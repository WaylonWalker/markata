"""Tests for FastAPI-owned persistent rendering without requiring a Go build."""
import asyncio
import os
import sys
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "python"))

from markata.pool import RendererPool
from markata.rendering import RenderError

# A minimal stdio peer verifies framing, process reuse and lifecycle.
FAKE_RENDERER = r'''
import json
import os
import sys
import time

for line in sys.stdin:
    request = json.loads(line)
    operation = request["op"]
    if operation == "capabilities":
        response = {"id": request["id"], "ok": True,
                    "capabilities": {"schema_version": 1}}
    elif operation == "render":
        source = request["source"]
        if source == "CRASH":
            sys.exit(7)
        if source == "SLOW":
            time.sleep(0.25)
        if source == "ERROR":
            response = {"id": request["id"], "ok": False,
                        "error": {"code": "render_failed", "message": "bad source"}}
        else:
            response = {"id": request["id"], "ok": True,
                        "result": {"schema_version": 1,
                                   "html": f"<p>{os.getpid()}:{source}</p>",
                                   "features": {"code": False},
                                   "assets": [], "external_urls": [],
                                   "head_html": "", "body_end_html": ""}}
    else:
        response = {"id": request["id"], "ok": True}
    print(json.dumps(response), flush=True)
'''


def pool(*, workers=1, max_pending=8):
    return RendererPool(workers=workers, max_pending=max_pending,
                        command=[sys.executable, "-u", "-c", FAKE_RENDERER])


class PoolTest(unittest.IsolatedAsyncioTestCase):
    async def test_reuses_warm_process(self):
        async with pool() as renderer:
            first = await renderer.render("first")
            second = await renderer.render("second")
            self.assertEqual(first.html.split(":")[0], second.html.split(":")[0])
            self.assertIn("second", second.html)
        with self.assertRaisesRegex(RenderError, "lifespan"):
            await renderer.render("closed")

    async def test_parallel_render_uses_multiple_processes(self):
        async with pool(workers=2) as renderer:
            a, b = await asyncio.gather(renderer.render("SLOW"), renderer.render("SLOW"))
            self.assertNotEqual(a.html.split(":")[0], b.html.split(":")[0])

    async def test_protocol_error_does_not_kill_worker(self):
        async with pool() as renderer:
            original = await renderer.render("before")
            with self.assertRaisesRegex(RenderError, "render_failed"):
                await renderer.render("ERROR")
            after = await renderer.render("after")
            self.assertEqual(original.html.split(":")[0], after.html.split(":")[0])

    async def test_crash_restarts_worker(self):
        async with pool() as renderer:
            with self.assertRaises(RenderError):
                await renderer.render("CRASH")
            recovered = await renderer.render("after", timeout=3)
            self.assertIn("after", recovered.html)

    async def test_timeout_restarts_worker(self):
        async with pool() as renderer:
            with self.assertRaisesRegex(RenderError, "timed out"):
                await renderer.render("SLOW", timeout=0.01)
            recovered = await renderer.render("after", timeout=3)
            self.assertIn("after", recovered.html)

    async def test_cancel_restarts_worker(self):
        async with pool() as renderer:
            task = asyncio.create_task(renderer.render("SLOW"))
            await asyncio.sleep(0.03)
            task.cancel()
            with self.assertRaises(asyncio.CancelledError):
                await task
            recovered = await renderer.render("after", timeout=3)
            self.assertIn("after", recovered.html)

    async def test_backpressure(self):
        async with pool(max_pending=0) as renderer:
            active = asyncio.create_task(renderer.render("SLOW"))
            await asyncio.sleep(0.03)
            with self.assertRaisesRegex(RenderError, "queue is full"):
                await renderer.render("overflow")
            await active

    async def test_size_limit(self):
        async with pool() as renderer:
            with self.assertRaisesRegex(RenderError, "200000"):
                await renderer.render("x" * 200001)

    async def test_subprocesses_reaped_at_shutdown(self):
        renderer = pool(workers=2)
        async with renderer:
            processes = [worker.process for worker in renderer._children]
            self.assertEqual(len(processes), 2)
        self.assertTrue(all(proc.returncode is not None for proc in processes))

    def test_bad_configuration(self):
        for kw in [{"workers": 0}, {"workers": 65}, {"max_pending": -1},
                   {"command": []}]:
            with self.assertRaises(ValueError):
                RendererPool(**kw)


if __name__ == "__main__":
    unittest.main()

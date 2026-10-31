"""An application-owned pool of persistent Markata renderer subprocesses.

The pool starts with the application's lifespan. Each Go process accepts
sequential JSONL requests; independent processes render concurrently.
"""
from __future__ import annotations

import asyncio
import json
from dataclasses import dataclass
from typing import Sequence

from .rendering import RenderError, RenderResult, _command, _result

_MAX_SOURCE_BYTES = 200_000
_MAX_RESPONSE_BYTES = 16 * 1024 * 1024


@dataclass(eq=False)
class _Worker:
    process: asyncio.subprocess.Process
    sequence: int = 0


class RendererPool:
    """Reusable async workers, owned by one application event loop.

    Use as an asynchronous context manager in a FastAPI lifespan or other
    application. Each Uvicorn process creates its own pool. max_pending bounds
    the number of callers waiting for a worker.
    """

    def __init__(self, *, workers: int = 2, max_pending: int = 64,
                 command: Sequence[str] | None = None):
        if not 1 <= workers <= 64:
            raise ValueError("workers must be between 1 and 64")
        if max_pending < 0:
            raise ValueError("max_pending must be nonnegative")
        if command is not None and not command:
            raise ValueError("command must be nonempty")
        self.workers = workers
        self.max_pending = max_pending
        self._command_override = tuple(command) if command is not None else None
        self._idle: asyncio.Queue[_Worker] = asyncio.Queue(maxsize=workers)
        self._slots = asyncio.Semaphore(workers + max_pending)
        self._children: set[_Worker] = set()
        self._repairs: set[asyncio.Task[None]] = set()
        self._started = False
        self._closed = False
        self._failure: str | None = None

    async def __aenter__(self) -> RendererPool:
        await self.start()
        return self

    async def __aexit__(self, *_exc: object) -> None:
        await self.aclose()

    async def start(self) -> None:
        """Spawn and verify all workers before admitting renders."""
        if self._started:
            raise RuntimeError("RendererPool already started")
        self._started = True
        try:
            for _ in range(self.workers):
                worker = await self._spawn()
                self._children.add(worker)
                self._idle.put_nowait(worker)
        except BaseException:
            await self.aclose()
            raise

    async def _spawn(self) -> _Worker:
        command = self._command_override
        if command is None:
            command = (_command(False, "/markata-assets", False)[0], "render", "--stdio")
        try:
            process = await asyncio.create_subprocess_exec(
                *command, stdin=asyncio.subprocess.PIPE,
                stdout=asyncio.subprocess.PIPE, stderr=asyncio.subprocess.DEVNULL,
                limit=_MAX_RESPONSE_BYTES,
            )
        except OSError as error:
            raise RenderError(f"Cannot start persistent Markata renderer: {error}") from error
        worker = _Worker(process)
        try:
            response = await asyncio.wait_for(
                self._exchange(worker, "capabilities"), timeout=5.0
            )
            caps = response.get("capabilities")
            if not response.get("ok") or not isinstance(caps, dict) or caps.get("schema_version") != 1:
                raise RenderError("Bundled Markata renderer does not support the v1 persistent stdio protocol")
            return worker
        except BaseException:
            self._kill(worker)
            await process.wait()
            raise

    @staticmethod
    async def _exchange(worker: _Worker, op: str, *,
                        source: str = "", options: dict | None = None) -> dict:
        worker.sequence += 1
        request_id = str(worker.sequence)
        request: dict = {"id": request_id, "op": op}
        if op == "render":
            request.update(source=source, options=options or {})
        data = (json.dumps(request, ensure_ascii=True, separators=(",", ":")) + "\n").encode("utf-8")
        if len(data) > 1 << 20:
            raise RenderError("Render request exceeds the 1 MiB protocol frame limit")
        if worker.process.stdin is None or worker.process.stdout is None:
            raise RenderError("Renderer process has no stdio pipes")
        try:
            worker.process.stdin.write(data)
            await worker.process.stdin.drain()
            line = await worker.process.stdout.readline()
            if not line:
                raise RenderError("Persistent Markata renderer terminated unexpectedly")
            response = json.loads(line)
        except (BrokenPipeError, ConnectionError, OSError, ValueError, asyncio.LimitOverrunError) as error:
            raise RenderError(f"Persistent renderer protocol failure: {error}") from error
        if not isinstance(response, dict) or response.get("id") != request_id or not isinstance(response.get("ok"), bool):
            raise RenderError("Persistent renderer returned a mismatched or invalid response")
        return response

    async def render(self, markdown: str, *, inline_assets: bool = False,
                     asset_base: str = "/markata-assets", unsafe: bool = False,
                     timeout: float = 30) -> RenderResult:
        """Render on a warm worker; cancellation/timeout replaces that worker."""
        if not self._started or self._closed:
            raise RenderError("RendererPool must be started inside an async lifespan")
        if self._failure is not None:
            raise RenderError(f"Renderer pool cannot recover a worker: {self._failure}")
        if timeout <= 0:
            raise ValueError("timeout must be positive")
        if not isinstance(markdown, str):
            raise TypeError("markdown must be a string")
        if len(markdown.encode("utf-8")) > _MAX_SOURCE_BYTES:
            raise RenderError("Markdown exceeds 200000 UTF-8 bytes")
        if self._slots.locked():
            raise RenderError("Renderer pool queue is full")
        await self._slots.acquire()
        try:
            try:
                async with asyncio.timeout(timeout):
                    worker = await self._idle.get()
                    try:
                        response = await self._exchange(
                            worker, "render", source=markdown,
                            options={"inline_assets": inline_assets,
                                     "asset_base": asset_base, "unsafe": unsafe},
                        )
                    except BaseException:
                        self._retire(worker)
                        raise
                    else:
                        self._idle.put_nowait(worker)
            except TimeoutError as error:
                raise RenderError(f"Markata render timed out after {timeout}s") from error
            if not response["ok"]:
                problem = response.get("error") or {}
                raise RenderError(f"Markata render failed ({problem.get('code', 'unknown')}): {problem.get('message', '')}")
            data = response.get("result")
            if not isinstance(data, dict):
                raise RenderError("Persistent renderer did not return a result object")
            return _result(json.dumps(data).encode("utf-8"), b"", 0)
        finally:
            self._slots.release()

    def _retire(self, worker: _Worker) -> None:
        self._children.discard(worker)
        self._kill(worker)
        if not self._closed:
            task = asyncio.create_task(self._replace(worker))
            self._repairs.add(task)
            task.add_done_callback(self._repairs.discard)

    async def _replace(self, worker: _Worker) -> None:
        try:
            await worker.process.wait()
            if self._closed:
                return
            new_worker = await self._spawn()
            if self._closed:
                self._kill(new_worker)
                await new_worker.process.wait()
            else:
                self._children.add(new_worker)
                self._idle.put_nowait(new_worker)
        except Exception as error:
            self._failure = str(error)

    @staticmethod
    def _kill(worker: _Worker) -> None:
        if worker.process.returncode is None:
            try:
                worker.process.kill()
            except ProcessLookupError:
                pass

    async def aclose(self) -> None:
        """Stop child processes and any pending replacement tasks."""
        if self._closed:
            return
        self._closed = True
        for task in list(self._repairs):
            task.cancel()
        for worker in self._children:
            self._kill(worker)
        await asyncio.gather(
            *(worker.process.wait() for worker in self._children),
            *self._repairs, return_exceptions=True,
        )
        self._children.clear()

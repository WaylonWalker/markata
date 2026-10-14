#!/usr/bin/env -S uv run --script
# /// script
# requires-python = ">=3.11"
# dependencies = ["playwright>=1.50,<2"]
# ///
"""Compare *rendering* under slow network conditions using CDP filmstrip frames.

Examples:
  uv run scripts/slow3g_compare.py --before http://127.0.0.1:8000/archive/ --after http://127.0.0.1:8000/archive/?variant=moved
  uv run scripts/slow3g_compare.py --before http://127.0.0.1:8000/archive/ --runs 5
  uv run scripts/slow3g_compare.py --before fixture:/abs/path/demo.html --runs 1 --observe-ms 3500

Important: filmstrip frames are timestamped at *receipt*. They are never labeled
with the checkpoint requested by an asynchronous page.screenshot() call.
Review actual screenshots to establish when main content became readable.
"""

from __future__ import annotations

import argparse
import asyncio
import base64
import csv
import json
import math
import os
from pathlib import Path
import statistics
import shutil
import time
from typing import Any

from playwright.async_api import async_playwright

CHECKPOINTS = (1000, 2000, 3000, 5000, 8000, 10000, 15000, 20000, 25000)
INIT_PERF = r"""
(() => {
  const metrics = { fcp: null, lcp: null, lcpTag: null, lcpURL: null, cls: 0 };
  Object.defineProperty(window, '__slow3gMetrics', { value: metrics });
  try {
    new PerformanceObserver(list => {
      for (const entry of list.getEntries()) {
        if (entry.name === 'first-contentful-paint') metrics.fcp = entry.startTime;
      }
    }).observe({ type: 'paint', buffered: true });
  } catch (error) { metrics.paintObserverError = String(error); }
  try {
    new PerformanceObserver(list => {
      for (const entry of list.getEntries()) {
        metrics.lcp = entry.startTime;
        metrics.lcpTag = entry.element?.tagName || null;
        metrics.lcpURL = entry.url || null;
      }
    }).observe({ type: 'largest-contentful-paint', buffered: true });
  } catch (error) { metrics.lcpObserverError = String(error); }
  try {
    new PerformanceObserver(list => {
      for (const entry of list.getEntries()) if (!entry.hadRecentInput) metrics.cls += entry.value;
    }).observe({ type: 'layout-shift', buffered: true });
  } catch (error) { metrics.clsObserverError = String(error); }
})();
"""


def ms(value: float | None) -> float | None:
    return round(value, 2) if value is not None else None


def asset_response_problem(kind: str, mime: str, status: int) -> str | None:
    """Flag a broken essential asset even if HTTP returned 200 with HTML."""
    if kind not in {"Stylesheet", "Script", "Image", "Font"}:
        return None
    if status >= 400:
        return f"HTTP {status}"
    content_type = (mime or "").lower().split(";", 1)[0].strip()
    if kind == "Stylesheet" and content_type != "text/css":
        return f"CSS returned {content_type or 'unknown MIME'}"
    if kind == "Script" and not ("javascript" in content_type or "ecmascript" in content_type):
        return f"JavaScript returned {content_type or 'unknown MIME'}"
    if kind == "Image" and not content_type.startswith("image/"):
        return f"image returned {content_type or 'unknown MIME'}"
    if kind == "Font" and not (
        content_type.startswith("font/")
        or content_type in {"application/font-woff", "application/x-font-woff", "application/octet-stream"}
    ):
        return f"font returned {content_type or 'unknown MIME'}"
    return None


def percentile_nearest_rank(values: list[float], percentile: float) -> float | None:
    if not values:
        return None
    sorted_values = sorted(values)
    return sorted_values[max(0, math.ceil(len(sorted_values) * percentile) - 1)]


async def one_run(browser: Any, args: argparse.Namespace, variant: str, run: int, url: str) -> dict[str, Any]:
    prefix = f"{variant}-{run:02d}"
    context = await browser.new_context(
        viewport={"width": args.width, "height": args.height},
        device_scale_factor=1,
        is_mobile=True,
        has_touch=True,
        ignore_https_errors=args.ignore_https_errors,
        service_workers="block",
        java_script_enabled=not args.no_js,
    )
    page = await context.new_page()
    page.set_default_timeout(args.timeout_ms)
    if not args.no_js:
        await context.add_init_script(script=INIT_PERF)
    cdp = await context.new_cdp_session(page)
    await cdp.send("Network.enable")
    await cdp.send("Network.setCacheDisabled", {"cacheDisabled": True})
    await cdp.send("Network.emulateNetworkConditions", {
        "offline": False,
        "latency": args.latency_ms,
        "downloadThroughput": args.download_bytes_per_sec,
        "uploadThroughput": args.upload_bytes_per_sec,
    })
    await cdp.send("Emulation.setCPUThrottlingRate", {"rate": args.cpu_slowdown})
    requests: dict[str, dict[str, Any]] = {}
    failures: list[dict[str, str]] = []
    invalid_assets: list[dict[str, Any]] = []
    document_delivery: dict[str, Any] = {}
    frames: list[dict[str, Any]] = []
    latest_frame: dict[str, Any] | None = None
    received_frame_count = 0
    filmstrip_dir = args.output / "filmstrips" / prefix
    filmstrip_dir.mkdir(parents=True, exist_ok=True)
    nav_started = None
    pending_acks: set[asyncio.Task] = set()

    def on_request(e: dict[str, Any]) -> None:
        requests[e["requestId"]] = {
            "url": e.get("request", {}).get("url", ""),
            "type": e.get("type", "Other"),
            "bytesDecodedChunks": 0,
            "bytesEncodedChunks": 0,
            "bytesEncodedFinished": None,
        }

    def on_data(e: dict[str, Any]) -> None:
        r = requests.get(e["requestId"])
        if r is not None:
            r["bytesDecodedChunks"] += e.get("dataLength", 0)
            r["bytesEncodedChunks"] += e.get("encodedDataLength", 0)

    def on_response(e: dict[str, Any]) -> None:
        kind = e.get("type", "")
        response = e.get("response", {})
        status = int(response.get("status", 0))
        mime = response.get("mimeType", "")
        request = requests.get(e.get("requestId"), {})
        request["http_status"] = status
        request["mime_type"] = mime
        if kind == "Document":
            headers = {str(k).lower(): v for k, v in response.get("headers", {}).items()}
            document_delivery.update({
                "status": status,
                "protocol": response.get("protocol"),
                "content_encoding": headers.get("content-encoding"),
                "mime_type": mime,
            })
        problem = asset_response_problem(kind, mime, status)
        if problem:
            invalid_assets.append({
                "url": request.get("url", "")[:400],
                "type": kind, "status": status, "mime_type": mime, "problem": problem,
            })

    def on_finished(e: dict[str, Any]) -> None:
        r = requests.get(e["requestId"])
        if r is not None:
            r["bytesEncodedFinished"] = e.get("encodedDataLength")

    def on_failed(e: dict[str, Any]) -> None:
        r = requests.get(e["requestId"], {})
        failures.append({"url": r.get("url", ""), "error": e.get("errorText", "unknown")})

    cdp.on("Network.requestWillBeSent", on_request)
    cdp.on("Network.dataReceived", on_data)
    cdp.on("Network.responseReceived", on_response)
    cdp.on("Network.loadingFinished", on_finished)
    cdp.on("Network.loadingFailed", on_failed)

    async def acknowledge(session_id: int) -> None:
        try:
            await cdp.send("Page.screencastFrameAck", {"sessionId": session_id})
        except Exception:
            # Browser/CDP can disconnect during teardown; the run carries the
            # resulting frame count and error status instead of crashing the loop.
            pass

    def on_frame(e: dict[str, Any]) -> None:
        nonlocal latest_frame, received_frame_count
        now = time.monotonic()
        received_frame_count += 1
        if nav_started is not None:
            latest_frame = {
                "data": e["data"],
                "received_at_ms": ms((now - nav_started) * 1000),
                "cdp_frame_timestamp": e.get("metadata", {}).get("timestamp"),
            }
        # Pyppeteer failed by passing a bare Future to create_task; always
        # schedule our async coroutine instead.
        task = asyncio.create_task(acknowledge(e["sessionId"]))
        pending_acks.add(task)
        task.add_done_callback(pending_acks.discard)

    async def snapshot_at_checkpoint(target: int) -> None:
        if nav_started is None:
            return
        remaining = target / 1000 - (time.monotonic() - nav_started)
        if remaining > 0:
            await asyncio.sleep(remaining)
        actual_ms = ms((time.monotonic() - nav_started) * 1000)
        current = latest_frame
        if not current:
            frames.append({"requested_checkpoint_ms": target, "snapshot_at_ms": actual_ms,
                           "error": "no CDP frame has arrived by checkpoint"})
            return
        filename = f"checkpoint-{target}ms-snapshot-{actual_ms:.0f}ms.jpg"
        try:
            (filmstrip_dir / filename).write_bytes(base64.b64decode(current["data"]))
            frames.append({"requested_checkpoint_ms": target,
                           "snapshot_at_ms": actual_ms,
                           "frame_received_at_ms": current["received_at_ms"],
                           "cdp_frame_timestamp": current["cdp_frame_timestamp"],
                           "file": str(Path("filmstrips") / prefix / filename)})
        except Exception as error:
            frames.append({"requested_checkpoint_ms": target, "error": str(error)})

    cdp.on("Page.screencastFrame", on_frame)
    error: str | None = None
    response_status: int | None = None
    metrics: dict[str, Any] = {}
    try:
        await cdp.send("Page.enable")
        await cdp.send("Page.startScreencast", {"format": "jpeg", "quality": 60, "everyNthFrame": 1})
        nav_started = time.monotonic()
        checkpoint_tasks = [asyncio.create_task(snapshot_at_checkpoint(t)) for t in args.checkpoints if t <= args.observe_ms]
        try:
            if url.startswith("fixture:"):
                # Offline harness self-test: bypasses network and does NOT test
                # HTTP transfer / throughput. Useful for checking CDP behavior.
                await page.set_content(Path(url.removeprefix("fixture:")).read_text(), wait_until="domcontentloaded")
            else:
                response = await page.goto(url, wait_until="commit", timeout=args.timeout_ms)
                response_status = response.status if response else None
                if response_status is not None and response_status >= 400:
                    error = f"navigation: HTTP {response_status}"
        except Exception as exc:
            error = f"navigation: {exc}"
        # A fixed observation window, regardless of onload or networkidle.
        remaining = args.observe_ms - (time.monotonic() - nav_started) * 1000
        if remaining > 0:
            await asyncio.sleep(remaining / 1000)
        if checkpoint_tasks:
            await asyncio.gather(*checkpoint_tasks, return_exceptions=True)
        if not args.no_js:
            try:
                metrics = await asyncio.wait_for(
                    page.evaluate("() => window.__slow3gMetrics || {}"),
                    timeout=min(10, args.timeout_ms / 1000),
                )
            except Exception as exc:
                error = (error + "; " if error else "") + f"metrics: {exc}"
    finally:
        try:
            await cdp.send("Page.stopScreencast")
        except Exception:
            pass
        if pending_acks:
            try:
                await asyncio.wait_for(asyncio.gather(*pending_acks, return_exceptions=True), timeout=3)
            except asyncio.TimeoutError:
                for task in pending_acks:
                    task.cancel()
        try:
            await context.close()
        except Exception:
            pass

    # Resource totals are *completed* resources, not bytes at paint. Actual
    # readability timestamps need screenshot review; the CSV includes FCP/LCP only.
    saved_frames = [f for f in frames if f.get("file")]
    return {
        "variant": variant,
        "run": run,
        "url": url,
        "http_status": response_status,
        "metrics": metrics,
        "frames": frames,
        "frame_count": len(saved_frames),
        "received_cdp_frame_count": received_frame_count,
        "missing_checkpoints_ms": [f["requested_checkpoint_ms"] for f in frames if not f.get("file")],
        "requests": list(requests.values()),
        "network_failures": failures,
        "invalid_asset_responses": invalid_assets,
        "document_delivery": document_delivery,
        "encoded_bytes_completed": sum(int(v["bytesEncodedFinished"] or 0) for v in requests.values()),
        "error": error,
        # A short smoke test may request fewer checkpoints than the usual 5.
        "valid_filmstrip": (
            len(saved_frames) >= min(args.min_frames, sum(t <= args.observe_ms for t in args.checkpoints))
            and not error and not invalid_assets and not failures
        ),
    }


def write_report(results: list[dict[str, Any]], output: Path, args: argparse.Namespace) -> None:
    (output / "results.json").write_text(json.dumps(results, indent=2), encoding="utf8")
    summary: list[dict[str, Any]] = []
    for variant in ("before", "after"):
        subset = [r for r in results if r["variant"] == variant]
        if not subset:
            continue
        usable = [r for r in subset if r["valid_filmstrip"]]
        fcps = [float(r["metrics"]["fcp"]) for r in usable if r["metrics"].get("fcp") is not None]
        lcps = [float(r["metrics"]["lcp"]) for r in usable if r["metrics"].get("lcp") is not None]
        clss = [float(r["metrics"]["cls"]) for r in usable if r["metrics"].get("cls") is not None]
        encoded_bytes = [float(r["encoded_bytes_completed"]) for r in usable]
        request_counts = [float(len(r.get("requests", []))) for r in usable]
        summary.append({
            "variant": variant,
            "valid_filmstrips": len(usable),
            "attempted": len(subset),
            "failed_network_requests": sum(len(r.get("network_failures", [])) for r in subset),
            "invalid_asset_responses": sum(len(r.get("invalid_asset_responses", [])) for r in subset),
            "median_fcp_ms": ms(statistics.median(fcps)) if fcps else None,
            "p90_fcp_ms": ms(percentile_nearest_rank(fcps, .9)),
            "median_lcp_ms": ms(statistics.median(lcps)) if lcps else None,
            "p90_lcp_ms": ms(percentile_nearest_rank(lcps, .9)),
            "median_cls": round(statistics.median(clss), 6) if clss else None,
            "p90_cls": round(percentile_nearest_rank(clss, .9), 6) if clss else None,
            "median_encoded_bytes_completed": ms(statistics.median(encoded_bytes)) if encoded_bytes else None,
            "p90_encoded_bytes_completed": ms(percentile_nearest_rank(encoded_bytes, .9)),
            "median_request_count": ms(statistics.median(request_counts)) if request_counts else None,
            "p90_request_count": ms(percentile_nearest_rank(request_counts, .9)),
            "note": "readable content requires filmstrip review; no inferred timing",
        })
    (output / "summary.json").write_text(json.dumps(summary, indent=2), encoding="utf8")
    with (output / "summary.csv").open("w", newline="", encoding="utf8") as fh:
        writer = csv.DictWriter(fh, fieldnames=list(summary[0].keys()) if summary else ["variant"])
        writer.writeheader()
        writer.writerows(summary)

    from html import escape
    lines = ["<!doctype html><html><head><meta charset='utf-8'><style>",
             "body{font:14px system-ui;background:#141921;color:#f0f0f0;padding:18px}",
             "article{margin:12px 0 24px;border:1px solid #59606b;padding:10px}",
             ".frames{display:flex;flex-wrap:wrap;gap:8px}.frame{width:175px}",
             "img{width:175px;height:auto;border:1px solid #555}.meta{color:#c6cede}",
             "</style></head><body><h1>Slow 3G filmstrip</h1>",
             "<p>Frame labels show the <strong>actual receipt time</strong>, not an assumed screenshot time.</p>",
             "<p>Inspect images manually for first <em>readable main content</em>.</p>"]
    for r in results:
        lines.append(f"<article><h2>{escape(r['variant'])} run {r['run']} – {len(r['frames'])} frames</h2>")
        lines.append(f"<p class='meta'>FCP: {r['metrics'].get('fcp')} ms; LCP: {r['metrics'].get('lcp')} ms; "
                     f"CLS: {r['metrics'].get('cls')}; error: {escape(str(r['error']))}</p><div class='frames'>")
        for frame in r["frames"]:
            if frame.get("file"):
                lines.append(f"<div class='frame'><img loading='lazy' src='{escape(frame['file'])}' alt='captured frame'>"
                             f"<p>Target {frame['requested_checkpoint_ms']}ms; snap {frame['snapshot_at_ms']}ms; last rendered frame arrived {frame['frame_received_at_ms']}ms</p></div>")
        lines.append("</div></article>")
    lines.append("</body></html>")
    (output / "filmstrip.html").write_text("\n".join(lines), encoding="utf8")


async def main(args: argparse.Namespace) -> int:
    args.output.mkdir(parents=True, exist_ok=True)
    results: list[dict[str, Any]] = []
    async with async_playwright() as p:
        browser = await p.chromium.launch(
            executable_path=args.chromium,
            headless=True,
            args=["--no-sandbox", "--disable-dev-shm-usage"],
        )
        try:
            for run in range(1, args.runs + 1):
                # Counterbalance ordering: A/B on odd pairs, B/A on even.
                # An always-before-first schedule can bias results if the host
                # warms up, throttles, or becomes resource constrained.
                pair = [("before", args.before), ("after", args.after)]
                if run % 2 == 0:
                    pair.reverse()
                for variant, url in pair:
                    if not url:
                        continue
                    print(f"[{variant} {run}/{args.runs}] {url}", flush=True)
                    try:
                        result = await one_run(browser, args, variant, run, url)
                    except Exception as error:
                        result = {"variant": variant, "run": run, "url": url, "error": str(error),
                                  "metrics": {}, "frames": [], "valid_filmstrip": False,
                                  "network_failures": []}
                    results.append(result)
                    write_report(results, args.output, args)
                    print(f"  frames={len(result['frames'])}, FCP={result['metrics'].get('fcp')}, "
                          f"error={result['error']}", flush=True)
        finally:
            await browser.close()
    if not any(r["valid_filmstrip"] for r in results):
        return 2
    if any(not r["valid_filmstrip"] for r in results):
        return 1
    return 0


if __name__ == "__main__":
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--before", required=True)
    ap.add_argument("--after", default=None)
    ap.add_argument("--runs", type=int, default=5)
    ap.add_argument("--output", type=Path, default=Path("slow3g-results"))
    ap.add_argument("--chromium", default=os.getenv("CHROMIUM_PATH", shutil.which("chromium")),
                    help="Chromium executable; defaults to system chromium if found, else Playwright-managed browser")
    ap.add_argument("--width", type=int, default=390)
    ap.add_argument("--height", type=int, default=844)
    ap.add_argument("--latency-ms", type=int, default=400)
    ap.add_argument("--download-bytes-per-sec", type=int, default=51200)
    ap.add_argument("--upload-bytes-per-sec", type=int, default=51200)
    ap.add_argument("--cpu-slowdown", type=float, default=4)
    ap.add_argument("--observe-ms", type=int, default=26000)
    ap.add_argument("--timeout-ms", type=int, default=35000)
    ap.add_argument("--min-frames", type=int, default=5)
    ap.add_argument("--checkpoints", type=int, nargs="+", default=CHECKPOINTS)
    ap.add_argument("--ignore-https-errors", action="store_true")
    ap.add_argument("--no-js", action="store_true", help="No-JS filmstrip only; performance observers unavailable")
    parsed = ap.parse_args()
    if parsed.runs < 1:
        ap.error("--runs must be >= 1")
    if not parsed.checkpoints or any(t <= 0 for t in parsed.checkpoints):
        ap.error("--checkpoints must contain positive millisecond offsets")
    if parsed.min_frames < 1:
        ap.error("--min-frames must be >= 1")
    if parsed.observe_ms < 1:
        ap.error("--observe-ms must be positive")
    raise SystemExit(asyncio.run(main(parsed)))

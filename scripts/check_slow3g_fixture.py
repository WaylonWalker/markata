#!/usr/bin/env python3
"""Assert the Slow 3G fixture actually fetched good gzip HTML and subresources."""

import argparse
import json
from pathlib import Path

parser = argparse.ArgumentParser()
parser.add_argument("result_dir", type=Path)
parser.add_argument("--expect-bad-css", action="store_true")
args = parser.parse_args()
runs = json.loads((args.result_dir / "results.json").read_text())

if args.expect_bad_css:
    assert len(runs) == 1, f"expected one bad-CSS run, got {len(runs)}"
    run = runs[0]
    assert not run.get("valid_filmstrip"), "broken stylesheet was accepted"
    issues = run.get("invalid_asset_responses", [])
    assert any(x.get("type") == "Stylesheet" and x.get("mime_type") == "text/html" for x in issues), issues
    print("Broken-CSS fixture correctly rejected.")
else:
    assert len(runs) == 2, f"expected one baseline/candidate pair, got {len(runs)}"
    for run in runs:
        assert run.get("valid_filmstrip"), (run.get("error"), run.get("invalid_asset_responses"), run.get("frames"))
        assert run.get("http_status") == 200, run.get("http_status")
        assert run.get("document_delivery", {}).get("content_encoding") == "gzip", run.get("document_delivery")
        assert not run.get("invalid_asset_responses"), run.get("invalid_asset_responses")
        assert not run.get("network_failures"), run.get("network_failures")
        kinds = {r.get("type") for r in run.get("requests", [])}
        assert {"Document", "Stylesheet", "Script", "Image"} <= kinds, kinds
        assert run.get("frame_count", 0) >= 3, run.get("frame_count")
    print("Both gzip HTTP fixture runs served real assets and valid CDP filmstrips.")

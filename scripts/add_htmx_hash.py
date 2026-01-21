#!/usr/bin/env python3
"""
Development script to add HTMX integrity hashes to feeds plugin.

This script should be used by developers when a new HTMX version is released.
It will:
1. Download/Get specified HTMX version
2. Calculate its SHA-256 hash
3. Add it to HTMX_INTEGRITY_HASHES dictionary in feeds.py

Usage:
    python scripts/add_htmx_hash.py 2.0.8
    python scripts/add_htmx_hash.py --list
    python scripts/add_htmx_hash.py --add-all
    python scripts/add_htmx_hash.py 1.9.10 --replace
"""

import argparse
import hashlib
import json
import re
from pathlib import Path
from typing import Optional
from urllib.request import Request, urlopen
from urllib.error import URLError, HTTPError


def get_htmx_versions():
    """Get list of all available HTMX versions from GitHub releases API."""
    try:
        url = "https://api.github.com/repos/bigskysoftware/htmx/releases"
        request = Request(url, headers={"User-Agent": "Markata-Dev/1.0"})
        with urlopen(request, timeout=10) as response:
            releases_data = json.loads(response.read().decode("utf-8"))

        # Extract version numbers from tag names (remove 'v' prefix)
        versions = []
        for release in releases_data:
            if "tag_name" in release and release["tag_name"].startswith("v"):
                version = release["tag_name"][1:]  # Remove 'v' prefix
                versions.append(version)

        # Filter out duplicate and sort by semantic version
        unique_versions = list(set(versions))
        unique_versions.sort(
            key=lambda v: [int(x) for x in re.findall(r"\d+", v)], reverse=True
        )

        return unique_versions
    except Exception as e:
        print(f"Error fetching HTMX versions: {e}")
        return []


def get_htmx_hash_from_github(version: str) -> Optional[str]:
    """Get HTMX hash directly from GitHub releases API."""
    try:
        url = "https://api.github.com/repos/bigskysoftware/htmx/releases"
        request = Request(url, headers={"User-Agent": "Markata-Dev/1.0"})
        with urlopen(request, timeout=10) as response:
            releases_data = json.loads(response.read().decode("utf-8"))

        # Find release with matching version
        for release in releases_data:
            if release["tag_name"] == f"v{version}":
                # Look for htmx.min.js asset
                for asset in release.get("assets", []):
                    if asset["name"] == "htmx.min.js":
                        # Extract hash from digest (remove 'sha256:' prefix)
                        digest = asset.get("digest", "")
                        if digest.startswith("sha256:"):
                            return digest[7:]  # Remove 'sha256:' prefix
                break

        return None
    except Exception as e:
        print(f"Error fetching hash from GitHub: {e}")
        return None


def add_htmx_hash(version: str, replace: bool = False, verbose: bool = False) -> bool:
    """Add HTMX version hash to feeds.py. Returns True if successful."""
    # First try to get hash from GitHub API (more reliable)
    sha256_hash = get_htmx_hash_from_github(version)

    if not sha256_hash:
        # Fall back to downloading from unpkg.com
        try:
            url = f"https://unpkg.com/htmx.org@{version}/dist/htmx.min.js"

            if verbose:
                print(f"Downloading HTMX {version} from {url}")

            request = Request(url, headers={"User-Agent": "Markata-Dev/1.0"})
            with urlopen(request, timeout=10) as response:
                content = response.read()
                sha256_hash = hashlib.sha256(content).hexdigest()
                if verbose:
                    print(f"SHA-256 hash: {sha256_hash}")
        except (URLError, HTTPError) as e:
            print(f"Error: Failed to download HTMX: {e}")
            return False
    else:
        if verbose:
            print(f"Got HTMX {version} hash from GitHub API")

    try:
        # Find and update the feeds.py file
        project_root = Path(__file__).parent.parent
        feeds_file = project_root / "markata" / "plugins" / "feeds.py"

        if not feeds_file.exists():
            print(f"Error: Could not find feeds.py at {feeds_file}")
            return False

        with open(feeds_file, "r") as f:
            file_content = f.read()

        # Find HTMX_INTEGRITY_HASHES dictionary
        pattern = r"(HTMX_INTEGRITY_HASHES = \{[^}]+)}"
        match = re.search(pattern, file_content, re.DOTALL)

        if not match:
            print("Error: Could not find HTMX_INTEGRITY_HASHES in feeds.py")
            return False

        # Add new hash
        new_hash_entry = f'        "{version}": "{sha256_hash}"'
        existing_dict = match.group(1)

        # Check if version already exists
        if f'"{version}":' in existing_dict:
            if verbose:
                print(f"Warning: HTMX version {version} already exists in hashes")
            if not replace:
                response = input("Replace existing hash? [y/N]: ")
                if response.lower() != "y":
                    print("Cancelled.")
                    return False

            # Replace existing entry
            new_dict = re.sub(
                rf'        "{version}": "[^"]*"', new_hash_entry, existing_dict
            )
        else:
            # Add new entry (before the closing brace)
            new_dict = existing_dict.rstrip() + f",\n{new_hash_entry}"

        # Update the file
        updated_content = file_content.replace(match.group(0), new_dict + "}")

        with open(feeds_file, "w") as f:
            f.write(updated_content)

        print(f"✅ Added HTMX {version} hash to feeds.py")
        if verbose:
            print(f"📝 File updated: {feeds_file}")
        return True

    except Exception as e:
        print(f"Error: {e}")
        return False


def add_all_htmx_versions(verbose: bool = False, skip_failed: bool = False) -> None:
    """Add hashes for all available HTMX versions."""
    versions = get_htmx_versions()
    if not versions:
        print("Could not fetch HTMX versions")
        return

    print(f"Found {len(versions)} HTMX versions")

    # Get existing versions to avoid duplicates
    project_root = Path(__file__).parent.parent
    feeds_file = project_root / "markata" / "plugins" / "feeds.py"

    with open(feeds_file, "r") as f:
        file_content = f.read()

    pattern = r"HTMX_INTEGRITY_HASHES = \{([^}]+)}"
    match = re.search(pattern, file_content, re.DOTALL)
    existing_versions = set()
    if match:
        existing_matches = re.findall(r'"([^"]+)":', match.group(1))
        existing_versions = set(existing_matches)

    if verbose:
        print(f"Existing versions: {sorted(existing_versions)}")

    # Filter out existing versions
    new_versions = [v for v in versions if v not in existing_versions]

    if not new_versions:
        print("All available versions already have hashes!")
        return

    print(f"Adding {len(new_versions)} new versions...")

    success_count = 0
    for version in new_versions:
        if verbose:
            print(f"\nProcessing {version}...")

        success = add_htmx_hash(version, replace=True, verbose=False)
        if success:
            success_count += 1
        elif not skip_failed:
            print(f"Failed to add {version}, stopping. Use --skip-failed to continue.")
            break

    print(f"\n✅ Successfully added {success_count}/{len(new_versions)} versions")


def list_htmx_versions(verbose: bool = False) -> None:
    """List all available HTMX versions."""
    versions = get_htmx_versions()
    if not versions:
        print("Could not fetch HTMX versions")
        return

    # Get existing versions
    project_root = Path(__file__).parent.parent
    feeds_file = project_root / "markata" / "plugins" / "feeds.py"

    with open(feeds_file, "r") as f:
        file_content = f.read()

    pattern = r"HTMX_INTEGRITY_HASHES = \{([^}]+)}"
    match = re.search(pattern, file_content, re.DOTALL)
    existing_versions = set()
    if match:
        existing_matches = re.findall(r'"([^"]+)":', match.group(1))
        existing_versions = set(existing_matches)

    print("HTMX Versions:")
    print("=" * 50)

    for version in versions[:20]:  # Show first 20 to avoid too much output
        status = "✅" if version in existing_versions else "❌"
        print(f"  {status} {version}")

    if len(versions) > 20:
        print(f"  ... and {len(versions) - 20} more versions")

    print(
        f"\nSummary: {len(existing_versions)} versions have hashes, {len(versions) - len(existing_versions)} missing"
    )

    if verbose:
        print(f"\nAll versions: {versions}")
        print(f"Existing versions: {sorted(existing_versions)}")
        print(
            f"Missing versions: {[v for v in versions if v not in existing_versions]}"
        )


def main():
    parser = argparse.ArgumentParser(
        description="Add HTMX integrity hash to feeds.py",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__,
    )

    parser.add_argument("version", nargs="?", help="HTMX version (e.g., 2.0.8, 1.9.10)")

    parser.add_argument(
        "--list",
        action="store_true",
        help="List all available HTMX versions and their hash status",
    )

    parser.add_argument(
        "--add-all",
        action="store_true",
        help="Add hashes for all missing HTMX versions",
    )

    parser.add_argument(
        "--replace", action="store_true", help="Replace existing hash without prompting"
    )

    parser.add_argument("--verbose", action="store_true", help="Show detailed output")

    parser.add_argument(
        "--skip-failed",
        action="store_true",
        help="Continue adding versions even if some fail (used with --add-all)",
    )

    args = parser.parse_args()

    if args.list:
        list_htmx_versions(args.verbose)
    elif args.add_all:
        add_all_htmx_versions(args.verbose, args.skip_failed)
    elif args.version:
        success = add_htmx_hash(args.version, args.replace, args.verbose)
        if not success:
            exit(1)
    else:
        parser.print_help()
        print("\nExamples:")
        print("  python add_htmx_hash.py 2.0.8              # Add specific version")
        print("  python add_htmx_hash.py --list             # List all versions")
        print(
            "  python add_htmx_hash.py --add-all           # Add all missing versions"
        )
        print("  python add_htmx_hash.py 2.0.8 --replace   # Replace existing hash")


if __name__ == "__main__":
    main()

# Development Scripts

This directory contains utility scripts for Markata developers.

## HTMX Hash Management

### `add_htmx_hash.py`

A utility script to add HTMX integrity hashes to feeds plugin when new HTMX versions are released.

**Usage:**
```bash
# Add a new HTMX version hash
python scripts/add_htmx_hash.py 2.0.8

# List all available versions and their hash status
python scripts/add_htmx_hash.py --list

# Add hashes for all missing versions
python scripts/add_htmx_hash.py --add-all

# Replace an existing hash without prompting
python scripts/add_htmx_hash.py 2.0.8 --replace

# Show detailed output
python scripts/add_htmx_hash.py --list --verbose

# Add all versions, skipping failed ones
python scripts/add_htmx_hash.py --add-all --skip-failed

# Show help
python scripts/add_htmx_hash.py --help
```

**Features:**
- **Smart Hash Retrieval**: First tries GitHub API for official hashes, falls back to unpkg.com
- **Version Management**: Lists all available HTMX versions from GitHub releases
- **Batch Operations**: Add all missing versions with `--add-all`
- **Status Tracking**: See which versions have hashes and which are missing
- **Safety**: Prompts before replacing existing hashes (unless `--replace` used)
- **Verbose Mode**: Detailed output for debugging and monitoring

**What it does:**
1. Fetches all HTMX versions from GitHub releases API
2. For single versions: Gets hash from GitHub API or downloads from unpkg.com and calculates SHA-256
3. Updates `HTMX_INTEGRITY_HASHES` dictionary in `markata/plugins/feeds.py`
4. For batch operations: Processes all missing versions automatically

**When to use:**
- **Single Version**: When a new HTMX version is released
- **List Mode**: To see current hash coverage and available versions
- **Batch Mode**: To populate hashes for many versions at once
- **Development**: When setting up a new development environment

**Examples:**
```bash
# Quick check of current status
python scripts/add_htmx_hash.py --list

# Add the latest version
python scripts/add_htmx_hash.py 2.0.7

# Populate all missing hashes (great for initial setup)
python scripts/add_htmx_hash.py --add-all

# Add with verbose output to see what's happening
python scripts/add_htmx_hash.py 2.0.7 --verbose
```

This ensures that HTMX files downloaded by Markata are verified for integrity and provides developers with easy tools to maintain the hash database.
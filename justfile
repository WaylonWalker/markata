
record-tapes:
    #!/usr/bin/env bash
    set -euxo pipefail
    # Export the function to be used by parallel
    convert_tape() {
        tape_file="$1"
        base_name=$(basename "$tape_file" .tape)
        vhs "tapes/$base_name.tape" -o "static/$base_name.gif"
    }
    export -f convert_tape

    # Find all .tape files and run the convert_tape function in parallel
    find tapes -name '*.tape' | parallel convert_tape


cov:
    #!/usr/bin/env bash
    set -euxo pipefail

    pytest --cov-report=term-missing --cov-config=pyproject.toml --cov=markata --cov=tests tests

no-cov:
    #!/usr/bin/env bash
    set -euxo pipefail
    pytest --cov-report=term-missing --cov-config=pyproject.toml --cov=markata --cov=tests tests --no-cov

lint:
    #!/usr/bin/env bash
    set -euxo pipefail
    # this is what runs in ci
    hatch run lint-format

build-docs:
    #!/usr/bin/env bash
    set -euxo pipefail
    markata build

clean-build:
    #!/usr/bin/env bash
    uv run markata clean
    uv run markata build

serve:
    #!/usr/bin/env bash
    set -euxo pipefail
    uv run markata server

ruff-fix:
    #!/usr/bin/env bash
    set -euxo pipefail
    ruff check markata --fix

lint-test: lint cov

get-tailwind:
    curl -L https://cdn.tailwindcss.com/3.4.17 -o static/tailwind.js
get-htmx:
    curl -L https://unpkg.com/htmx.org@1.9.2 -o static/htmx.js

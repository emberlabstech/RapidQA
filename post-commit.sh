#!/usr/bin/env bash

BASE="/D1/DevCode/Labs"

REPO="$(git rev-parse --show-toplevel)"

cd "$REPO" || exit 1

files=()

if git rev-parse --verify HEAD^ >/dev/null 2>&1; then
    mapfile -d '' files < <(
        git diff \
            --name-only \
            --diff-filter=ACMRT \
            -z \
            HEAD^ \
            HEAD
    )
else
    mapfile -d '' files < <(
        git diff-tree \
            --root \
            --no-commit-id \
            --name-only \
            --diff-filter=ACMRT \
            -r \
            -z \
            HEAD
    )
fi

if [ "${#files[@]}" -eq 0 ]; then
    exit 0
fi

absolute_files=()

for file in "${files[@]}"; do
    absolute_files+=("$REPO/$file")
done

/opt/commit-review/commit-review -- "${absolute_files[@]}"

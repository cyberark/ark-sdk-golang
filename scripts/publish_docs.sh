#!/usr/bin/env bash
set -euo pipefail

# Read version (e.g., 2.3.1)
VERSION="$(tr -d ' \n' < VERSION)"
if [[ -z "${VERSION}" ]]; then
  echo "VERSION file is empty"; exit 1
fi
MAJOR="${VERSION%%.*}"
MAJOR_RE="^${MAJOR//./\\.}\\."
ALIAS="latest"

echo "Current VERSION: ${VERSION}"
echo "Major: ${MAJOR}  Alias: ${ALIAS}"

# Make sure gh-pages is up to date
echo "Making sure gh-pages is latest"
git fetch origin gh-pages:gh-pages

# Get existing deployed list from gh-pages
echo "Reading deployed versions..."
LIST="$(poetry run mike list -b gh-pages || true)"
echo "---- mike list ----"
echo "${LIST}"
echo "-------------------"

# Find if alias exists (e.g., 'v2 -> 2.2.0')
existing_alias_target="$(echo "${LIST}" | awk -v a="^${ALIAS} -> " '$0 ~ a {sub(a,""); print $0}' | head -n1 || true)"

# collect same-major versions into a space-separated string
same_major_versions="$(printf '%s\n' "$LIST" \
  | awk -v re="${MAJOR_RE}" '$0 ~ re {print $1}' \
  | paste -sd' ' -)"

delete_args=""
if [ -n "$existing_alias_target" ]; then
  echo "Found alias $existing_alias_target; will delete."
  delete_args+="$existing_alias_target"
fi
if [ -n "$same_major_versions" ]; then
  echo "Found versions with the same major: $same_major_versions"
  delete_args+="$same_major_versions"
fi

# Delete previous same-major deployments/alias if present
if [ -n "$delete_args" ]; then
  # shellcheck disable=SC2086  # we intentionally split $delete_args
  echo "Deleting previous deployments: $delete_args"
  poetry run mike delete -b gh-pages -p $delete_args
fi

# Deploy new version and (re)create the major alias
# -u updates alias to point to this version; -p pushes; -b selects branch
echo "Deploying ${VERSION} with alias ${ALIAS}..."
poetry run mike deploy -b gh-pages -u -p "${VERSION}" "${ALIAS}"

echo "Done."

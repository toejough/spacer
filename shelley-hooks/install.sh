#!/usr/bin/env bash
# Install Shelley slash hooks for this project.
# Run from the repo root.
set -euo pipefail

HOOKS_DIR="$HOME/.config/shelley/hooks"
REPO_HOOKS_DIR="$(cd "$(dirname "$0")" && pwd)"

mkdir -p "$HOOKS_DIR/slash"

for hook in "$REPO_HOOKS_DIR"/slash/*; do
  name=$(basename "$hook")
  target="$HOOKS_DIR/slash/$name"
  if [ -e "$target" ]; then
    echo "Backing up existing $target"
    mv "$target" "$target.bak.$(date +%s)"
  fi
  cp "$hook" "$target"
  chmod +x "$target"
  echo "Installed: $target"
done

echo "Shelley hooks installed. You can now use /opsx in Shelley."

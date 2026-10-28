#!/usr/bin/env bash
# this is an inconsistent workaround for
# (✓) Executing command [ command=go run . ]
# ERROR: The process "11260" not found.

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

while true; do
  # go run .
  dlv debug $SCRIPT_DIR/../cmd/server --listen=127.0.0.1:2345 \
    --continue --headless --accept-multiclient \
    --output bin/debug-roommates
  echo "command exited (code $?), retrying in 2s..."
  sleep 2
done
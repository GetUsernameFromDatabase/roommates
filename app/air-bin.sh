#!/bin/bash
# load env
set -o allexport
source "../.env"
set +o allexport

# run binary
./tmp/main
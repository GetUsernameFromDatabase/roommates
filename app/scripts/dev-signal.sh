#!/usr/bin/env bash
# templ's --cmd: templ kills and reruns this on every Go change, it only tells
# scripts/dev-run.sh to restart the server
#
# the server is not templ's child because on Windows a failed TASKKILL of any
# process in the tree makes templ keep the dead PID and never run --cmd again,
# cmd.exe -> sleep always dies cleanly

mkdir -p bin
date +%s%N > bin/.dev-restart
exec sleep 2147483647

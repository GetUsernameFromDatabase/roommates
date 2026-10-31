#!/usr/bin/env bash
# runs the server under dlv (127.0.0.1:2345) and restarts it whenever the
# Taskfile's restart-watch writes bin/.dev-restart, on a Go or embedded
# web/assets change

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
TRIGGER=bin/.dev-restart
DLV_ADDR=127.0.0.1:2345

# a signal left from an earlier run only causes one extra restart, removing
# the file here could race with templ's first signal
mkdir -p bin

pid=""
winpid=""

dlv_port_busy() {
  netstat -ano | grep -q "$DLV_ADDR .*LISTENING"
}

# by Windows PID, the bash PID of a native program can be gone while it still runs
dlv_alive() {
  tasklist //FI "PID eq $winpid" //FO CSV //NH | grep -q "\"$winpid\""
}

# a process of the tree can refuse to die while dlv is building or starting
# the server, so the kill is retried until dlv is gone and its port is free,
# returns 1 if it is still busy after ~5s
stop() {
  [ -z "$pid" ] && return 0
  for _ in $(seq 10); do
    taskkill //T //F //PID "$winpid" > /dev/null 2>&1
    if ! dlv_alive && ! dlv_port_busy; then
      pid=""
      return 0
    fi
    sleep 0.5
  done
  echo "dev-run: dlv or $DLV_ADDR is still busy, retrying"
  return 1
}

start() {
  dlv debug "$SCRIPT_DIR/../cmd/server" --listen=$DLV_ADDR \
    --continue --headless --accept-multiclient \
    --output bin/debug-roommates &
  pid=$!
  winpid=$(cat "/proc/$pid/winpid")
}

trap stop EXIT
trap 'exit 130' INT TERM

last=""
while true; do
  current=$(cat "$TRIGGER" 2> /dev/null)
  if [ -n "$current" ] && [ "$current" != "$last" ]; then
    # still busy: keep last unchanged so the next poll tries again
    stop && start && last=$current
  elif [ -n "$pid" ] && ! dlv_alive; then
    # e.g. a compile error, the next change starts it again
    echo "dev-run: server exited, waiting for the next change"
    pid=""
  fi
  sleep 0.5
done

#!/usr/bin/env bash
set -euo pipefail

usage() {
  echo "usage: $0 [workspace]" >&2
  echo "environment: ULTRAPLAN_BIN, ULTRAPLAN_SERVER_UNIT, ULTRAPLAN_SERVER_URL" >&2
}

if (( $# > 1 )); then
  usage
  exit 2
fi

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
workspace_input=${1:-${ULTRAPLAN_WORKSPACE:-$script_dir/..}}
if [[ ! -d "$workspace_input" ]]; then
  echo "workspace directory does not exist: $workspace_input" >&2
  exit 4
fi
workspace=$(cd -- "$workspace_input" && pwd)
if [[ ! -f "$workspace/ultraplan.yml" ]]; then
  echo "workspace does not contain ultraplan.yml: $workspace" >&2
  exit 4
fi

ultraplan_bin=${ULTRAPLAN_BIN:-}
if [[ -z "$ultraplan_bin" ]]; then
  ultraplan_bin=$(command -v ultraplan || true)
fi
if [[ -z "$ultraplan_bin" || ! -x "$ultraplan_bin" ]]; then
  echo "ultraplan executable not found; set ULTRAPLAN_BIN to the installed executable" >&2
  exit 4
fi
ultraplan_bin=$(realpath -- "$ultraplan_bin")

unit=${ULTRAPLAN_SERVER_UNIT:-ultraplan-serve.service}
server_url=${ULTRAPLAN_SERVER_URL:-http://127.0.0.1:8080/}

echo "stopping $unit"
systemctl --user stop "$unit" 2>/dev/null || true

for _ in {1..50}; do
  load_state=$(systemctl --user show "$unit" --property=LoadState --value 2>/dev/null || true)
  active_state=$(systemctl --user show "$unit" --property=ActiveState --value 2>/dev/null || true)
  if [[ -z "$load_state" || "$load_state" == "not-found" || "$active_state" == "inactive" || "$active_state" == "failed" ]]; then
    break
  fi
  sleep 0.1
done

systemctl --user reset-failed "$unit" 2>/dev/null || true
echo "starting $unit in $workspace"
systemd-run \
  --user \
  --collect \
  --unit="${unit%.service}" \
  --description="UltraPlan local server" \
  --working-directory="$workspace" \
  "$ultraplan_bin" serve >/dev/null

for _ in {1..100}; do
  if systemctl --user is-active --quiet "$unit"; then
    if ! command -v curl >/dev/null 2>&1 || curl --fail --silent --output /dev/null "$server_url"; then
      echo "restarted $unit"
      echo "dashboard: $server_url"
      exit 0
    fi
  fi
  sleep 0.1
done

echo "server did not become ready: $unit" >&2
systemctl --user --no-pager --full status "$unit" >&2 || true
exit 1

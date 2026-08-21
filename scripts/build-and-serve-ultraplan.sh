#!/usr/bin/env bash
set -euo pipefail

usage() {
  echo "usage: $0 [workspace]" >&2
  echo "environment: ULTRAPLAN_SOURCE, ULTRAPLAN_DEV_LISTEN, ULTRAPLAN_DEV_UNIT, ULTRAPLAN_DEV_BUILD_DIR" >&2
}

if (( $# > 1 )); then
  usage
  exit 2
fi

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
workspace_input=${1:-${ULTRAPLAN_WORKSPACE:-$script_dir/..}}
source_input=${ULTRAPLAN_SOURCE:-$script_dir/../../ultraplan-go}
listen=${ULTRAPLAN_DEV_LISTEN:-127.0.0.1:8081}
unit=${ULTRAPLAN_DEV_UNIT:-ultraplan-dev-serve.service}
build_dir=${ULTRAPLAN_DEV_BUILD_DIR:-${TMPDIR:-/tmp}/ultraplan-dev-$UID}

if [[ ! -d "$workspace_input" ]]; then
  echo "workspace directory does not exist: $workspace_input" >&2
  exit 4
fi
workspace=$(cd -- "$workspace_input" && pwd)
if [[ ! -f "$workspace/ultraplan.yml" ]]; then
  echo "workspace does not contain ultraplan.yml: $workspace" >&2
  exit 4
fi

if [[ ! -d "$source_input" ]]; then
  echo "UltraPlan source directory does not exist: $source_input" >&2
  exit 4
fi
source=$(cd -- "$source_input" && pwd)
if [[ ! -f "$source/go.mod" || ! -d "$source/cmd/ultraplan" ]]; then
  echo "UltraPlan source checkout is invalid: $source" >&2
  exit 4
fi

if [[ "$build_dir" != /* ]]; then
  echo "temporary build directory must be absolute: $build_dir" >&2
  exit 4
fi
mkdir -p -- "$build_dir"
chmod 700 "$build_dir"

new_binary=$(mktemp "$build_dir/.ultraplan.XXXXXX")
cleanup() {
  if [[ -n "${new_binary:-}" && -f "$new_binary" ]]; then
    rm -f -- "$new_binary"
  fi
}
trap cleanup EXIT

echo "building $source/cmd/ultraplan"
(
  cd -- "$source"
  go build -trimpath -o "$new_binary" ./cmd/ultraplan
)
chmod 755 "$new_binary"

echo "stopping prior $unit if present"
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

binary="$build_dir/ultraplan"
mv -f -- "$new_binary" "$binary"
new_binary=

echo "starting $unit on $listen"
systemd-run \
  --user \
  --collect \
  --unit="${unit%.service}" \
  --description="UltraPlan local build server" \
  --working-directory="$workspace" \
  "$binary" --workspace "$workspace" serve --listen "$listen" >/dev/null

server_url="http://$listen/"
for _ in {1..150}; do
  if systemctl --user is-active --quiet "$unit"; then
    if ! command -v curl >/dev/null 2>&1 || curl --fail --silent --output /dev/null "$server_url"; then
      echo "running local build: $binary"
      echo "dashboard: $server_url"
      exit 0
    fi
  fi
  sleep 0.1
done

echo "local build server did not become ready: $unit" >&2
systemctl --user --no-pager --full status "$unit" >&2 || true
exit 1

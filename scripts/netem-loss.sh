#!/bin/sh
set -eu

if [ "$(id -u)" -ne 0 ]; then
  echo "Run as root, for example: sudo ./scripts/netem-loss.sh eth0 50" >&2
  exit 1
fi

INTERFACE=${1:-}
LOSS=${2:-50}

if [ -z "$INTERFACE" ]; then
  echo "Usage: sudo $0 <interface> [loss-percent] [command ...]" >&2
  exit 1
fi

case "$LOSS" in
  ''|*[!0-9]*) echo "Loss must be an integer from 0 to 100" >&2; exit 1 ;;
esac
if [ "$LOSS" -lt 0 ] || [ "$LOSS" -gt 100 ]; then
  echo "Loss must be an integer from 0 to 100" >&2
  exit 1
fi

cleanup() {
  tc qdisc del dev "$INTERFACE" root 2>/dev/null || true
  echo "netem removed from $INTERFACE"
}
trap cleanup EXIT INT TERM

# Remove an old test rule so repeated demonstrations remain safe.
tc qdisc del dev "$INTERFACE" root 2>/dev/null || true
tc qdisc add dev "$INTERFACE" root netem loss "$LOSS"%
echo "Applying $LOSS% packet loss on $INTERFACE. Cleanup is automatic."

shift 2 2>/dev/null || shift $#
if [ "$#" -gt 0 ]; then
  "$@"
else
  echo "Press Enter to stop the test and restore networking."
  read -r _
fi

#!/bin/sh
set -eu

mkdir -p data/raft

cleanup() {
  kill "${NODE1:-}" "${NODE2:-}" "${NODE3:-}" 2>/dev/null || true
}
trap cleanup EXIT INT TERM

go run ./cmd/logger-service --id log-1 --http 127.0.0.1:9101 --raft 127.0.0.1:9201 --data data/raft/log-1 --bootstrap & NODE1=$!
go run ./cmd/logger-service --id log-2 --http 127.0.0.1:9102 --raft 127.0.0.1:9202 --data data/raft/log-2 & NODE2=$!
go run ./cmd/logger-service --id log-3 --http 127.0.0.1:9103 --raft 127.0.0.1:9203 --data data/raft/log-3 & NODE3=$!

sleep 4
curl --fail --silent --show-error -X POST -H 'Content-Type: application/json' -d '{"id":"log-2","address":"127.0.0.1:9202"}' http://127.0.0.1:9101/cluster/join
curl --fail --silent --show-error -X POST -H 'Content-Type: application/json' -d '{"id":"log-3","address":"127.0.0.1:9203"}' http://127.0.0.1:9101/cluster/join

echo
echo "Three-node logger is ready. Leader API: http://127.0.0.1:9101"
echo "Run game with: go run ./cmd/client --ui gui --log-server http://127.0.0.1:9101"
echo "Press Ctrl+C to stop all logger nodes."
wait

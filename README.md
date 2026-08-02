# HoloDam: Javanese Strategy Chess Game

Two independent Go processes play HoloDam through UDP sockets. One player acts as the
authoritative host, while another player joins as a client. The game can be controlled through a
responsive browser GUI or a terminal interface.

## Feature Matrix

| Feature | Status |
| --- | --- |
| HoloDam engine | Implemented |
| Host-authoritative two-process gameplay | Implemented |
| Responsive browser GUI | Implemented |
| Terminal interface | Implemented |
| Reliable UDP transport | Implemented |
| ACK, retry, and duplicate suppression | Implemented |
| JSONL local gameplay logs | Implemented |
| Netem packet-loss helper script | Implemented |
| LAN and Tailscale usage guide | Implemented |
| Non-linear Elo rating system (logistic) | Implemented |
| Separate Raft logging service | Implemented |
| Three-node Raft deployment | Implemented |
| Two-person video recording guide | TBA |

## Quick Start

Run the browser GUI:

```bash
go run ./cmd/client --ui gui
```

Open:

```text
http://localhost:8080
```

Build standalone programs:

```bash
go build -o client ./cmd/client
go build -o logger-service ./cmd/logger-service
```

Run the built game client:

```bash
./client --ui gui
```

For one-machine testing, open two terminals.

Player A:

```bash
go run ./cmd/client --ui gui --addr :8080 --listen :9001
```

Player B:

```bash
go run ./cmd/client --ui gui --addr :8081 --listen :9002
```

Open `http://localhost:8080` and choose **Host game**. Open `http://localhost:8081`, enter
`127.0.0.1:9001` as the opponent address, and choose **Join game**.

## Two Machines

On Player A's machine, run:

```bash
go run ./cmd/client --ui gui --addr :8080 --listen :9001
```

Open `http://localhost:8080`, select a profile, and choose **Host game**. Share Player A's LAN IP
and UDP port with Player B, for example:

```text
192.168.1.20:9001
```

On Player B's machine, run:

```bash
go run ./cmd/client --ui gui --addr :8080 --listen :9002
```

Open `http://localhost:8080`, enter Player A's address, and choose **Join game**.

Both devices must allow the configured UDP ports through their firewalls. TCP port `8080` is only
used to open the browser interface on each device.

Two devices on different networks can use a VPN such as Tailscale. Install Tailscale on both
devices, join the same tailnet, and enter Player A's Tailscale address instead of the LAN address:

```text
100.x.x.x:9001
```

Tailscale encrypts and routes the connection across different networks. The game's reliable UDP
layer still handles packet loss through acknowledgements and retransmission.

## Game Controls

In the browser GUI, select one of your pieces and then select a destination point. The interface
shows the current turn, both players, captured pieces, move number, elapsed time, and gameplay log.

Available actions include:

- Select or change the active piece
- Move to a valid destination
- Continue a chained capture with the same piece, or choose **End turn** to stop voluntarily
- Apply Dam Ora Mangan when a player skips an available capture: the opponent removes any three offending pieces
- View the gameplay history
- Surrender the match
- Return to the main menu

In terminal mode, the available commands are:

```text
move <source> <target>
board
help
surrender
quit
```

## Reliable UDP

The custom transport is implemented above Go UDP sockets. Reliable packets contain a session ID
and sequence number. The receiver sends an ACK, while the sender retries when an ACK is lost or
delayed. Duplicate packets are acknowledged but are not applied to the game twice.

To test the game with 50% packet loss on Linux:

```bash
sudo ./scripts/netem-loss.sh eth0 50
```

Replace `eth0` with the network interface used by the game. The script removes the `tc netem` rule
automatically when it exits.

## Raft Logger

The Raft logger runs separately from the game processes. Start a local three-node cluster with:

```bash
./scripts/start-raft-local.sh
```

Run a game client and send its gameplay events to the Raft leader:

```bash
go run ./cmd/client --ui gui --log-server http://127.0.0.1:9101
```

Check the leader and replicated events:

```bash
curl http://127.0.0.1:9101/health
curl http://127.0.0.1:9101/events
```

A three-node cluster can continue committing events when one logger node fails. If the cluster
loses its majority, new writes are rejected to prevent conflicting gameplay histories.

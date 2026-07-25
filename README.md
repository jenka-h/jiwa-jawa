# JavaC: Javanese-Styled Chess

Jiwa Jawa is a Go project scaffold for a two-player networked `javanese-styled chess` game. The project is intended to run each player as a separate process and communicate over UDP with a custom reliability layer.

At the moment, this repository contains the planned directory structure for the application. Source packages can be added under `src/`, executables under `cmd/`, and supporting materials under `docs/`, `assets/`, and `scripts/`.

## Current Directory Structure

```text
jiwa-jawa/
├── assets/          # Images, icons, demo media, or other static assets
├── cmd/             # Go executable entry points, for example cmd/player
├── docs/            # Architecture notes, protocol notes, and test plans
├── scripts/         # Helper scripts for running, testing, or network simulation
├── src/
│   ├── app/         # Application orchestration and game flow
│   ├── logging/     # Game/event logging
│   ├── protocol/    # Message format, encoding, and decoding
│   ├── rating/      # Optional player rating calculation
│   ├── session/     # Player sessions, turns, and game lifecycle state
│   ├── transport/   # UDP transport and reliability layer
│   ├── types/       # Shared domain types such as board, piece, move, and player
│   └── ui/          # CLI/TUI/GUI presentation layer
└── README.md
```

The `cmd/player` package should contain only the program entry point and command-line parsing. Most reusable code should live under `src/` packages.

## Package Responsibilities

| Directory | Responsibility |
| --- | --- |
| `cmd/` | Main Go executables. A likely first executable is `cmd/player`, used to start one player process. |
| `src/app/` | Coordinates the game loop, UI, session state, transport, protocol, and logging. |
| `src/types/` | Core game data structures such as board, position, move, piece, player, and game state. |
| `src/session/` | Manages local player identity, turns, game lifecycle, and opponent state. |
| `src/protocol/` | Defines network messages and handles serialization/deserialization. |
| `src/transport/` | Wraps UDP sockets and implements custom reliable delivery. |
| `src/logging/` | Stores moves, received messages, network events, and game history. |
| `src/ui/` | User interaction layer, such as CLI, TUI, or future GUI. |
| `src/rating/` | Optional rating system for players. |
| `docs/` | Design documents, protocol explanation, and packet-loss testing notes. |
| `scripts/` | Helper scripts for development and testing. |
| `assets/` | Static files used by documentation, UI, or demos. |

## Planned Features

### Main Features

| Feature | Description | Planned Area | Status |
| --- | --- | --- | --- |
| Separate player programs | Each player runs their own process. | `cmd/player`, `src/app`, `src/session` | Planned |
| UDP communication | Players communicate using UDP sockets. | `src/transport` | Planned |
| Reliable protocol over UDP | Messages are resent and acknowledged to tolerate packet loss. | `src/transport`, `src/protocol` | Planned |
| Catur jawa rules | Moves are validated according to the game rules. | `src/types`, `src/session`, `src/app` | Planned |
| Packet-loss testing | Test reliability using Linux `netem` or similar tools. | `docs/`, `scripts/` | Planned |
| Game logging | Store moves and important game/network events. | `src/logging` | Planned |

### Optional / Bonus Features

| Feature | Description | Planned Area | Status |
| --- | --- | --- | --- |
| GUI | Add a graphical interface after the core game works. | `src/ui`, `assets` | Not started |
| Separate logging service | Run logging as a separate service, potentially with Raft. | future `cmd/logger-service` | Not started |
| Rating system | Calculate player ratings using a linear algebra based approach. | `src/rating` | Not started |
| Demo video | Record a gameplay/networking demo. | `assets`, `docs` | Not started |

## Getting Started

This repository does not currently include Go source files or a `go.mod` file. To initialize the Go module, run this from the `jiwa-jawa` directory:

```sh
go mod init jiwa-jawa
```

Then add a player executable, for example:

```text
cmd/player/main.go
```

After Go files are added, common commands will be:

```sh
go run ./cmd/player
```

```sh
go test ./...
```

```sh
go fmt ./...
```
## Packet-Loss Testing

On Linux, packet loss can be simulated with `tc netem`. For example:

```sh
sudo tc qdisc add dev lo root netem loss 50%
```

Remove the rule after testing:

```sh
sudo tc qdisc del dev lo root
```

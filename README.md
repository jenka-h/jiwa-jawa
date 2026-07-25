# JavaC: A Javanese-Style Chess

A Go foundation for a multiplayer `catur jawa` game over UDP.

This project is planned as a two-player networked board game where each player runs a separate program. Communication uses UDP, with a custom reliability layer built above UDP to handle packet loss.

## Feature checklist

### Main features

| Feature | Requirement | Planned implementation area | Status |
| --- | --- | --- | --- |
| Separate player programs | Player A and Player B run as two different programs/processes. | `cmd/player`, `internal/app`, `internal/network` | Planned |
| UDP communication | All player-to-player communication uses UDP sockets. | `internal/network` | Planned |
| Custom reliability protocol | Build a custom protocol above UDP to ensure messages are delivered despite packet loss. | `internal/network`, `internal/protocol` | Planned |
| Valid catur jawa gameplay | Moves must follow the real rules of catur jawa. | `internal/domain`, `internal/rules` | Planned |
| Packet-loss testing | Test with Linux `netem`, for example `tc qdisc add dev eth0 root netem loss 50%`, then remove the rule after testing. | manual test plan/docs | Planned |
| Game logging | Store moves and/or opponent moves so the current game condition can be understood later. | `internal/logging` | Planned |

### Bonus features

| Bonus | Requirement | Planned implementation area | Status |
| --- | --- | --- | --- |
| GUI | Create a graphical interface for the game. | future `internal/gui` or separate frontend | Not started |
| Separate logging service with Raft | Logging is handled by a separate program, not the game program, and uses Raft to ensure correctness. | future `cmd/logger-service`, `internal/raftlog` | Not started |
| Rating system | Create a rating system for all players. The calculation can be free-form but should use linear algebra. | future `internal/rating` | Not started |
| Demo video | Record a demo of the program, ideally with one partner. | documentation/demo asset | Not started |

## Project Structure

```text
jiwa-jawa
├── cmd
│   ├── player          # main executable for Player A / Player B
│   └── logger-service  # optional bonus logging service
├── internal
│   ├── app             # application orchestration
│   ├── domain          # board, pieces, moves, game state
│   ├── rules           # catur jawa rule validation
│   ├── protocol        # network message format and encoding/decoding
│   ├── network         # UDP socket and custom reliable transport
│   ├── logging         # local game event logging
│   ├── rating          # optional rating system
│   └── gui             # optional GUI layer
├── docs                # architecture notes and test plans
├── go.mod
└── README.md
```

## Package responsibilities

| Package | Responsibility |
| --- | --- |
| `cmd/player` | Entry point for running a player process. Parses CLI arguments such as local port, peer address, and player name. |
| `internal/app` | Coordinates the game controller, rules engine, network transport, protocol, and logger. |
| `internal/domain` | Contains pure game data: board, position, piece, move, player, and game state. No UDP or GUI code. |
| `internal/rules` | Validates and applies legal catur jawa moves. |
| `internal/protocol` | Converts game messages into bytes and converts received bytes back into structured messages. |
| `internal/network` | Wraps Go UDP sockets and implements reliable delivery with sequence numbers, ACKs, retries, and duplicate detection. |
| `internal/logging` | Persists moves and important game/network events. |
| `internal/gui` | Optional GUI layer. Should call the app/controller layer instead of containing game logic. |
| `internal/rating` | Optional player rating calculation. |``
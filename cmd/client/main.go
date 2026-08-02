package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"time"

	"jiwa-jawa/internal/engine"
	"jiwa-jawa/internal/gamelog"
	"jiwa-jawa/internal/gui"
	"jiwa-jawa/internal/transport/protocol"
	"jiwa-jawa/internal/transport/rudp"
)

type cliConfig struct {
	mode      string
	name      string
	id        string
	listen    string
	peer      string
	sessionID uint64
	logPath   string
	logServer string
}

type cliGame struct {
	cfg        cliConfig
	conn       *rudp.Connection
	logger     *gamelog.Logger
	state      *engine.GameState
	localSide  engine.Side
	remoteSide engine.Side
	localID    string
	remoteID   string
	sessionID  uint64
}

func main() {
	ui := flag.String("ui", "cli", "interface to run: cli or gui")
	addr := flag.String("addr", ":8080", "HTTP address for the GUI; :8080 allows access from other devices")
	assets := flag.String("assets", "assets", "path to GUI assets directory")

	mode := flag.String("mode", "host", "CLI mode: host or join")
	name := flag.String("name", "Player", "player display name")
	id := flag.String("id", "", "stable player id; defaults to --name")
	listen := flag.String("listen", ":9001", "UDP listen address for CLI mode")
	peer := flag.String("peer", "127.0.0.1:9002", "peer UDP address for CLI mode")
	sessionID := flag.Uint64("session", uint64(time.Now().UnixNano()), "match session id used by the reliable UDP protocol")
	logPath := flag.String("log", "logs/game.jsonl", "local JSONL game log path")
	logServer := flag.String("log-server", "", "Raft logger leader HTTP URL, for example http://127.0.0.1:9101")
	profilePath := flag.String("profiles", "data/profiles.json", "persistent GUI profile file")
	flag.Parse()

	switch *ui {
	case "gui":
		runGUI(*addr, *assets, gui.Options{
			ListenAddr:  strings.TrimSpace(*listen),
			PeerAddr:    strings.TrimSpace(*peer),
			SessionID:   *sessionID,
			LogPath:     strings.TrimSpace(*logPath),
			LogServer:   strings.TrimSpace(*logServer),
			ProfilePath: strings.TrimSpace(*profilePath),
		})
	case "cli":
		playerID := strings.TrimSpace(*id)
		if playerID == "" {
			playerID = strings.TrimSpace(*name)
		}
		cfg := cliConfig{
			mode:      strings.ToLower(strings.TrimSpace(*mode)),
			name:      strings.TrimSpace(*name),
			id:        playerID,
			listen:    strings.TrimSpace(*listen),
			peer:      strings.TrimSpace(*peer),
			sessionID: *sessionID,
			logPath:   strings.TrimSpace(*logPath),
			logServer: strings.TrimSpace(*logServer),
		}
		if err := runCLI(cfg); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown --ui %q; use cli or gui", *ui)
	}
}

func runGUI(addr, assets string, options gui.Options) {
	server := gui.NewServerWithOptions(assets, options)
	displayAddr := addr
	if strings.HasPrefix(addr, ":") {
		displayAddr = "localhost" + addr
	}
	fmt.Printf("HoloDam is running at http://%s\n", displayAddr)
	if strings.HasPrefix(addr, ":") || strings.HasPrefix(addr, "0.0.0.0:") {
		fmt.Println("Other devices: open http://<this-device-LAN-IP>" + addr[strings.LastIndex(addr, ":"):])
	}
	fmt.Println("Press Ctrl+C to stop.")

	if err := server.ListenAndServe(addr); err != nil {
		log.Fatal(err)
	}
}

func runCLI(cfg cliConfig) error {
	if cfg.mode != "host" && cfg.mode != "join" {
		return fmt.Errorf("invalid --mode %q: use host or join", cfg.mode)
	}
	if cfg.id == "" {
		return fmt.Errorf("--id or --name is required")
	}
	if cfg.mode == "join" && cfg.peer == "" {
		return fmt.Errorf("--peer is required in join mode")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger, err := gamelog.NewWithRemote(cfg.logPath, cfg.logServer)
	if err != nil {
		return err
	}
	defer logger.Close()

	peerAddr := ""
	if cfg.mode == "join" {
		peerAddr = cfg.peer
	}
	conn, err := rudp.NewConnection(rudp.Config{
		ListenAddr: cfg.listen,
		PeerAddr:   peerAddr,
		SessionID:  cfg.sessionID,
	})
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.Start(ctx); err != nil {
		return err
	}

	board, err := engine.NewDefaultBoard()
	if err != nil {
		return err
	}
	game := &cliGame{
		cfg:       cfg,
		conn:      conn,
		logger:    logger,
		state:     engine.NewGameState(board, engine.SideOne),
		localID:   cfg.id,
		sessionID: cfg.sessionID,
	}

	fmt.Printf("Listening on %s\n", conn.LocalAddr())
	if err := logger.Log(gamelog.Event{Type: "start", SessionID: cfg.sessionID, Message: fmt.Sprintf("%s listening on %s", cfg.mode, conn.LocalAddr())}); err != nil {
		return err
	}

	if cfg.mode == "host" {
		game.localSide = engine.SideOne
		game.remoteSide = engine.SideTwo
		if err := game.waitForJoin(ctx); err != nil {
			return err
		}
	} else {
		game.localSide = engine.SideTwo
		game.remoteSide = engine.SideOne
		if err := game.joinHost(ctx); err != nil {
			return err
		}
	}

	fmt.Printf("Ready. You are %s (%s).\n", game.localSide, cfg.name)
	fmt.Println("Commands: move <source> <target>, end, penalty <point>, board, help, surrender, quit")
	game.printBoard()
	return game.loop(ctx)
}

func (g *cliGame) waitForJoin(ctx context.Context) error {
	fmt.Println("Waiting for player B to join...")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-g.conn.Errors():
			fmt.Printf("transport error: %v\n", err)
		case received := <-g.conn.Incoming():
			if received.Packet.Header.MessageType != protocol.MessageJoin {
				continue
			}
			var join protocol.JoinPayload
			if err := json.Unmarshal(received.Packet.Payload, &join); err != nil {
				return err
			}
			g.remoteID = join.PlayerID
			if err := g.conn.SetPeer(received.Addr.String()); err != nil {
				return err
			}
			fmt.Printf("Player joined from %s: %s (%s)\n", received.Addr, join.PlayerName, join.PlayerID)
			_ = g.logger.Log(gamelog.Event{Type: "join", SessionID: g.sessionID, Side: g.remoteSide, Message: fmt.Sprintf("%s (%s)", join.PlayerName, join.PlayerID), Seq: received.Packet.Header.Sequence})
			accept := protocol.AcceptPayload{PlayerID: g.localID, PlayerName: g.cfg.name, SessionID: g.sessionID}
			return g.sendJSON(protocol.MessageAccept, accept)
		}
	}
}

func (g *cliGame) joinHost(ctx context.Context) error {
	join := protocol.JoinPayload{PlayerID: g.localID, PlayerName: g.cfg.name}
	fmt.Printf("Joining host at %s...\n", g.cfg.peer)
	if err := g.sendJSON(protocol.MessageJoin, join); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-g.conn.Errors():
			fmt.Printf("transport error: %v\n", err)
		case received := <-g.conn.Incoming():
			if received.Packet.Header.MessageType != protocol.MessageAccept {
				continue
			}
			var accept protocol.AcceptPayload
			if err := json.Unmarshal(received.Packet.Payload, &accept); err != nil {
				return err
			}
			g.remoteID = accept.PlayerID
			g.sessionID = accept.SessionID
			g.conn.SetSessionID(accept.SessionID)
			fmt.Printf("Accepted by host %s (%s), session %d\n", accept.PlayerName, accept.PlayerID, accept.SessionID)
			_ = g.logger.Log(gamelog.Event{Type: "accept", SessionID: g.sessionID, Side: g.remoteSide, Message: fmt.Sprintf("%s (%s)", accept.PlayerName, accept.PlayerID), Seq: received.Packet.Header.Sequence})
			return nil
		}
	}
}

func (g *cliGame) loop(ctx context.Context) error {
	inputCh := make(chan string)
	go scanInput(inputCh)

	for !g.state.Finished {
		if g.state.CurrentTurn == g.localSide {
			fmt.Print("your turn> ")
		} else {
			fmt.Printf("waiting for %s> ", g.state.CurrentTurn)
		}

		select {
		case <-ctx.Done():
			_ = g.sendJSON(protocol.MessageLeave, protocol.FinishPayload{Reason: "interrupted"})
			return nil
		case err := <-g.conn.Errors():
			fmt.Printf("transport error: %v\n", err)
			_ = g.logger.Log(gamelog.Event{Type: "transport_error", SessionID: g.sessionID, Error: err.Error()})
		case received := <-g.conn.Incoming():
			if err := g.handlePacket(received); err != nil {
				fmt.Printf("packet error: %v\n", err)
				_ = g.logger.Log(gamelog.Event{Type: "packet_error", SessionID: g.sessionID, Error: err.Error(), Seq: received.Packet.Header.Sequence})
			}
		case line := <-inputCh:
			if err := g.handleInput(strings.TrimSpace(line)); err != nil {
				fmt.Println(err)
			}
		}
	}

	fmt.Printf("Game finished. Winner: %s\n", g.state.Winner)
	_ = g.logger.Log(gamelog.Event{Type: "finish", SessionID: g.sessionID, Side: g.state.Winner, Message: "game finished"})
	return nil
}

func scanInput(out chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		out <- scanner.Text()
	}
}

func (g *cliGame) handleInput(line string) error {
	if line == "" {
		return nil
	}
	fields := strings.Fields(line)
	switch strings.ToLower(fields[0]) {
	case "help":
		fmt.Println("Commands:")
		fmt.Println("  move <source> <target>  apply and send a move")
		fmt.Println("  m <source> <target>     short form of move")
		fmt.Println("  board                   show board and point ids")
		fmt.Println("  end                     decline an optional chained capture")
		fmt.Println("  penalty <point>         remove a piece during Dam Ora Mangan")
		fmt.Println("  surrender               finish the game and give the win to the peer")
		fmt.Println("  quit                    leave without declaring a winner")
	case "board":
		g.printBoard()
	case "quit":
		_ = g.sendJSON(protocol.MessageLeave, protocol.FinishPayload{Reason: "quit"})
		g.state.Finish(engine.SideNone)
	case "surrender":
		winner := engine.Opponent(g.localSide)
		g.state.Finish(winner)
		_ = g.sendJSON(protocol.MessageFinish, protocol.FinishPayload{WinnerID: string(winner), Reason: "surrender"})
	case "end", "end-turn":
		if err := g.state.EndTurn(g.localSide); err != nil {
			return err
		}
		if err := g.sendJSON(protocol.MessageEndTurn, protocol.EndTurnPayload{PlayerID: g.localID}); err != nil {
			return err
		}
		fmt.Println("Turn ended; optional chained capture declined.")
	case "penalty":
		if len(fields) != 2 {
			return fmt.Errorf("usage: penalty <point>")
		}
		target, err := strconv.Atoi(fields[1])
		if err != nil {
			return fmt.Errorf("invalid target %q", fields[1])
		}
		if err := g.state.ApplyPenalty(g.localSide, engine.PointID(target)); err != nil {
			return err
		}
		if err := g.sendJSON(protocol.MessagePenalty, protocol.PenaltyPayload{PlayerID: g.localID, Target: target}); err != nil {
			return err
		}
		fmt.Printf("Dam Ora Mangan: removed opponent piece at %d (%d remaining)\n", target, g.state.PenaltiesRemaining)
		g.printBoard()
	case "move", "m":
		if len(fields) != 3 {
			return fmt.Errorf("usage: move <source> <target>")
		}
		if g.state.CurrentTurn != g.localSide {
			return fmt.Errorf("not your turn; current turn is %s", g.state.CurrentTurn)
		}
		source, err := strconv.Atoi(fields[1])
		if err != nil {
			return fmt.Errorf("invalid source %q", fields[1])
		}
		target, err := strconv.Atoi(fields[2])
		if err != nil {
			return fmt.Errorf("invalid target %q", fields[2])
		}
		move := engine.Move{Source: engine.PointID(source), Target: engine.PointID(target)}
		result, err := g.state.ApplyMoveResult(move)
		if err != nil {
			return err
		}
		payload := protocol.MovePayload{PlayerID: g.localID, Source: source, Target: target}
		if err := g.sendJSON(protocol.MessageMove, payload); err != nil {
			return err
		}
		g.logMove("move_sent", g.localSide, move, result)
		g.printMoveResult(result)
		g.printBoard()
	default:
		return fmt.Errorf("unknown command %q; type help", fields[0])
	}
	return nil
}

func (g *cliGame) handlePacket(received rudp.ReceivedPacket) error {
	switch received.Packet.Header.MessageType {
	case protocol.MessageMove:
		var payload protocol.MovePayload
		if err := json.Unmarshal(received.Packet.Payload, &payload); err != nil {
			return err
		}
		move := engine.Move{Source: engine.PointID(payload.Source), Target: engine.PointID(payload.Target)}
		result, err := g.state.ApplyMoveResult(move)
		if err != nil {
			return err
		}
		g.logMove("move_received", g.remoteSide, move, result)
		fmt.Printf("\nRemote move: %s\n", move)
		g.printMoveResult(result)
		g.printBoard()
	case protocol.MessageEndTurn:
		var payload protocol.EndTurnPayload
		if err := json.Unmarshal(received.Packet.Payload, &payload); err != nil {
			return err
		}
		if err := g.state.EndTurn(g.remoteSide); err != nil {
			return err
		}
		fmt.Println("\nOpponent ended their turn and declined the chained capture.")
	case protocol.MessagePenalty:
		var payload protocol.PenaltyPayload
		if err := json.Unmarshal(received.Packet.Payload, &payload); err != nil {
			return err
		}
		if err := g.state.ApplyPenalty(g.remoteSide, engine.PointID(payload.Target)); err != nil {
			return err
		}
		fmt.Printf("\nDam Ora Mangan: opponent removed piece at %d (%d remaining)\n", payload.Target, g.state.PenaltiesRemaining)
		g.printBoard()
	case protocol.MessageLeave:
		fmt.Println("\nPeer left the match.")
		g.state.Finish(g.localSide)
	case protocol.MessageFinish:
		var payload protocol.FinishPayload
		_ = json.Unmarshal(received.Packet.Payload, &payload)
		winner := engine.Side(payload.WinnerID)
		if winner == engine.SideNone {
			winner = g.remoteSide
		}
		fmt.Printf("\nPeer finished the match: %s\n", payload.Reason)
		g.state.Finish(winner)
	case protocol.MessageHeartbeat, protocol.MessageStateRequest, protocol.MessageStateSnapshot, protocol.MessageJoin, protocol.MessageAccept:
		// Valid protocol messages, but not part of the active move loop.
	default:
		return fmt.Errorf("unsupported message type %s", received.Packet.Header.MessageType)
	}
	return nil
}

func (g *cliGame) sendJSON(messageType protocol.MessageType, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	packet := protocol.NewPacket(g.sessionID, 0, messageType, data, true)
	if err := g.conn.SendReliable(packet); err != nil {
		_ = g.logger.Log(gamelog.Event{Type: "send_error", SessionID: g.sessionID, Message: messageType.String(), Error: err.Error()})
		return err
	}
	_ = g.logger.Log(gamelog.Event{Type: "send", SessionID: g.sessionID, Message: messageType.String()})
	return nil
}

func (g *cliGame) logMove(eventType string, side engine.Side, move engine.Move, result engine.MoveResult) {
	message := "normal move"
	if len(result.Captured) > 0 {
		message = fmt.Sprintf("captured %v", result.Captured)
	}
	if result.MustContinue {
		message += "; optional chained capture available"
	}
	_ = g.logger.Log(gamelog.Event{
		Type:      eventType,
		SessionID: g.sessionID,
		Side:      side,
		Message:   message,
		Move:      &gamelog.MoveEvent{Source: move.Source, Target: move.Target},
	})
}

func (g *cliGame) printMoveResult(result engine.MoveResult) {
	if len(result.Captured) > 0 {
		fmt.Printf("Captured: %v\n", result.Captured)
	}
	if result.MustContinue {
		fmt.Println("Another capture is available from the landing point; continue or use 'end'.")
	}
	if result.DamOraMangan {
		fmt.Println("Dam Ora Mangan! Opponent may remove any three of your pieces.")
	}
	if result.GameOver {
		fmt.Printf("Winner: %s\n", result.Winner)
	}
}

func (g *cliGame) printBoard() {
	fmt.Println()
	fmt.Println("Board:")
	positions := map[[2]int]engine.PointID{}
	for id, pos := range g.state.Board.Points {
		positions[[2]int{pos.X, pos.Y}] = id
	}
	for y := 0; y <= 4; y++ {
		var line strings.Builder
		for x := 0; x <= 8; x++ {
			id, ok := positions[[2]int{x, y}]
			if !ok {
				line.WriteString("   ")
				continue
			}
			marker := "."
			if piece, occupied := g.state.Board.PieceAt(id); occupied {
				switch piece.Owner {
				case engine.SideOne:
					marker = "A"
				case engine.SideTwo:
					marker = "B"
				}
			}
			line.WriteString(fmt.Sprintf("%2s ", marker))
		}
		fmt.Println(line.String())
	}

	ids := make([]int, 0, len(g.state.Board.Points))
	for id := range g.state.Board.Points {
		ids = append(ids, int(id))
	}
	sort.Ints(ids)
	fmt.Println("Point IDs:")
	for _, rawID := range ids {
		id := engine.PointID(rawID)
		pos := g.state.Board.Points[id]
		marker := "."
		if piece, occupied := g.state.Board.PieceAt(id); occupied {
			if piece.Owner == engine.SideOne {
				marker = "A"
			} else {
				marker = "B"
			}
		}
		fmt.Printf("%2d:(%d,%d)=%s  ", id, pos.X, pos.Y, marker)
		if rawID%4 == 0 {
			fmt.Println()
		}
	}
	fmt.Printf("\nTurn: %s | Move: %d | Pieces A=%d B=%d\n\n", g.state.CurrentTurn, g.state.MoveNumber, engine.CountPieces(g.state.Board, engine.SideOne), engine.CountPieces(g.state.Board, engine.SideTwo))
}

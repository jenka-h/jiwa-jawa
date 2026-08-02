package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"jiwa-jawa/internal/raftlog"
)

func main() {
	id := flag.String("id", "log-1", "unique Raft node ID")
	httpAddr := flag.String("http", ":9101", "HTTP API listen address")
	raftAddr := flag.String("raft", "127.0.0.1:9201", "Raft address advertised to other nodes")
	dataDir := flag.String("data", "data/raft/log-1", "durable Raft data directory")
	bootstrap := flag.Bool("bootstrap", false, "bootstrap a new cluster; use only on the first node")
	flag.Parse()

	store, err := raftlog.Open(*id, *raftAddr, *dataDir, *bootstrap)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	server := &http.Server{Addr: *httpAddr, Handler: raftlog.NewHTTPHandler(store), ReadHeaderTimeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	fmt.Printf("Raft logger %s: HTTP %s, Raft %s, data %s\n", *id, *httpAddr, *raftAddr, *dataDir)
	if *bootstrap {
		fmt.Println("Bootstrapping cluster; wait for this node to become leader before joining others.")
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

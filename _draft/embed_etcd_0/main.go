package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/coreos/etcd/embed"
)

var (
	name    = flag.String("name", "", "player identity")
	connect = flag.String("connect", "", "IP:port to connect with another player")
	dataDir = flag.String("data-dir", filepath.Join(os.TempDir(), "tiwe-game.data"), "path to game data storage")
)

func main() {
	flag.Parse()

	if err := start(); err != nil {
		log.Fatal(err)
	}
}

func start() error {
	cfg := embed.NewConfig()
	cfg.Dir = *dataDir
	cfg.EnableV2 = false
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		return err
	}
	defer e.Close()
	select {
	case <-e.Server.ReadyNotify():
		log.Printf("Server is ready!")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Printf("Server took too long to start!")
	}
	return <-e.Err()
}

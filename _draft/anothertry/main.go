package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var (
	addr  = flag.String("addr", "", "address to listen for peers")
	peers = flag.String("peers", "", "comma-separated list of peer addresses")
)

var (
	mu    sync.RWMutex
	conns map[string]net.Conn
)

func init() {
	conns = make(map[string]net.Conn)
}

func main() {
	flag.Parse()
	var flags []string
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		flags = append(flags, f.Name, f.Value.String())
	})
	log.Printf("tiwe-experimental build\ncmd args: %v\ncmd flags: %v", flag.Args(), flags)

	var wg sync.WaitGroup
	defer wg.Wait()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on %v", ln.Addr())
	defer ln.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		serve(ctx, ln)
	}()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				var b bytes.Buffer
				mu.RLock()
				fmt.Fprintf(&b, "connection count: %d\n", len(conns))
				for who, conn := range conns {
					fmt.Fprintf(&b, "self (%v) <-> %s (%v)\n", conn.LocalAddr(), who, conn.RemoteAddr())
				}
				mu.RUnlock()
				log.Print(b.String())
			}
		}
	}()

	if *peers != "" {
		time.Sleep(500 * time.Millisecond)
		startGame(ln.Addr().String(), strings.Split(*peers, ","))
	}

	ch := make(chan os.Signal, 1)
	defer close(ch)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

func serve(ctx context.Context, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			return
		}
		go handleConnection(ctx, conn)
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	log.Printf("connected: %v <-> %v", conn.LocalAddr(), conn.RemoteAddr())
	defer func() {
		conn.Close()
		log.Printf("connection closed: %v <-> %v", conn.LocalAddr(), conn.RemoteAddr())
	}()
	dec, enc := gob.NewDecoder(conn), gob.NewEncoder(conn)

	var who string
	err := dec.Decode(&who)
	if err != nil {
		log.Print(err)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	if _, ok := conns[who]; ok {
		log.Printf("already connected to %v", who)
		return
	}
	conns[who] = conn

	var peers []string
	err = dec.Decode(&peers)
	if err != nil {
		log.Print(err)
		return
	}
	self := conn.LocalAddr().String()
	startGame(self, peers)

	_ = enc
}

func startGame(self string, peers []string) {
	mu.Lock()
	defer mu.Unlock()
	var d net.Dialer
	for _, addr := range peers {
		if _, ok := conns[addr]; addr == self || ok {
			continue
		}
		conn, err := d.DialContext(context.TODO(), "tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		conns[addr] = conn
		enc := gob.NewEncoder(conn)
		enc.Encode(self)
		enc.Encode(append([]string{self}, peers...))
	}
}

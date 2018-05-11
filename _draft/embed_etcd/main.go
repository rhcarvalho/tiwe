package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/coreos/etcd/embed"
)

func main() {
	cfg := embed.NewConfig()
	var err error
	cfg.Dir, err = ioutil.TempDir("", "tiwe-experiment-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(cfg.Dir)

	var e *embed.Etcd

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		if e != nil {
			log.Printf("Shutting down.")
			e.Server.Stop()
		}
	}()

	e, err = embed.StartEtcd(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	select {
	case <-e.Server.ReadyNotify():
		log.Printf("Server is ready!")
	case <-time.After(1 * time.Nanosecond):
		e.Server.Stop()                              // trigger a shutdown
		log.Printf("Server took too long to start!") // blocks forever
	}
	log.Fatal(<-e.Err())
}

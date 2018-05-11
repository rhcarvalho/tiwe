package main

import (
	"log"
	"net"
)

type ObservableConn struct {
	net.Conn
}

func (oc *ObservableConn) Read(b []byte) (n int, err error) {
	n, err = oc.Conn.Read(b)
	log.Printf("read -> %x", b)
	return
}

func (oc *ObservableConn) Write(b []byte) (n int, err error) {
	log.Printf("write <- %x", b)
	return oc.Conn.Write(b)
}

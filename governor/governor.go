package main

import (
	"fmt"
	"log"
	"net"
)

type providerStatus struct {
	status string
}

func main() {
	addr := "127.0.0.1:1234"

	status := make(chan providerStatus)
	go providerListener(addr, status)

	for {
		s := <-status
		log.Println("Status from", s.status)
	}

}

func providerListener(addr string, status chan providerStatus) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Exit(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		} else {
			providerHandler(conn, status)
		}
	}
}

func providerHandler(conn net.Conn, status chan providerStatus) {
	buf := make ([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		log.Print(conn.RemoteAddr(), err)
	}
	status <- providerStatus{string(buf[0:n])}

	fmt.Fprintf(conn, "bye\n")
	conn.Close()
}

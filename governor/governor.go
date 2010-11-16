package main

import (
	"fmt"
	"http"
	"log"
	"net"
	"strings"
)

type status struct {
	name string
	info provider
}

type provider struct {
	online bool
}

type providers map[string]provider

func main() {
	govAddr := "127.0.0.1:1234"
	httpAddr := "127.0.0.1:8080"

	p := make(providers, 10)
	http.Handle("/", &p)
	go httpServer(httpAddr)

	updates := make(chan status)
	go providersListener(govAddr, updates)

	for {
		s := <-updates
		log.Println("Status from", s.name)
		p[s.name] = s.info
	}
}

func httpServer(addr string) {
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Exit(err)
	}
}

func (p providers) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello from display status")
	fmt.Fprintf(w, "<ul>")

	for name, info := range p {
		fmt.Fprintf(w, "<li>%v - %#v</li>", name, info)
	}
	fmt.Fprintf(w, "</ul>")
}

func providersListener(addr string, updates chan status) {
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
			providerHandler(conn, updates)
		}
	}
}

func providerHandler(conn net.Conn, updates chan status) {
	buf := make([]byte, 1024)

	status := *new(status)
	status.info.online = true

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Print(conn.RemoteAddr(), err)
			status.info.online = false
			updates <-status
			return
		}

		req := string(buf[0:n])
		log.Println("Recieved:", req)
		parts := strings.Split(req, " ", 2)
		command := parts[0]
		arguments := parts[1]

		ok := true
		switch command {
		case "NAME": status.name = arguments
		default: ok = false
		}

		if ok {
			fmt.Fprintf(conn, "OK")
		} else {
			fmt.Fprintf(conn, "ERROR")
		}

		updates <- status
	}
	conn.Close()
}

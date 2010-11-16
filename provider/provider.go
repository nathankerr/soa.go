package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	addr := "127.0.0.1:1234"

	conn, err := net.Dial("tcp", "", addr)
	if err != nil {
		log.Exit(err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "PROVIDER %v", conn.LocalAddr())
}

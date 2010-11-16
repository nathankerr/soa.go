package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	addr := "127.0.0.1:1234"

	conn, err := net.Dial("tcp", "", addr)
	if err != nil {
		log.Exit(err)
	}
	defer conn.Close()

	var name string
	if len(os.Args) > 1 {
		name = os.Args[1]
	} else {
		parts := strings.Split(conn.LocalAddr().String(), ":", -1)
		name = "Provider " + parts[len(parts) - 1]
	}

	fmt.Fprintf(conn, "NAME %v", name)

	time.Sleep(5000000000)
}

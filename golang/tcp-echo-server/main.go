package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
)

var addr = "127.0.0.01:8080"

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err != io.EOF {
				slog.Error(err.Error())
				os.Exit(1)
			}
			return
		}
		conn.Write(bytes)
		slog.Info(fmt.Sprintf("Finished connection from %s", conn.RemoteAddr()))
	}
}

func main() {
	if len(os.Args) == 2 {
		addr = os.Args[1]
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	slog.Info(fmt.Sprintf("Listening on: %s", addr))
	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		slog.Info(fmt.Sprintf("Accepting connection from %s", conn.RemoteAddr()))
		go handleConn(conn)
	}
}

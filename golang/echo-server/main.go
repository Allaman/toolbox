package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"sync"
)

var tcpAddr = "127.0.0.1:8080"
var httpAddr = "127.0.0.1:8081"

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err == io.EOF {
				slog.Info(fmt.Sprintf("TCP connection closed by client %s", conn.RemoteAddr()))
				return
			}
			// Log the error but don't exit - just close this connection
			slog.Warn(fmt.Sprintf("TCP read error from %s: %s", conn.RemoteAddr(), err.Error()))
			return
		}
		_, writeErr := conn.Write(bytes)
		if writeErr != nil {
			slog.Warn(fmt.Sprintf("TCP write error to %s: %s", conn.RemoteAddr(), writeErr.Error()))
			return
		}
		slog.Info(fmt.Sprintf("Echoed data to TCP client %s", conn.RemoteAddr()))
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info(fmt.Sprintf("HTTP request from %s: %s %s", r.RemoteAddr, r.Method, r.URL.Path))

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(body)
		slog.Info(fmt.Sprintf("Echoed %d bytes to %s", len(body), r.RemoteAddr))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Echo server is running. Send POST request to echo content.\n")
	}
}

func startTCPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	listener, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		slog.Error(fmt.Sprintf("TCP server error: %s", err.Error()))
		os.Exit(1)
	}
	defer listener.Close()

	slog.Info(fmt.Sprintf("TCP server listening on: %s", tcpAddr))

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Warn(fmt.Sprintf("TCP accept error: %s", err.Error()))
			continue
		}
		slog.Info(fmt.Sprintf("Accepting TCP connection from %s", conn.RemoteAddr()))
		go handleConn(conn)
	}
}

func startHTTPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	http.HandleFunc("/", echoHandler)
	http.HandleFunc("/echo", echoHandler)

	slog.Info(fmt.Sprintf("HTTP server listening on: %s", httpAddr))

	err := http.ListenAndServe(httpAddr, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("HTTP server error: %s", err.Error()))
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) >= 2 {
		tcpAddr = os.Args[1]
	}
	if len(os.Args) >= 3 {
		httpAddr = os.Args[2]
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go startTCPServer(&wg)

	wg.Add(1)
	go startHTTPServer(&wg)

	slog.Info("Server started with TCP and HTTP support")
	slog.Info(fmt.Sprintf("TCP echo server: %s", tcpAddr))
	slog.Info(fmt.Sprintf("HTTP echo server: %s", httpAddr))

	wg.Wait()
}

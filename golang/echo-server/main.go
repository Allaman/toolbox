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
	clientIP := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	for {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err == io.EOF {
				slog.Info("TCP connection closed by client", "source_ip", clientIP)
				return
			}
			// Log the error but don't exit - just close this connection
			slog.Warn("TCP read error", "source_ip", clientIP, "error", err.Error())
			return
		}
		_, writeErr := conn.Write(bytes)
		if writeErr != nil {
			slog.Warn("TCP write error", "source_ip", clientIP, "error", writeErr.Error())
			return
		}
		slog.Info("Echoed data to TCP client", "source_ip", clientIP)
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := r.RemoteAddr
	slog.Info("HTTP request received", "source_ip", clientIP, "method", r.Method, "path", r.URL.Path)

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Warn("Error reading request body", "source_ip", clientIP, "error", err.Error())
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(body)
		slog.Info("Echoed data to HTTP client", "source_ip", clientIP, "bytes", len(body))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Echo server is running. Send POST request to echo content.\n")
		slog.Info("Served HTTP info page", "source_ip", clientIP)
	}
}

func startTCPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	listener, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		slog.Error("TCP server error", "error", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	slog.Info("TCP server listening", "address", tcpAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Warn("TCP accept error", "error", err.Error())
			continue
		}
		clientIP := conn.RemoteAddr().String()
		slog.Info("Accepting TCP connection", "source_ip", clientIP)
		go handleConn(conn)
	}
}

func startHTTPServer(wg *sync.WaitGroup) {
	defer wg.Done()

	http.HandleFunc("/", echoHandler)
	http.HandleFunc("/echo", echoHandler)

	slog.Info("HTTP server listening", "address", httpAddr)

	err := http.ListenAndServe(httpAddr, nil)
	if err != nil {
		slog.Error("HTTP server error", "error", err.Error())
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
	slog.Info("TCP echo server", "address", tcpAddr)
	slog.Info("HTTP echo server", "address", httpAddr)

	wg.Wait()
}

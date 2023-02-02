package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	startServer(80)
}

func startServer(port int) {
	address := ":" + fmt.Sprintf("%v", port)
	fmt.Println("Starting tcp server on", address)
	listener, err := net.Listen("tcp", address)
	handleError(err)
	defer listener.Close()
	for {
		connection, err := listener.Accept()
		go handleConnection(connection, err)
	}
}

func handleConnection(connection net.Conn, err error) {
	handleError(err)
	defer connection.Close()

	request := getRequest(connection)
	fmt.Println("Intercepted request", request)

	route(request, connection)
}

func route(request request, connection net.Conn) {
	if request.method == "GET" && request.path == "/hi" {
		_, _ = io.WriteString(connection, ok().toHttpResponseString())
		return
	}
	_, _ = io.WriteString(connection, notFund(request.path).toHttpResponseString())
	return
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"net"
	"strings"
)

type request struct {
	method  string
	path    string
	headers map[string][]string
	body    string
}

func getRequest(connection net.Conn) request {
	finishedReadingHeaders := false
	request := &request{headers: make(map[string][]string)}
	payload := readBytes(connection)
	lines := strings.Split(payload, "\r\n")
	for i, line := range lines {
		if i == 0 {
			parseRequestLine(line, request)
		} else {
			if line == "" {
				finishedReadingHeaders = true
				continue
			}
			if !finishedReadingHeaders {
				parseHeader(line, request)
			} else {
				request.body = strings.ReplaceAll(line, "\n", "")
				break
			}
		}
	}
	return *request
}

func readBytes(connection net.Conn) string {
	bufferSize := 1024
	payloadBuilder := new(strings.Builder)
	for {
		buf := make([]byte, bufferSize)
		readLen, _ := connection.Read(buf)
		payloadBuilder.WriteString(string(buf[:readLen]))
		if readLen < bufferSize {
			break
		}
	}
	return payloadBuilder.String()
}

func parseRequestLine(line string, request *request) {
	fields := strings.Fields(line)
	request.method = fields[0]
	request.path = fields[1]
}

func parseHeader(line string, request *request) {
	header := strings.Split(line, ": ")
	values := strings.Split(header[1], ",")
	request.headers[header[0]] = values
}

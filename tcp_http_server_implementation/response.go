package main

import (
	"fmt"
	"strings"
)

type response struct {
	status     string
	statusCode int
	headers    map[string][]string
	body       *string
}

func ok() *response {
	return &response{status: "OK", statusCode: 200}
}

func notFund(path string) *response {
	body := fmt.Sprintf("Requested resource (%s) was not found on this server", path)
	return &response{
		status:     "NOT FOUND",
		body:       &body,
		statusCode: 404,
		headers: map[string][]string{
			"Content-Type": {"application/text"},
		},
	}
}

func (r *response) toHttpResponseString() string {
	stringBuilder := new(strings.Builder)
	stringBuilder.WriteString(fmt.Sprintf("HTTP/1.1 %v %s\r\n", r.statusCode, r.status))
	for header, value := range r.headers {
		stringBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", header, strings.Join(value, ",")))
	}
	if r.body != nil {
		stringBuilder.WriteString(fmt.Sprintf("Content-Length: %v\r\n", len(*r.body)))
		stringBuilder.WriteString("\r\n")
		stringBuilder.WriteString(*r.body)
	}
	return stringBuilder.String()
}

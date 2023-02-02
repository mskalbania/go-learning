package main

import (
	"fmt"
	"net/http"
	"time"
)

type status struct {
	host          string
	executionTime int64
	statusCode    int
	status        string
	errorMessage  string
}

func main() {
	hosts := []string{
		"https://google.com",
		"https://facebook.com",
	}
	channel := make(chan status)

	for _, host := range hosts {
		go check(host, channel)
	}

	for hostStatus := range channel {
		if hostStatus.status == "UP" {
			fmt.Printf("Host %s is UP (status - %v) and responded in %v ms\n", hostStatus.host, hostStatus.statusCode, hostStatus.executionTime)
		} else {
			fmt.Println(hostStatus)
		}
		go func(host string) {
			time.Sleep(3 * time.Second)
			check(host, channel)
		}(hostStatus.host)

	}
}

func check(host string, c chan status) {
	timeStartMs := time.Now().UnixMilli()
	rs, err := http.Get(host)
	timeTakenMs := time.Now().UnixMilli() - timeStartMs
	if err != nil {
		c <- status{
			host:         host,
			status:       "DOWN",
			errorMessage: fmt.Sprint(err),
		}
		return
	}
	if rs.StatusCode > 399 {
		c <- status{
			host:          host,
			status:        "DOWN",
			executionTime: timeTakenMs,
			statusCode:    rs.StatusCode,
		}
		return
	}
	c <- status{
		host:          host,
		executionTime: timeTakenMs,
		statusCode:    rs.StatusCode,
		status:        "UP",
	}
}

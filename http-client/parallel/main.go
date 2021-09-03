package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	numberOfRequests = 10
	baseURL          = "https://postman-echo.com"
	c                = http.Client{Timeout: time.Duration(1) * time.Second}
)

func main() {
	wg := sync.WaitGroup{}
	for no := 1; no <= numberOfRequests; no++ {
		wg.Add(1)
		go func(no int) {
			resp := get("/get?foo1=bar1&foo2=bar2")
			fmt.Printf("Requested number %d with status %v\n", no, resp.Status)
			wg.Done()
		}(no)
	}
	wg.Wait()
}

func get(path string) *http.Response {
	resp, err := c.Get(fmt.Sprintf("%s%s", baseURL, path))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return resp
}

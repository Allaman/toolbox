package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	aux "github.com/allaman/toolbox/http-client/auxiliary"
)

var (
	baseURL = "https://postman-echo.com"
	c       = http.Client{Timeout: time.Duration(1) * time.Second}
)

func main() {
	get("/get?foo1=bar1&foo2=bar2")
	post("/post", map[string]string{"foo": "bar"})
}

func get(path string) {
	resp, err := c.Get(fmt.Sprintf("%s%s", baseURL, path))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	aux.PrintResponse(resp)
}

func post(path string, body interface{}) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	resp, err := c.Post(fmt.Sprintf("%s%s", baseURL, path), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	aux.PrintResponse(resp)
}

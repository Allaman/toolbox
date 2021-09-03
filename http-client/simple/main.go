package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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
	printResponse(resp)
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
	printResponse(resp)
}

func printResponse(resp *http.Response) {
	if resp.StatusCode == 200 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		prettyJSON, err := formatJSON(b)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", prettyJSON)
	} else {
		fmt.Printf("Status is: %s\n", resp.Status)
	}
}

func formatJSON(data []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "    ")
	if err == nil {
		return out.Bytes(), err
	}
	return data, nil
}

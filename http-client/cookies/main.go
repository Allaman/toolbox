package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	aux "github.com/allaman/toolbox/http-client/auxiliary"
)

var (
	baseURL = "https://postman-echo.com"
	c       = http.Client{Timeout: time.Duration(1) * time.Second}
)

func main() {
	cookies := get("/get?foo1=bar1&foo2=bar2")
	post("/post", map[string]string{"foo": "bar"}, cookies)
}

func get(path string) []*http.Cookie {
	resp, err := c.Get(fmt.Sprintf("%s%s", baseURL, path))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	for _, cookie := range resp.Cookies() {
		fmt.Printf("Found a cookie named '%s' with value '%s'\n", cookie.Name, cookie.Value)
	}
	return resp.Cookies()
}

func post(path string, body interface{}, cookies []*http.Cookie) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", baseURL, path), bytes.NewBuffer((jsonData)))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	for i := range cookies {
		req.AddCookie(cookies[i])
	}
	fmt.Println("Posting with cookie...")
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Result with posted cookie")
	aux.PrintResponse(resp)
}

// Read converts a io.Reader object to a string
func read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}

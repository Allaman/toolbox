package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

var (
	baseURL = "https://example.de/api/v4/"
	c       = http.Client{Timeout: time.Duration(1) * time.Second}
)

func searchProjectID(id int, token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", baseURL, "projects/", strconv.Itoa(id)), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("PRIVATE-TOKEN", token)
	log.Printf("Request is: %v", req)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if !isHTTPResponseCode200(resp) {
		log.Printf("statuscode was '%d'", resp.StatusCode)
	}
	return resp, nil
}

func isHTTPResponseCode200(resp *http.Response) bool {
	return resp.StatusCode == 200
}

func getJSONFromHTTPResponse(resp *http.Response) ([]byte, error) {
	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func getJSONKey(json string, key string) string {
	return gjson.Get(json, key).String()
}

func printProjectName(resp *http.Response) error {
	json, err := getJSONFromHTTPResponse(resp)
	if err != nil {
		return err
	}
	project := getJSONKey(string(json), "name_with_namespace")
	fmt.Println(project)
	return nil
}

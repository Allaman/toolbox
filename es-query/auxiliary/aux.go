package aux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/tidwall/gjson"
)

// NewESClient returns an simple opiniated ES client
func NewESClient() *elasticsearch.Client {
	cfg := elasticsearch.Config{
		RetryOnStatus: []int{429, 502, 503, 504},
		RetryBackoff: func(i int) time.Duration {
			// A simple exponential delay
			d := time.Duration(math.Exp2(float64(i))) * time.Second
			log.Printf("Attempt: %d | Sleeping for %s...\n", i, d)
			return d
		},
		Transport: &http.Transport{
			ResponseHeaderTimeout: (time.Millisecond * 30000),
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalln("error while creating ES client: ", err)
	}
	return es
}

// Read converts a io.Reader object to a string
func Read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}

// ColoredLog prints simple colored log messages
func ColoredLog(message string) {
	colorReset := "\033[0m"
	colorCyan := "\033[36m"
	log.Println(string(colorCyan), message, string(colorReset))
}

// GetKey simplyfies gjson.Get
func GetKey(json *string, key string) string {
	return gjson.Get(*json, key).String()
}

// PrintStats prints Statistics from a ES search query result
func PrintStats(json *string) {
	values := gjson.GetMany(*json, "hits.total.value", "took")
	log.Printf(
		"%d hits; took: %dms",
		values[0].Int(),
		values[1].Int(),
	)
	docs := gjson.Get(*json, "hits.hits")
	fmt.Println("Printing first document")
	for _, doc := range docs.Array() {
		fmt.Println(doc.String())
		break
	}
}

// ConstructBody builds a valid ES search body
// inspired from https://kb.objectrocket.com/elasticsearch/how-to-construct-elasticsearch-queries-from-a-string-using-golang-550
func ConstructBody(body string) *strings.Reader {
	isValid := json.Valid([]byte(body)) // returns bool
	// Default query is "{}" if JSON is invalid
	if !isValid {
		log.Fatalf("Not a valid json - can not construct body: %s", body)
	}
	// Build a new string from JSON query
	var b strings.Builder
	b.WriteString(body)
	// Instantiate a *strings.Reader object from string
	bodyReader := strings.NewReader(b.String())
	// Return a *strings.Reader object
	return bodyReader
}

// CheckForESSearchResultError validates if a ES search response errored
func CheckForESSearchResultError(res *esapi.Response) {
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s with statuscode %d", err, res.StatusCode)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
}

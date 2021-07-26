package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/tidwall/gjson"
)

func main() {
	es, _ := elasticsearch.NewDefaultClient()

	log.SetFlags(0)

	var index = "kubernetes_cluster*"

	// all documents within the last hour
	query := map[string]interface{}{
		"size": 2000,
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"@timestamp": map[string]interface{}{
					"format": "strict_date_optional_time",
					"gte":    "now-1h",
					"lt":     "now",
				},
			},
		},
	}
	{
	}

	// all documents
	// query := map[string]interface{}{
	// 	"size": 2000,
	// 	"query": map[string]interface{}{
	// 		"match_all": map[string]interface{}{},
	// 	},
	// }

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	json := read(res.Body)
	if !gjson.Valid(json) {
		log.Fatal("invalid json")
	}

	printStats(&json)

	createCSV(&json)

}

func read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}

func createCSV(json *string) {
	coloredLog("Printing first document's values as csv")
	hits := gjson.Get(*json, "hits.hits")
	csvWriter := csv.NewWriter(os.Stdout)
	hits.ForEach(func(key, value gjson.Result) bool {
		raw := value.Raw
		timestamp := getKey(&raw, "_source.@timestamp")
		index := getKey(&raw, "_index")
		id := getKey(&raw, "_id")
		l := getKey(&raw, "_source.log")
		records := []string{timestamp, index, id, l}
		csvWriter.Write(records)
		if err := csvWriter.Error(); err != nil {
			log.Fatalln("error writing csv:", err)
		}
		return false // keep iterating
	})
	csvWriter.Flush()
}

func printStats(json *string) {
	values := gjson.GetMany(*json, "hits.total.value", "took")
	log.Printf(
		"%d hits; took: %dms",
		values[0].Int(),
		values[1].Int(),
	)
	docs := gjson.Get(*json, "hits.hits")
	coloredLog("Printing first document")
	for _, doc := range docs.Array() {
		println(doc.String())
		break
	}
}

func coloredLog(message string) {
	colorReset := "\033[0m"
	colorCyan := "\033[36m"
	log.Println(string(colorCyan), message, string(colorReset))
}
func getKey(json *string, key string) string {
	return gjson.Get(*json, key).String()
}

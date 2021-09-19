package main

import (
	"context"
	"log"

	aux "github.com/allaman/toolbox/es-query/auxiliary"
	"github.com/tidwall/gjson"
)

func main() {
	es := aux.NewESClient()

	var index = "" // TODO:
	var size = 100

	// all documents within the last hour
	var query = `"bool": {
            "filter": [
                {
                "range": {
                    "time": {
                        "format": "strict_date_optional_time",
                        "gte":    "now-1h",
                        "lt":     "now"
                        }
                    }
                }
            ]
        }`

	// all documents in an alternative form for an query
	// query := map[string]interface{}{
	// 	"size": 2000,
	// 	"query": map[string]interface{}{
	// 		"match_all": map[string]interface{}{},
	// 	},
	// }

	if index == "" {
		log.Fatalln("index should not be empty")
	}

	body := aux.ConstructBodyFromQuery(query, size)

	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(body),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	aux.CheckForESSearchResultError(res)

	json := aux.Read(res.Body)
	if !gjson.Valid(json) {
		log.Fatal("invalid json")
	}
	aux.PrintStats(&json)
}

package main

import (
	"context"
	"fmt"
	"log"

	aux "github.com/allaman/toolbox/es-query/auxiliary"
	"github.com/tidwall/gjson"
)

func main() {
	es := aux.NewESClient()

	var index = "kubernetes_cluster*"

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

	// all documents and alternative raw form of an query
	// query := map[string]interface{}{
	// 	"size": 2000,
	// 	"query": map[string]interface{}{
	// 		"match_all": map[string]interface{}{},
	// 	},
	// }

	body := aux.ConstructBodyFromQuery(query, 100)

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
	if res.IsError() {
		fmt.Printf("err", res)
	}

	aux.CheckForESSearchResultError(res)

	json := aux.Read(res.Body)
	if !gjson.Valid(json) {
		log.Fatal("invalid json")
	}
	aux.PrintStats(&json)
}

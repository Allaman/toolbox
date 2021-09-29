package main

import (
	"fmt"
	"log"
	"time"

	aux "github.com/allaman/toolbox/es-query/auxiliary"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/tidwall/gjson"
)

func main() {
	es := aux.NewESClient()

	var index = "" // TODO:
	var size = 1000

	var query = `"bool": {
            "filter": [
                {
                    "match_phrase": {
                        "kubernetes.container_name": {
                            "query": "FooBar"
                        }
                    }
                },
                {
                "range": {
                    "time": {
                        "format": "strict_date_optional_time",
                        "gte": "2021-09-21T00:00:00.000Z",
                        "lte": "2021-09-21T23:59:59.000Z"
                        }
                    }
                }
            ]
        }`

	if index == "" {
		log.Fatalln("index should not be empty")
	}

	rawData := scroll(es, index, query, size)
	printRawData(rawData)

}

func printRawData(rawData []gjson.Result) {
	for _, gJsonResult := range rawData {
		gJsonResult.ForEach(func(key, value gjson.Result) bool {
			fmt.Println(value.Raw)
			return false // NOTE: to keep iterating set this to true
		})
		return // NOTE: remove this to loop through all elements
	}
}

func scroll(es *elasticsearch.Client, index string, query string, size int) []gjson.Result {
	var results []gjson.Result
	body := aux.ConstructBody(`{"query": {`, query, size)
	res, err := es.Search(
		es.Search.WithIndex(index),
		es.Search.WithBody(body),
		es.Search.WithScroll(time.Minute),
		es.Search.WithErrorTrace(),
	)
	defer res.Body.Close()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	aux.CheckForESSearchResultError(res)
	log.Println("Scrolling ...")
	// Extract the first batch of documents and extract scrollID
	json := aux.Read(res.Body)
	hits := gjson.Get(json, "hits.hits")
	results = append(results, hits)
	scrollID := gjson.Get(json, "_scroll_id").String()
	for {
		log.Println("Scrolling...")
		res, err := es.Scroll(es.Scroll.WithScrollID(scrollID), es.Scroll.WithScroll(time.Minute))
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		aux.CheckForESSearchResultError(res)
		json := aux.Read(res.Body)
		res.Body.Close()
		scrollID = gjson.Get(json, "_scroll_id").String()
		hits := gjson.Get(json, "hits.hits")
		// Stop scrolling when no hits are returned
		if len(hits.Array()) < 1 {
			return results
		}
		results = append(results, hits)
	}
}

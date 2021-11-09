package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	aux "github.com/allaman/toolbox/es-query/auxiliary"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/tidwall/gjson"
)

func main() {
	es := aux.NewESClient()
	bQuery, err := ioutil.ReadFile("query.json")
	if err != nil {
		log.Fatal(err)
	}

	var index = "" // TODO:

	if index == "" {
		log.Fatalln("index should not be empty")
	}

	rawData := scroll(es, index, string(bQuery))
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

func scroll(es *elasticsearch.Client, index string, query string) []gjson.Result {
	var results []gjson.Result
	body := aux.ConstructBody(query)
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

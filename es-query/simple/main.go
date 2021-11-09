package main

import (
	"context"
	"io/ioutil"
	"log"

	aux "github.com/allaman/toolbox/es-query/auxiliary"
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

	body := aux.ConstructBody(string(bQuery))

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

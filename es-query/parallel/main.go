package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
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
	var maxID = 8  // defies the number of slices and therefore the number of cuncurrent go routines

	if index == "" {
		log.Fatalln("index should not be empty")
	}

	wg := &sync.WaitGroup{}
	rawData := make(chan []gjson.Result)
	for i := 0; i < maxID; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			rawData <- slicedSearch(es, index, string(bQuery), j, maxID)
		}(i)
	}
	// Wait untill all routines finished AND entries are fetched from channel
	go func() {
		wg.Wait()
		close(rawData)
	}()
	// Loop over channel and fetch entries
	for gJsonResults := range rawData {
		for _, gJsonResult := range gJsonResults {
			gJsonResult.ForEach(func(key, value gjson.Result) bool {
				fmt.Println(value.Raw)
				return true // TODO: to keep iterating set this to true
			})
		}
	}
}
func slicedSearch(es *elasticsearch.Client, index string, query string, ID, maxID int) []gjson.Result {
	// Pass the query string to the function and have it return a Reader object
	log.Printf("Routine %d is %s", ID, "starting 🚀")
	var results []gjson.Result
	body := aux.ConstructBody(fmt.Sprintf(query, ID, maxID))
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
	// Extract the first batch of documents and extract scrollID
	json := aux.Read(res.Body)
	hits := gjson.Get(json, "hits.hits")
	total := gjson.Get(json, "hits.total.value").Int()
	log.Printf("Routine %d is getting '%d' total hits", ID, total)
	if total == 0 {
		log.Fatalf("%d: %s %s", ID, "total hits 0:", res)
	}
	results = append(results, hits)
	scrollID := gjson.Get(json, "_scroll_id").String()
	for {
		log.Printf("Routine %d is %s", ID, "Scrolling... ")
		// Perform the scroll request and pass the scrollID and scroll duration
		res, err := es.Scroll(es.Scroll.WithScrollID(scrollID), es.Scroll.WithScroll(time.Minute))
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		aux.CheckForESSearchResultError(res)
		json := aux.Read(res.Body)
		res.Body.Close()
		scrollID = gjson.Get(json, "_scroll_id").String()
		hits := gjson.Get(json, "hits.hits")
		// Stop scrolling when no hits are returned
		if len(hits.Array()) < 1 {
			log.Printf("Routine %d is %s", ID, "closing ... ")
			return results
		}
		results = append(results, hits)
	}
}

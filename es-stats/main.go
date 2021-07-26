package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/fatih/color"
	"github.com/tidwall/gjson"
)

var (
	faint = color.New(color.Faint)
	bold  = color.New(color.Bold)
)

func main() {
	verCmd := flag.NewFlagSet("version", flag.ExitOnError)
	statusCmd := flag.NewFlagSet("status", flag.ExitOnError)
	idxCmd := flag.NewFlagSet("indices", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'version', 'status', or 'indicies' subcommands")
		os.Exit(1)
	}

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	switch os.Args[1] {
	case "version":
		verCmd.Parse(os.Args[2:]) // prevent 'declared but not used error'
		faint.Print(elasticsearch.Version)
		faint.Print(es.Info())
	case "status":
		statusCmd.Parse(os.Args[2:])
		res, err := es.Cluster.Stats(es.Cluster.Stats.WithHuman())
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()
		json := read(res.Body)

		faint.Print("cluster ")
		bold.Print(gjson.Get(json, "cluster_name"))
		faint.Print(" status=")
		status := gjson.Get(json, "status")
		switch status.Str {
		case "green":
			bold.Add(color.FgHiGreen).Print(status)
		case "yellow":
			bold.Add(color.FgHiYellow).Print(status)
		case "red":
			bold.Add(color.FgHiRed).Print(status)
		default:
			bold.Add(color.FgHiRed, color.Underline).Print(status)
		}
		fmt.Println("\n" + strings.Repeat("─", 50))

		stats := []string{
			"indices.count",
			"indices.docs.count",
			"indices.store.size",
			"nodes.count.total",
			"nodes.os.mem.used_percent",
			"nodes.process.cpu.percent",
			"nodes.jvm.versions.#.version",
			"nodes.jvm.mem.heap_used",
			"nodes.jvm.mem.heap_max",
			"nodes.fs.free",
		}

		var maxwidth int
		for _, item := range stats {
			if len(item) > maxwidth {
				maxwidth = len(item)
			}
		}

		for _, item := range stats {
			pad := maxwidth - len(item)
			fmt.Print(strings.Repeat(" ", pad))
			faint.Printf("%s |", item)
			fmt.Printf(" %s\n", gjson.Get(json, item))
		}
	case "indices":
		idxCmd.Parse(os.Args[2:])
		// https://pkg.go.dev/github.com/elastic/go-elasticsearch/esapi#CatIndices
		cat := es.Cat.Indices
		// res, err := es.Cat.Indices(es.Cat.Indices.WithFormat("json"))
		// sort by size and show column headings
		res, err := cat(cat.WithS("store.size"), cat.WithV(true))
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()
		data := read(res.Body)
		fmt.Print(data)
	}
}

func read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var version = "0.1.0"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	address := flag.String("address", "http://localhost:2020", "Cistern node address")
	collection := flag.String("collection", "", "Collection")
	columns := flag.String("columns", "", "Columns")
	group := flag.String("group", "", "Group")
	filters := flag.String("filters", "", "Filters")
	start := flag.Int64("start", time.Now().Unix()-3600, "Start ts")
	end := flag.Int64("end", time.Now().Unix(), "End ts")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	queryDesc, err := parseQuery(*columns, *group, *filters)
	if err != nil {
		log.Fatalln(err)
	}

	queryDesc.TimeRange.Start = time.Unix(*start, 0)
	queryDesc.TimeRange.End = time.Unix(*end, 0)

	buf := &bytes.Buffer{}
	err = json.NewEncoder(buf).Encode(queryDesc)
	if err != nil {
		log.Fatalln(err)
	}
	response, err := http.Post(fmt.Sprintf("%s/collections/%s/query", *address, *collection), "application/json", buf)
	if err != nil {
		log.Fatalln(err)
	}

	queryResult := QueryResult{}
	err = json.NewDecoder(response.Body).Decode(&queryResult)
	if err != nil {
		log.Fatalln(err)
	}

	pretty, err := json.MarshalIndent(queryResult, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s", pretty)
}

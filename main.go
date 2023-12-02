package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	url = "https://jsonmock.hackerrank.com/api/medical_records"
)

func main() {
	client := Client{URL: url}
	var mode = flag.Int("mode", 1, "0 for sync or 1 for async")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *mode == 0 {
		fmt.Printf("************* SYNC MODE *************")
		records, err := client.Fetch()
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		for _, record := range records {
			for _, data := range record.Data {
				fmt.Printf("data: %v\n", data)
			}
		}
		return
	}

	fmt.Printf("************* ASYNC MODE *************")
	channel := client.FetchAsync()
	for record := range channel {
		for _, data := range record.Data {
			fmt.Printf("data: %v\n", data)
		}
	}
}

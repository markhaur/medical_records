package main

import "fmt"

var (
	url = "https://jsonmock.hackerrank.com/api/medical_records"
)

func main() {
	client := Client{URL: url}

	// sync
	// records, err := client.Fetch()
	// if err != nil {
	// 	fmt.Printf("error: %v\n", err)
	// 	os.Exit(1)
	// }

	// for _, record := range records {
	// 	fmt.Printf("record: %v\n", record)
	// }

	// async
	channel := client.FetchAsync()
	for record := range channel {
		fmt.Printf("record: %v\n", record.Data)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type MedicalRecord struct {
	Page       int      `json:"page"`
	PerPage    int      `json:"per_page"`
	Total      int      `json:"total"`
	TotalPages int      `json:"total_pages"`
	Data       []Record `json:"data"`
}

type Record struct {
	ID        int32 `json:"id"`
	Timestamp int64 `json:"timestamp"`
	Diagnosis struct {
		ID       int32  `json:"id"`
		Name     string `json:"name"`
		Severity int32  `json:"severity"`
	} `json:"diagnosis"`
	Vitals struct {
		BloodPressureDiastole int32   `json:"bloodPressureDiastole"`
		BloodPressureSystole  int32   `json:"bloodPressureSystole"`
		Pulse                 int32   `json:"pulse"`
		BreathingRate         int32   `json:"breathingRate"`
		BodyTemperature       float32 `json:"bodyTemperature"`
	} `json:"vitals"`
	Doctor struct {
		ID   int32  `json:"id"`
		Name string `json:"name"`
	} `json:"doctor"`
	UserID   int32  `json:"userId"`
	UserName string `json:"userName"`
	UserDob  string `json:"userDob"`
	Meta     struct {
		Height int32 `json:"height"`
		Weight int32 `json:"weight"`
	} `json:"meta"`
}

type Client struct {
	URL string
}

func (c *Client) Fetch() ([]MedicalRecord, error) {
	records := make([]MedicalRecord, 0)

	page := 1
	for {
		resp, err := fetchPage(fmt.Sprintf("%s?page=%d", c.URL, page))
		if err != nil {
			return nil, fmt.Errorf("couldn't fetch page %d: %v", page, err)
		}

		records = append(records, *resp)
		page++
		if page > resp.TotalPages {
			break
		}
	}

	return records, nil
}

func (c *Client) FetchAsync() <-chan MedicalRecord {
	recordChannel := make(chan MedicalRecord, 10)

	go func() {
		defer close(recordChannel)

		resp, err := fetchPage(c.URL)
		if err != nil {
			return
		}
		recordChannel <- *resp

		var wg sync.WaitGroup
		wg.Add(resp.TotalPages - 1)

		for i := 2; i <= resp.TotalPages; i++ {
			go func(page int) {
				defer wg.Done()
				resp, err := fetchPage(fmt.Sprintf("%s?page=%d", c.URL, page))
				if err != nil {
					log.Printf("couldn't fetch page %d: %v\n", page, err)
				}
				recordChannel <- *resp
			}(i)
		}

		wg.Wait()
	}()

	return recordChannel
}

func fetchPage(url string) (*MedicalRecord, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch records: %v", err)
	}
	defer resp.Body.Close()

	var medicalRecord MedicalRecord
	if err = json.NewDecoder(resp.Body).Decode(&medicalRecord); err != nil {
		return nil, fmt.Errorf("couldn't decode response: %v", err)
	}

	return &medicalRecord, nil
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	Search string `json:"Search"`
}

type Response struct {
	Chinese          string   `json:"chinese"`
	Pinyin           string   `json:"pinyin"`
	PinyinNormalized string   `json:"pinyin_normalized"`
	Meanings         []string `json:"meanings"`
}

func main() {
	// Тестовые запросы
	testCases := []string{
		"中国",     // Иероглифы
		"zhong1", // Пиньинь
		"река",   // Русский перевод
		"hello",  // Английский (для проверки отсутствия результатов)
	}

	for _, searchTerm := range testCases {
		fmt.Printf("\n=== Testing search for: %s ===\n", searchTerm)

		reqBody := Request{Search: searchTerm}
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			fmt.Printf("Error marshaling request: %v\n", err)
			continue
		}

		resp, err := http.Post("http://localhost:8080/", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			continue
		}

		fmt.Printf("Status: %s\n", resp.Status)
		fmt.Printf("Response body:\n%s\n", string(body))

		// Пытаемся разобрать JSON, если это массив
		if len(body) > 0 && body[0] == '[' {
			var responses []Response
			if err := json.Unmarshal(body, &responses); err != nil {
				fmt.Printf("Error unmarshaling response: %v\n", err)
			} else {
				fmt.Printf("Found %d results:\n", len(responses))
				for i, r := range responses {
					fmt.Printf("  %d. %s [%s] - %v\n", i+1, r.Chinese, r.Pinyin, r.Meanings)
				}
			}
		}
	}
}

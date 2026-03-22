package main

import (
	"backend/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)



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

		reqBody := model.Request{Search: searchTerm}
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

		
	}
}

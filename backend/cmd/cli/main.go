package main

import (
	"backend/internal/search"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./cli <search_query>")
		os.Exit(1)
	}

	query := os.Args[1]
	url := "http://localhost:8080/api/entries?word=" + query

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	
	var result search.Result
	json.Unmarshal(body, &result)
	
	data := result.Data

	
	fmt.Printf("Found %d results:\n\n", len(data))

	for i, entry := range data {
		
		fmt.Printf("%d. %s [%s]\n", i+1, entry.Hanzi, entry.Pinyin)

		meanings := entry.Meanings
		for _, m := range meanings {
			fmt.Printf("   %d. %s\n", int(m.Index), m.Text)
		}
		fmt.Println()
	}
}

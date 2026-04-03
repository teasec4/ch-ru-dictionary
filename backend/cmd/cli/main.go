package main

import (
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
	url := "http://localhost:8080/api/entries?hanzi=" + query + "&limit=5"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	data := result["data"].([]interface{})
	fmt.Printf("Found %d results:\n\n", len(data))

	for i, entry := range data {
		e := entry.(map[string]interface{})
		fmt.Printf("%d. %s [%s]\n", i+1, e["hanzi"], e["pinyin"])

		meanings := e["meanings"].([]interface{})
		for _, m := range meanings {
			meaning := m.(map[string]interface{})
			fmt.Printf("   %d. %s\n", int(meaning["index"].(float64)), meaning["text"])
		}
		fmt.Println()
	}
}

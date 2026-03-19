package main

import (
	"backend/internal/handler"
	"backend/internal/service"
	"net/http"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Get the path to the database relative to the project root
	dbPath := filepath.Join("..", "db", "dictionary.db")
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	defer db.Close()

	sDB := service.NewDbService(db)
	h := handler.NewResponseHandler(sDB)

	http.HandleFunc("/", h.GetMeaning)

	port := ":8080"
	println("Server starting on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

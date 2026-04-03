package main

import (
	"backend/internal/handler"
	"backend/internal/search"
	"backend/internal/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sqlx.Connect("sqlite3", "dictionary.db")
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	defer db.Close()

	svcDB := service.NewDB(db)
	if err := svcDB.InitFTS(); err != nil {
		log.Printf("FTS init warning: %v", err)
	}

	svc := search.New(db)
	h := handler.NewHandler(svc)

	r := gin.Default()

	r.GET("/api/entries", h.SearchEntries)

	port := ":8080"
	fmt.Println("Server starting on port", port)
	if err := r.Run(port); err != nil {
		panic(err)
	}
}

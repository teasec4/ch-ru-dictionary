package handler

import (
	"backend/internal/model"
	"backend/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseHandler struct {
	Service *service.DbService
}

func NewResponseHandler(s *service.DbService) *ResponseHandler {
	return &ResponseHandler{
		Service: s,
	}
}

func (h *ResponseHandler) GetMeaning(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, что поисковый запрос не пустой
	if req.Search == "" {
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	// Логируем поисковый запрос
	fmt.Printf("Search request: %q\n", req.Search)

	// Выполняем поиск в базе данных
	results, err := h.Service.SearchWords(req.Search)
	if err != nil {
		fmt.Printf("Search error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Логируем количество найденных результатов
	fmt.Printf("Found %d results for query: %q\n", len(results), req.Search)
	if len(results) > 0 {
		fmt.Printf("First result: %s [%s]\n", results[0].Chinese, results[0].Pinyin)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		fmt.Printf("JSON encode error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

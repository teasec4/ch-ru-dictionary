package handler

import (
	"backend/internal/model"
	"backend/internal/service"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	body, _ := io.ReadAll(r.Body)
	fmt.Println("RAW BODY:", string(body))
	
	r.Body = io.NopCloser(bytes.NewBuffer(body)) 
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

	// Выполняем поиск в базе данных h.Service....

	response := model.Response{
		Chinese: "你好",
		Pinyin: "ni2hao3",
		PinyinNormalized: "nihao",
		Meanings: []string{"привет" , "здороваться"},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("JSON encode error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

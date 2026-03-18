package handler

import (
	"backend/internal/model"
	"backend/internal/service"
	"encoding/json"
	"net/http"
)

type ResponseHandler struct{
	Service *service.DbService
}

func NewResponseHandler(s *service.DbService) *ResponseHandler{
	return &ResponseHandler{
		Service: s,
	}
}

func (h *ResponseHandler) GetMeaning(w http.ResponseWriter, r *http.Request){
	var req model.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// db service logic
	
	w.Header().Set("Content-Type", "application/json")
	// send encoded data 
	json.NewEncoder(w).Encode()
}
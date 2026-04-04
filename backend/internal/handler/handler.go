package handler

import (
	"net/http"

	"backend/internal/search"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *search.Searcher
}

func NewHandler(svc *search.Searcher) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) SearchEntries(c *gin.Context) {
	query := c.Query("word")
	if query == "" {
		query = c.Query("q")
	}

	result, _ := h.svc.Search(query)

	c.JSON(http.StatusOK, result)
}

func detectSearchType(word string) (hanzi, pinyin, meaning string) {
	return word, "", ""
}

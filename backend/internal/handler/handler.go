package handler

import (
	"net/http"
	"strconv"

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
	hanzi := c.Query("hanzi")
	pinyin := c.Query("pinyin")
	meaning := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, _ := h.svc.Search(hanzi, pinyin, meaning, page, limit)

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Autocomplete(c *gin.Context) {
	prefix := c.Query("prefix")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	entries, err := h.svc.Autocomplete(prefix, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": entries})
}

package handler

import (
	"net/http"
	"strconv"
	"unicode/utf8"

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
	word := c.Query("word")
	hanzi := c.Query("hanzi")
	pinyin := c.Query("pinyin")
	meaning := c.Query("q")
	mode := c.DefaultQuery("mode", "exact")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	if word != "" && hanzi == "" && pinyin == "" && meaning == "" {
		hanzi, pinyin, meaning = detectSearchType(word)
	}

	result, _ := h.svc.Search(hanzi, pinyin, meaning, page, limit, mode == "extended")

	c.JSON(http.StatusOK, result)
}

func detectSearchType(word string) (hanzi, pinyin, meaning string) {
	if !utf8.ValidString(word) {
		return word, "", ""
	}

	hasCJK := false
	hasLatin := false
	hasCyrillic := false

	for _, r := range word {
		switch {
		case r >= 0x4E00 && r <= 0x9FFF:
			hasCJK = true
		case r >= 0x0041 && r <= 0x005A || r >= 0x0061 && r <= 0x007A:
			hasLatin = true
		case r >= 0x0400 && r <= 0x04FF:
			hasCyrillic = true
		}
	}

	if hasCJK {
		return word, "", ""
	}
	if hasLatin {
		return "", word, ""
	}
	if hasCyrillic {
		return "", "", word
	}

	return word, "", ""
}



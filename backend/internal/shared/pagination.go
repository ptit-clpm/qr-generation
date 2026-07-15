package shared

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Page struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func Pagination(c *gin.Context) Page {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 12
	}
	return Page{Page: page, Limit: limit}
}

func (p Page) Offset() int {
	return (p.Page - 1) * p.Limit
}

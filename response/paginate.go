package response

import (
	"net/http"
	"strconv"
)

type PaginateResult struct {
	PerPage     int         `json:"per_page"`
	CurrentPage int         `json:"current_page"`
	Total       int64       `json:"total"`
	Data        interface{} `json:"data"`
}

func Page(r *http.Request) int {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}
	return page
}

func PageSize(r *http.Request) int {
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 20
	}
	return pageSize
}

func Paginate(r *http.Request, count int64, data interface{}) (p PaginateResult) {
	p.Total = count
	p.Data = data
	p.CurrentPage = Page(r)
	p.PerPage = PageSize(r)
	return p
}

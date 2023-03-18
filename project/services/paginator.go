package services

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Object holding request queries for pagination (page?, limit, offset, dir?)
type PaginationParams struct {
	Page   int
	Limit  int
	Offset int
	dir    string
}

// Full and elegant pagination obejct, inspired by laravel pagination
type Pagination struct {
	CurrentPageUrl   string `json:"current_page_url"`
	CurrentPage      int    `json:"current_page"`
	TotalPages       int    `json:"total_pages"`
	PerPage          int    `json:"per_page"`
	Limit            int    `json:"limit"`
	PreviousPage     int    `json:"previous_page"`
	NextPage         int    `json:"next_page"`
	CurrentPageTotal int    `json:"current_page_total"`
	PreviousPageUrl  string `json:"previous_page_url"`
	NextPageUrl      string `json:"next_page_url"`
	LastPageUrl      string `json:"last_page_url"`
	Total            int    `json:"total"`
}

type Db struct {
	*gorm.DB
}

// Extract default pagination values
func (pagiParam *PaginationParams) ParseQuery(c *gin.Context) {
	var (
		err                 error
		limit, offset, page int64
	)
	/** pagination size */
	limit, err = strconv.ParseInt(strings.TrimSpace(c.DefaultQuery("limit", "100")), 10, 64)
	if err != nil {
		limit = 100
	}

	page, err = strconv.ParseInt(strings.TrimSpace(c.DefaultQuery("page", "1")), 10, 64)
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}

	offset = (page - 1) * (limit)

	pagiParam.Limit = int(limit)
	pagiParam.Offset = int(offset)
	pagiParam.Page = int(page)
}

func ArrayContains(arr []string, check interface{}) bool {
	for _, item := range arr {
		if item == check {
			return true
		}
	}
	return false
}

// Custom gorm model pagination, inspired by laravel pagination. model interface{} is a pointer
func Paginate(req *gin.Context, params PaginationParams, total int) (pagination *Pagination, err error) {
	var (
		scheme string
	)
	page := params.Page
	limit := params.Limit
	if (page-1)*limit > total {
		return nil, errors.New("Invalid paginator params")
	}

	host := req.Request.Host
	if strings.Contains(host, "localhost") || strings.Contains(host, "127.0.0.1") {
		scheme = "http://"
	} else {
		scheme = "https://"
	}

	paramsToIgnore := []string{"limit", "page"}
	current_page_host_and_path := scheme + req.Request.Host + req.Request.URL.Path + "?"
	for key, element := range req.Request.URL.Query() {
		value := strings.Join(element, fmt.Sprintf("&%s=", key))
		if !ArrayContains(paramsToIgnore, key) {
			current_page_host_and_path += fmt.Sprintf("%s=%s&", key, value)
		}
	}

	pagination = &Pagination{}

	// ceil is used to give a whole value even if size of remaining items
	// sometimes, last page's items may not be up to limit size
	total_pages := int(math.Ceil(float64(total) / float64((limit))))
	pagination.TotalPages = total_pages

	// if for instance last page has items not up to limit size, page*limit is
	// more than total and this has to be resolved by reducing page by one
	// before multiply by limit to the remain items for the last page
	if page*limit > total {
		pagination.CurrentPageTotal = total - (limit * (page - 1))
	} else {
		pagination.CurrentPageTotal = limit
	}

	pagination.CurrentPageUrl =
		fmt.Sprintf("%spage=%v&limit=%v", current_page_host_and_path, page, limit)
	if page > 1 {
		previous_page := page - 1
		pagination.PreviousPageUrl =
			fmt.Sprintf("%spage=%v&limit=%v", current_page_host_and_path, previous_page, limit)
		pagination.PreviousPage = previous_page
	}
	if page < total_pages && total_pages > 0 {
		next_page := page + 1
		pagination.NextPageUrl =
			fmt.Sprintf("%spage=%v&limit=%v", current_page_host_and_path, next_page, limit)
		pagination.NextPage = next_page

		pagination.LastPageUrl =
			fmt.Sprintf("%spage=%v&limit=%v", current_page_host_and_path, total_pages, limit)
	}

	pagination.Limit = limit
	pagination.PerPage = limit
	pagination.Total = total
	pagination.CurrentPage = page

	return pagination, nil
}

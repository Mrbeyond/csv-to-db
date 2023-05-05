package controller

import (
	"csvapi-test/model"
	"csvapi-test/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Fetch(c *gin.Context) {
	var (
		ohlcs            []model.Ohcl
		db               = model.DB
		search           = c.Query("search")
		pagination       services.Pagination
		isFullPagination = strings.ToLower(c.Query("ptype")) == "full" // pagination type
	)

	// Compute full search on the request search query if set
	if search != "" && !strings.HasPrefix(db.Dialector.Name(), "sqlite") {
		db = db.Where(
			`to_tsvector(
					'english', unix || ' ' || symbol || ' ' || open || ' ' || high  || ' ' ||  low  || ' ' || close  
					) @@ to_tsquery('english', ?)`,
			search)
	}

	paginationQueries := &services.PaginationParams{}
	paginationQueries.ParseQuery(c)

	// Add full pagination object to the result
	if isFullPagination {
		var total int64

		if err := db.Model(model.Ohcl{}).Count(&total).Error; err != nil {
			services.ServerErrror(c, err, "")
			return
		}

		paginationP, err := services.Paginate(c, *paginationQueries, int(total))
		if err != nil {
			services.ServerErrror(c, err, "")
			return
		}
		pagination = *paginationP
	}

	if err := db.Limit(paginationQueries.Limit).
		Offset(paginationQueries.Offset).
		Find(&ohlcs).Error; err != nil {
		services.ServerErrror(c, err, "")
		return
	}

	response := gin.H{
		"status":  "success",
		"message": "Data successfully fetched",
		"data":    ohlcs,
	}
	if isFullPagination {
		response["pagination"] = pagination
	}
	c.JSON(http.StatusOK, response)
}

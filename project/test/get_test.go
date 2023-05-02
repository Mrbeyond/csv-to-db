package test

import (
	"csvapi-test/model"
	"csvapi-test/router"
	"csvapi-test/services"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SimpleResponse struct {
	Data    []OHLC `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (res *SimpleResponse) MustContainSearch(t *testing.T, search string) bool {
	valid := true
	for _, row := range res.Data {
		record := []string{
			fmt.Sprintf("%d", row.UNIX),
			row.SYMBOL,
			fmt.Sprintf("%f", row.OPEN),
			fmt.Sprintf("%f", row.HIGH),
			fmt.Sprintf("%f", row.LOW),
			fmt.Sprintf("%f", row.CLOSE),
		}
		valid = assert.Contains(t, record, search)
		if !valid {
			break
		}
	}
	return valid
}

type FullPaginationResponse struct {
	Data       []OHLC              `json:"data"`
	Message    string              `json:"message"`
	Status     string              `json:"status"`
	Pagination services.Pagination `json:"pagination"`
}

func TestGetMax100(t *testing.T) {
	model.DbConfig("TESTING")
	router := router.AppInstance()

	var response SimpleResponse

	// client := &http.Client{
	// 	Timeout: time.Second * 10,
	// }

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/data?limit=100", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(w, req)

	err = json.NewDecoder(w.Body).Decode(&response)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code, "Status code must be 200")
	assert.Equal(t, response.Status, "success", "Response status must be success")
	assert.LessOrEqual(t, len(response.Data), 100, "Length data must not be greater than 100")

	// resp, err := client.Do(req)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer resp.Body.Close()

	// assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code must be 200")
	// assert.Equal(t, response.Status, "success", "Response status must be success")
	// assert.LessOrEqual(t, len(response.Data), 100, "Length data must not be greater than 100")
}

// func TestGetWithSearch(t *testing.T) {
// 	search := fmt.Sprintf("%d", 1644719460000)
// 	url := fmt.Sprintf("%s:%s/data?limit=20&search=%s", baseUrl, port, search)
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	client := &http.Client{
// 		Timeout: time.Second * 10,
// 	}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer resp.Body.Close()

// 	var response SimpleResponse
// 	err = json.NewDecoder(resp.Body).Decode(&response)
// 	assert.Nil(t, err)
// 	assert.True(t, response.MustContainSearch(t, search))
// 	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code must be 200")
// 	assert.Equal(t, response.Status, "success", "Response status must be success")
// 	assert.LessOrEqual(t, len(response.Data), 20, "Length data must not be greater than 20")
// }

// func TestGetWithFullPagination(t *testing.T) {
// 	url := fmt.Sprintf("%s:%s/data?limit=10&ptype=full", baseUrl, port)
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	client := &http.Client{
// 		Timeout: time.Second * 10,
// 	}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer resp.Body.Close()
// 	var response FullPaginationResponse

// 	err = json.NewDecoder(resp.Body).Decode(&response)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, response.Pagination)
// 	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code must be 200")
// 	assert.Equal(t, response.Status, "success", "Response status must be success")
// 	assert.LessOrEqual(t, len(response.Data), 100, "Length data must not be greater than 100")
// }

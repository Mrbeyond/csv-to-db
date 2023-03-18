package test

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	port    = "8090"             // Port must be specific port used on the app container
	baseUrl = "http://127.0.0.1" // Base url for the app instance on the host the container
)

type OHLC struct {
	UNIX   uint64  `json:"unix"`
	SYMBOL string  `json:"symbol"`
	OPEN   float32 `json:"open"`
	HIGH   float32 `json:"high"`
	LOW    float32 `json:"low"`
	CLOSE  float32 `json:"close"`
}

var (
	// Expected data sample
	fields = []OHLC{
		{UNIX: 1644719700000, SYMBOL: "BTCUSDT", OPEN: 42123.29000000, HIGH: 42148.32000000, LOW: 42120.82000000, CLOSE: 42146.06000000},
		{UNIX: 1644719640000, SYMBOL: "BTCUSDT", OPEN: 42113.08000000, HIGH: 42126.32000000, LOW: 42113.07000000, CLOSE: 42123.30000000},
		{UNIX: 1644719580000, SYMBOL: "BTCUSDT", OPEN: 42120.80000000, HIGH: 42130.23000000, LOW: 42111.01000000, CLOSE: 42113.07000000},
		{UNIX: 1644719520000, SYMBOL: "BTCUSDT", OPEN: 42114.47000000, HIGH: 42123.31000000, LOW: 42102.22000000, CLOSE: 42120.80000000},
		{UNIX: 1644719460000, SYMBOL: "BTCUSDT", OPEN: 42148.23000000, HIGH: 42148.24000000, LOW: 42114.04000000, CLOSE: 42114.48000000},
	}
	//Expected data header
	csvHeader = []string{"UNIX", "SYMBOL", "OPEN", "HIGH", "LOW", "CLOSE"}
)

type CreateResponse struct {
	Data struct {
		CsvLinesRead   string `json:"csvLinesRead"`
		TotalSavedRows string `json:"totalSavedRows"`
	} `json:"data"`
	Status string `json:"status"`
}

func TestSaveSCV(t *testing.T) {
	file, err := os.CreateTemp("", "csv_file*.csv")
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	defer file.Close()
	defer os.Remove(file.Name())

	max := 50000 // max number of rows to generate

	// write csv into the temp file
	csvWriter := csv.NewWriter(file)
	// Write header
	if err = csvWriter.Write(csvHeader); err != nil {
		t.Fatal(err)
	}
	// Add rows
	for i := 0; i < max/5; i++ {
		for _, field := range fields {
			record := []string{
				fmt.Sprintf("%d", field.UNIX),
				field.SYMBOL,
				fmt.Sprintf("%f", field.OPEN),
				fmt.Sprintf("%f", field.HIGH),
				fmt.Sprintf("%f", field.LOW),
				fmt.Sprintf("%f", field.CLOSE),
			}
			if err = csvWriter.Write(record); err != nil {
				t.Fatal()
			}
		}
	}
	csvWriter.Flush()
	file.Seek(0, 0)
	// Create a new multipart/form-data request
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	mPart, err := writer.CreateFormFile("csv_file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}

	//copy csv content into mPart
	_, err = io.Copy(mPart, file)
	if err != nil {
		t.Fatalf("Error copying file to form file: %v", err)
	}
	writer.Close()

	url := fmt.Sprintf("%s:%s/data", baseUrl, port)
	// Create a new request to the create endpoint with the test file attached
	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: time.Minute * 3,
	}
	// Use the endpoint handler to process the request and capture the response
	beforeRequest := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	timeDiff := time.Now().Sub(beforeRequest).Seconds()
	defer resp.Body.Close()

	var response CreateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	t.Log(response, timeDiff)
	assert.Nil(t, err, "Invalid response type")
	assert.NotNil(t, response, "Response must not ne nil")
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Status code must be 201")
	assert.Equal(t, response.Status, "success", "Response status must be success")
	assert.NotEmpty(t, response.Data.CsvLinesRead, "Lines read must be returned")
	assert.Equal(t, response.Data.CsvLinesRead, response.Data.TotalSavedRows, "Not all files are saved")

	assert.LessOrEqual(t, timeDiff, float64(15), "Request time should not be more than 15 seconds")
}

// Test for non csv file uploaded
func TestNonCsvFileType(t *testing.T) {
	// Create a temp test file
	file, err := os.CreateTemp("", "csv_file*.txt")
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	defer os.Remove(file.Name())

	_, err = file.WriteString("Just random string")
	if err != nil {
		t.Fatalf("Error creating writing string into test file: %v", err)
	}
	defer file.Close()

	// Create a new multipart/form-data request
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	mPart, err := writer.CreateFormFile("csv_file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}

	//copy csv content into mPart
	_, err = io.Copy(mPart, file)
	if err != nil {
		t.Fatalf("Error copying file to form file: %v", err)
	}
	writer.Close()
	url := fmt.Sprintf("%s:%s/data", baseUrl, port)
	req, err := http.NewRequest("POST", url, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	bodyString, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var response CreateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	assert.Contains(t, string(bodyString), "Expected a csv file",
		"Invalid response for wrong file type ",
	)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// Test for invalid csv header against expected format
func TestUnexpectedCvsFormat(t *testing.T) {
	// Create a temp test file
	file, err := os.CreateTemp("", "csv_file.csv")
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	max := 50 //Max number of csv rows to create

	// write csv into the temp file
	csvWriter := csv.NewWriter(file)
	// Write header Wrong header
	csvWrongHeader := []string{"UNI", "SYMBOL", "OPN", "HGH", "LOW", "CLOSE"}
	if err = csvWriter.Write(csvWrongHeader); err != nil {
		t.Fatal()
	}
	// Add rows
	for i := 0; i < max/5; i++ {
		for _, field := range fields {
			record := []string{
				fmt.Sprintf("%d", field.UNIX),
				field.SYMBOL,
				fmt.Sprintf("%f", field.OPEN),
				fmt.Sprintf("%f", field.HIGH),
				fmt.Sprintf("%f", field.LOW),
				fmt.Sprintf("%f", field.CLOSE),
			}
			if err = csvWriter.Write(record); err != nil {
				t.Fatal()
			}
		}
	}
	defer csvWriter.Flush()
	file.Seek(0, 0)

	// Create a new multipart/form-data request
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	mPart, err := writer.CreateFormFile("csv_file", filepath.Base(file.Name()))
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}

	//copy csv content into mPart
	_, err = io.Copy(mPart, file)
	if err != nil {
		t.Fatalf("Error copying file to form file: %v", err)
	}
	defer writer.Close()
	// Create a new request to the create endpoint with the test file attached
	url := fmt.Sprintf("%s:%s/data", baseUrl, port)
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		t.Fatalf("Error From client request %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	bodyString, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bodyString))

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
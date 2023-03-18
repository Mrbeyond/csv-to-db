package controller

import (
	"bufio"
	"context"
	"csvapi-test/model"
	"csvapi-test/services"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	chunkVolume = 4000    // Size of slice of ohcl to save to the db as batches
	MB4         = 4 << 20 // File size factor for worker pool
)

// Object holding csv rows in chunks and channel flow for saving the chunks in batches
type ProcessPool struct {
	chunk           []model.Ohcl
	dataChan        chan []model.Ohcl
	wg              *sync.WaitGroup
	numWorkers      int //Number of worker to process db insertion from dbChannel
	errorMessage    string
	csvLinesRead    int // Total number of lines read
	totalChunkSaved int
	done            bool // Checks if every lines has been saved
	mutex           sync.Mutex
}

// Read through the entire csv rows and append them in chunks into the db channel
func (processPool *ProcessPool) generateCsvChunk(csvReader *csv.Reader) {
	// Read csv so far there's no error saving or reading into the chunks
	for {
		// Read the line of csv reader
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			processPool.errorMessage = err.Error()
			close(processPool.dataChan) // close data channel
			break
		}
		// Check and valid number of columns, strictly based on the expected data
		if len(row) != 6 {
			processPool.errorMessage = "Invalid row detected"
			processPool.chunk = make([]model.Ohcl, 0) // empty the chunk
			close(processPool.dataChan)               // close data channel
			break
		}

		// Convert row to Ohcl
		uix, err := strconv.ParseUint(row[0], 10, 64)
		symbol := row[1]
		open, err := strconv.ParseFloat(row[2], 32)
		high, err := strconv.ParseFloat(row[3], 32)
		low, err := strconv.ParseFloat(row[4], 32)
		close, err := strconv.ParseFloat(row[5], 32)

		ohlc := model.Ohcl{
			UNIX:   uix,
			SYMBOL: symbol,
			OPEN:   float32(open),
			HIGH:   float32(high),
			LOW:    float32(low),
			CLOSE:  float32(close),
		}
		processPool.chunk = append(processPool.chunk, ohlc)
		// check if the lenght of the rows equal to chunkVolume then send it to db channel
		if len(processPool.chunk) == chunkVolume {
			processPool.csvLinesRead += chunkVolume
			processPool.dataChan <- processPool.chunk
			// empty the chunk
			processPool.chunk = make([]model.Ohcl, 0)
		}
	}
	// check for remant of rows if not up to and checked by chunkVolume
	if len(processPool.chunk) > 0 {
		processPool.csvLinesRead += len(processPool.chunk)
		if processPool.errorMessage == "" {
			processPool.dataChan <- processPool.chunk
		}
	}
}

// Immplement worker pool to save csv chunks into the db
func (processPool *ProcessPool) processCsvChunk(db *gorm.DB) {

	for i := 0; i < processPool.numWorkers; i++ {
		go func() {
			defer processPool.wg.Done()
			for rows := range processPool.dataChan {
				processPool.mutex.Lock()
				if err := db.Create(rows).Error; err != nil {
					processPool.errorMessage = err.Error()
					close(processPool.dataChan) // close data channel
					processPool.mutex.Unlock()
					return
				}
				processPool.mutex.Unlock()

				// check if all the scv lines has been saved
				processPool.mutex.Lock()
				processPool.totalChunkSaved += len(rows)
				fmt.Println(processPool.totalChunkSaved)
				if processPool.csvLinesRead == processPool.totalChunkSaved {
					processPool.done = true
					close(processPool.dataChan) //close data channel
					processPool.mutex.Unlock()
					return // All csv lines are saved
				}
				processPool.mutex.Unlock()

			}
		}()
	}
}

func Create(c *gin.Context) {
	// Increase the context timeout incase of a very large csv file to.
	_, cancel := context.WithTimeout(c.Request.Context(), 3*time.Minute)
	defer cancel()

	file, err := c.FormFile("csv_file")
	if err != nil {
		services.BadRequestErrror(c, err, ": Form cannot be parsed")
		return
	}

	if strings.ToLower(filepath.Ext(file.Filename)) != ".csv" {
		services.BadRequestErrror(c, nil, "Expected a csv file")
		return
	}

	// Get an io.Reader for the file contents using file.Open()
	fileContent, err := file.Open()
	if err != nil {
		services.ServerErrror(c, err, "")
		return
	}
	defer fileContent.Close()

	// Bufio is used to efficiently handle reading of a large file
	bufioReader := bufio.NewReader(fileContent)
	csvReader := csv.NewReader(bufioReader)

	// Read first row to remove the header
	header, err := csvReader.Read()
	if err != nil {
		services.ServerErrror(c, err, ":occurs when reading the file header")
		return
	} else if !validateSCVHeader(header) {
		services.BadRequestErrror(c, nil, "Csv file is not valid")
		return
	}

	// Set max of 20 additional workers for big files
	ff := file.Size / MB4
	worKerFactor := math.Min(float64(ff), 20)

	// number of workers for worker pool from system CPU available
	numWorkers := runtime.NumCPU() + int(worKerFactor)

	var (
		db = model.DB     // Database instance
		wg sync.WaitGroup // wait group syncer for workpool
	)

	// Disable logging and default transaction for the database session
	db = db.Session(&gorm.Session{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true, //disable default transaction to help speed
	})

	tx := db.Begin() // Perform the saving with explicit transaction to enable rollback

	wg.Add(numWorkers)
	processPool := ProcessPool{
		wg:         &wg,
		chunk:      make([]model.Ohcl, 0, chunkVolume),
		dataChan:   make(chan []model.Ohcl, numWorkers),
		numWorkers: numWorkers,
	}

	// go processPool.processCsvChunk(db) //Use worker pool to save csv in chunks
	go processPool.processCsvChunk(tx)      //Use worker pool to save csv in chunks
	processPool.generateCsvChunk(csvReader) // Read word scv rows into chunks

	// lock flow until all workers are done
	wg.Wait()

	// Check if theres no error for worker pool and commit transaction
	if processPool.errorMessage == "" && processPool.done {
		if err := tx.Commit().Error; err != nil {
			services.ServerErrror(c, err, "")
			return
		}
	} else {
		tx.Rollback() // rollback the transaction
		services.ServerErrror(c, nil, processPool.errorMessage)
		return
	}

	data := map[string]string{
		"csvLinesRead":   fmt.Sprintf("%d", processPool.csvLinesRead),
		"totalSavedRows": fmt.Sprintf("%d", processPool.totalChunkSaved),
	}
	response := gin.H{
		"status": "success",
		"data":   data,
	}
	c.JSON(http.StatusCreated, response)
}

func validateSCVHeader(header []string) (valid bool) {
	expectedHeader := []string{"UNIX", "SYMBOL", "OPEN", "HIGH", "LOW", "CLOSE"}
	valid = true
	for index, value := range expectedHeader {
		if value != strings.ToUpper(header[index]) {
			valid = false
			break
		}
	}
	return valid && len(header) == len(expectedHeader)
}

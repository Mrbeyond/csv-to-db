package services

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
)

// Checks if error is not nil and return error message else return empty string
func ErrorExists(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

// Compute 500 Server Error response
func ServerErrror(c *gin.Context, err error, extra string) {
	response := gin.H{
		"status":  "failed",
		"error":   true,
		"message": ErrorExists(err) + " " + extra,
	}

	c.JSON(http.StatusInternalServerError, response)
}

// Compute 400 Bad Request Error response
func BadRequestErrror(c *gin.Context, err error, extra string) {
	response := gin.H{
		"status":  "failed",
		"error":   true,
		"message": ErrorExists(err) + " " + extra,
	}

	c.JSON(http.StatusBadRequest, response)
}

// Compute 404 Not Found Error response
func NotFoundError(c *gin.Context, err error, extra string) {
	response := gin.H{
		"status":  "failed",
		"error":   true,
		"message": ErrorExists(err) + " " + extra,
	}

	c.JSON(http.StatusNotFound, response)
}

// Compute 401 Unauthorized Error response
func UnauthorizedError(c *gin.Context, err error, extra string) {
	response := gin.H{
		"status":  "failed",
		"error":   true,
		"message": ErrorExists(err) + " " + extra,
	}

	c.JSON(http.StatusUnauthorized, response)
}

// Compute 40 Forbidden Error response
func ForbiddenError(c *gin.Context, err error, extra string) {
	response := gin.H{
		"status":  "failed",
		"error":   true,
		"message": ErrorExists(err) + " " + extra,
	}

	c.JSON(http.StatusForbidden, response)
}

/** Check error from *gorm.DB operation or query, send error response if found and return true */
func GormQueryErrorCheck(c *gin.Context, result *gorm.DB, model, customMessage string) (found bool) {
	var errorMessage string
	if customMessage != "" {
		errorMessage = customMessage
	} else {
		errorMessage = model + "not found"
	}

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			NotFoundError(c, nil, errorMessage)
		} else {
			ServerErrror(c, result.Error, "")
		}
		return true
	} else if result.RowsAffected == 0 {
		NotFoundError(c, nil, errorMessage)
		return true
	}
	return
}

// Validates the request payload with the model struct
func AbortWithRequestError(c *gin.Context, model interface{}, err error) {
	//model is a pointer
	var ers validator.ValidationErrors
	errorMessage := make([]map[string]string, 0)
	if errors.As(err, &ers) {
		for _, er := range err.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(model).Elem().FieldByName(er.Field())
			jsonTag := string(field.Tag.Get("json"))
			errorTag := er.Tag()
			errorMessage = append(errorMessage, BindErrorResolver(jsonTag, errorTag, er.Param()))
		}
	} else {
		field := reflect.TypeOf(model).Elem().Name()
		errorMessage = append(errorMessage, map[string]string{field: err.Error()})
	}

	response := gin.H{
		"status":  "failed",
		"error":   true,
		"message": errorMessage,
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, response)
}

// Bind appropriate message text to error type based on the struct tag validation
func BindErrorResolver(jsonTagName, errorTag, errorParam string) (message map[string]string) {
	message = map[string]string{}
	switch errorTag {
	case "gte":
		message[jsonTagName] = fmt.Sprintf("%s must be greater than or equal to %s", jsonTagName, errorParam)
		return
	case "lte":
		message[jsonTagName] = fmt.Sprintf("%s must be less than or equal to %s", jsonTagName, errorParam)
		return
	case "required":
		message[jsonTagName] = fmt.Sprintf("%s is required", jsonTagName)
		return
	case "email":
		message[jsonTagName] = fmt.Sprintf("%s must be a valid email address", jsonTagName)
		return
	case "oneof":
		message[jsonTagName] = fmt.Sprintf("%s must be one of %s", jsonTagName, strings.Join(strings.Split(errorParam, " "), ", "))
		return

	}
	return
}

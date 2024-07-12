package twilio

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type jsonResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

var validate = validator.New()

// validateBody method declaration

func (app *Config) validateBody(c *gin.Context, data any) error {

	// validates the JSON request body (c.BindJSON(&data)) and the struct (validate.Struct(&data)) using the validate instance.

	if err := c.BindJSON(&data); err != nil {
		return err
	}

	if err := validate.Struct(&data); err != nil {
		return err
	}

	return nil
}

// writeJSON method declaration
// writeJSON is a method that sends a JSON response (c.JSON) using the provided HTTP

func (app *Config) writeJSON(c *gin.Context, status int, data any) {
	c.JSON(status, jsonResponse{Status: status, Message: "success", Data: data})
}

// errorJSON method declaration
// err parameter indicating the error that occurred and an optional status parameter to specify the HTTP status code
// formats the response using the jsonResponse struct with the provided status code and error message

func (app *Config) errorJSON(c *gin.Context, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	c.JSON(statusCode, jsonResponse{Status: statusCode, Message: err.Error()})
}

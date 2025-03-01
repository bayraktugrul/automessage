package errors

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	required = "required"
	oneOf    = "oneof"
)

func ValidationError(c *gin.Context, err error) {
	respondWithError(c, http.StatusBadRequest, err)
}

func InternalServerError(c *gin.Context, err error) {
	respondWithError(c, http.StatusInternalServerError, err)
}

func handleValidationError(err error) ValidationErrorResponse {
	var response ValidationErrorResponse

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()

			switch tag {
			case required:
				response.Message = field + " field is required"
			case oneOf:
				param := e.Param()
				response.Message = field + " must be one of: " + param
			default:
				response.Message = err.Error()
			}

			break
		}
	} else {
		response.Message = err.Error()
	}

	return response
}

func respondWithError(c *gin.Context, statusCode int, err error) {
	var response interface{}

	switch {
	case statusCode == http.StatusInternalServerError:
		response = ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "Internal server error",
		}
	default:
		if _, ok := err.(validator.ValidationErrors); ok || strings.Contains(err.Error(), "json") {
			validationResponse := handleValidationError(err)
			response = ErrorResponse{
				Code:  statusCode,
				Error: validationResponse.Message,
			}
		} else {
			response = ErrorResponse{
				Code:  statusCode,
				Error: err.Error(),
			}
		}
	}

	c.JSON(statusCode, response)
}

package utils

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"
)

func UserDataValidationError(err error) (string, int) {
	var errorMessage string
	errorStatus := http.StatusBadRequest
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			switch e.Field() {
			case "Email":
				if e.Tag() == "required" {
					errorMessage = "Email is required"
				} else if e.Tag() == "email" {
					errorMessage = "Invalid email format"
				}
			case "Username":
				if e.Tag() == "required" {
					errorMessage = "Username is required"
				}
			default:
				errorMessage = fmt.Sprintf("Validation failed on field '%s' for condition '%s'", e.Field(), e.Tag())
			}
		}
	}
	return errorMessage, errorStatus
}

func UserDataInsertError(err error) (string, int) {
	var errorStatus int
	var errorMessage string
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == "23505" { // Unique violation code
			errorStatus = http.StatusConflict
			errorMessage = "Username already exists"
		} else {
			errorStatus = http.StatusInternalServerError
			errorMessage = fmt.Sprintf("Database error: %v", pgErr.Code)
		}
	} else {
		log.Printf("Unexpected error: %v (type: %v)", err, reflect.TypeOf(err))
		errorStatus = http.StatusInternalServerError
		errorMessage = "Could not create user"
	}
	return errorMessage, errorStatus
}

package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type CreateEventRequest struct {
	Title       string    `json:"title" validate:"required,max=100" example:"This is a title"`
	Description *string   `json:"description" example:"This is a description"`
	StartTime   time.Time `json:"start_time" validate:"required" example:"2025-10-17T00:00:00Z"`
	EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime" example:"2025-10-20T10:00:00Z"`
}

func (e *CreateEventRequest) Validate() map[string]string {
	validate := validator.New()
	err := validate.Struct(e)

	errors := make(map[string]string)

	if err != nil {
		errors["struct"] = fmt.Sprintf("Invalid struct %s", err)
	}

	if err == nil {
		return nil
	}

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()

		switch tag {
		case "required":
			errors[field] = fmt.Sprintf("Field %s is required", field)
		case "max":
			errors[field] = fmt.Sprintf("Field %s cannot exceed %s characters", field, err.Param())
		case "gtfield":
			errors[field] = fmt.Sprintf("Field %s must be greater than %s", field, err.Param())
		default:
			errors[field] = fmt.Sprintf("Validation error %s: %s", field, tag)
		}
	}

	return errors
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

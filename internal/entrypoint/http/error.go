package http

import (
	"errors"
	"event-api/internal/domain/model"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func handlerErrorResponse(w http.ResponseWriter, err error) {
	var httpStatus int
	var validationError *model.ValidationError
	var notFoundError *model.NotFoundError
	log.Println(err)
	switch {
	case errors.As(err, &validationError):
		httpStatus = http.StatusBadRequest
	case errors.As(err, &notFoundError):
		httpStatus = http.StatusNotFound
	default:
		httpStatus = http.StatusInternalServerError
	}

	respondWithJSON(w, httpStatus, ErrorResponse{
		Error: err.Error(),
	})
}

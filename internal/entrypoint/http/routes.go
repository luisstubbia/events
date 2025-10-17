package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	createEvent(w http.ResponseWriter, r *http.Request)
	getEvent(w http.ResponseWriter, r *http.Request)
	getAllEvents(w http.ResponseWriter, r *http.Request)
	healthCheck(w http.ResponseWriter, r *http.Request)
}

func SetupRoutes(h Handler) *mux.Router {
	router := mux.NewRouter()
	router.Use(corsMiddleware)
	router.HandleFunc("/events", h.getAllEvents).Methods("GET")
	router.HandleFunc("/events", h.createEvent).Methods("POST")
	router.HandleFunc("/events/{id}", h.getEvent).Methods("GET")
	router.HandleFunc("/health", h.healthCheck).Methods("GET")

	router.PathPrefix("/swagger/").Handler(swaggerHandler())
	router.HandleFunc("/swagger/doc.json", swaggerJSON)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
	})
	return router
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"bugtracker-backend/internal/db"
	"bugtracker-backend/internal/models"

	"github.com/gorilla/mux"
)

func RegisterCommentRoutes(r *mux.Router) {
	r.HandleFunc("/bugs/{id}/comments", GetComments).Methods("GET")
	r.HandleFunc("/bugs/{id}/comments", CreateComment).Methods("POST")
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("Getting comments for bug %s", id)

	if _, err := strconv.Atoi(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid bug ID format",
		})
		return
	}

	comments, err := db.GetComments(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "bug not found" {
			status = http.StatusNotFound
		} else if err.Error() == "invalid bug ID format" {
			status = http.StatusBadRequest
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("Creating comment for bug %s", id)

	if _, err := strconv.Atoi(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid bug ID format",
		})
		return
	}

	var req models.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}

	if err := req.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	comment := &models.Comment{
		Content: req.Content,
		Author:  req.Author,
	}

	if err := db.CreateComment(id, comment); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "bug not found" {
			status = http.StatusNotFound
		} else if err.Error() == "invalid bug ID format" {
			status = http.StatusBadRequest
		}
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

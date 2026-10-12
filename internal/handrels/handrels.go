package handrels

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dinbgv/todoAppDanil/internal/database"
	"github.com/Dinbgv/todoAppDanil/internal/models"
)

type Handrels struct {
	store *database.TaskStore
}

func NewHandrels(store *database.TaskStore) *Handrels {
	return &Handrels{store: store}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, massage string) {
	respondWithJSON(w, statusCode, map[string]string{"error": massage})
}

func (h *Handrels) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		log.Printf("GetAll error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error get all tasks")
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handrels) GetTaskById(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idString := pathParts[0]

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Task ID")
		return
	}

	task, err := h.store.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handrels) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Missing title")
		return
	}

	task, err := h.store.Create(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func (h *Handrels) UpdateTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idString := pathParts[0]

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Task ID")
		return
	}

	var input models.UpdateTaskInput

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Missing title")
		return
	}

	task, err := h.store.Update(id, input)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusNotFound, "Task not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handrels) DeleteTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idString := pathParts[0]
	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Task ID")
		return
	}
	err = h.store.Delete(id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusNotFound, "Task not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

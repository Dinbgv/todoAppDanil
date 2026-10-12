package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Dinbgv/todoAppDanil/internal/database"
	"github.com/Dinbgv/todoAppDanil/internal/handrels"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://tasksuser:taskspass@localhost:5432/tasksbd?sslmode=disable"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Starting server on port %s", serverPort)

	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Database connected")

	taskstore := database.NewTaskStore(db)

	handrel := handrels.NewHandrels(taskstore)

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", methodHandrel(handrel.GetAllTasks, "GET"))
	mux.HandleFunc("/tasks/create", methodHandrel(handrel.CreateTask, "POST"))

	mux.HandleFunc("/tasks/", taskIDHandler(handrel))

	loggerMux := loggingMiddleware(mux)

	serverAddr := ":" + serverPort

	err = http.ListenAndServe(serverAddr, loggerMux)
	if err != nil {
		log.Fatal(err)
	}
}

func methodHandrel(handrelFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		handrelFunc(w, r)
	}
}

func taskIDHandler(handrel *handrels.Handrels) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handrel.GetTaskById(w, r)
		case http.MethodPut:
			handrel.UpdateTask(w, r)
		case http.MethodDelete:
			handrel.DeleteTask(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

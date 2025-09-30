package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"task_manager_server/internal/application/dto/request"
	"task_manager_server/internal/application/usecase"
	"task_manager_server/internal/domain/entites"
	"task_manager_server/internal/infrastructure/middlewares"
	"task_manager_server/internal/infrastructure/persistence/repository"
	"task_manager_server/pkg/security"

	"github.com/gorilla/mux"
)

type TaskHandler struct {
	useCase *usecase.TaskUseCase
}

func NewTaskHandler(useCase *usecase.TaskUseCase) *TaskHandler { return &TaskHandler{useCase: useCase} }

func (t *TaskHandler) InitRoutes(r *mux.Router) {
	taskAPI := r.PathPrefix("/tasks").Subrouter()
	taskAPI.Use(middlewares.AuthJWTMiddleware)

	taskAPI.HandleFunc("/create", t.Create).Methods("POST")
	taskAPI.HandleFunc("", t.GetTaskList).Methods("GET")
	taskAPI.HandleFunc("", t.UpdateStatus).Methods("PATCH")
	taskAPI.HandleFunc("", t.DeleteTask).Methods("DELETE")
	//taskAPI.HandleFunc("/drop", t.).Methods("DELETE")

}

func (t *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	reqBody := new(request.CreateTaskResponse)

	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		log.Println(err)
		http.Error(w, "invalid request body, failed create task", http.StatusBadRequest)
		return
	}
	userID, err := security.GetIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	err = t.useCase.Create(
		r.Context(),
		reqBody.Title,
		reqBody.Description,
		reqBody.StatusID,
		userID,
	)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "task successfully created")
}

func (t *TaskHandler) GetTaskList(w http.ResponseWriter, r *http.Request) {
	var taskSlice []*entites.Task
	userID, err := security.GetIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	if status := r.URL.Query().Get("status"); status != "" {
		statusID, err := strconv.Atoi(status)
		if err != nil {
			http.Error(w, "invalid query param", http.StatusBadRequest)
			return
		}
		taskSlice, err = t.useCase.GetFilteredList(r.Context(), userID, statusID)
	} else {
		taskSlice, err = t.useCase.GetList(r.Context(), userID)
	}

	if err != nil {
		if errors.Is(err, repository.RepoErr) {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(taskSlice)
}

func (t *TaskHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID, err := security.GetIDFromContext(r.Context())

	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	title := r.URL.Query().Get("title")
	statusID, err := strconv.Atoi(r.URL.Query().Get("status"))

	if title == "" || err != nil {
		http.Error(w, "not found valid url param", http.StatusBadRequest)
		return
	}
	if err = t.useCase.UpdateStatus(r.Context(), title, userID, statusID); err != nil {
		if errors.Is(err, entites.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Task's status successfully updated"))
}

func (t *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID, err := security.GetIDFromContext(r.Context())
	if err != nil {
		log.Printf("context err: %v\n", err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "task name should not be empty", http.StatusBadRequest)
		return
	}

	id, err := t.useCase.DeleteTask(r.Context(), title, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if id == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("task with that title just not exist"))
		return
	}
	w.Write([]byte("task successfully deleted"))
}

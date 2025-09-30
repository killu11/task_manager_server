package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"task_manager_server/internal/application/dto/request"
	"task_manager_server/internal/application/dto/response"
	"task_manager_server/internal/application/usecase"
	"task_manager_server/internal/domain/entites"
	"task_manager_server/internal/infrastructure/persistence/repository"
	"task_manager_server/pkg/security"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	useCase *usecase.UserUseCase
}

// go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/main.go -d ./cmd,./internal/infrastructure/handlers --parseDependency --parseInternal
func NewUserHandler(uc *usecase.UserUseCase) *UserHandler { return &UserHandler{useCase: uc} }

func (h *UserHandler) InitRoutes(r *mux.Router) {
	r.HandleFunc("/sign-up", h.Register).Methods("POST")
	r.HandleFunc("/sign-in", h.Login).Methods("POST")
}

// Register godoc
// @Summary      Регистрация пользователя
// @Description  Создаёт пользователя и возвращает JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      request.RegisterRequest  true  "Данные регистрации"
// @Success      201      {object}  response.RegisterResponse
// @Failure      400      {string}  string "Invalid request body"
// @Failure      405      {string}  string "Unsupported method"
// @Failure      409      {string}  string "User already exists"
// @Router       /sign-up [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported method, resource supported only POST-method", http.StatusMethodNotAllowed)
		return
	}

	userData := new(request.RegisterRequest)

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(userData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.useCase.Register(
		r.Context(),
		userData.Username,
		userData.Password,
	)

	if err != nil {
		matchError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(&response.RegisterResponse{
		Token:    token,
		Username: userData.Username,
	})
}

// Login godoc
// @Summary      Вход пользователя
// @Description  Проверяет входные данные и возвращает access-token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      request.LoginRequest  true  "Данные входа"
// @Success      200     {object}  response.LoginResponse
// @Failure      400      {string}  string "Invalid request body"
// @Failure 	 401 	  {string}  string ""User with this nickname not found""
// @Failure      405      {string}  string "Unsupported method"
// @Failure      409      {string}  string "User already exists"
// @Router       /sign-in [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
		return
	}
	requestBody := new(request.LoginRequest)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(requestBody)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	token, err := h.useCase.Auth(r.Context(), requestBody.Username, requestBody.Password)

	if err != nil {
		matchError(err, w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response.LoginResponse{Token: token, Username: requestBody.Username})
}

func matchError(err error, w http.ResponseWriter) {
	switch {
	case errors.As(err, &entites.InvalidCredErr):
		if errors.Is(err, entites.NameAlreadyUsingErr) {
			http.Error(w, "Name already using", http.StatusConflict)
			return
		}
		http.Error(w, "Invalid credentials", http.StatusBadRequest)

	case errors.As(err, &repository.RepoErr):
		http.Error(w, "Register service error", http.StatusInternalServerError)

	case errors.Is(err, security.HashPassErr):
		http.Error(w, "Register service error", http.StatusInternalServerError)

	case errors.Is(err, security.VerifyPassErr):
		http.Error(w, "Invalid password", http.StatusUnauthorized)

	case errors.Is(err, security.GenerateJWTErr):
		http.Error(w, "Register service error", http.StatusInternalServerError)

	case errors.Is(err, entites.UserNotExistErr):
		http.Error(w, "User with this nickname not found", http.StatusUnauthorized)

	default:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

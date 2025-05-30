package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/your-username/golang-ecommerce-app/models"
	"github.com/your-username/golang-ecommerce-app/services"
	"github.com/your-username/golang-ecommerce-app/utils"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) SignupUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserId   string `json:"userId"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	newUser := models.User{
		UserId:    body.UserId,
		Password:  body.Password,
		Email:     body.Email,
		Role:      "user",
		CreatedAt: time.Now(),
	}
	if body.UserId == "" || body.Password == "" || body.Email == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "userId, password and email are required")
		return
	}

	result, err := uc.userService.Signup(r.Context(), newUser)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, result)
}

func (uc *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserId   string `json:"userId"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.UserId == "" || body.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "userId and password are required")
		return
	}

	result, err := uc.userService.Login(r.Context(), body.UserId, body.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, result)
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	if userId == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "userId is required")
		return
	}

	var updates map[string]any
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid update data")
		return
	}

	result, err := uc.userService.UpdateUserService(r.Context(), userId, updates)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, result)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	if userId == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "userId is required")
		return
	}

	if err := uc.userService.DeleteUserService(r.Context(), userId); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (uc *UserController) DeleteUserAllAccess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	if userId == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "userId is required")
		return
	}

	if err := uc.userService.DeleteUserServiceAllAccess(r.Context(), userId); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// updateUserRole updates the role of a user
func (uc *UserController) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	if userId == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "userId is required")
		return
	}

	var body struct {
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.Role == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "role is required")
		return
	}

	result, err := uc.userService.UpdateUserRoleService(r.Context(), userId, body.Role)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, result)
}
package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"restaurant-backend/src/models"
	"restaurant-backend/src/repositories"
	"restaurant-backend/src/utils"
	"strings"
	"time"
)

type AuthController struct {
	userRepo *repositories.UserRepository
}

func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{
		userRepo: repositories.NewUserRepository(db),
	}
}

func (ac *AuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := models.RegisterResponse{
			Success: false,
			Error:   "Invalid JSON format",
		}

		w.WriteHeader((http.StatusBadRequest))

		json.NewEncoder(w).Encode(response)
	}

	if err := ac.validateRegisterRequest(&req); err != nil {
		response := models.RegisterResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	exists, err := ac.userRepo.UserExists(req.Email)
	if err != nil {
		response := models.RegisterResponse{
			Success: false,
			Error:   "Database error",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if exists {
		response := models.RegisterResponse{
			Success: false,
			Error:   "User with this email already exists",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		response := models.RegisterResponse{
			Success: false,
			Error:   "Error processing password",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	user := &models.User{
		Name:     strings.TrimSpace(req.Name),
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: hashedPassword,
	}

	if err := ac.userRepo.CreateUser(user); err != nil {
		response := models.RegisterResponse{
			Success: false,
			Error:   "Error creating user",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenData := fmt.Sprintf("%s:%d:%s", user.Id.String(), time.Now().Unix(), os.Getenv("COOKIE_SECRET_KEY"))

	hashedToken := utils.HashString(tokenData)

	http.SetCookie(w, &http.Cookie{
		Name:     "dashboard-cookie",
		Value:    hashedToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400 * 7,
	})

	response := models.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		User:    user,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

func (ac *AuthController) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := models.RegisterResponse{
			Success: false,
			Error:   "Invalid JSON format",
		}

		w.WriteHeader((http.StatusBadRequest))

		json.NewEncoder(w).Encode(response)
	}

	if err := ac.validateLoginRequest(&req); err != nil {
		response := models.RegisterResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	exists, err := ac.userRepo.UserExists(req.Email)
	if err != nil {
		response := models.RegisterResponse{
			Success: false,
			Error:   "Database error",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !exists {
		response := models.RegisterResponse{
			Success: false,
			Error:   "User with this email not exists",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	existedUser, err := ac.userRepo.GetUserByEmail(req.Email)

	if err != nil {
		response := models.RegisterResponse{
			Success: false,
			Error:   "Error processing existed user password",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	isEqualPass := utils.IsPasswordEqualHash(existedUser.Password, req.Password)
	if !isEqualPass {
		response := models.RegisterResponse{
			Success: false,
			Error:   "Error invalid email or password",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenData := fmt.Sprintf("%s:%d:%s", existedUser.Id.String(), time.Now().Unix(), os.Getenv("COOKIE_SECRET_KEY"))

	hashedToken := utils.HashString(tokenData)

	http.SetCookie(w, &http.Cookie{
		Name:     "dashboard-cookie",
		Value:    hashedToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400 * 7,
	})

	response := models.RegisterResponse{
		Success: true,
		Message: "User login success",
		User:    existedUser,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (ac *AuthController) validateRegisterRequest(req *models.RegisterRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if len(strings.TrimSpace(req.Name)) < 2 {
		return fmt.Errorf("name must be at least 2 characters long")
	}
	if len(strings.TrimSpace(req.Name)) > 50 {
		return fmt.Errorf("name must be no more than 50 characters long")
	}

	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		return fmt.Errorf("invalid email format")
	}

	if strings.TrimSpace(req.Password) == "" {
		return fmt.Errorf("password is required")
	}
	if len(strings.TrimSpace(req.Password)) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}

	return nil
}

func (ac *AuthController) validateLoginRequest(req *models.LoginRequest) error {
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		return fmt.Errorf("invalid email format")
	}
	if strings.TrimSpace(req.Password) == "" {
		return fmt.Errorf("password is required")
	}
	if len(strings.TrimSpace(req.Password)) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}

	return nil
}

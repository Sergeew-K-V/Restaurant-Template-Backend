package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"restaurant-backend/src/repositories"
	"restaurant-backend/src/types"
	"restaurant-backend/src/utils"
	"strings"
)

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterResponseUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Id    string `json:"id"`
}

type RegisterResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	User    *RegisterResponseUser `json:"user,omitempty"`
	Error   string                `json:"error,omitempty"`
}

func (as *AuthService) RegisterUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var userRepo = repositories.NewUserRepository(as.db)
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := RegisterResponse{
			Success: false,
			Error:   "Invalid JSON format",
		}

		w.WriteHeader((http.StatusBadRequest))

		json.NewEncoder(w).Encode(response)
	}

	if err := as.validateRegisterRequest(&req); err != nil {
		response := RegisterResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	exists, err := userRepo.UserExists(req.Email)
	if err != nil {
		response := RegisterResponse{
			Success: false,
			Error:   "Database error",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if exists {
		response := RegisterResponse{
			Success: false,
			Error:   "User with this email already exists",
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		response := RegisterResponse{
			Success: false,
			Error:   "Error processing password",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	user := &repositories.User{
		Name:     strings.TrimSpace(req.Name),
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: hashedPassword,
	}

	if err := userRepo.CreateUser(user); err != nil {
		response := RegisterResponse{
			Success: false,
			Error:   "Error creating user",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenData := fmt.Sprintf("%s:%s", user.Id.String(), ctx.Value(types.CookieSecretKey))

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

	response := RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		User: &RegisterResponseUser{
			Name:  user.Name,
			Email: user.Email,
			Id:    user.Id.String(),
		},
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

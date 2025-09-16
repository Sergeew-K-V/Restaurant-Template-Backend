package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"restaurant-backend/src/repositories"
	"restaurant-backend/src/types"
	"restaurant-backend/src/utils"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponseUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Id    string `json:"id"`
}

type LoginResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	User    *LoginResponseUser `json:"user,omitempty"`
	Error   string             `json:"error,omitempty"`
}

func (as *AuthService) LoginUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginUser")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var userRepo = repositories.NewUserRepository(as.db)
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := LoginResponse{
			Success: false,
			Error:   "Invalid JSON format",
		}

		w.WriteHeader((http.StatusBadRequest))

		json.NewEncoder(w).Encode(response)
	}

	if err := as.validateLoginRequest(&req); err != nil {
		response := LoginResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	exists, err := userRepo.UserExists(req.Email)
	if err != nil {
		response := LoginResponse{
			Success: false,
			Error:   "Database error",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !exists {
		response := LoginResponse{
			Success: false,
			Error:   "User with this email not exists",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	existedUser, err := userRepo.GetUserByEmail(req.Email)

	if err != nil {
		response := LoginResponse{
			Success: false,
			Error:   "Error processing existed user password",
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	isEqualPass := utils.IsPasswordEqualHash(existedUser.Password, req.Password)
	if !isEqualPass {
		response := LoginResponse{
			Success: false,
			Error:   "Error invalid email or password",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	tokenData := fmt.Sprintf("%s:%s", existedUser.Id.String(), ctx.Value(types.CookieSecretKey))

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

	response := LoginResponse{
		Success: true,
		Message: "User login success",
		User: &LoginResponseUser{
			Name:  existedUser.Name,
			Email: existedUser.Email,
			Id:    existedUser.Id.String(),
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

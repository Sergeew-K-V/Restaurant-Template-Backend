package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"restaurant-backend/src/repositories"
)

func (as *AuthService) Profile(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Profile")
	path := r.URL.Path
	id := path[len("/api/profile/"):]

	userRepo := repositories.NewUserRepository(as.db)

	user, err := userRepo.GetUserById(id)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

}

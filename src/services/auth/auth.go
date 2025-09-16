package auth

import (
	"context"
	"database/sql"
	"net/http"
)

type AuthService struct {
	ctx context.Context
	db  *sql.DB
}

type AuthServiceInterface interface {
	RegisterUser(ctx context.Context, w http.ResponseWriter, r *http.Request)
	LoginUser(ctx context.Context, w http.ResponseWriter, r *http.Request)
	Me(ctx context.Context, w http.ResponseWriter, r *http.Request)
}

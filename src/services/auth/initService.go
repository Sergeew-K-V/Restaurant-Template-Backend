package auth

import (
	"context"
	"database/sql"
	"fmt"
)

func InitService(ctx context.Context, db *sql.DB) (AuthServiceInterface, error) {
	if ctx == nil {
		return nil, fmt.Errorf("error to init auth service")
	}

	if db == nil {
		return nil, fmt.Errorf("error to init auth service")
	}

	return &AuthService{ctx: ctx, db: db}, nil
}

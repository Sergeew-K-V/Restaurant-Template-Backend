package auth

import (
	"fmt"
	"strings"
)

func (as *AuthService) validateLoginRequest(req *LoginRequest) error {
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

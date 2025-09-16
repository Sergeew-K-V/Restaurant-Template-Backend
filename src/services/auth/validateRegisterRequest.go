package auth

import (
	"fmt"
	"strings"
)

func (as *AuthService) validateRegisterRequest(req *RegisterRequest) error {
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

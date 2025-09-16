package types

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

// Predefined context keys
const (
	CookieSecretKey ContextKey = "cookieSecretKey"
)

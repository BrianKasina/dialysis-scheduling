package utils

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// JWT structure for decoding
type JWTClaims struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	jwt.StandardClaims
}

// AuthorizationMiddleware validates JWT tokens and session
func AuthorizationMiddleware(next http.Handler, secretKey string, db *sql.DB, requiredRole string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			ErrorHandler(w, http.StatusUnauthorized, errors.New("missing authorization header"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			ErrorHandler(w, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		// Validate the session in the database
		if !validateSession(claims.ID, claims.Type, db) {
			ErrorHandler(w, http.StatusUnauthorized, errors.New("session validation failed"))
			return
		}

		// Check role if required
		if requiredRole != "" && !hasRole(claims.ID, requiredRole, db) {
			ErrorHandler(w, http.StatusForbidden, errors.New("insufficient role privileges"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// validateSession checks if session exists in the database (simplified)
func validateSession(userID, userType string, db *sql.DB) bool {
	// Custom logic here, example with a dummy SQL query
	query := `SELECT COUNT(*) FROM sessions WHERE user_id = ? AND expires_at > NOW()`
	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	if err != nil || count == 0 {
		return false
	}
	return true
}

// hasRole checks if a user has a specific role (simplified)
func hasRole(userID, roleName string, db *sql.DB) bool {
	// Custom logic here, example with a dummy SQL query
	query := `SELECT role_name FROM users WHERE user_id = ?`
	var userRole string
	err := db.QueryRow(query, userID).Scan(&userRole)
	if err != nil || userRole != roleName {
		return false
	}
	return true
}

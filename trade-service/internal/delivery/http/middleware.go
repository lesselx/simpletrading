package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const userContextKey = contextKey("userEmail")

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "mysecret"
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// Validate the algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Get email from claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid claims", http.StatusUnauthorized)
			return
		}

		email, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "Invalid subject claim", http.StatusUnauthorized)
			return
		}

		// Set email into context
		ctx := context.WithValue(r.Context(), userContextKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Optionally, export this for handlers
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(userContextKey).(string)
	return email, ok
}

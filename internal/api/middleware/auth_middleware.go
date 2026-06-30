package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"1000words-game/internal/api/responses"
	"1000words-game/internal/services"
	dbModels "1000words-game/models/db"
)

type contextKey string

const currentUserContextKey contextKey = "current_user"

func RequireAuth(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := BearerTokenFromRequest(r)
			if err != nil {
				responses.Unauthorized(w, "Authentication token is required")
				return
			}

			user, err := authService.Authenticate(token)
			if err != nil {
				responses.Unauthorized(w, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), currentUserContextKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CurrentUser(r *http.Request) (*dbModels.User, bool) {
	user, ok := r.Context().Value(currentUserContextKey).(*dbModels.User)

	return user, ok
}

func BearerTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.Fields(header)

	if len(parts) != 2 {
		return "", errors.New("authorization header is not valid")
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("authorization scheme must be bearer")
	}

	if strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("token is required")
	}

	return parts[1], nil
}

package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"

	"github.com/divijg19/physiolink/backend/internal/config"
)

type ctxKey string

const UserIDKey ctxKey = "user_id"
const UserRoleKey ctxKey = "user_role"

func JWTAuth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			var tokenStr string
			if auth != "" && strings.HasPrefix(auth, "Bearer ") {
				tokenStr = strings.TrimPrefix(auth, "Bearer ")
			} else {
				// Support React Native client using x-auth-token
				tokenStr = r.Header.Get("x-auth-token")
			}
			if tokenStr == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			// Our JWT payload: { user: { id, role }, exp }
			var sub string
			var role string
			if u, ok := claims["user"].(map[string]interface{}); ok && u != nil {
				if v, ok := u["id"].(string); ok {
					sub = v
				}
				if v, ok := u["role"].(string); ok {
					role = v
				}
			}
			if sub == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, sub)
			if role != "" {
				ctx = context.WithValue(ctx, UserRoleKey, role)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CookieAuth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			tokenStr := cookie.Value
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			var sub string
			var role string
			if u, ok := claims["user"].(map[string]interface{}); ok && u != nil {
				if v, ok := u["id"].(string); ok {
					sub = v
				}
				if v, ok := u["role"].(string); ok {
					role = v
				}
			}
			if sub == "" {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, sub)
			if role != "" {
				ctx = context.WithValue(ctx, UserRoleKey, role)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OptionalCookieAuth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			tokenStr := cookie.Value
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(cfg.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				next.ServeHTTP(w, r)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}
			var sub string
			var role string
			if u, ok := claims["user"].(map[string]interface{}); ok && u != nil {
				if v, ok := u["id"].(string); ok {
					sub = v
				}
				if v, ok := u["role"].(string); ok {
					role = v
				}
			}
			if sub != "" {
				ctx := context.WithValue(r.Context(), UserIDKey, sub)
				if role != "" {
					ctx = context.WithValue(ctx, UserRoleKey, role)
				}
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

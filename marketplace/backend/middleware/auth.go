package middleware

import (
	"context"
	"net/http"
	"strings"

	"marketplace/utils"
)

type contextKey string

const (
	ContextUserID contextKey = "user_id"
	ContextEmail  contextKey = "email"
	ContextRole   contextKey = "role"
	ContextShopID contextKey = "shop_id"
)

func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.Error(w, http.StatusUnauthorized, "Токени воридшавӣ ёфт нашуд (Missing Authorization header)")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Error(w, http.StatusUnauthorized, "Формати нодурусти токен (Invalid Authorization format)")
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			utils.Error(w, http.StatusUnauthorized, "Токени нодуруст ё муҳлаташ гузашта (Invalid or expired token)")
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextEmail, claims.Email)
		ctx = context.WithValue(ctx, ContextRole, claims.Role)
		if claims.ShopID != nil {
			ctx = context.WithValue(ctx, ContextShopID, *claims.ShopID)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(ContextRole).(string)
			if !ok {
				utils.Error(w, http.StatusForbidden, "Дастрасӣ маҳдуд аст (Access denied)")
				return
			}

			allowed := false
			for _, role := range allowedRoles {
				if userRole == role || userRole == "admin" {
					allowed = true
					break
				}
			}

			if !allowed {
				utils.Error(w, http.StatusForbidden, "Шумо ба ин қисмат ҳуқуқи дастрасӣ надоред (Insufficient permissions)")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireSeller(next http.Handler) http.Handler {
	return RequireRoles("seller", "admin")(next)
}

func RequireModerator(next http.Handler) http.Handler {
	return RequireRoles("moderator", "admin")(next)
}

func GetUserIDFromContext(r *http.Request) int64 {
	val, _ := r.Context().Value(ContextUserID).(int64)
	return val
}

func GetShopIDFromContext(r *http.Request) (int64, bool) {
	val, ok := r.Context().Value(ContextShopID).(int64)
	return val, ok
}

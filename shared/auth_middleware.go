package shared

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"thoughts_backend_api/models"
	"thoughts_backend_api/types"
)

func AuthMiddleware(db *gorm.DB, jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				RespondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "missing authorization header",
				})
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				RespondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization header format",
				})
				return
			}

			tokenString := parts[1]
			claims := &types.Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
				return jwtSecret, nil
			})
			if err != nil || !token.Valid {
				RespondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "invalid token",
				})
				return
			}

			var user models.User
			if err := db.First(&user, claims.UserID).Error; err != nil {
				RespondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "user not found",
				})
				return
			}

			ctx := context.WithValue(r.Context(), types.UserContextKey, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(types.UserContextKey).(*models.User)
	return user, ok
}

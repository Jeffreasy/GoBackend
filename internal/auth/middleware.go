// middleware.go (auth directory)

// Deze middleware functie valideert of een inkomend verzoek een geldige JWT bevat voordat toegang wordt verleend tot beschermde routes.

package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func JWTAuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Geen token gevonden", http.StatusUnauthorized)
				return
			}

			// Expect "Authorization: Bearer <token>"
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Ongeldig token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Ongeldige claims", http.StatusUnauthorized)
				return
			}

			// Voeg de user_id uit de claims toe aan de context zodat handlers deze kunnen gebruiken
			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

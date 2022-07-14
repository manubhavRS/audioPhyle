package middlewares

import (
	"audioPhile/models"
	"net/http"
)

func AccessControlMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var signedUser *models.UserModel
		signedUser = UserFromContext(r.Context())
		if signedUser.Role != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

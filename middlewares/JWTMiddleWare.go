package middlewares

import (
	"audioPhile/database/helper"
	"audioPhile/models"
	"audioPhile/utilities"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func GenerateJWT(user *models.UserModel) (string, error) {
	var mySigningKey = []byte(utilities.Secretkey)
	claims := &models.Claims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Printf("GenerateJWT Error: %v", err)
		return "", err
	}
	return tokenString, nil
}

func JWTAuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header[utilities.TokenString] == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var mySigningKey = []byte(utilities.Secretkey)
		var claims models.Claims
		token, err := jwt.ParseWithClaims(r.Header[utilities.TokenString][0], &claims, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("there was an error in parsing")
			}
			return mySigningKey, nil
		})
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID := fmt.Sprint(claims.ID)
		if claims.StandardClaims.ExpiresAt-time.Now().Unix() < 0 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		users, err := helper.FetchUserDetailsHelper(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("JWT Token verified.")
		ctx := context.WithValue(r.Context(), utilities.ContextUserKey, users)
		rec := statusRecorder{w, 200}
		handler.ServeHTTP(&rec, r.WithContext(ctx))
		if rec.status == 200 {
			if claims.StandardClaims.ExpiresAt-time.Now().Unix() < 30 && claims.StandardClaims.ExpiresAt-time.Now().Unix() > 0 {
				refreshToken, err := GenerateJWT(&users)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				log.Printf("Refresh Token: " + refreshToken)
				w.Header().Set("Content-Type", "application/json")
				resp := make(map[string]string)
				resp["RefreshToken"] = refreshToken
				jsonResponse, err := json.Marshal(resp)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write(jsonResponse)
			}

		}

	})
}

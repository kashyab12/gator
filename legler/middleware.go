package legler

import (
	"github.com/kashyab12/gator/internal/database"
	"net/http"
	"strings"
)

type AuthedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		writer.Header().Set("Access-Control-Allow-Headers", "*")
		if request.Method == "OPTIONS" {
			writer.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func (config *ApiConfig) AuthMiddleware(authHandler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authHeader := r.Header.Get("Authorization"); authHeader == "" {
			_ = RespondWithError(w, http.StatusUnauthorized, "Please provide ApiKey")
		} else if splitAuthHeader := strings.Split(authHeader, "ApiKey "); len(splitAuthHeader) < 2 {
			_ = RespondWithError(w, http.StatusUnauthorized, "Please provide ApiKey in format Authorization: ApiKey <key>")
		} else {
			apiKey := splitAuthHeader[1]
			if queryUser, queryErr := config.DB.GetUser(r.Context(), apiKey); queryErr != nil {
				_ = RespondWithError(w, http.StatusInternalServerError, queryErr.Error())
			} else {
				authHandler(w, r, queryUser)
			}
		}
	}
}

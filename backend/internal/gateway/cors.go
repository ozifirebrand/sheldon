package gateway

import (
	"net/http"
)

func CORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		responseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE")
		responseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if request.Method == http.MethodOptions {
			responseWriter.WriteHeader(http.StatusNoContent)
			return
		}
		handler.ServeHTTP(responseWriter, request)
	})
}

package gateway

import (
	"encoding/json"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func writeJSON(responseWriter http.ResponseWriter, statusCode int, value interface{}) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(statusCode)
	if value != nil {
		_ = json.NewEncoder(responseWriter).Encode(value)
	}
}

func writeError(responseWriter http.ResponseWriter, request *http.Request, err error) {
	statusErr, ok := status.FromError(err)
	if !ok {
		http.Error(responseWriter, "internal server error", http.StatusInternalServerError)
		return
	}
	message := statusErr.Message()
	switch statusErr.Code() {
	case codes.AlreadyExists:
		http.Error(responseWriter, message, http.StatusConflict)
	case codes.InvalidArgument:
		http.Error(responseWriter, message, http.StatusBadRequest)
	default:
		http.Error(responseWriter, "internal server error", http.StatusInternalServerError)
	}
}

func appointmentIDFromPath(path, prefix string) string {
	return strings.TrimPrefix(path, prefix)
}

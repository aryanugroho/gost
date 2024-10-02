package middleware

import (
	"context"
	"net/http"

	"[base_project]/internal/model"
	"[base_project]/pkg/logger"
	"[base_project]/pkg/uuid"
)

func CorrelationLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := uuid.New().String()
		ctx := context.WithValue(r.Context(), model.ContextCorrelationID, correlationID)
		r = r.WithContext(ctx)
		logger.Http(r)
		next.ServeHTTP(w, r)
	})
}

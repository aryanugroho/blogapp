package middleware

import (
	"context"
	"net/http"

	"github.com/aryanugroho/blogapp/internal/logger"
	"github.com/aryanugroho/blogapp/internal/uuid"
	"github.com/aryanugroho/blogapp/model"
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

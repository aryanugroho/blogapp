package httpapi

import (
	"net/http"

	"github.com/aryanugroho/blogapp/internal/httpapi/server"
)

func Ping() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		server.WriteJSON(wr, r, http.StatusOK, map[string]string{"ping": "pong"})
	})
}

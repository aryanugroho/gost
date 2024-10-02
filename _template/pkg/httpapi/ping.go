package httpapi

import (
	"net/http"

	"[base_project]/pkg/httpapi/server"
)

func Ping() http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, r *http.Request) {
		server.WriteJSON(wr, r, http.StatusOK, map[string]string{"ping": "pong"})
	})
}

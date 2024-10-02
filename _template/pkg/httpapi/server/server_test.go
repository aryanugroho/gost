package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"[base_project]/pkg/httpapi"
	"[base_project]/pkg/httpapi/server"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	srv := server.New()
	srv.Register(httpapi.Routes(httpapi.Dependencies{}))

	wr := httptest.NewRecorder()
	srv.ServeHTTP(wr, httptest.NewRequest(http.MethodGet, "/ping", nil))

	res := wr.Result()
	assert.Equal(t, res.StatusCode, http.StatusOK)
}

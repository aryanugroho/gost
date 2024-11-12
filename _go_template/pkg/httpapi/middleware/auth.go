package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"github.com/aryanugroho//internal/model"
	"github.com/aryanugroho//pkg/helpers/errors"
	"github.com/aryanugroho//pkg/httpapi/server"
	"github.com/aryanugroho//pkg/logger"
)

func AuthAPI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		if len(apiKey) == 0 {
			server.WriteResponseError(w, r, errors.Unauthorized.New("missing api key"))
			return
		}

		haser := sha256.New()
		haser.Write([]byte(apiKey))
		hashedAPIKey := base64.StdEncoding.EncodeToString(haser.Sum(nil))

		clientID, err := getClient(hashedAPIKey)
		if err != nil {
			server.WriteResponseError(w, r, errors.Unauthorized.New("invalid api key"))
			return
		}

		ctx := context.WithValue(r.Context(), model.ClientID, clientID)
		r = r.WithContext(ctx)
		logger.Http(r)
		next.ServeHTTP(w, r)
	})
}

var clients = make(map[string]string)

func RegisterClientID(clientID string, clientName string) {
	clients[clientID] = clientName
}

func getClient(clientID string) (string, error) {

	v, ok := clients[clientID]
	if !ok {
		return "", errors.NotFound.New("client not found")
	}
	return v, nil
}

package openapi

import (
	"net/http"

	"github.com/swaggest/swgui/v5cdn"
)

func (o *OpenAPIV3) ServeDocs(openAPIPath string) http.Handler {
	handler := v5cdn.New(
		o.reflector.Spec.Info.Title,
		openAPIPath,
		"",
	)

	return handler
}

package openapi

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/swaggest/openapi-go/openapi3"
)

type OpenAPIV3 struct {
	reflector openapi3.Reflector
}

func NewOpenApiV3(title, version, description string) OpenAPIV3 {
	o := OpenAPIV3{}
	o.reflector = openapi3.Reflector{}
	o.reflector.Spec = &openapi3.Spec{Openapi: "3.0.3"}
	o.reflector.Spec.Info.
		WithTitle(title).
		WithVersion(version).
		WithDescription(description)

	return o
}

func (o *OpenAPIV3) AddOperation(req, resp interface{}, path string, method string, description string) error {
	op, err := openapi3.NewReflector().NewOperationContext(method, path)
	if err != nil {
		return errors.New("failed to create openapi3 operation context")
	}
	op.AddReqStructure(req)
	op.AddRespStructure(resp)
	op.SetDescription(description)

	err = o.reflector.AddOperation(op)
	if err != nil {
		return fmt.Errorf("failed to add operation to openapi3 spec: %w", err)
	}
	return nil
}

func (o *OpenAPIV3) Build(w http.ResponseWriter, r *http.Request) {
	schema, err := o.reflector.Spec.MarshalJSON()
	if err != nil {
		http.Error(w, "failed to marshal openapi3 spec", http.StatusInternalServerError)
		return
	}
	w.Write(schema)
}

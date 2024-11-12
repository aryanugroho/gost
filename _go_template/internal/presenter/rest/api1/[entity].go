package api1

import (
	"net/http"

	"github.com/aryanugroho//internal/usecase"
	"github.com/aryanugroho//pkg/helpers/errors"
	"github.com/aryanugroho//pkg/helpers/validation"
	"github.com/aryanugroho//pkg/httpapi/server"
)

func (api *API) InitAPI() {
	api.BaseRoutes. = api.BaseRoutes.ApiRoot.PathPrefix("/").Subrouter()

	api.BaseRoutes..Handle("/", http.HandlerFunc(api.create)).Methods("POST")
	api.docs.AddOperation(usecase.Payload{}, usecase.Payload{}, "/api/v1/create", http.MethodPut, "Create something")
}

func (api *API) create(w http.ResponseWriter, r *http.Request) {
	dto := usecase.Payload{}
	err := server.ReadJSON(r, &dto)
	if err != nil {
		server.WriteResponseError(w, r, errors.BadRequest.New(err.Error()))
		return
	}

	err = validation.Validate.Struct(dto)
	if err != nil {
		server.WriteResponseError(w, r, errors.BadRequest.New(err.Error()))
		return
	}

	result, err := api.usecase.Create(r.Context(), &dto)
	if err != nil {
		server.WriteResponseError(w, r, err)
		return
	}

	server.WriteJSON(w, r, http.StatusOK, result)
}

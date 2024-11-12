package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aryanugroho//internal/model"
	"github.com/aryanugroho//pkg/helpers/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// WriteJSON writes the given value as JSON encoded data to the response writer with given
// status. 'v' must be compatible with JSON Marshal.
func WriteResponseError(wr http.ResponseWriter, r *http.Request, err error) {
	var defaultMessage string
	var statusCode, errorCode int

	errType := errors.GetType(err)

	switch errType {
	case 42201:
		defaultMessage = "Default message to be showed."
		statusCode = http.StatusUnprocessableEntity
		errorCode = 42201
	case errors.Unauthorized:
		defaultMessage = "Unauthorized."
	case errors.NotFound:
		defaultMessage = "Resource not found."
	case errors.UnprocessableEntity:
		defaultMessage = "Unprocessable entity error."
	case errors.BadRequest:
		defaultMessage = "Bad request."
	default:
		defaultMessage = "Internal server error."
	}

	message := err.Error()
	if message == "" {
		message = defaultMessage
	}
	if errorCode == 0 {
		// Get error code from err type value
		errorCode = errors.GetErrorCode(err)
	}
	if statusCode == 0 {
		statusCode = int(errType)
	}

	response := struct {
		Errors []ErrorResponse `json:"errors"`
	}{
		Errors: []ErrorResponse{
			{
				Code:    errorCode,
				Message: message,
			},
		},
	}

	var errsResponse []ErrorResponse
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		statusCode = 400
		for _, e := range validationErrs {
			errsResponse = append(errsResponse, ErrorResponse{
				Code:    int(errors.InvalidInputParameterCode),
				Message: e.Error(),
			})
		}
		response.Errors = errsResponse
	}

	wr.Header().Set("Content-type", "application/json; charset=utf-8")
	wr.WriteHeader(statusCode)
	req := r.WithContext(context.WithValue(r.Context(), model.StatusCode, statusCode))
	*r = *req
	resp, _ := json.Marshal(response)
	wr.Write(resp)
}

// WriteJSON writes the given value as JSON encoded data to the response writer with given
// status. 'v' must be compatible with JSON Marshal.
func WriteJSON(wr http.ResponseWriter, r *http.Request, status int, v interface{}) {
	req := r.WithContext(context.WithValue(r.Context(), model.StatusCode, status))
	*r = *req

	wr.Header().Set("Content-type", "application/json; charset=utf-8")
	wr.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(wr).Encode(v); err != nil {
		panic(fmt.Errorf("writeJSON failed: %s", err.Error()))
	}
}

// ReadJSON reads the body of the given request as JSON encoded data and unmarshalls it into
// the 'into' pointer. If the `into` has a `Validate() error` it will be used to validate the
// data after unmarshal.
func ReadJSON(req *http.Request, into interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(into); err != nil {
		return fmt.Errorf("failed to decode body: %s", err.Error())
	}
	if v, ok := into.(interface{ Validate() error }); ok {
		return v.Validate()
	}
	return nil
}

func staticHandler(status int, message string) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		WriteJSON(wr, req, status, map[string]string{"message": message})
	})
}

// New initialises a new API server with a router and default ping handler.
func New() Server {
	router := mux.NewRouter()
	router.MethodNotAllowedHandler = staticHandler(http.StatusMethodNotAllowed, "method not allowed")
	router.NotFoundHandler = staticHandler(http.StatusNotFound, "path not found")
	return Server{router: router}
}

// Server wraps an HTTP router and acts as an HTTP server
type Server struct {
	router *mux.Router
}

// Register registers the given routes to the server.
func (s *Server) Register(routes ...[]Route) *Server {
	for _, list := range routes {
		for _, route := range list {
			s.router.Handle(route.Path, route.Handler).Methods(route.Method)
		}
	}
	return s
}

// Serve starts the server on given addr and blocks until the underlying server
// exits with error or the context is cancelled. If the context is cancelled, a
// graceful shutdown is performed with configured grace period.
func (s *Server) Serve(baseCtx context.Context, addr string) error {
	ctx, cancel := context.WithCancel(baseCtx)
	defer cancel()

	srv := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	go func() {
		defer cancel()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("server exited with error: %s\n", err.Error())
		}
	}()

	return s.waitForInterrupt(ctx, srv)
}

func (s *Server) waitForInterrupt(ctx context.Context, srv *http.Server) error {
	<-ctx.Done()

	const gracefulPeriod = 1 * time.Second
	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulPeriod)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

func (s *Server) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(wr, req)
}

// Route represents an HTTP endpoint and the handler to be used.
type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}

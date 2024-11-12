package middleware

import "net/http"

// Middleware wraps a given http Handler to intercept request.
type Middleware func(http.Handler) http.Handler

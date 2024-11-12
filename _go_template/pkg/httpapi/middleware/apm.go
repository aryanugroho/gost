package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aryanugroho//internal/model"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func APM(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opts := []ddtrace.StartSpanOption{
			tracer.ServiceName(os.Getenv("DD_SERVICE")),
			tracer.ResourceName(r.Method + " " + r.Host + r.URL.Path),
			tracer.SpanType(ext.SpanTypeWeb),
			tracer.Tag(ext.HTTPMethod, r.Method),
			tracer.Tag(ext.HTTPURL, r.URL.Path),
			tracer.Tag(ext.EventSampleRate, 1.0),
			tracer.Measured(),
		}
		if spanctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header)); err == nil {
			opts = append(opts, tracer.ChildOf(spanctx))
		}

		clientID := r.Context().Value(ClientID)
		if clientID != nil {
			opts = append(opts, tracer.Tag("client.id", clientID))
		}

		span, ctx := tracer.StartSpanFromContext(r.Context(), "http.request", opts...)
		defer span.Finish()

		ctx = context.WithValue(ctx, TraceID, span.Context().TraceID())
		ctx = context.WithValue(ctx, SpanID, span.Context().SpanID())

		// pass the span through the request context
		r = r.WithContext(ctx)

		// serve the request to the next middleware
		next.ServeHTTP(w, r)

		status := r.Context().Value(model.StatusCode)
		span.SetTag(ext.HTTPCode, fmt.Sprintf("%v", status))

		if clientID != nil {
			span.SetTag("client.id", clientID)
		}

		statusCode := status.(int)
		if statusCode >= 500 && statusCode < 600 {
			span.SetTag(ext.Error, fmt.Errorf("%d: %s", status, http.StatusText(statusCode)))
		}
	})
}

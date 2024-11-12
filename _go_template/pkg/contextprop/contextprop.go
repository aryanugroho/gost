package contextprop

import "context"

type ContextKey string

const (
	ContextCorrelationID ContextKey = "cid"
	TraceID              ContextKey = "trace_id"
	SpanID               ContextKey = "span_id"
	ClientID             ContextKey = "client_id"
	StatusCode           ContextKey = "status_code"
)

func GetContextValue(ctx context.Context, key ContextKey) string {
	value, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return value
}

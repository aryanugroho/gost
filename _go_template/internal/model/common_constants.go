package model

type ContextKey string

const (
	ContextCorrelationID ContextKey = "cid"
	TraceID              ContextKey = "trace_id"
	SpanID               ContextKey = "span_id"
	ClientID             ContextKey = "client_id"
	StatusCode           ContextKey = "status_code"
)

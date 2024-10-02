package contextprop

import "context"

type ContextKey string

func GetContextValue(ctx context.Context, key ContextKey) string {
	value, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return value
}

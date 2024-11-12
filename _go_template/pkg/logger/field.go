package logger

import "go.uber.org/zap"

type Field struct {
	zap.Field
}

func (f *Field) format() zap.Field {
	return zap.Field{
		Key: f.Key,
	}
}

func convertFields(fields []Field) []zap.Field {
	zFields := []zap.Field{}
	for _, v := range fields {
		zF := zap.Field{
			Key:       v.Key,
			Type:      v.Type,
			Integer:   v.Integer,
			String:    v.String,
			Interface: v.Interface,
		}
		zFields = append(zFields, zF)
	}
	return zFields
}

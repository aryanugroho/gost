package main

// mapCustomTypeToGoType maps custom types to Go types
func mapCustomTypeToGoType(customType string) string {
    switch customType {
    case "string":
        return "string"
    case "int":
        return "int64"
    case "decimal":
        return "float64"
    case "time": // assuming 'time' is the custom type for time.Time
        return "time.Time"
    default:
        return customType
    }
}

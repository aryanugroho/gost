package model

import "github.com/aryanugroho//pkg/uuid"

// NewId is a globally unique identifier.  It is a [A-Z0-9] string 26
// characters long.  It is a UUID version 4 Guid that is zbased32 encoded
// without the padding.
func NewId() string {
	return uuid.New().String()
}

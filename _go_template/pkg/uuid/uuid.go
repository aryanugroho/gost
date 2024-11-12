package uuid

import (
	"github.com/google/uuid"
)

// UUID custom common uuid
type UUID struct {
	uuid uuid.UUID
}

// New generates new uuid
func New() UUID {
	return UUID{
		uuid: uuid.New(),
	}
}

// String return new uuid of string
func (u UUID) String() string {
	return u.uuid.String()
}

// check uuid
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	if err != nil {
		return false
	}
	return true
}

package uuidv7

import "github.com/google/uuid"

func NewE() uuid.UUID {
	// Generate a new UUIDv7
	val, error := uuid.NewV7()

	if error != nil {
		panic(error)
	}

	return val
}

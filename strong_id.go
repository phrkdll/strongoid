package strongoid

import "github.com/google/uuid"

type StrongIdConstraints interface {
	int64 | string | uuid.UUID
}

// A base struct for strongly typed IDs
// To trigger the marshaller below, instead new type, define an alias like this:
// "type MyStronglyTypedId = StrongId[int]"
type StrongId[T StrongIdConstraints] struct {
	Inner T
}

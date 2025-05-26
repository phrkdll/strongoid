package strongoid

import "github.com/google/uuid"

type IdConstraints interface {
	int64 | string | uuid.UUID
}

// A base struct for strongly typed IDs
// "type MyStronglyTypedId StrongId[int]"
type Id[T IdConstraints] struct {
	Inner T
}

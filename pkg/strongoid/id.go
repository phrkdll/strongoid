package strongoid

// Allowed inner types for strongly typed IDs
type IdConstraints interface {
	int64 | string | any
}

// Base struct for strongly typed IDs
type Id[T IdConstraints] struct {
	Inner T
}

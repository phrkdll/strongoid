package strongoid

import (
	"database/sql/driver"
	"errors"
)

var (
	ErrUnsupportedData = errors.New("unsupported data")
)

func (id *Id[T]) Scan(dbValue any) (err error) {
	switch value := dbValue.(type) {
	case T:
		id.Inner = value
	default:
		return ErrUnsupportedData
	}
	return nil
}

func (id Id[T]) Value() (driver.Value, error) {
	return id.Inner, nil
}

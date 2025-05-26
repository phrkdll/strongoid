package strongoid

import (
	"database/sql/driver"
	"fmt"
)

func (id *Id[T]) Scan(dbValue any) (err error) {
	switch value := dbValue.(type) {
	case T:
		id.Inner = value
	default:
		return fmt.Errorf("unsupported data %#v", dbValue)
	}
	return nil
}

func (id Id[T]) Value() (driver.Value, error) {
	return id.Inner, nil
}

package strongoid

import "encoding/json"

func (id Id[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.Inner)
}

func (id *Id[T]) UnmarshalJSON(data []byte) error {
	var value T

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	id.Inner = value

	return nil
}

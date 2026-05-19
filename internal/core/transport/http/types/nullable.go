package core_http_types

import (
	"encoding/json"

	"github.com/odelshchwank/BigProjectLesson/internal/core/domain"
)

// Создаем ещё один уровень Nullable на уровне транспорта, чтобы именно к нему
// "прикручивать" метод по декодированию JSON
type Nullable[T any] struct {
	domain.Nullable[T]
}

// Если был вызван этот метод, то значит значение было передано, т.е. Set == true
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	if string(b) == "null" {
		n.Value = nil

		return nil
	}

	var value T

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value

	return nil
}

func (n *Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}

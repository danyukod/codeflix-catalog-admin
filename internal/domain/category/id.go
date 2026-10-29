package category

import (
	"errors"

	"github.com/google/uuid"
)

type Id struct {
	id string
}

func (i Id) GetValue() string {
	return i.id
}

func Unique() Id {
	return Id{uuid.New().String()}
}

func FromString(id string) (Id, error) {
	categoryID := Id{id}
	if err := categoryID.Validate(); err != nil {
		return Id{}, err
	}
	return categoryID, nil
}

func (i Id) Validate() error {
	if i.id == "" {
		return errors.New("id cannot be empty")
	}
	if _, err := uuid.Parse(i.id); err != nil {
		return errors.New("id must be a valid UUID")
	}
	return nil
}

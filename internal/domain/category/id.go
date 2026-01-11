package category

import (
	"errors"

	"github.com/google/uuid"
)

type Id struct {
	id string
}

func (i Id) GetIdentifier() string {
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

func FromUUID(uuid uuid.UUID) Id {
	return Id{uuid.String()}
}

func (i Id) Validate() error {
	if i.id == "" {
		return errors.New("id cannot be empty")
	}
	return nil
}

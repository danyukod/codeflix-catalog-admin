package update

type CategoryCommand struct {
	id          string
	name        string
	description string
	isActive    bool
}

func (c CategoryCommand) With(id string, name, description string, isActive bool) *CategoryCommand {
	return &CategoryCommand{id, name, description, isActive}
}

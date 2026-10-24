package create

type CategoryCommand struct {
	name        string
	description string
	isActive    bool
}

func (c CategoryCommand) With(name, description string, isActive bool) *CategoryCommand {
	return &CategoryCommand{name, description, isActive}
}

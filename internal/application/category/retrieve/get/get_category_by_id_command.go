package get

type CategoryCommand struct {
	id string
}

func (c CategoryCommand) With(id string) *CategoryCommand {
	return &CategoryCommand{id}
}

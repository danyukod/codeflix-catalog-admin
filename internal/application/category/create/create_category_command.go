package create

type CreateCategoryCommand struct {
	name        string
	description string
	isActive    bool
}

func (c CreateCategoryCommand) With(name, description string, isActive bool) *CreateCategoryCommand {
	return &CreateCategoryCommand{name, description, isActive}
}

package update

type UpdateCategoryCommand struct {
	id          string
	name        string
	description string
	isActive    bool
}

func (c UpdateCategoryCommand) With(id string, name, description string, isActive bool) *UpdateCategoryCommand {
	return &UpdateCategoryCommand{id, name, description, isActive}
}

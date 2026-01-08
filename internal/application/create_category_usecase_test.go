package application

import "testing"

func TestCreateCategoryUseCase_Execute(t *testing.T) {
	categoryCreateUsecase := new(CreateCategoryUseCase)
	if categoryCreateUsecase == nil {
		t.Error("categoryCreateUsecase is nil")
	}

	category, err := categoryCreateUsecase.Execute()
	if err != nil {
		t.Error(err)
	}
	if category == nil {
		t.Error("category is nil")
	}
}

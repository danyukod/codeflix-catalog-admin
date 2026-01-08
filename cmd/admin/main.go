package main

import (
	"fmt"

	"github.com/danyukod/codeflix-catalog-admin/internal/application"
)

func main() {
	fmt.Println("Hello, World!")
	createCategoryUsecase := application.CreateCategoryUseCase{}
	fmt.Println(createCategoryUsecase.Execute())
}

package services

import "gorm.io/gorm"

type ExampleService struct {
	DB *gorm.DB
}

func (es *ExampleService) Index() {
	//
}

func (es *ExampleService) Store() {
	//
}

func (es *ExampleService) Show() {
	//
}

func (es *ExampleService) Update() {
	//
}

func (es *ExampleService) Delete() {
	//
}

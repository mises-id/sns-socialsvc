package models

import "github.com/go-playground/validator"

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

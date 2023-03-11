package configuration

import (
	"log"

	"github.com/go-playground/validator/v10"
)

type GoAppTools struct {
	ErrorLogger log.Logger
	InfoLogger  log.Logger
	Validate    *validator.Validate
}

func Check(erro error, app GoAppTools) {
	if erro != nil {
		app.ErrorLogger.Fatal(erro)
	}
}

package configuration

import (
	"log"
)

type GoAppTools struct {
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
}

func Check(erro error, app *GoAppTools) {
	if erro != nil {
		app.ErrorLogger.Fatal(erro)
	}
}

// Package utils fornece utilitários que auxiliam na manipulação e diagnóstico de problemas.
package utils

import (
	"github.com/sirupsen/logrus"
)

// GoAppTools encapsula ferramentas de logging para diagnóstico e monitoramento.
type GoAppTools struct {
	// Logger é a instância do logger Logrus.
	Logger *logrus.Logger
}

// NewGoAppTools cria uma nova instância de GoAppTools com configurações padrão para o logger.
func NewGoAppTools() *GoAppTools {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	return &GoAppTools{
		Logger: logger,
	}
}

// Check avalia se um erro ocorreu e, se verdadeiro, registra o erro usando Logrus.
// Diferente de outras implementações, essa função não encerra o programa, mas apenas registra o erro.
func (app *GoAppTools) Check(err error) {
	if err != nil {
		app.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("An error occurred")
	}
}

// CheckAndPanic avalia se um erro ocorreu e, se verdadeiro, registra o erro e causa um panic.
func (app *GoAppTools) CheckAndPanic(err error) {
	if err != nil {
		app.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Panic("An error occurred")
	}
}

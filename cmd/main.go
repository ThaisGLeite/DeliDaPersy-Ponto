// Package main é o ponto de entrada da aplicação deli-ponto.
package main

import (
	"deli-ponto/pkg/driver"
	"deli-ponto/pkg/handlers"
	"deli-ponto/pkg/utils"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

var (
	dynamoClient *dynamodb.Client
	logs         *utils.GoAppTools
)

// inLambda verifica se o código está sendo executado em um ambiente Lambda.
func inLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}

// setupRouter configura as rotas da aplicação.
func setupRouter() *gin.Engine {
	r := gin.New()
	r.GET("/", func(ctx *gin.Context) {
		handlers.ResponseOK(ctx, logs)
	})
	r.GET("/pontos", func(ctx *gin.Context) {
		handlers.GetPunches(ctx, dynamoClient, logs)
	})
	r.POST("/ponto/:nome", func(ctx *gin.Context) {
		nome := ctx.Param("nome")
		handlers.PostPunch(nome, ctx, dynamoClient, logs)
	})
	return r
}

func main() {
	// Inicializa o logger
	logs = utils.NewGoAppTools()

	// Configura a conexão com o DynamoDB
	var err error
	dynamoClient, err = driver.ConfigAws(logs)
	logs.Check(err) // Verifica e registra erros usando o método Check de utils

	// Configura e inicia o servidor
	router := setupRouter()
	if inLambda() {
		err = gateway.ListenAndServe(":8080", router)
	} else {
		err = http.ListenAndServe(":8080", router)
	}
	logs.CheckAndPanic(err) // Verifica e registra erros críticos, e encerra o programa se necessário
}

// Para compilar o binario do sistema usamos:
//
//	GOARCH=arm64 GOOS=linux  CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap .
//
// para criar o zip do projeto comando:
//
// zip lambda.zip bootstrap

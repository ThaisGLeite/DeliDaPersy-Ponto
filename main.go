package main

import (
	"deli-ponto/configuration"
	"deli-ponto/driver"
	"deli-ponto/handlers"
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

var (
	dynamoClient *dynamodb.Client
	logs         configuration.GoAppTools
)

func inLambda() bool {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

func setupRouter() *gin.Engine {
	appRouter := gin.New()
	appRouter.GET("/", func(ctx *gin.Context) {
		logs.InfoLogger.Println("Servidor Ok")
		handlers.ResponseOK(ctx, logs)
	})

	appRouter.GET("/pontos", func(ctx *gin.Context) {
		handlers.GetPunches(ctx, dynamoClient, logs)
	})

	appRouter.POST("/ponto/:nome", func(ctx *gin.Context) {
		nome := ctx.Param("nome")
		handlers.PostPunch(nome, ctx, dynamoClient, logs)
	})
	return appRouter
}

// Para compilar o binario do sistema usamos: GOOS=linux GOARCH=amd64 go build -o ponto-app .
// para criar o zip do projeto comando: zip lambda.zip ponto-app
func main() {
	InfoLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)
	ErrorLogger := log.New(os.Stdout, " ", log.LstdFlags|log.Lshortfile)

	logs.InfoLogger = *InfoLogger
	logs.ErrorLogger = *ErrorLogger
	var err error
	// chamada de função para a criação da sessao de login com o banco
	dynamoClient, err = driver.ConfigAws()
	//chamada da função para revificar o erro retornado
	configuration.Check(err, logs)

	if inLambda() {

		log.Fatal(gateway.ListenAndServe(":8080", setupRouter()))
	} else {

		log.Fatal(http.ListenAndServe(":8080", setupRouter()))
	}
}

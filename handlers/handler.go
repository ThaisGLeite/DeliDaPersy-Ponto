package handlers

import (
	"deli-ponto/configuration"
	"deli-ponto/database/query"
	"deli-ponto/model"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func GetPunches(c *gin.Context, dynamoClient *dynamodb.Client, logs *configuration.GoAppTools) {
	punches := make([]model.Punch, 0)

	ponto := query.SelectPunch("Bianca", *dynamoClient, logs)
	punches = append(punches, ponto)
	ponto = query.SelectPunch("Danilo", *dynamoClient, logs)
	punches = append(punches, ponto)
	ponto = query.SelectPunch("paty", *dynamoClient, logs)
	punches = append(punches, ponto)
	c.IndentedJSON(http.StatusOK, punches)

}

func ResponseOK(c *gin.Context, app *configuration.GoAppTools) {
	c.IndentedJSON(http.StatusOK, "Servidor up")
}

func PostPunch(nome string, c *gin.Context, dynamoClient *dynamodb.Client, logs *configuration.GoAppTools) {
	query.InsertPunch(dynamoClient, nome, logs)
	response := ("ponto do colaborador " + nome + " batido")
	c.IndentedJSON(http.StatusCreated, response)
}

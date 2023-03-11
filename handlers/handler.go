package handlers

import (
	"deli-ponto/configuration"
	"deli-ponto/database/query"
	"deli-ponto/model"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func GetPunches(c *gin.Context, dynamoClient *dynamodb.Client, logs configuration.GoAppTools) {
	punches := make([]model.Punch, 0)

	ponto := query.SelectPunch("Bianca", *dynamoClient, logs)
	punches = append(punches, ponto)
	ponto = query.SelectPunch("Danilo", *dynamoClient, logs)
	punches = append(punches, ponto)
	ponto = query.SelectPunch("Patricia", *dynamoClient, logs)
	punches = append(punches, ponto)
	c.IndentedJSON(http.StatusOK, punches)

}

func ResponseOK(c *gin.Context, app configuration.GoAppTools) {
	c.IndentedJSON(http.StatusOK, "Servidor up")
}

func PostPunch(c *gin.Context, dynamoClient *dynamodb.Client, logs configuration.GoAppTools) {
	var newPunch model.Punch

	//configue o model punh with the retorn of context gin
	err := c.BindJSON(&newPunch)
	//faz a chacagem de errode forma unificada
	configuration.Check(err, logs)
	//calling the quiry package to mount the request for a DB
	query.InsertPunch(dynamoClient, newPunch.Nome, logs)
	name := ("ponto do colaborador " + newPunch.Nome + " batido")
	c.IndentedJSON(http.StatusCreated, (name))
}

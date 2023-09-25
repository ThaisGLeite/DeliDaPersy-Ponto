// Package handlers provides the HTTP handlers for the deli-ponto service.
package handlers

import (
	"deli-ponto/pkg/model"
	"deli-ponto/pkg/query"
	"deli-ponto/pkg/utils"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
)

// GetPunches is an HTTP handler that retrieves punch records for a set of names.
func GetPunches(c *gin.Context, dynamoClient *dynamodb.Client, logs *utils.GoAppTools) {
	// Initialize an empty slice to store punch records
	punches := make([]model.Punch, 0)

	// Fetch punch records for specific names and append them to the slice
	for _, name := range []string{"Bianca", "Danilo", "paty"} {
		punch := query.SelectPunch(name, dynamoClient, logs)
		punches = append(punches, punch)
	}

	// Respond with the collected punch records
	c.IndentedJSON(http.StatusOK, punches)
}

// ResponseOK is an HTTP handler that responds with a 200 OK status.
func ResponseOK(c *gin.Context, app *utils.GoAppTools) {
	c.IndentedJSON(http.StatusOK, "Servidor up")
}

// PostPunch is an HTTP handler that inserts a new punch record for a given name.
func PostPunch(nome string, c *gin.Context, dynamoClient *dynamodb.Client, logs *utils.GoAppTools) {
	// Insert a new punch record
	query.InsertPunch(dynamoClient, nome, logs)

	// Create a response message
	response := ("Ponto do colaborador " + nome + " batido")

	// Respond with a 201 Created status
	c.IndentedJSON(http.StatusCreated, response)
}

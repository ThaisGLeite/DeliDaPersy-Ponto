// Package handlers provides the HTTP handlers for the deli-ponto service.
package handlers

import (
	"deli-ponto/pkg/model"
	"deli-ponto/pkg/query"
	"deli-ponto/pkg/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// GetPunches is an HTTP handler that retrieves punch records for a set of names.
func GetPunches(w http.ResponseWriter, r *http.Request, dynamoClient *dynamodb.Client, logs *utils.GoAppTools) {
	// Initialize an empty slice to store punch records
	punches := make([]model.Punch, 0)

	// Fetch punch records for specific names and append them to the slice
	for _, name := range []string{"Bianca", "Danilo", "paty"} {
		punch := query.SelectPunch(name, dynamoClient, logs)
		punches = append(punches, punch)
	}

	// Encode the punches slice to JSON and write it to the ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(punches)
}

// ResponseOK is an HTTP handler that responds with a 200 OK status.
func ResponseOK(w http.ResponseWriter, r *http.Request, app *utils.GoAppTools) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Servidor up")
}

// PostPunch is an HTTP handler that inserts a new punch record for a given name.
func PostPunch(nome string, w http.ResponseWriter, r *http.Request, dynamoClient *dynamodb.Client, logs *utils.GoAppTools) {
	// Insert a new punch record
	query.InsertPunch(dynamoClient, nome, logs)

	// Create a response message
	response := ("Ponto do colaborador " + nome + " batido")

	// Set content type and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Encode the response to JSON and write it to the ResponseWriter
	json.NewEncoder(w).Encode(response)
}

// GetReport is an HTTP handler that retrieves a report for a given name and month.
func GetReport(w http.ResponseWriter, r *http.Request, dynamoClient *dynamodb.Client, logs *utils.GoAppTools, nome string, mes string) {
	// Modify the names as per the requirements
	if nome == "Bianca" {
		nome = "Bia"
	} else if nome == "Patricia" {
		nome = "paty"
	}

	ano := time.Now().Year()
	periodo := strconv.Itoa(ano) + "-" + mes

	// Fetch the report
	report := query.SelectReport(nome, periodo, *dynamoClient, logs)

	// Set the content type and send the report as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

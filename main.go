// Package main é o ponto de entrada da aplicação deli-ponto.
package main

import (
	"deli-ponto/pkg/driver"
	"deli-ponto/pkg/handlers"
	"deli-ponto/pkg/utils"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/apex/gateway"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	dynamoClient *dynamodb.Client
	logs         *utils.GoAppTools
)

// inLambda verifica se o código está sendo executado em um ambiente Lambda.
func inLambda() bool {
	return os.Getenv("LAMBDA_TASK_ROOT") != ""
}

// SecureMiddleware adds security-related headers to responses.
func SecureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Inicializa o logger
	logs = utils.NewGoAppTools()

	// Configura a conexão com o DynamoDB
	var err error
	dynamoClient, err = driver.ConfigAws(logs)

	logs.Check(err) // Verifica e registra erros usando o método Check de utils

	// Create a new HTTP Mux (router)
	mux := http.NewServeMux()

	// Define your routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.ResponseOK(w, r, logs)
	})
	mux.HandleFunc("/pontos", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetPunches(w, r, dynamoClient, logs)
	})
	mux.HandleFunc("/ponto/", func(w http.ResponseWriter, r *http.Request) {
		nome := r.URL.Path[len("/ponto/"):]
		handlers.PostPunch(nome, w, r, dynamoClient, logs)
	})
	// Add the new "/report/" route
	mux.HandleFunc("/report/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 4 {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		nome := parts[2]
		mes := parts[3]

		handlers.GetReport(w, r, dynamoClient, logs, nome, mes)
	})

	// Wrap your router with the SecureMiddleware
	secureMux := SecureMiddleware(mux)

	// Start the server
	if inLambda() {
		log.Fatal(gateway.ListenAndServe(":8080", secureMux))
	} else {
		log.Fatal(http.ListenAndServe(":8080", secureMux))
	}
}

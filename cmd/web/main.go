package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	appRouter := gin.New()
	appRouter.GET("/", func(ctx *gin.Context) {
		log.Println("Creating a scalable web application with Ginalabi")
	})

	err := appRouter.Run()
	check(err)
}

func check(erro error) {
	if erro != nil {
		log.Fatal(erro)
	}
}

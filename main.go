package main

import (
	"log"

	"github.com/SathvikPN/codepulse/services"
)

func main() {
	err := services.InitDB()
	if err != nil {
		log.Fatalf("database initialization failed, [err] %v", err)
	}

	router := services.Router()
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("REST api server failed, [err] %v", err)
	}

}

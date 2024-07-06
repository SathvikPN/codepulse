package main

import (
	"log"

	"github.com/SathvikPN/codepulse/services"
)

func main() {
	services.InitDB()
	router := services.Router()
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("REST api server failed, [err] %v", err)
	}

}

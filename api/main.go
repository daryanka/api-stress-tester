package main

import (
	"fmt"
	"github.com/daryanka/api-stress-tester/api/clients"
	"github.com/daryanka/api-stress-tester/api/controllers"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	utils.InitLogger()
	clients.InitDatabase()
	engine := controllers.StartRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	if err = engine.Run(fmt.Sprintf(":%v", port)); err != nil {
		panic(err.Error())
	}
}

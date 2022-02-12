package main

import (
	"go-ethereum-demo/config"
	"go-ethereum-demo/databases"
	"go-ethereum-demo/routes"
	"log"
	"os"
)

func main() {
	var err error

	// Init viper config
	err = config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Init database
	err = databases.InitMySql()
	if err != nil{
		log.Fatal(err)
	}
	defer databases.Close()

	// Register the routes
	r := routes.SetRouter()

	// start server on 8080 port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}

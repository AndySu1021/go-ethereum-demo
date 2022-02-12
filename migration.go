package main

import (
	"go-ethereum-demo/config"
	"go-ethereum-demo/databases"
	"go-ethereum-demo/migrations"
	"log"
)

func main()  {
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

	tables := []interface{}{
		&migrations.Block{},
		&migrations.BlockLog{},
		&migrations.Transaction{},
		&migrations.TransactionLog{},
	}

	err = databases.MySqlClient.AutoMigrate(tables...)
	if err != nil {
		log.Fatal(err)
	}
}

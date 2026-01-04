package main

import (
	"e-commerce-api/internal/app"
	"e-commerce-api/internal/config"
	"fmt"
	"log"
)

func main() {
	config, err := config.LoadENV()
	if err != nil {
		log.Fatalln(err)
	}

	app, err := app.New(config)
	if err != nil {
		log.Fatalf("error while trying to init app: %s", err.Error())
	}

	app.Engine.Run(fmt.Sprintf(":%d", config.PORT))
}

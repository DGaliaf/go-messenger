package main

import (
	"context"
	"log"
	"messenger-rest-api/app/internal/app"
	"messenger-rest-api/app/internal/config"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample messenger service.

// @host      localhost:9000
// @BasePath  /
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("config initializing")
	cfg := config.GetConfig()

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, err)
	}

	log.Println("Running Application")
	a.Run(ctx)
}

package main

import (
	"context"
	"log"
	"os"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/pkg/database"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())

	// connect DB
	db := database.DbConn(ctx, &cfg)

	defer db.Disconnect(ctx)

	log.Println(db)
}

package main

import (
	"context"
	"log"
	"os"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/pkg/database"
	"github.com/guatom999/Go-MicroService/server"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		log.Printf("choosen env is :%v", os.Args[1])
		return os.Args[1]
	}())

	// connect DB
	db := database.DbConn(ctx, &cfg)

	defer db.Disconnect(ctx)

	server.Start(ctx, &cfg, db)

	log.Println(db)

}

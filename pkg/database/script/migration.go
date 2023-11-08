package main

import (
	"context"
	"log"
	"os"

	"github.com/guatom999/Go-MicroService/config"
	"github.com/guatom999/Go-MicroService/pkg/database/migration"
)

func main() {
	ctx := context.Background()
	// _ = ctx
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())

	log.Printf("cfg is :%v", cfg)

	switch cfg.App.Name {
	case "auth":
		migration.AuthMigrate(ctx, &cfg)
	case "player":
	case "inventory":
	case "item":
	case "payment":
	}

}

package main

import (
	"flag"

	"context"
	"time"

	"github.com/mises-id/socialsvc/app/models"
	"github.com/mises-id/socialsvc/lib/db"
	_ "github.com/mises-id/socialsvc/lib/mises"

	// This Service
	"github.com/mises-id/socialsvc/handlers"
	"github.com/mises-id/socialsvc/svc/server"
)

func main() {
	// Update addresses if they have been overwritten by flags
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	db.SetupMongo(ctx)
	models.EnsureIndex()

	cfg := server.DefaultConfig
	cfg = handlers.SetConfig(cfg)

	server.Run(cfg)
}

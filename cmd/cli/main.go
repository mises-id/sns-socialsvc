package main

import (
	"flag"
	"fmt"

	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/services/session"
	_ "github.com/mises-id/sns-socialsvc/config"
	"github.com/mises-id/sns-socialsvc/lib/airdrop"
	"github.com/mises-id/sns-socialsvc/lib/db"
	_ "github.com/mises-id/sns-socialsvc/lib/mises"

	// This Service
	"github.com/mises-id/sns-socialsvc/handlers"
	"github.com/mises-id/sns-socialsvc/svc/server"
)

func main() {
	// Update addresses if they have been overwritten by flags
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	fmt.Println("setup mongo...")
	db.SetupMongo(ctx)
	models.EnsureIndex()

	fmt.Println("setup mises...")
	session.SetupMisesClient()
	fmt.Println("setup airdrop...")
	airdrop.SetAirdropClient()

	cfg := server.DefaultConfig
	cfg = handlers.SetConfig(cfg)

	server.Run(cfg)
}

package db

import (
	"context"

	"github.com/mises-id/socialsvc/config/env"
	"github.com/mises-id/socialsvc/lib/db/odm"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoDB   *mongo.Database
	odmClient *odm.Client
)

func SetupMongo(ctx context.Context) {
	clientOpts := options.Client().SetMaxPoolSize(30).ApplyURI(env.Envs.MongoURI)
	if env.Envs.DBUser != "" {
		clientOpts = clientOpts.SetAuth(options.Credential{
			Username: env.Envs.DBUser,
			Password: env.Envs.DBPass,
		})
	}
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		panic(err)
	}
	mongoDB = client.Database(env.Envs.DBName)
	odmClient = odm.NewClient(mongoDB)
}

func DB() *mongo.Database {
	return mongoDB
}

func ODM(ctx context.Context) *odm.DB {
	return odmClient.NewSession(ctx)
}

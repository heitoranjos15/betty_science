package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"betty/science/app/league_of_legends/client/riot"
	"betty/science/app/league_of_legends/core"
	"betty/science/app/league_of_legends/integrations"
	"betty/science/app/league_of_legends/repo"
	"betty/science/config"
)

func main() {
	cfg := config.LoadConfig()
	api := integrations.NewLeagueOfLegendsAPI(cfg.RiotAPIKey)
	client := riot.NewClientMatch(api)

	db := initMongoDB(cfg.MongoURI, cfg.MongoDB)
	matchRepository := repo.NewMongoRepo("matches", db)
	teamRepository := repo.NewMongoRepo("teams", db)

	coreMatch := core.NewMatchCore(client, matchRepository, teamRepository)
	err := coreMatch.Load()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func initMongoDB(uri, dbName string) *mongo.Database {
	// monitor := &event.CommandMonitor{
	// 	Started: func(_ context.Context, evt *event.CommandStartedEvent) {
	// 		fmt.Printf("MongoDB Command: %s %v\n", evt.CommandName, evt.Command)
	// 	},
	// }
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	return client.Database(dbName)
}

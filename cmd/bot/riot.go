package main

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"betty/science/app/league_of_legends/client/riot"
	"betty/science/app/league_of_legends/core"
	"betty/science/app/league_of_legends/integrations"
	"betty/science/app/league_of_legends/repo"
	"betty/science/config"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalln("Usage: riot <command>\nCommands:\n  next-match\n  load-frames")
		panic("Not enough arguments")
	}

	cfg := config.LoadConfig()
	db := initMongoDB(cfg.MongoURI, cfg.MongoDB)

	bot := Bot{
		Name:  "Riot Bot",
		cfg:   cfg,
		mongo: db,
	}

	commands := map[string][]func(){
		"schedule":       {bot.nextMatchBot},
		"update_matches": {bot.nextMatchBot, bot.loadGamesBot},
		"complete_games": {bot.loadGamesBot, bot.loadFramesBot},
		"load_frames":    {bot.loadFramesBot},
	}

	cmd := args[1]
	if len(cmd) > 4 && cmd[:4] == "CMD=" {
		cmd = cmd[4:]
	}
	cmdFunc, exists := commands[cmd]
	if !exists {
		log.Fatalf("Unknown command: %s\n", args[1])
		panic("Unknown command")
	}

	for _, cfunc := range cmdFunc {
		cfunc()
	}

}

type Bot struct {
	Name  string
	cfg   *config.Config
	mongo *mongo.Database
}

func (b Bot) nextMatchBot() {
	api := integrations.NewLeagueOfLegendsAPI(b.cfg.RiotAPIKey)
	client := riot.NewClientMatch(api)
	matchRepository := repo.NewMongoRepo("matches", b.mongo)
	teamRepository := repo.NewMongoRepo("teams", b.mongo)
	coreMatch := core.NewMatchCore(client, matchRepository, teamRepository)
	err := coreMatch.Load()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func (b Bot) loadGamesBot() {
	api := integrations.NewLeagueOfLegendsAPI(b.cfg.RiotAPIKey)
	client := riot.NewClientGame(api)
	gameRepository := repo.NewMongoRepo("games", b.mongo)
	teamRepository := repo.NewMongoRepo("teams", b.mongo)
	matchRepository := repo.NewMongoRepo("matches", b.mongo)
	coreGame := core.NewGameCore(b.cfg, client, gameRepository, teamRepository, matchRepository)
	err := coreGame.Load()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func (b Bot) loadFramesBot() {
	api := integrations.NewLeagueOfLegendsAPI(b.cfg.RiotAPIKey)
	client := riot.NewFramesClient(api)
	frameRepository := repo.NewMongoRepo("frames", b.mongo)
	gameRepository := repo.NewMongoRepo("games", b.mongo)
	playersRepository := repo.NewMongoRepo("players", b.mongo)
	coreFrame := core.NewFrameCore(client, frameRepository, gameRepository, playersRepository)
	err := coreFrame.Load()
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

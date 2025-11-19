package main

import (
	"context"
	"os"

	"betty/science/app/league_of_legends/client/riot"
	"betty/science/app/league_of_legends/core"
	"betty/science/app/league_of_legends/integrations"
	"betty/science/app/league_of_legends/machine"
	"betty/science/app/league_of_legends/repo"
	"betty/science/app/league_of_legends/worker"
	"betty/science/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.LoadConfig()
	mongoDB := initMongoDB(cfg.MongoURI, cfg.MongoDB)

	command := &Command{
		cfg:   cfg,
		mongo: mongoDB,
	}

	commands := map[string][]func(){
		"schedule": {command.RunMatch},
		"game":     {command.RunGame},
		"frame":    {command.RunFrame},
	}

	args := os.Args
	cmd := args[1]
	if len(cmd) > 4 && cmd[:4] == "CMD=" {
		cmd = cmd[4:]
	}
	cmdFunc, exists := commands[cmd]
	if !exists {
		panic("Unknown command")
	}

	for _, cfunc := range cmdFunc {
		cfunc()
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

type Command struct {
	cfg   config.Config
	mongo *mongo.Database
}

type CmdSetup struct {
	api             *integrations.LeagueOfLegendsClient
	clientMatch     *riot.ClientMatch
	matchRepository *repo.MongoRepo
	teamRepository  *repo.MongoRepo
	gameRepository  *repo.MongoRepo
}

func (c *Command) defaultSetup() *CmdSetup {
	api := integrations.NewLeagueOfLegendsAPI(c.cfg.RiotAPIKey)

	clientMatch := riot.NewClientMatch(api)

	matchRepository := repo.NewMongoRepo("matches", c.mongo)
	teamRepository := repo.NewMongoRepo("teams", c.mongo)
	gameRepository := repo.NewMongoRepo("games", c.mongo)
	return &CmdSetup{
		api:             api,
		clientMatch:     clientMatch,
		matchRepository: matchRepository,
		teamRepository:  teamRepository,
		gameRepository:  gameRepository,
	}
}

func (c *Command) RunMatch() {
	setup := c.defaultSetup()
	coreMatch := core.NewMatchCore(setup.matchRepository)
	coreTeam := core.NewTeamCore(setup.teamRepository)
	workerMatch := worker.NewWorkerMatch(setup.clientMatch, coreMatch, coreTeam)
	botMatch := machine.NewBot(workerMatch, c.cfg)
	botMatch.Execute()
}

func (c *Command) RunGame() {
	setup := c.defaultSetup()
	clientGame := riot.NewClientGame(setup.api)
	coreMatch := core.NewMatchCore(setup.matchRepository)
	coreGame := core.NewGameCore(&c.cfg, clientGame, setup.gameRepository, setup.teamRepository, setup.matchRepository, repo.NewMongoRepo("players", c.mongo))
	workerGame := worker.NewWorkerGame(clientGame, coreGame, coreMatch)
	botGame := machine.NewBot(workerGame, c.cfg)
	botGame.Execute()
}

func (c *Command) RunFrame() {
	setup := c.defaultSetup()
	clientFrame := riot.NewFramesClient(setup.api)
	frameRepository := repo.NewMongoRepo("frames", c.mongo)
	playersRepository := repo.NewMongoRepo("players", c.mongo)
	coreFrame := core.NewFrameCore(clientFrame, frameRepository, setup.gameRepository, playersRepository)
	clientGame := riot.NewClientGame(setup.api)
	coreGame := core.NewGameCore(&c.cfg, clientGame, setup.gameRepository, setup.teamRepository, setup.matchRepository, playersRepository)
	workerFrame := worker.NewWorkerFrame(clientFrame, coreFrame, coreGame)
	botFrame := machine.NewBot(workerFrame, c.cfg)
	botFrame.Execute()
}

package core

import (
	"context"
	"log"
	"slices"
	"time"

	"betty/science/app/league_of_legends/client"
	"betty/science/app/league_of_legends/models"
	"betty/science/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type frameClient interface {
	LoadData(game models.Game) (client.FrameResponse, error)
}

type FrameChannel struct {
	frame client.FrameResponse
	g     models.Game
	err   error
}

type FrameCore struct {
	client    frameClient
	db        frameDB
	gameDB    gameDB
	playersDB playersDB
}

func NewFrameCore(client frameClient, db frameDB, gameDB gameDB, playersDB playersDB) *FrameCore {
	return &FrameCore{
		client:    client,
		db:        db,
		gameDB:    gameDB,
		playersDB: playersDB,
	}
}

func (ec *FrameCore) Load() error {
	ctx := context.Background()

	games, err := ec.gameDB.GetGames(ctx, bson.M{
		// "external_id": "113475871524050784",
		// "state": "completed",
		// "load_state": bson.M{
		// 	"$ne": "loaded",
		// },
		"load_state": "without_frames", // TODO: better filter to not overload Riot API
	})
	if err != nil {
		log.Println("[core-frame] Error fetching games:", err)
		return err
	}

	botCfg := config.SetupBot()
	bots := botCfg.Bots
	frameChan := make(chan FrameChannel)

	limitBots := len(games)
	if limitBots >= botCfg.Workers {
		limitBots = limitBots / botCfg.Workers
	}

	for i := 0; i < botCfg.Workers; i++ {
		bot := bots[i]
		log.Printf("[%s] hey", bot.Name)
		if i*limitBots >= len(games) {
			log.Printf("[%s] i dont need to work", bot.Name)
		}

		loadGames := games[i*limitBots : (i+1)*limitBots]
		go ec.loadGamesWorker(bot, loadGames, frameChan)
	}

	ec.processLoadedFrames(ctx, frameChan)

	log.Printf("[core-frame] loaded frames for %d games", len(games))

	return nil
}

func (ec *FrameCore) processLoadedFrames(ctx context.Context, frameChan <-chan FrameChannel) {
	for frameResp := range frameChan {
		ec.saveGame(ctx, frameResp)
		ec.savePlayers(ctx, frameResp)
		ec.saveFrames(ctx, frameResp)
	}

}

func (ec *FrameCore) saveGame(ctx context.Context, chanResp FrameChannel) {
	log.Printf("[core-frame] Saving game data for game %s", chanResp.g.ExternalID)
	game := chanResp.g
	frame := chanResp.frame
	if chanResp.err != nil {
		log.Printf("[core-frame] Error loading frames for game %s: %v", game.ExternalID, chanResp.err)
	}
	game.StartTime = frame.GameStart
	game.Duration = frame.GameEnd.Sub(chanResp.frame.GameStart).Seconds()
	game.Winner = frame.WinnerID
	game.LoadState = "loaded"

	players := ec.getPlayersData(frame.Players) // TODO: method that injects PlayerID
	gamePlayers := []models.GamePlayer{}
	for i, p := range frame.Players {
		gamePlayers = append(gamePlayers, models.GamePlayer{
			PlayerID: players[i].ID,
			Role:     p.Role,
			Champion: p.Champion,
			Side:     p.Side,
		})
	}
	game.Players = gamePlayers

	err := ec.gameDB.SaveGame(ctx, game)
	if err != nil {
		log.Printf("[core-frame] Error updating game %s: %v", game.ExternalID, err)
	}

}

func (ec *FrameCore) savePlayers(ctx context.Context, chanResp FrameChannel) {
	log.Printf("[core-frame] Saving players data for game %s", chanResp.g.ExternalID)
	frame := chanResp.frame
	players := ec.getPlayersData(frame.Players) // TODO: method that injects PlayerID
	game := chanResp.g
	if err := ec.playersDB.SaveBulkPlayers(ctx, players); err != nil {
		log.Printf("[core-frame] Error saving players for game %s: %v", game.ExternalID, err)
	}
}

func (ec *FrameCore) saveFrames(ctx context.Context, chanResp FrameChannel) {
	log.Printf("[core-frame] Saving frames data for game %s", chanResp.g.ExternalID)
	frame := chanResp.frame
	players := ec.getPlayersData(frame.Players) // TODO: method that injects PlayerID
	for i, pData := range players {
		frame.Frame.Players[i].PlayerID = pData.ID
	}

	game := chanResp.g
	err := ec.db.SaveFrame(ctx, frame.Frame)
	if err != nil {
		log.Printf("[core-frame] Error saving frame for game %s: %v", game.ExternalID, err)
	}
}

func (ec *FrameCore) loadGamesWorker(bot *config.Bot, games []models.Game, outChan chan<- FrameChannel) {
	defer close(outChan)
	log.Printf("[%s] loading frames for %d games", bot.Name, len(games))

	for _, game := range games {
		frame, err := ec.client.LoadData(game)

		outChan <- FrameChannel{frame: frame, g: game, err: err}

		log.Printf("[%s] loaded frames for game %s", bot.Name, game.ExternalID)
		time.Sleep(time.Duration(bot.DelaySeconds) * time.Second) // to avoid rate limiting
	}
}

func (ec *FrameCore) getPlayersData(playerFrame []client.PlayerFrame) []models.Player {
	var players []models.Player

	for _, pf := range playerFrame {
		existPlayer, err := ec.playersDB.GetPlayerByExternalID(context.Background(), pf.ExternalID)
		if err != nil {
			log.Printf("[core-frame] Player with external ID %s not found, creating new one", pf.ExternalID)
			player := models.Player{
				ExternalID: pf.ExternalID,
				Name:       pf.Name,
				Roles:      []string{pf.Role},
				Teams:      []primitive.ObjectID{pf.TeamID},
				ActualRole: pf.Role,
				ActualTeam: pf.TeamID,
			}
			players = append(players, player)
			continue
		}

		if len(existPlayer.Roles) == 0 || !slices.Contains(existPlayer.Roles, pf.Role) {
			existPlayer.ActualRole = pf.Role
			existPlayer.Roles = append(existPlayer.Roles, pf.Role)
		}
		if len(existPlayer.Teams) == 0 || !slices.Contains(existPlayer.Teams, pf.TeamID) {
			existPlayer.ActualTeam = pf.TeamID
			existPlayer.Teams = append(existPlayer.Teams, pf.TeamID)
		}
		players = append(players, existPlayer)

	}
	return players
}

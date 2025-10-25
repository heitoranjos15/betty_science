package core

import (
	"context"
	"log"
	"slices"
	"time"

	"betty/science/app/league_of_legends/client"
	"betty/science/app/league_of_legends/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type frameTeamClient interface {
	LoadData(game models.Game) (client.FrameResponse, error)
}

type TeamFrameCore struct {
	client    frameTeamClient
	db        frameDB
	gameDB    gameDB
	playersDB playersDB
}

func NewTeamFrameCore(client frameTeamClient, db frameDB, gameDB gameDB, playersDB playersDB) *TeamFrameCore {
	return &TeamFrameCore{
		client:    client,
		db:        db,
		gameDB:    gameDB,
		playersDB: playersDB,
	}
}

func (ec *TeamFrameCore) Load() error {
	ctx := context.Background()

	games, err := ec.gameDB.GetGames(ctx, bson.M{
		// "external_id": "113475798006664639",
		"state": "completed",
		"load_state": bson.M{
			"$ne": "loaded",
		},
	})
	if err != nil {
		log.Println("[core-frame] Error fetching games:", err)
		return err
	}

	for i, game := range games {
		time.Sleep(2 * time.Second) // to avoid rate limiting

		frame, err := ec.client.LoadData(game)
		if err != nil {
			log.Printf("[core-frame] Error loading frames for game %s: %v", game.ExternalID, err)
			games[i].LoadState = "error"
			games[i].ErrorMsg = err.Error()
			err = ec.gameDB.SaveGame(ctx, games[i])
			if err != nil {
				log.Printf("[core-frame] Error saving error state for game %s: %v", game.ExternalID, err)
			}
			continue
		}

		games[i].StartTime = frame.GameStart
		games[i].Duration = frame.GameEnd.Sub(frame.GameStart).Seconds()
		games[i].Winner = frame.WinnerID

		players := ec.getPlayersData(frame.Players)
		if err := ec.playersDB.SaveBulkPlayers(ctx, players); err != nil {
			log.Printf("[core-frame] Error saving players for game %s: %v", game.ExternalID, err)
		}

		gamePlayers := []models.GamePlayer{}
		for j, p := range frame.Players {
			gamePlayers = append(gamePlayers, models.GamePlayer{
				PlayerID: players[j].ID,
				Role:     p.Role,
				Champion: p.Champion,
				Side:     p.Side,
			})
			log.Printf("[core-frame] Player %d: %s, Role: %s, TeamID: %s", j+1, p.Name, p.Role, p.TeamID.Hex())
		}
		games[i].Players = gamePlayers

		games[i].LoadState = "loaded"
		err = ec.gameDB.SaveGame(ctx, games[i])
		if err != nil {
			log.Printf("[core-frame] Error updating game %s: %v", game.ExternalID, err)
			continue
		}

		err = ec.db.SaveFrame(ctx, frame.Frame)
		if err != nil {
			log.Printf("[core-frame] Error saving frame for game %s: %v", game.ExternalID, err)
			continue
		}

		log.Printf("[core-frame] loaded frames for game %s", game.ExternalID)
	}

	log.Printf("[core-frame] loaded frames for %d games", len(games))

	return nil
}

func (ec *TeamFrameCore) getPlayersData(playerFrame []client.PlayerFrame) []models.Player {
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

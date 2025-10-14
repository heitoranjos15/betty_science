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
	db        gameDB
	playersDB playersDB
}

func NewTeamFrameCore(client frameTeamClient, db gameDB, playersDB playersDB) *TeamFrameCore {
	return &TeamFrameCore{
		client:    client,
		db:        db,
		playersDB: playersDB,
	}
}

func (ec *TeamFrameCore) Load() error {
	ctx := context.Background()

	games, err := ec.db.GetGames(ctx, bson.M{"state": "completed"})
	if err != nil {
		log.Println("[core-frame] Error fetching games:", err)
		return err
	}

	for i, game := range games {
		frame, err := ec.client.LoadData(game)
		if err != nil {
			log.Printf("[core-frame] Error loading frames for game %s: %v", game.ExternalID, err)
			continue
		}
		games[i].Frames = append(games[i].Frames, frame.Frame)

		// if i == 1 {
		// 	players := ec.playersEnrich(frame.PlayerGamesDetails)
		// }

		time.Sleep(2 * time.Second) // to avoid rate limiting
	}

	return nil
}

func (ec *TeamFrameCore) playersEnrich(playerFrame []models.GamePlayer) []models.Player {
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

		if !slices.Contains(existPlayer.Roles, pf.Role) {
			existPlayer.Roles = append(existPlayer.Roles, pf.Role)
			existPlayer.ActualRole = pf.Role
		}
		if !slices.Contains(existPlayer.Teams, pf.TeamID) {
			existPlayer.Teams = append(existPlayer.Teams, pf.TeamID)
			existPlayer.ActualTeam = pf.TeamID
		}

	}
	return players
}

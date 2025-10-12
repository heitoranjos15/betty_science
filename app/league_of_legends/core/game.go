package core

import (
	"betty/science/app/league_of_legends/models"
	"context"
	"log"
)

type gameClient interface {
	LoadData(match models.Match) ([]models.Game, error)
}

type GameCore struct {
	client gameClient
	db     gameDB
	teamDB teamDB
}

func NewGameCore(client gameClient, db gameDB, teamDB teamDB) *GameCore {
	return &GameCore{
		client: client,
		db:     db,
		teamDB: teamDB,
	}
}

func (ec *GameCore) Load(matchExternalID string) error {
	ctx := context.Background()
	game, err := ec.client.LoadData(models.Match{ExternalID: matchExternalID})

	if err != nil {
		log.Println("[core-game] Error loading game data:", err)
		return err
	}

	for _, g := range game {
		for i, t := range g.Teams {
			existingTeam, err := ec.teamDB.GetTeamByName(ctx, t.Name)
			if err != nil {
				log.Println("[core-game] team not found, creating new one:", t.Name)
				continue
			}
			if existingTeam.ExternalID == "" {
				err = ec.teamDB.UpdateTeamExternalID(ctx, existingTeam.ID, t.ExternalID)
			}
			g.Teams[i].ID = existingTeam.ID
		}
	}

	if err := ec.db.SaveBulkGames(ctx, game); err != nil {
		log.Println("[core-game] Error saving game data:", err)
		return err
	}
	return nil
}

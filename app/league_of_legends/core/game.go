package core

import (
	"betty/science/app/league_of_legends/models"
	"betty/science/config"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type gameClient interface {
	LoadData(match models.Match) ([]models.Game, error)
}

type GameCore struct {
	cfg     *config.Config
	client  gameClient
	db      gameDB
	teamDB  teamDB
	matchDB matchDB
}

func NewGameCore(cfg *config.Config, client gameClient, db gameDB, teamDB teamDB, matchDB matchDB) *GameCore {
	return &GameCore{
		cfg:     cfg,
		client:  client,
		db:      db,
		teamDB:  teamDB,
		matchDB: matchDB,
	}
}

func (ec *GameCore) Load() error {
	ctx := context.Background()

	matches, err := ec.matchDB.GetMatches(ctx, bson.M{"load_state": "without_games"})

	if err != nil {
		log.Println("[core-game] Error fetching matches:", err)
	}

	for _, m := range matches {
		games, err := ec.loadGamesForMatch(ctx, m)
		if err != nil {
			log.Printf(fmt.Sprintf("[core-game] Error loading games for match %s:", m.ExternalID), err)
			continue
		}

		if err := ec.db.SaveBulkGames(ctx, games); err != nil {
			log.Println("[core-game] Error saving game data:", err)
			return err
		}
		m.LoadState = "loaded"
		if err := ec.matchDB.SaveMatch(ctx, m); err != nil {
			log.Println("[core-game] Error updating match load state:", err)
		}

		log.Printf("[core-game] loaded %d games for match %s", len(games), m.ExternalID)
		time.Sleep(2 * time.Second) // to avoid rate limiting
	}

	return nil
}

func (ec *GameCore) loadGamesForMatch(ctx context.Context, match models.Match) ([]models.Game, error) {
	games, err := ec.client.LoadData(match)

	if err != nil {
		log.Println("[core-game] Error loading game data:", err)
		return games, err
	}

	for i, g := range games {
		gameTeams, err := ec.setGameTeamIDs(ctx, g)
		if err != nil {
			log.Println("[core-game] Error setting game team IDs:", err)
			return games, err
		}
		games[i].Teams = gameTeams
		games[i].MatchID = match.ID
	}

	return games, nil
}

func (ec *GameCore) setGameTeamIDs(ctx context.Context, game models.Game) ([]models.GameTeam, error) {
	var gameTeams []models.GameTeam
	for _, t := range game.Teams {
		existingTeam, err := ec.teamDB.GetTeamByName(ctx, t.Name)
		if err != nil {
			log.Println("[core-game] team not found", t.Name)
			continue
		}
		if existingTeam.ExternalID == "" {
			// fallback to updating external ID if not set
			err = ec.teamDB.UpdateTeamExternalID(ctx, existingTeam.ID, t.ExternalID)
		}
		gameTeams = append(gameTeams, models.GameTeam{
			ID:         existingTeam.ID,
			ExternalID: t.ExternalID,
			Name:       t.Name,
			Side:       t.Side,
		})
	}
	return gameTeams, nil
}

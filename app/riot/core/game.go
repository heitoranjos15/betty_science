package core

import (
	models "betty/science/app/riot"
	"betty/science/app/riot/clients/league"
	"betty/science/app/riot/repo"
	"context"
	"log"
	"slices"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type gameClient interface {
	LoadData(match models.Match) ([]models.Game, error)
}

type GameCore struct {
	client    gameClient
	db        gameDB
	teamDB    teamDB
	matchDB   matchDB
	playersDB playersDB
}

func NewGameCore(client gameClient, db gameDB, teamDB teamDB, matchDB matchDB, playersDB playersDB) *GameCore {
	return &GameCore{
		client:    client,
		db:        db,
		teamDB:    teamDB,
		matchDB:   matchDB,
		playersDB: playersDB,
	}
}

func (ec *GameCore) UpdateGameByFrameResp(game models.Game, resp league.FrameResponse) error {
	ctx := context.Background()

	gamePlayers := []repo.GamePlayer{}
	for _, p := range resp.Players {
		player := ec.getPlayerData(ctx, p)
		if err := ec.playersDB.SaveBulkPlayers(ctx, []repo.Player{player}); err != nil { // TODO NO NEED TO USE BULK
			log.Printf("[core-frame] Error saving player %s: %v", player.ExternalID, err)
			return err
		}
		gamePlayers = append(gamePlayers, repo.GamePlayer{
			PlayerID: player.ID,
			Role:     p.Role,
			Champion: p.Champion,
			Side:     p.Side,
		})
	}

	updates := bson.M{
		"load_state":     "loaded_frames",
		"start_time":     resp.GameStart,
		"duration":       resp.GameEnd.Sub(resp.GameStart).Seconds(),
		"winner_team_id": resp.WinnerID,
		"players":        gamePlayers,
	}

	return ec.db.UpdateGameByExternalID(ctx, game.ExternalID, updates)
}

func (ec *GameCore) getPlayerData(ctx context.Context, pf models.GamePlayer) repo.Player {
	existPlayer, err := ec.playersDB.GetPlayerByExternalID(ctx, pf.ExternalID)
	if err != nil {
		log.Printf("[core-frame] Player with external ID %s not found, creating new one", pf.ExternalID)
		player := repo.Player{
			ExternalID: pf.ExternalID,
			Name:       pf.Name,
			Roles:      []string{pf.Role},
			Teams:      []primitive.ObjectID{pf.TeamID},
			ActualRole: pf.Role,
			ActualTeam: pf.TeamID,
		}
		return player
	}

	if len(existPlayer.Roles) == 0 || !slices.Contains(existPlayer.Roles, pf.Role) {
		existPlayer.ActualRole = pf.Role
		existPlayer.Roles = append(existPlayer.Roles, pf.Role)
	}
	if len(existPlayer.Teams) == 0 || !slices.Contains(existPlayer.Teams, pf.TeamID) {
		existPlayer.ActualTeam = pf.TeamID
		existPlayer.Teams = append(existPlayer.Teams, pf.TeamID)
	}
	return existPlayer
}

func (ec *GameCore) LoadBulk() ([]models.Match, error) {
	ctx := context.Background()
	matches := []models.Match{}

	log.Println("[core-game] fetching matches from DB to load games for")

	query := bson.M{
		// "external_id": bson.M{"$in": []string{
		// 	"113475871524050783",
		// 	"115016265550803595",
		// 	"113475871523985235",
		// 	"113475798006664622",
		// 	"115016265550803589",
		// }},
		"load_state": "without_games",
	}
	dbResult, err := ec.matchDB.GetMatches(ctx, query)
	log.Println("[core-game] fetched matches from DB:", len(dbResult))
	if err != nil {
		log.Println("[core-game] Error fetching matches:", err)
		return matches, err
	}

	for _, result := range dbResult {
		data := models.Match{
			ID:         result.ID,
			ExternalID: result.ExternalID,
		}
		matches = append(matches, data)
	}

	log.Printf("[core-game] fetched %d matches to load games for", len(matches))

	return matches, nil
}

func (ec *GameCore) SaveBulk(games []models.Game, match models.Match) error {
	ctx := context.Background()
	saveData := []repo.Game{}

	for _, g := range games {
		teams := []repo.GameTeam{}
		for _, t := range g.Teams {
			team, err := ec.teamDB.GetTeamByExternalID(ctx, t.ExternalID)
			if err != nil {
				log.Printf("[core-game] Error fetching team by external ID %s: %v", t.ExternalID, err)
				return err
			}
			teams = append(teams, repo.GameTeam{
				ID:   team.ID,
				Side: t.Side,
			})
		}

		data := repo.Game{
			MatchID:      match.ID,
			Number:       g.Number,
			ExternalID:   g.ExternalID,
			State:        g.State,
			Duration:     g.Duration,
			WinnerTeamID: g.WinnerTeamID,
			StartTime:    g.StartTime,
			Teams:        teams,
			LoadState:    "without_frames",
		}
		saveData = append(saveData, data)
	}
	return ec.db.SaveBulkGames(ctx, saveData)
}

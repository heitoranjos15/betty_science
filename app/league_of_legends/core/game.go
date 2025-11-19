package core

import (
	"betty/science/app/league_of_legends/models"
	"betty/science/app/league_of_legends/repo"
	"betty/science/config"
	"context"
	"log"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type gameClient interface {
	LoadData(match models.Match) ([]models.Game, error)
}

type outChanType struct {
	game repo.Game
	m    repo.Match
	err  error
}

type GameCore struct {
	cfg       *config.Config
	client    gameClient
	db        gameDB
	teamDB    teamDB
	matchDB   matchDB
	playersDB playersDB
}

func NewGameCore(cfg *config.Config, client gameClient, db gameDB, teamDB teamDB, matchDB matchDB, playersDB playersDB) *GameCore {
	return &GameCore{
		cfg:       cfg,
		client:    client,
		db:        db,
		teamDB:    teamDB,
		matchDB:   matchDB,
		playersDB: playersDB,
	}
}

func (ec *GameCore) UpdateGameByFrameResp(game models.Game, resp models.FrameResponse) error {
	ctx := context.Background()

	gamePlayers := []models.GamePlayer{}
	for _, p := range resp.Players {
		player := ec.getPlayerData(ctx, p)
		if err := ec.playersDB.SaveBulkPlayers(ctx, []repo.Player{player}); err != nil { // TODO NO NEED TO USE BULK
			log.Printf("[core-frame] Error saving player %s: %v", player.ExternalID, err)
		}
		gamePlayers = append(gamePlayers, models.GamePlayer{
			PlayerID: player.ID,
			Role:     p.Role,
			Champion: p.Champion,
			Side:     p.Side,
		})
	}

	updates := bson.M{
		"load_state": "loaded_frames",
		"start_time": game.StartTime,
		"duration":   resp.GameEnd.Sub(resp.GameStart).Seconds(),
		"winner":     game.Winner,
		"players":    gamePlayers,
	}

	return ec.db.UpdateGameByExternalID(ctx, game.ExternalID, updates)
}

func (ec *GameCore) getPlayerData(ctx context.Context, pf models.PlayerGameInfo) repo.Player {
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

func (ec *GameCore) SaveBulk(games []models.Game) error {
	ctx := context.Background()
	saveData := []repo.Game{}

	for _, g := range games {
		data := repo.Game{
			ExternalID: g.ExternalID,
			State:      g.State,
			Duration:   g.Duration,
			Winner:     g.Winner,
			StartTime:  g.StartTime,
			LoadState:  "without_frames",
		}
		saveData = append(saveData, data)
	}
	return ec.db.SaveBulkGames(ctx, saveData)
}

func (ec *GameCore) Load() error {
	ctx := context.Background()

	matches, err := ec.matchDB.GetMatches(ctx, bson.M{"load_state": "without_games"})

	if err != nil {
		log.Println("[core-game] Error fetching matches:", err)
	}

	if len(matches) == 0 {
		log.Println("[core-game] No matches to load games for")
		return nil
	}
	bots := ec.cfg.Workers
	outChan := make(chan outChanType)

	type bot struct {
		name  string
		emote string
	}
	botsNames := []bot{
		{name: "Woody", emote: "🤠"},
		{name: "Sasuke", emote: ":|"},
		{name: "Gash", emote: ":D"},
	}
	limitBots := len(matches)
	if len(matches) >= bots {
		limitBots = len(matches) / bots
	}

	for i := 0; i < ec.cfg.Workers; i++ {
		bot := botsNames[i]
		log.Printf("[%s] hey", bot.name)
		if i*limitBots >= len(matches) {
			log.Printf("[%s] i dont need to work %s", bot.name, bot.emote)
			break
		}

		loadMatches := matches[i*limitBots : (i+1)*limitBots]
		log.Printf("[%s] loading %d matches %s", bot.name, len(loadMatches), bot.emote)

		go ec.loadGamesForMatch(loadMatches, outChan)
	}

	matchesToSave := []repo.Match{}
	gamesToSave := []repo.Game{}

	for matchLoaded := range outChan {
		log.Printf("[core-game] processed match %s", matchLoaded.m.ExternalID)
		if matchLoaded.err != nil {
			log.Printf("[core-game]error processing %s: %v", matchLoaded.m.ExternalID, matchLoaded.err)
			continue
		}
		g := matchLoaded.game
		m := matchLoaded.m

		if g.State == "completed" {
			m.LoadState = "loaded"
		}
		gamesToSave = append(gamesToSave, g)
		matchesToSave = append(matchesToSave, m)
	}

	if err := ec.db.SaveBulkGames(ctx, gamesToSave); err != nil {
		log.Println("[core-game] Error saving game data:", err)
		return err
	}

	if err := ec.matchDB.SaveBulkMatches(ctx, matchesToSave); err != nil {
		log.Println("[core-game] Error updating match load state:", err)
		return err
	}

	log.Printf("[core-game] loaded %d games", len(gamesToSave))
	log.Printf("[core-game] updated %d matches", len(matchesToSave))
	return nil
}

func (ec *GameCore) loadGamesForMatch(matches []repo.Match, outChan chan<- outChanType) {
	defer close(outChan)
	log.Printf("[core-game] loading games for %d matches", len(matches))

	for _, match := range matches {
		matchLoadData := models.Match{
			ID:         match.ID,
			ExternalID: match.ExternalID,
		}
		games, err := ec.client.LoadData(matchLoadData)

		if err != nil {
			log.Println("[core-game] Error loading game data:", err)
			outChan <- outChanType{game: repo.Game{}, m: match, err: err}
			continue
		}

		for _, g := range games {
			// gameTeams, err := ec.setGameTeamIDs(ctx, g)
			gameData := repo.Game{
				ExternalID: g.ExternalID,
				State:      g.State,
				Duration:   g.Duration,
				Winner:     g.Winner,
				StartTime:  g.StartTime,
				LoadState:  "without_frames",
			}

			if err != nil {
				log.Println("[core-game] Error setting game team IDs:", err)
				continue
			}
			// g.Teams = gameTeams
			// g.MatchID = match.ID
			outChan <- outChanType{game: gameData, m: match, err: nil}
			log.Printf("[core-game] loaded game %s for match %s", g.ExternalID, match.ExternalID)
			time.Sleep(2 * time.Second) // to avoid rate limiting
		}
	}
}

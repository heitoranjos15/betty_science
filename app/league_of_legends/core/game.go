package core

import (
	"betty/science/app/league_of_legends/models"
	"betty/science/config"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type gameClient interface {
	LoadData(match models.Match) ([]models.Game, error)
}

type outChanType struct {
	game models.Game
	m    models.Match
	err  error
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

		go ec.loadGamesForMatch(ctx, loadMatches, outChan)
	}

	matchesToSave := []models.Match{}
	gamesToSave := []models.Game{}

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

func (ec *GameCore) loadGamesForMatch(ctx context.Context, matches []models.Match, outChan chan<- outChanType) {
	defer close(outChan)
	log.Printf("[core-game] loading games for %d matches", len(matches))

	for _, match := range matches {
		games, err := ec.client.LoadData(match)

		if err != nil {
			log.Println("[core-game] Error loading game data:", err)
			outChan <- outChanType{game: models.Game{}, m: match, err: err}
			continue
		}

		for _, g := range games {
			gameTeams, err := ec.setGameTeamIDs(ctx, g)

			if err != nil {
				log.Println("[core-game] Error setting game team IDs:", err)
				continue
			}
			g.Teams = gameTeams
			g.MatchID = match.ID
			outChan <- outChanType{game: g, m: match, err: nil}
			log.Printf("[core-game] loaded game %s for match %s", g.ExternalID, match.ExternalID)
			time.Sleep(2 * time.Second) // to avoid rate limiting
		}
	}
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

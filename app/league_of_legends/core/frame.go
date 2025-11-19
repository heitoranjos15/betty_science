package core

import (
	"context"
	"log"

	"betty/science/app/league_of_legends/models"
	"betty/science/app/league_of_legends/repo"

	"go.mongodb.org/mongo-driver/bson"
)

type frameClient interface {
	LoadData(game models.Game) (models.FrameResponse, error)
}

type FrameChannel struct {
	frame models.FrameResponse
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

func (ec *FrameCore) LoadBulk() ([]models.Game, error) {
	ctx := context.Background()

	games, err := ec.gameDB.GetGames(ctx, bson.M{
		"load_state": "without_frames",
	})
	if err != nil {
		log.Println("[core-frame] Error fetching games:", err)
		return nil, err
	}

	gamesData := []models.Game{}
	for _, g := range games {
		teams := []models.GameTeam{}
		for _, gt := range g.Teams {
			teams = append(teams, models.GameTeam{
				ID: gt.ID,
			})
		}

		game := models.Game{
			ID:         g.ID,
			ExternalID: g.ExternalID,
			Teams:      teams,
		}

		gamesData = append(gamesData, game)
	}

	return gamesData, nil
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

	log.Printf("[core-frame] loaded frames for %d games", len(games))

	return nil
}

func (ec *FrameCore) Save(game models.Game, frame models.Frame) error {
	ctx := context.Background()

	teams := []repo.FrameTeam{}
	for _, t := range frame.Teams {
		teams = append(teams, repo.FrameTeam{
			TeamID:     t.TeamID,
			Gold:       t.Gold,
			Towers:     t.Towers,
			Inhibitors: t.Inhibitors,
			Dragons:    t.Dragons,
			Barons:     t.Barons,
			TotalKills: t.TotalKills,
		})
	}
	players := []repo.FramePlayer{}
	for _, p := range frame.Players {
		players = append(players, repo.FramePlayer{
			PlayerID:            p.PlayerID,
			ExternalID:          p.ExternalID,
			Level:               p.Level,
			Kills:               p.Kills,
			Deaths:              p.Deaths,
			Assists:             p.Assists,
			TotalGoldEarned:     p.TotalGoldEarned,
			CreepScore:          p.CreepScore,
			KillParticipation:   p.KillParticipation,
			ChampionDamageShare: p.ChampionDamageShare,
			WardsPlaced:         p.WardsPlaced,
			WardsDestroyed:      p.WardsDestroyed,
			AttackDamage:        p.AttackDamage,
			AbilityPower:        p.AbilityPower,
			CriticalChance:      p.CriticalChance,
			AttackSpeed:         p.AttackSpeed,
			LifeSteal:           p.LifeSteal,
			Armor:               p.Armor,
			MagicResistance:     p.MagicResistance,
			Tenacity:            p.Tenacity,
			Items:               p.Items,
			Runes: repo.Runes{
				Main:      p.Runes.Main,
				Secondary: p.Runes.Secondary,
				Perks:     p.Runes.Perks,
			},
			Abilities: p.Abilities,
		})
	}

	frameData := repo.Frame{
		GameID:    game.ID,
		TimeStamp: frame.TimeStamp,
		Teams:     teams,
		Players:   players,
		Time:      frame.Time,
	}
	err := ec.db.SaveFrame(ctx, frameData)
	if err != nil {
		log.Printf("[core-frame] Error saving frame for game %s: %v", frame.GameID.Hex(), err)
	}

	log.Printf("[core-frame] Saved frame at %v for game %s", frame.TimeStamp, game.ExternalID)
	return nil
}

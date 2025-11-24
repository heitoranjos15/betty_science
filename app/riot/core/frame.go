package core

import (
	"context"
	"log"

	models "betty/science/app/riot"
	"betty/science/app/riot/clients/league"
	"betty/science/app/riot/repo"

	"go.mongodb.org/mongo-driver/bson"
)

type frameClient interface {
	LoadData(game models.Game) (league.FrameResponse, error)
}

type FrameChannel struct {
	frame models.Frame
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
		// "external_id": "113475871523985236",
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
				ID:   gt.ID,
				Side: gt.Side,
			})
		}

		game := models.Game{
			ID:         g.ID,
			ExternalID: g.ExternalID,
			Teams:      teams,
		}

		gamesData = append(gamesData, game)
	}

	log.Printf("[core-frame] fetched %d games to load frames for", len(gamesData))
	return gamesData, nil
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
		playerData, err := ec.playersDB.GetPlayerByExternalID(ctx, p.ExternalID)
		if err != nil {
			log.Printf("[core-frame] Error fetching player by external ID %s: %v", p.ExternalID, err)
			return err
		}
		players = append(players, repo.FramePlayer{
			PlayerID:            playerData.ID,
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
		TimeStamp: frame.Timestamp,
		Teams:     teams,
		Players:   players,
		Time:      frame.Time,
	}
	err := ec.db.SaveFrame(ctx, frameData)
	if err != nil {
		log.Printf("[core-frame] Error saving frame for game %s: %v", frame.GameID.Hex(), err)
	}

	return nil
}

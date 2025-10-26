package riot

import (
	"betty/science/app/league_of_legends/client"
	"betty/science/app/league_of_legends/models"
	"errors"
	"log"
	"time"
)

var ErrorGameFrameNoContent = errors.New("game frame has no content")

type clientFrame struct {
	api api
}

func NewFramesClient(api api) *clientFrame {
	return &clientFrame{
		api: api,
	}
}

func (c *clientFrame) LoadData(game models.Game) (client.FrameResponse, error) {
	var frames client.FrameResponse
	log.Println("[client-team-frame] Loading frames for game:", game.ExternalID)

	now := time.Now().UTC().Add(-2 * time.Minute)
	gameFrames, err := c.api.GetFrames(game.ExternalID, now)
	if err != nil {
		return frames, err
	}

	frame := gameFrames.Frames[len(gameFrames.Frames)-1]
	frames.GameStart, frames.GameEnd = c.findGameStartAndEnd(game, frame)

	playerFrames, err := c.api.GetPlayerFrames(game.ExternalID, now)
	if err != nil {
		return frames, err
	}

	timestamp, _ := time.Parse(time.RFC3339, frame.Rfc460Timestamp)

	playerFrame := playerFrames.Frames[len(playerFrames.Frames)-1]

	teams := c.frameTeams(game.Teams, frame)
	winner := teams[0]
	for _, t := range teams {
		if t.Gold > winner.Gold {
			winner = t
		}
	}
	frames.WinnerID = winner.TeamID

	frames.Frame = models.Frame{
		GameID:    game.ID,
		TimeStamp: timestamp,
		Teams:     teams,
		Players:   c.framePlayers(game.Teams, gameFrames.GameMetadata, playerFrame),
	}

	frames.Players = c.playersDetails(game.Teams, gameFrames.GameMetadata)

	return frames, nil
}

func (c *clientFrame) findGameStartAndEnd(game models.Game, frame GameFrame) (time.Time, time.Time) {
	gameEnd := time.Time{}
	gameStart := game.StartTime

	timestamp, _ := time.Parse(time.RFC3339, frame.Rfc460Timestamp)
	gameEnd = timestamp

	if game.StartTime.IsZero() {
		gameStart = c.gameStart(game.ExternalID, timestamp)
		return gameStart, gameEnd
	}

	return gameStart, gameEnd
}

func (c clientFrame) gameStart(gameExternalID string, endGameTime time.Time) time.Time {
	pastTime := endGameTime.Add(-20 * time.Minute) // assuming average game duration of 20 minutes
	startTime := pastTime

	for {
		log.Println("[client-team-frame] Fetching first frame for game:", gameExternalID, "at", pastTime)
		firstFrame, err := c.api.GetFrames(gameExternalID, pastTime)
		if err != nil {
			if errors.Is(err, ErrorGameFrameNoContent) {
				log.Println("[client-team-frame] Game start time found for game:", gameExternalID, "at", startTime)
				break
			}

			log.Println("[client-team-frame] Error fetching first frame for game:", gameExternalID, err)
			break
		}

		time.Sleep(1 * time.Second) // to avoid rate limiting

		pastTime = pastTime.Add(-1 * time.Minute)
		startTime, _ = time.Parse(time.RFC3339, firstFrame.Frames[0].Rfc460Timestamp)
	}

	return startTime
}

func (c clientFrame) framePlayers(teams []models.GameTeam, gameMetadata GameMetadata, playerFrame ParticipantFrame) []models.FramePlayer {
	var frame []models.FramePlayer

	for _, team := range teams {
		playerMeta := c.findParticipantMetadata(team.Side, gameMetadata)
		for _, pm := range playerMeta {
			playerID := pm.ParticipantID
			pf := playerFrame.Participants[playerID-1]
			frame = append(frame, models.FramePlayer{
				ExternalID:          pm.EsportsPlayerID,
				Level:               pf.Level,
				Kills:               pf.Kills,
				Deaths:              pf.Deaths,
				Assists:             pf.Assists,
				TotalGoldEarned:     pf.TotalGoldEarned,
				CreepScore:          pf.CreepScore,
				KillParticipation:   pf.KillParticipation,
				ChampionDamageShare: pf.ChampionDamageShare,
				WardsPlaced:         pf.WardsPlaced,
				WardsDestroyed:      pf.WardsDestroyed,
				AttackDamage:        pf.AttackDamage,
				AbilityPower:        pf.AbilityPower,
				CriticalChance:      pf.CriticalChance,
				AttackSpeed:         pf.AttackSpeed,
				LifeSteal:           pf.LifeSteal,
				Armor:               pf.Armor,
				MagicResistance:     pf.MagicResistance,
				Tenacity:            pf.Tenacity,
				Items:               pf.Items,
				Runes: models.Runes{
					Main:      pf.PerkMetadata.StyleID,
					Secondary: pf.PerkMetadata.SubStyleID,
					Perks:     pf.PerkMetadata.Perks,
				},
				Abilities: pf.Abilities,
			})
		}
	}
	return frame
}

func (c clientFrame) frameTeams(teams []models.GameTeam, gameMetadata GameFrame) []models.FrameTeam {
	var frame []models.FrameTeam

	for _, team := range teams {
		teamMeta := c.findTeamMetadata(team.Side, gameMetadata)
		frame = append(frame, models.FrameTeam{
			TeamID:     team.ID,
			Gold:       teamMeta.TotalGold,
			Towers:     teamMeta.Towers,
			Inhibitors: teamMeta.Inhibitors,
			Dragons:    teamMeta.Dragons,
			Barons:     teamMeta.Barons,
			TotalKills: teamMeta.TotalKills,
		})
	}

	return frame
}

func (c clientFrame) playersDetails(team []models.GameTeam, gameMetadata GameMetadata) []client.PlayerFrame {
	var players []client.PlayerFrame

	for _, team := range team {
		teamMeta := c.findParticipantMetadata(team.Side, gameMetadata)
		for _, player := range teamMeta {
			players = append(players, client.PlayerFrame{
				ExternalID: player.EsportsPlayerID,
				Role:       player.Role,
				Name:       player.SummonerName,
				Side:       team.Side,
				Champion:   player.ChampionID,
				TeamID:     team.ID,
			})
		}
	}

	return players
}

func (c clientFrame) findTeamMetadata(side string, frame GameFrame) TeamFrame {
	if side == "blue" {
		return frame.BlueTeam
	}
	return frame.RedTeam
}

func (c clientFrame) findParticipantMetadata(side string, metadata GameMetadata) []ParticipantMeta {
	if side == "blue" {
		return metadata.BlueTeamMetadata.ParticipantMetadata
	}
	return metadata.RedTeamMetadata.ParticipantMetadata
}

package riot

import (
	"errors"
	"log"
	"sort"
	"sync"
	"time"

	"betty/science/app/league_of_legends/models"
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

func (c *clientFrame) collectFrames(frames []GameFrame, gameExternalID string) ([]GameFrame, error) {
	endGameFrame := frames[len(frames)-1]

	endTime, err := time.Parse(time.RFC3339, endGameFrame.Rfc460Timestamp)
	if err != nil {
		return frames, err
	}

	timeCounter := 0
	lastFrame := false
	for !lastFrame {
		type chanResult struct {
			LastFrame bool
			Frames    []GameFrame
		}
		channel := make(chan chanResult, 3)
		wg := sync.WaitGroup{}

		for range make([]int, 3) {
			wg.Add(1)
			timeCounter++
			go func(channel chan<- chanResult, timeCounter int, startTime time.Time) {
				defer wg.Done()
				frameTime := startTime.Add(time.Duration(-timeCounter) * time.Minute)
				log.Println("[client-team-frame] Fetching remaining frames for game:", gameExternalID, "at", frameTime)
				frames, err := c.api.GetFrames(gameExternalID, frameTime)
				if err != nil {
					if errors.Is(err, ErrorGameFrameNoContent) {
						log.Println("[client-team-frame] No more frames available for game:", gameExternalID)
						lastFrame = true
						return
					}
				}

				channel <- chanResult{
					LastFrame: false,
					Frames:    frames.Frames,
				}
			}(channel, timeCounter, endTime)
		}

		wg.Wait()
		close(channel)
		time.Sleep(1 * time.Second) // to avoid rate limiting

		for result := range channel {
			frames = append(frames, result.Frames...)
		}

		if lastFrame {
			log.Println("[client-team-frame] Completed fetching all frames for game:", gameExternalID)
			break
		}
	}

	sort.SliceStable(frames, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, frames[i].Rfc460Timestamp)
		timeJ, _ := time.Parse(time.RFC3339, frames[j].Rfc460Timestamp)
		return timeI.Before(timeJ)
	})
	return frames, nil
}

func (c *clientFrame) LoadData(game models.Game) (models.FrameResponse, error) {
	var frames models.FrameResponse
	log.Println("[client-team-frame] Loading frames for game:", game.ExternalID)

	now := time.Now().UTC().Add(-2 * time.Minute)
	gameFrames, err := c.api.GetFrames(game.ExternalID, now)
	if err != nil {
		log.Println("[client-team-frame] Error fetching game frames for game:", game.ExternalID, err)
		return frames, err
	}

	frame := gameFrames.Frames[len(gameFrames.Frames)-1]
	allFrames, err := c.collectFrames(gameFrames.Frames, game.ExternalID)
	if err != nil {
		return frames, err
	}

	gameFrames.Frames = allFrames

	firstFrameTime, _ := time.Parse(time.RFC3339, gameFrames.Frames[0].Rfc460Timestamp)
	lastFrameTime, _ := time.Parse(time.RFC3339, gameFrames.Frames[len(gameFrames.Frames)-1].Rfc460Timestamp)
	frames.GameStart = firstFrameTime
	frames.GameEnd = lastFrameTime

	playerFrames, err := c.api.GetPlayerFrames(game.ExternalID, lastFrameTime)
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

	frames.Frames = []models.Frame{
		{
			GameID:    game.ID,
			TimeStamp: timestamp,
			Teams:     teams,
			Players:   c.framePlayers(game.Teams, gameFrames.GameMetadata, playerFrame),
		},
	}

	frames.Players = c.playersDetails(game.Teams, gameFrames.GameMetadata)

	return frames, nil
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

func (c clientFrame) playersDetails(team []models.GameTeam, gameMetadata GameMetadata) []models.PlayerGameInfo {
	var players []models.PlayerGameInfo

	for _, team := range team {
		teamMeta := c.findParticipantMetadata(team.Side, gameMetadata)
		for _, player := range teamMeta {
			players = append(players, models.PlayerGameInfo{
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

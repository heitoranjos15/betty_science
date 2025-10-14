package riot

import (
	"betty/science/app/league_of_legends/client"
	"betty/science/app/league_of_legends/models"
	"time"
)

type clientTeamFrame struct {
	api api
}

func NewTeamFramesClient(api api) *clientTeamFrame {
	return &clientTeamFrame{
		api: api,
	}
}

func (c *clientTeamFrame) LoadData(game models.Game) (client.FrameResponse, error) {
	var frames client.FrameResponse
	now := time.Now()
	gameFrames, err := c.api.GetFrames(game.ExternalID, now)
	if err != nil {
		return frames, err
	}

	frame := gameFrames.Frames[0]
	startTime, _ := time.Parse(time.RFC3339, frame.Rfc460Timestamp)

	frames.Frame = models.Frame{
		GameID:    game.ID,
		TimeStamp: startTime,
		Teams:     c.frameTeams(game.Teams, frame),
	}

	frames.PlayerGamesDetails = c.playersDetails(game.Teams, gameFrames.GameMetadata)

	return frames, nil
}

func (c clientTeamFrame) frameTeams(teams []models.GameTeam, gameMetadata GameFrame) []models.FrameTeam {
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

func (c clientTeamFrame) playersDetails(team []models.GameTeam, gameMetadata GameMetadata) []models.GamePlayer {
	var players []models.GamePlayer

	for _, team := range team {
		teamMeta := c.findParticipantMetadata(team.Side, gameMetadata)
		for _, player := range teamMeta {
			players = append(players, models.GamePlayer{
				Side:       team.Side,
				Champion:   player.ChampionID,
				Role:       player.Role,
				Name:       player.SummonerName,
				ExternalID: player.EsportsPlayerID,
				TeamID:     team.ID,
			})
		}
	}

	return players
}

func (c clientTeamFrame) findTeamMetadata(side string, frame GameFrame) TeamFrame {
	if side == "blue" {
		return frame.BlueTeam
	}
	return frame.RedTeam
}

func (c clientTeamFrame) findParticipantMetadata(side string, metadata GameMetadata) []ParticipantMeta {
	if side == "blue" {
		return metadata.BlueTeamMetadata.ParticipantMetadata
	}
	return metadata.RedTeamMetadata.ParticipantMetadata
}

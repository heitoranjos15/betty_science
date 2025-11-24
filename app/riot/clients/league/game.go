package league

import (
	models "betty/science/app/riot"
	"log"
)

type gameClient struct {
	api api
}

func NewClientGame(api api) *gameClient {
	return &gameClient{
		api: api,
	}
}

func (c *gameClient) LoadData(match models.Match) ([]models.Game, error) {
	var response []models.Game
	data, err := c.api.GetGameDetails(match.ExternalID)
	if err != nil {
		log.Println("[client-riot-game] Error fetching game details:", err)
		return response, err
	}

	for _, details := range data.Data.Event.Match.Games {
		if details.State == "unneeded" {
			continue
		}
		game := c.parseGameDetails(details, data.Data.Event.Match.Teams)
		game.MatchID = match.ID

		response = append(response, game)
	}

	return response, nil
}

func (c *gameClient) parseGameDetails(data GamesDetails, teamDetails []GameDetailsTeam) models.Game {
	game := models.Game{
		ExternalID: data.ID,
		Number:     data.Number,
		State:      data.State,
	}

	for _, details := range teamDetails {
		gameTeam := models.GameTeam{
			ExternalID: details.ID,
		}
		for _, side := range data.Teams {
			if side.ID == details.ID {
				gameTeam.Side = side.Side
			}
		}
		game.Teams = append(game.Teams, gameTeam)
	}

	return game
}

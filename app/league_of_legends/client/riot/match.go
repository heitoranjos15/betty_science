package riot

import (
	"betty/science/app/league_of_legends/models"
	"errors"
	"fmt"
	"log"
	"time"
)

type ClientMatch struct {
	api api
}

func NewClientMatch(api api) *ClientMatch {
	return &ClientMatch{
		api: api,
	}
}

func (c *ClientMatch) Load() (models.MatchResponse, error) {
	return c.LoadData(nil)
}

func (c *ClientMatch) LoadData(_ any) (models.MatchResponse, error) {
	var response models.MatchResponse
	data, err := c.api.GetSchedule()
	if err != nil {
		log.Println("[client-riot] Error fetching schedule:", err)
		return response, err
	}

	for _, event := range data.Data.Schedule.Events {
		error := c.validateEvent(event)
		if error != nil {
			log.Printf("[client-riot] Skipping event ID %s due to validation error: %v", event.Match.ID, error)
			continue
		}
		match, err := c.match(event)
		if err != nil {
			log.Printf("[client-riot] Skipping event ID %s due to match parsing error: %v", event.Match.ID, err)
			continue
		}
		response.Match = append(response.Match, match)

		teams, err := c.team(event)
		if err != nil {
			log.Printf("[client-riot] Skipping event ID %s due to team parsing error: %v", event.Match.ID, err)
			continue
		}
		response.TeamsDetails = append(response.TeamsDetails, teams...)
	}
	return response, nil
}

func (c ClientMatch) validateEvent(event Event) error {
	for _, team := range event.Match.Teams {
		if team.Name == "" || team.Name == "TBD" {
			return errors.New("invalid team name")
		}
	}

	return nil
}

func (c ClientMatch) match(event Event) (models.Match, error) {
	timeParsed, err := time.Parse(time.RFC3339, event.StartTime)
	if err != nil {
		return models.Match{}, errors.New("invalid start time format")
	}

	return models.Match{
		ExternalID: event.Match.ID,
		StartTime:  timeParsed,
		State:      event.State,
		BestOf:     event.Match.Strategy.Count,
		Format:     event.Type,
		League:     event.League.Name,
		LoadState:  "without_games",
	}, nil
}

func (c ClientMatch) team(event Event) ([]models.TournamentTeam, error) {
	teams := []models.TournamentTeam{}
	for _, team := range event.Match.Teams {
		teams = append(teams, models.TournamentTeam{
			TournamentName: fmt.Sprintf("%s %d", event.League.Name, time.Now().Year()),
			Team: models.Team{
				Name:     team.Name,
				ImageURL: team.Image,
				Code:     team.Code,
			},
		})
	}
	return teams, nil
}

// func (ep eventParser) matchResults(event iModels.Event) []models.MatchResult {
//   results := []models.MatchResult{}
//   for _, team := range event.Match.Teams {
//     if team.Result != nil {
//       results = append(results, models.MatchResult{
//         TeamID:   primitive.NilObjectID, // TODO: map team external ID to internal ID
//         GameWins: team.Result.GameWins,
//       })
//     }
//   }
//   return results
// }

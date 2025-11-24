package league

import (
	models "betty/science/app/riot"
	"errors"
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

func (c *ClientMatch) Load() ([]models.Match, error) {
	var matches []models.Match
	data, err := c.api.GetSchedule()
	if err != nil {
		log.Println("[client-riot] Error fetching schedule:", err)
		return matches, err
	}

	for _, event := range data.Data.Schedule.Events {
		match, err := c.parse(event)
		if err != nil {
			log.Printf("[client-riot] Skipping event ID %s due to match parsing error: %v", event.Match.ID, err)
			continue
		}
		matches = append(matches, match)
	}
	return matches, nil
}

func (c ClientMatch) parse(event Event) (models.Match, error) {
	if err := c.validateEvent(event); err != nil {
		return models.Match{}, err
	}

	timeParsed, err := time.Parse(time.RFC3339, event.StartTime)
	if err != nil {
		return models.Match{}, errors.New("invalid start time format")
	}

	matchTeams := []models.Team{}
	for _, team := range event.Match.Teams {
		matchTeams = append(matchTeams, models.Team{
			Name:     team.Name,
			Code:     team.Code,
			ImageURL: team.Image,
		})
	}

	return models.Match{
		ExternalID: event.Match.ID,
		StartTime:  timeParsed,
		State:      event.State,
		BestOf:     event.Match.Strategy.Count,
		Format:     event.Type,
		League:     event.League.Name,
		Teams:      matchTeams,
	}, nil
}

func (c ClientMatch) validateEvent(event Event) error {
	for _, team := range event.Match.Teams {
		if team.Name == "" || team.Name == "TBD" {
			return errors.New("invalid team name")
		}
	}

	return nil
}

package core

import (
	"context"

	models "betty/science/app/riot"
	"betty/science/app/riot/repo"
)

type MatchCore struct {
	db matchDB
}

func NewMatchCore(db matchDB) *MatchCore {
	return &MatchCore{
		db: db,
	}
}

func (ec *MatchCore) Save(matches []models.Match) error {
	ctx := context.Background()

	saveData := []repo.Match{}
	for _, m := range matches {
		data := repo.Match{
			StartTime:  m.StartTime,
			State:      m.State,
			BestOf:     m.BestOf,
			Format:     m.Format,
			League:     m.League,
			ExternalID: m.ExternalID,
			Tournament: m.Tournament,
			TeamsID:    m.TeamsID,
			LoadState:  "without_games",
		}
		saveData = append(saveData, data)
	}

	return ec.db.SaveBulkMatches(ctx, saveData)
}

func (ec *MatchCore) LoadBulk() ([]models.Game, error) {
	ctx := context.Background()
	matches := []models.Game{}

	dbResult, err := ec.db.GetMatches(ctx, map[string]any{"load_state": "without_games"})
	if err != nil {
		return matches, err
	}

	for _, result := range dbResult {
		data := models.Game{
			ExternalID: result.ExternalID,
		}
		matches = append(matches, data)
	}

	return matches, nil
}

package worker

import (
	channels "betty/science/app"
	models "betty/science/app/riot"
	"log"
	"sync"
)

type clientMatch interface {
	Load() ([]models.Match, error)
}

type coreTeam interface {
	SetTeamID(team *models.Team) error
}

type coreMatch interface {
	LoadBulk() ([]models.Game, error)
	Save(matches []models.Match) error
}

type WorkerMatch struct {
	client   clientMatch
	core     coreMatch
	coreTeam coreTeam
}

func NewWorkerMatch(
	client clientMatch,
	core coreMatch,
	coreTeam coreTeam,
) *WorkerMatch { // TODO: use of pointer?
	return &WorkerMatch{
		client:   client,
		core:     core,
		coreTeam: coreTeam,
	}
}

func (w *WorkerMatch) LoadBulk() ([]any, error) {
	// No bulk loading needed for this worker
	return []any{}, nil
}

func (w *WorkerMatch) Run(data []any, workerName string, delay int, botChan chan<- channels.BotResponse, wg *sync.WaitGroup) { // TODO: group parameter about the worker
	defer wg.Done()

	resp, err := w.client.Load()
	if err != nil {
		botChan <- channels.BotResponse{
			Error: err,
		}
		return
	}

	for i := range resp {
		for _, team := range resp[i].Teams {
			err := w.coreTeam.SetTeamID(&team)
			if err != nil {
				log.Printf("[worker-match] Error setting team ID for team %s: %v", team.Name, err)
				botChan <- channels.BotResponse{
					Error: err,
				}
				return
			}
			resp[i].TeamsID = append(resp[i].TeamsID, team.ID)
		}
	}

	err = w.core.Save(resp)

	if err != nil {
		botChan <- channels.BotResponse{
			Error: err,
		}
		return
	}

	botChan <- channels.BotResponse{
		Error:          err,
		TotalProcessed: len(resp),
	}

	log.Println("WorkerMatch: Match loading process completed")
}

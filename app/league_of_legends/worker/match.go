package worker

import (
	channels "betty/science/app"
	"betty/science/app/league_of_legends/models"
	"log"
	"sync"
)

type clientMatch interface {
	Load() (models.MatchResponse, error)
}

type coreMatch interface {
	LoadBulk() ([]models.Game, error)
	Save(matches []models.Match) error
}

type coreTeam interface {
	UpdateTeamsTournament(teams []models.TournamentTeam) error
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
) *WorkerMatch {
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

func (w *WorkerMatch) Run(data []any, workerName string, delay int, botChan chan<- channels.BotResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := w.client.Load()
	if err != nil {
		botChan <- channels.BotResponse{
			Error: err,
		}
		return
	}

	err = w.core.Save(resp.Match)

	if err != nil {
		botChan <- channels.BotResponse{
			Error: err,
		}
		return
	}

	err = w.coreTeam.UpdateTeamsTournament(resp.TeamsDetails)

	botChan <- channels.BotResponse{
		Error:          err,
		TotalProcessed: len(resp.Match),
	}

	log.Println("WorkerMatch: Match loading process completed")
}

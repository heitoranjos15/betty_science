package worker

import (
	channels "betty/science/app"
	"betty/science/app/league_of_legends/models"
	"log"
	"sync"
	"time"
)

type clientGame interface {
	LoadData(models.Match) ([]models.Game, error)
}

type coreGame interface {
	SaveBulk(games []models.Game) error
	UpdateGameByFrameResp(game models.Game, resp models.FrameResponse) error
	LoadBulk() ([]models.Match, error)
}

type WorkerGame struct {
	client    clientGame
	core      coreGame
	coreMatch coreMatch
}

func NewWorkerGame(
	client clientGame,
	core coreGame,
	coreMatch coreMatch,
) *WorkerGame {
	return &WorkerGame{
		client:    client,
		core:      core,
		coreMatch: coreMatch,
	}
}

func (w *WorkerGame) LoadBulk() ([]models.Match, error) {
	return w.core.LoadBulk()
}

func (w *WorkerGame) Run(data []models.Match, workerName string, delay int, botChan chan<- channels.BotResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("[%s] [worker-game] Starting processing %d matches", workerName, len(data))

	matchProcessed := []models.Match{}
	for _, match := range data {
		resp, err := w.client.LoadData(match)
		if err != nil {
			botChan <- channels.BotResponse{
				Error: err,
			}
			continue
		}
		log.Printf("[%s] [worker game] Match %s games loaded: %d", workerName, match.ExternalID, len(resp))

		err = w.core.SaveBulk(resp)
		if err != nil {
			botChan <- channels.BotResponse{
				Error: err,
			}
			continue
		}

		match.LoadState = "loaded"
		matchProcessed = append(matchProcessed, match)
		time.Sleep(time.Duration(delay) * time.Second)
	}

	if len(matchProcessed) > 0 {
		err := w.coreMatch.Save(matchProcessed)
		if err != nil {
			log.Printf("[%s] [worker game] Error updating matches after games loaded: %v", workerName, err)
		}
	}

	botChan <- channels.BotResponse{
		TotalProcessed: len(matchProcessed),
		Error:          nil,
	}
}

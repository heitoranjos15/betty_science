package worker

import (
	channels "betty/science/app"
	models "betty/science/app/riot"
	"betty/science/app/riot/clients/league"
	"log"
	"sync"
	"time"
)

type clientFrame interface {
	LoadData(models.Game) (league.FrameResponse, error)
}

type coreFrame interface {
	Save(game models.Game, frames models.Frame) error
	LoadBulk() ([]models.Game, error)
}

type WorkerFrame struct {
	client   clientFrame
	core     coreFrame
	coreGame coreGame
}

func NewWorkerFrame(
	client clientFrame,
	core coreFrame,
	coreGame coreGame,
) *WorkerFrame {
	return &WorkerFrame{
		client:   client,
		core:     core,
		coreGame: coreGame,
	}
}

func (w *WorkerFrame) LoadBulk() ([]models.Game, error) {
	return w.core.LoadBulk()
}

func (w *WorkerFrame) Run(data []models.Game, workerName string, delay int, botChan chan<- channels.BotResponse, wg *sync.WaitGroup) {
	defer wg.Done()

	gamesProcessed := []models.Game{}
	for _, game := range data {
		resp, err := w.client.LoadData(game)
		if err != nil {
			botChan <- channels.BotResponse{
				Error: err,
			}
			continue
		}

		for _, frame := range resp.Frames {
			err = w.core.Save(game, frame)
			if err != nil {
				botChan <- channels.BotResponse{
					Error: err,
				}
				continue
			}
		}

		err = w.coreGame.UpdateGameByFrameResp(game, resp)
		if err != nil {
			log.Printf("[worker-frame] Error updating game by frame response: %v", err)
		}
		log.Printf("[%s] [worker-frame] Processed game %s with %d frames", workerName, game.ExternalID, len(resp.Frames))

		time.Sleep(time.Duration(delay) * time.Second)
	}

	botChan <- channels.BotResponse{
		TotalProcessed: len(gamesProcessed),
	}
}

package bot

import (
	channels "betty/science/app"
	"betty/science/config"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
)

type botConfig struct {
	workersTotal int
	workers      []*worker
}

type worker struct {
	Name         string
	Emote        string
	DelaySeconds int
}

func SetupBot(cfg config.Config) botConfig {
	if cfg.Workers <= 0 {
		panic("Workers on config must be greater than 0")
	}

	workers := []*worker{
		{Name: "Woody", Emote: "🤠", DelaySeconds: rand.IntN(6) + 2},
		{Name: "Sasuke", Emote: ":|", DelaySeconds: rand.IntN(6) + 2},
		{Name: "Gash", Emote: ":D", DelaySeconds: rand.IntN(6) + 2},
	}

	if cfg.Workers > len(workers) {
		for i := len(workers) - 1; i >= cfg.Workers; i++ {
			workers = append(workers, &worker{Name: fmt.Sprintf("Default bot %d", i), Emote: ":)", DelaySeconds: rand.IntN(6) + 2})
		}
	}
	return botConfig{
		workersTotal: cfg.Workers,
		workers:      workers,
	}
}

type WorkerInterface[T any] interface {
	LoadBulk() ([]T, error)
	Run(data []T, botName string, delay int, botChan chan<- channels.BotResponse, wg *sync.WaitGroup)
}

type Bot[T any] struct {
	worker WorkerInterface[T]
	config botConfig
}

func NewBot[T any](worker WorkerInterface[T], cfg config.Config) *Bot[T] {
	cfgBot := SetupBot(cfg)
	return &Bot[T]{
		worker: worker,
		config: cfgBot,
	}
}

func (b Bot[T]) Execute() {
	wg := sync.WaitGroup{}

	worker := b.worker
	config := b.config

	data, err := b.worker.LoadBulk()
	log.Printf("[bot] Loaded %d items to process", len(data))
	if err != nil {
		log.Printf("[bot] Error loading data: %v", err)
		return
	}

	botChan := make(chan channels.BotResponse, config.workersTotal)
	for i := 0; i < config.workersTotal; i++ {
		name := config.workers[i].Name
		emote := config.workers[i].Emote
		delay := config.workers[i].DelaySeconds

		if i >= len(data) {
			log.Printf("[%s] not enough data %s", name, emote)
			break
		}

		wg.Add(1)
		log.Printf("[%s] starting worker %s", name, emote)

		totalData := len(data)
		if totalData >= config.workersTotal {
			totalData = totalData / config.workersTotal
		}
		dataChunk := data[i*totalData : (i+1)*totalData]
		go worker.Run(dataChunk, name, delay, botChan, &wg)
	}

	log.Println("[bot] Workers started")

	totalProcessed := 0
	go func() {
		wg.Wait()
		close(botChan)
	}()

	for resp := range botChan {
		if resp.Error != nil {
			log.Printf("[bot] Error processing: %v", resp.Error)
			continue
		}
		totalProcessed += resp.TotalProcessed
	}
	log.Printf("[bot] All workers finished. Total processed: %d", totalProcessed)
}

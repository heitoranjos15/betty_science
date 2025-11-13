package config

import "math/rand/v2"

type BotConfig struct {
	Workers int
	Bots    []*Bot
}

type Bot struct {
	Name         string
	Emote        string
	DelaySeconds int
}

func SetupBot() BotConfig {
	cfg := LoadConfig()

	if cfg.Workers <= 0 {
		panic("Workers on config must be greater than 0")
	}

	bots := []*Bot{
		{Name: "Woody", Emote: "🤠", DelaySeconds: rand.IntN(5)},
		{Name: "Sasuke", Emote: ":|", DelaySeconds: rand.IntN(5)},
		{Name: "Gash", Emote: ":D", DelaySeconds: rand.IntN(5)},
	}

	if cfg.Workers > len(bots) {
		for i := len(bots) - 1; i >= cfg.Workers; i++ {
			bots = append(bots, &Bot{Name: "Default Bot" + string(i), Emote: ":)", DelaySeconds: rand.IntN(5)})
		}
	}
	return BotConfig{
		Workers: cfg.Workers,
		Bots:    bots,
	}
}

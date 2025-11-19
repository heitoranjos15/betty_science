package channels

type BotResponse struct {
	Error          error
	TotalProcessed int
}

type BotConfig struct {
	Name         string
	Emote        string
	DelaySeconds int
}

type BotInputChan <-chan BotConfig

type BotChan chan BotResponse

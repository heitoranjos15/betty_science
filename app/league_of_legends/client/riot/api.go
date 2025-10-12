package riot

type api interface {
	GetSchedule() (ScheduleResponse, error)
	GetGameDetails(string) (GameDetailsResponse, error)
}

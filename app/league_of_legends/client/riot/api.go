package riot

import (
	"time"
)

type api interface {
	GetSchedule() (ScheduleResponse, error)
	GetGameDetails(string) (GameDetailsResponse, error)
	GetFrames(string, time.Time) (FramesResponse, error)
	GetPlayerFrames(string, time.Time) (PlayerFramesResponse, error)
	GetFirstFrame(string, time.Time) (GameFrame, error)
}

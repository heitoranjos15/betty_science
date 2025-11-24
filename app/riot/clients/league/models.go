package league

import (
	"time"

	models "betty/science/app/riot"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FrameResponse struct {
	Frames    []models.Frame
	Players   []models.GamePlayer
	GameStart time.Time
	GameEnd   time.Time
	WinnerID  primitive.ObjectID
}

package client

import (
	"betty/science/app/league_of_legends/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamDetails struct {
	TournamentName string
	Team           models.Team
}

type MatchResponse struct {
	Match        []models.Match
	TeamsDetails []TeamDetails
}

type PlayerFrame struct {
	ExternalID string
	Name       string
	Team       string
	Role       string
	Side       string
	Champion   string
	TeamID     primitive.ObjectID
}

type FrameResponse struct {
	Frame     models.Frame
	Players   []PlayerFrame
	GameStart time.Time
	GameEnd   time.Time
	WinnerID  primitive.ObjectID
}

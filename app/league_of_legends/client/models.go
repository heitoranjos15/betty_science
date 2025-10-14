package client

import (
	"betty/science/app/league_of_legends/models"
)

type TeamTournament struct {
	TournamentName string
	Team           models.Team
}

type MatchResponse struct {
	Match            []models.Match
	TeamsTournaments []TeamTournament
}

type playerFrame struct {
	Name       string
	Team       string
	ExternalID string
	Role       string
}

type FrameResponse struct {
	Frame              models.Frame
	PlayerGamesDetails []models.GamePlayer
}

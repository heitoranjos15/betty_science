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

type GameResponse struct {
	Teams []models.Team
	Games []models.Game
}

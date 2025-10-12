package core

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"betty/science/app/league_of_legends/client"
	"betty/science/app/league_of_legends/models"
)

type clientAPI interface {
	LoadData(_ any) (client.MatchResponse, error)
}

type matchDB interface {
	SaveBulkMatches(context.Context, []models.Match) error
}

type teamDB interface {
	UpdateTeamTournaments(context.Context, primitive.ObjectID, string) error
	GetTeamByName(context.Context, string) (models.Team, error)
	SaveTeamByName(context.Context, models.Team) error
}

type gameDB interface {
	SaveBulkGames(context.Context, []models.Game) error
}

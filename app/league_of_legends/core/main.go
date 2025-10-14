package core

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"betty/science/app/league_of_legends/client"
	"betty/science/app/league_of_legends/models"
)

type clientAPI interface {
	LoadData(_ any) (client.MatchResponse, error)
}

type matchDB interface {
	SaveBulkMatches(context.Context, []models.Match) error
	GetMatches(context.Context, bson.M) ([]models.Match, error)
}

type teamDB interface {
	UpdateTeamTournaments(context.Context, primitive.ObjectID, string) error
	GetTeamByName(context.Context, string) (models.Team, error)
	SaveTeamByName(context.Context, models.Team) error
	UpdateTeamExternalID(context.Context, primitive.ObjectID, string) error
}

type gameDB interface {
	SaveBulkGames(context.Context, []models.Game) error
	GetGames(context.Context, bson.M) ([]models.Game, error)
}

type playersDB interface {
	GetPlayerByExternalID(ctx context.Context, externalID string) (models.Player, error)
	SaveBulkPlayers(ctx context.Context, players []models.Player) error
}

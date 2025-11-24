package core

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"betty/science/app/riot/repo"
)

type matchDB interface {
	SaveBulkMatches(context.Context, []repo.Match) error
	SaveMatch(context.Context, repo.Match) error
	GetMatches(context.Context, bson.M) ([]repo.Match, error)
}

type teamDB interface {
	UpdateTeamTournaments(context.Context, primitive.ObjectID, string) error
	GetTeamByName(context.Context, string) (repo.Team, error)
	GetTeamByExternalID(context.Context, string) (repo.Team, error)
	SaveTeamByName(context.Context, repo.Team) error
	UpdateTeamExternalID(context.Context, primitive.ObjectID, string) error
}

type gameDB interface {
	SaveBulkGames(context.Context, []repo.Game) error
	SaveGame(context.Context, repo.Game) error
	GetGames(context.Context, bson.M) ([]repo.Game, error)
	UpdateGameByExternalID(context.Context, string, bson.M) error
}

type playersDB interface {
	GetPlayerByExternalID(ctx context.Context, externalID string) (repo.Player, error)
	SaveBulkPlayers(ctx context.Context, players []repo.Player) error
}

type frameDB interface {
	SaveFrame(context.Context, repo.Frame) error
}

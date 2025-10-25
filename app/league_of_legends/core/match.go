package core

import (
	"context"
	"log"
	"slices"

	"betty/science/app/league_of_legends/client"
	"betty/science/app/league_of_legends/models"
)

type matchClient interface {
	LoadData(_ any) (client.MatchResponse, error)
}

type MatchCore struct {
	client matchClient
	db     matchDB
	teamDB teamDB
}

func NewMatchCore(client matchClient, db matchDB, teamDB teamDB) *MatchCore {
	return &MatchCore{
		client: client,
		db:     db,
		teamDB: teamDB,
	}
}

func (ec *MatchCore) Load() error {
	ctx := context.Background()
	data, err := ec.client.LoadData(nil)
	if err != nil {
		return err
	}

	if err := ec.db.SaveBulkMatches(ctx, data.Match); err != nil {
		return err
	}

	for _, tt := range data.TeamsDetails {
		if err := ec.saveTeamTournament(ctx, tt.Team, tt.TournamentName); err != nil {
			log.Println("[core-match] couldnt save team tournament", err)
		}
	}

	log.Printf("[core-match] loaded %d matches", len(data.Match))
	log.Printf("[core-match] loaded %d team-tournament associations", len(data.TeamsDetails))

	return nil
}

func (ec *MatchCore) saveTeamTournament(ctx context.Context, team models.Team, tournament string) error {
	existingTeam, err := ec.teamDB.GetTeamByName(ctx, team.Name)

	if err != nil {
		log.Println("[core-match] team not found, creating new one:", team.Name)
		team.Tournaments = []string{tournament}
		err = ec.teamDB.SaveTeamByName(ctx, team)
		if err != nil {
			return err
		}
	}

	if slices.Contains(existingTeam.Tournaments, tournament) {
		return nil
	}

	return ec.teamDB.UpdateTeamTournaments(ctx, existingTeam.ID, tournament)
}

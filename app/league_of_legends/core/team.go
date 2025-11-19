package core

import (
	"betty/science/app/league_of_legends/models"
	"betty/science/app/league_of_legends/repo"
	"context"
	"slices"
)

type TeamCore struct {
	db teamDB
}

func NewTeamCore(db teamDB) *TeamCore {
	return &TeamCore{
		db: db,
	}
}

func (ec *TeamCore) UpdateTeamsTournament(teams []models.TournamentTeam) error {
	ctx := context.Background()

	for _, t := range teams {
		team := t.Team
		existingTeam, err := ec.db.GetTeamByName(ctx, team.Name)

		if err != nil {
			teamData := repo.Team{
				Name:        team.Name,
				Code:        team.Code,
				ImageURL:    team.ImageURL,
				Tournaments: []string{t.TournamentName},
			}
			if err := ec.db.SaveTeamByName(ctx, teamData); err != nil {
				return err
			}
		}

		if !slices.Contains(existingTeam.Tournaments, t.TournamentName) {
			return ec.db.UpdateTeamTournaments(ctx, existingTeam.ID, t.TournamentName)
		}
	}

	return nil
}

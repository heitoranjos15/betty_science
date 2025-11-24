package core

import (
	models "betty/science/app/riot"
	"betty/science/app/riot/repo"
	"context"
)

type TeamCore struct {
	db teamDB
}

func NewTeamCore(db teamDB) *TeamCore {
	return &TeamCore{
		db: db,
	}
}

func (ec *TeamCore) SetTeamID(team *models.Team) error {
	ctx := context.Background()
	savedTeam, err := ec.db.GetTeamByName(ctx, team.Name)
	if err != nil {
		insertData := repo.Team{
			Name:     team.Name,
			Code:     team.Code,
			ImageURL: team.ImageURL,
		}
		if err := ec.db.SaveTeamByName(ctx, insertData); err != nil {
			return err
		}

		savedTeam, err = ec.db.GetTeamByName(ctx, team.Name)
		if err != nil {
			return err
		}
	}

	team.ID = savedTeam.ID
	return nil
}

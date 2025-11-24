package repo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Match struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	StartTime  time.Time            `bson:"start_time"`
	State      string               `bson:"state"`
	BestOf     int                  `bson:"best_of"`
	Format     string               `bson:"format"`
	League     string               `bson:"league"`
	Tournament string               `bson:"tournament"`
	ExternalID string               `bson:"external_id"`
	LoadState  string               `bson:"load_state"`
	TeamsID    []primitive.ObjectID `bson:"teams_id"`
}

type Team struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ExternalID  string             `bson:"external_id"`
	Name        string             `bson:"name"`
	Code        string             `bson:"code"`
	ImageURL    string             `bson:"image"`
	Tournaments []string           `bson:"tournaments"`
}

type Player struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	ExternalID string               `bson:"external_id"`
	Name       string               `bson:"name"`
	ActualTeam primitive.ObjectID   `bson:"actual_team"`
	ActualRole string               `bson:"actual_role"`
	Roles      []string             `bson:"roles"`
	Teams      []primitive.ObjectID `bson:"teams"`
	ImageURL   string               `bson:"image"`
}

type Game struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	MatchID      primitive.ObjectID `bson:"match_id"`
	ExternalID   string             `bson:"external_id"`
	Number       int                `bson:"game_number"`
	WinnerTeamID primitive.ObjectID `bson:"winner_team_id,omitempty"`
	Duration     float64            `bson:"duration,omitempty"`
	ScheduleTime time.Time          `bson:"schedule_time,omitempty"`
	StartTime    time.Time          `bson:"start_time,omitempty"`
	Teams        []GameTeam         `bson:"teams,omitempty"`
	Players      []GamePlayer       `bson:"players,omitempty"`
	State        string             `bson:"state"`
	LoadState    string             `bson:"load_state"`
	ErrorMsg     string             `bson:"error_message,omitempty"`
}

type GameTeam struct {
	ID   primitive.ObjectID `bson:"team_id"`
	Side string             `bson:"side"`
}

type GamePlayer struct {
	PlayerID primitive.ObjectID `bson:"player_id"`
	Side     string             `bson:"side"`
	Champion string             `bson:"champion"`
	Role     string             `bson:"role"`
}

type Frame struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	GameID    primitive.ObjectID `bson:"game_id"`
	Teams     []FrameTeam        `bson:"teams"`
	Players   []FramePlayer      `bson:"players"`
	Time      int                `bson:"time"`
	TimeStamp time.Time          `bson:"timestamp"`
}

type FrameTeam struct {
	TeamID     primitive.ObjectID `bson:"team_id"`
	Gold       int                `bson:"gold"`
	Towers     int                `bson:"towers"`
	Dragons    []string           `bson:"dragons"`
	Barons     int                `bson:"barons"`
	Inhibitors int                `bson:"inhibitors"`
	TotalKills int                `bson:"total_kills"`
}

type FramePlayer struct {
	PlayerID            primitive.ObjectID `bson:"player_id"`
	ExternalID          string             `bson:"external_id"`
	Level               int                `bson:"level"`
	Kills               int                `bson:"kills"`
	Deaths              int                `bson:"deaths"`
	Assists             int                `bson:"assists"`
	TotalGoldEarned     int                `bson:"total_gold_earned"`
	CreepScore          int                `bson:"creep_score"`
	KillParticipation   float64            `bson:"kill_participation"`
	ChampionDamageShare float64            `bson:"champion_damage_share"`
	WardsPlaced         int                `bson:"wards_placed"`
	WardsDestroyed      int                `bson:"wards_destroyed"`
	AttackDamage        int                `bson:"attack_damage"`
	AbilityPower        int                `bson:"ability_power"`
	CriticalChance      float64            `bson:"critical_chance"`
	AttackSpeed         int                `bson:"attack_speed"`
	LifeSteal           int                `bson:"life_steal"`
	Armor               int                `bson:"armor"`
	MagicResistance     int                `bson:"magic_resistance"`
	Tenacity            float64            `bson:"tenacity"`
	Items               []int              `bson:"items"`
	Runes               Runes              `bson:"runes"`
	Abilities           []string           `bson:"abilities"`
}

type Runes struct {
	Main      int   `bson:"main"`
	Secondary int   `bson:"secondary"`
	Perks     []int `bson:"perks"`
}

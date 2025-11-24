package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Match struct {
	ID         primitive.ObjectID
	ExternalID string
	StartTime  time.Time
	State      string
	BestOf     int
	Format     string
	League     string
	Tournament string
	Teams      []Team
	TeamsID    []primitive.ObjectID
}

type Team struct {
	ID         primitive.ObjectID
	ExternalID string
	Name       string
	Code       string
	ImageURL   string
}

type Player struct {
	ID         primitive.ObjectID
	ExternalID string
}

type PlayerDetails struct {
	Name     string
	Team     primitive.ObjectID
	Role     string
	ImageURL string
}

type Game struct {
	ID           primitive.ObjectID
	ExternalID   string
	MatchID      primitive.ObjectID
	Number       int
	WinnerTeamID primitive.ObjectID
	Duration     float64
	ScheduleTime time.Time
	StartTime    time.Time
	Teams        []GameTeam
	Players      []GamePlayer
	State        string
	ErrorMsg     string
}

type GameTeam struct {
	ID         primitive.ObjectID
	ExternalID string
	Side       string
}

type GamePlayer struct {
	TeamID     primitive.ObjectID
	Name       string
	ExternalID string
	Side       string
	Role       string
	Champion   string
}

type Frame struct {
	GameID    primitive.ObjectID
	Teams     []FrameTeam
	Players   []FramePlayer
	Time      int
	Timestamp time.Time
}

type FrameTeam struct {
	TeamID     primitive.ObjectID
	Gold       int
	Towers     int
	Dragons    []string
	Barons     int
	Inhibitors int
	TotalKills int
}

type FramePlayer struct {
	ExternalID          string
	Level               int
	Kills               int
	Deaths              int
	Assists             int
	TotalGoldEarned     int
	CreepScore          int
	KillParticipation   float64
	ChampionDamageShare float64
	WardsPlaced         int
	WardsDestroyed      int
	AttackDamage        int
	AbilityPower        int
	CriticalChance      float64
	AttackSpeed         int
	LifeSteal           int
	Armor               int
	MagicResistance     int
	Tenacity            float64
	Items               []int
	Runes               Runes
	Abilities           []string
}

type Runes struct {
	Main      int
	Secondary int
	Perks     []int
}

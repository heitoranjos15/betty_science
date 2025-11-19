package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TournamentTeam struct {
	TournamentName string
	Team           Team
}

type MatchResponse struct {
	Match        []Match
	TeamsDetails []TournamentTeam
}

type Match struct {
	ID         primitive.ObjectID
	StartTime  time.Time
	State      string
	BestOf     int
	Format     string
	League     string
	Tournament string
	ExternalID string
	LoadState  string
	TeamsID    []primitive.ObjectID
}

type Team struct {
	ID          primitive.ObjectID
	ExternalID  string
	Name        string
	Code        string
	ImageURL    string
	Tournaments []string
}

type Player struct {
	ID         primitive.ObjectID
	ExternalID string
	Name       string
	ActualTeam primitive.ObjectID
	ActualRole string
	Roles      []string
	Teams      []primitive.ObjectID
	ImageURL   string
}

type Game struct {
	ID           primitive.ObjectID
	MatchID      primitive.ObjectID
	ExternalID   string
	Number       int
	Winner       primitive.ObjectID
	Duration     float64
	ScheduleTime time.Time
	StartTime    time.Time
	Teams        []GameTeam
	Players      []GamePlayer
	State        string
	LoadState    string
	ErrorMsg     string
}

type GameTeam struct {
	ID   primitive.ObjectID
	Side string

	// fallback fields to teams collection
	ExternalID string
	Name       string

	// TODO: winner already on GameTeam, so on front end i dont need to check both
	// Winner bool
}

type GamePlayer struct {
	PlayerID primitive.ObjectID
	Side     string
	Champion string
	Role     string
}

type PlayerGameInfo struct {
	PlayerID   primitive.ObjectID
	ExternalID string
	Name       string
	Team       string
	Role       string
	Side       string
	Champion   string
	TeamID     primitive.ObjectID
}

type FrameResponse struct {
	Frames    []Frame
	Players   []PlayerGameInfo
	GameStart time.Time
	GameEnd   time.Time
	WinnerID  primitive.ObjectID
}

type Frame struct {
	ID        primitive.ObjectID
	GameID    primitive.ObjectID
	Teams     []FrameTeam
	Players   []FramePlayer
	Time      int
	TimeStamp time.Time
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
	PlayerID            primitive.ObjectID
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

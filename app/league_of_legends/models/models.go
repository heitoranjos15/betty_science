package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Match struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StartTime  time.Time          `bson:"start_time" json:"start_time"`
	State      string             `bson:"state" json:"state"`
	BestOf     int                `bson:"best_of" json:"best_of"`
	Format     string             `bson:"format" json:"format"`
	League     string             `bson:"league" json:"league"`
	ExternalID string             `bson:"external_id" json:"external_id"`
}

type Team struct {
	ID          primitive.ObjectID `bson:"id,omitempty" json:"id"`
	ExternalID  string             `bson:"external_id" json:"external_id"`
	Name        string             `bson:"name" json:"name"`
	Code        string             `bson:"code" json:"code"`
	ImageURL    string             `bson:"image" json:"image"`
	Tournaments []string           `bson:"tournaments" json:"tournaments"`
}

type Player struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	ExternalID string               `bson:"external_id" json:"external_id"`
	Name       string               `bson:"name" json:"name"`
	ActualTeam primitive.ObjectID   `bson:"actual_team" json:"actual_team"`
	ActualRole string               `bson:"actual_role" json:"actual_role"`
	Roles      []string             `bson:"roles" json:"roles"`
	Teams      []primitive.ObjectID `bson:"teams" json:"teams"`
	ImageURL   string               `bson:"image" json:"image"`
}

type Game struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MatchID    primitive.ObjectID `bson:"match_id" json:"match_id"`
	ExternalID string             `bson:"external_id" json:"external_id"`
	Number     int                `bson:"game_number" json:"game_number"`
	Winner     primitive.ObjectID `bson:"winner" json:"winner"`
	Frames     []Frame            `bson:"frames" json:"frames"`
	Duration   int                `bson:"duration" json:"duration"` // in seconds
	Teams      []GameTeam         `bson:"teams" json:"teams"`
	Players    []GamePlayer       `bson:"players" json:"players"`
	StartTime  time.Time          `bson:"start_time" json:"start_time"`
	State      string             `bson:"state" json:"state"` // unloaded, loaded, error
	ErrorMsg   string             `bson:"error_message,omitempty" json:"error_message,omitempty"`
}

type GameTeam struct {
	ID     primitive.ObjectID `bson:"team_id" json:"team_id"`
	Side   string             `bson:"side" json:"side"`
	Winner bool               `bson:"winner" json:"winner"`

	// fallback fields to teams collection
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
}

type GamePlayer struct {
	PlayerID primitive.ObjectID `bson:"player_id" json:"player_id"`
	Side     string             `bson:"side" json:"side"`
	Champion string             `bson:"champion" json:"champion"`
}

type Frame struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GameID    primitive.ObjectID `bson:"game_id" json:"game_id"`
	Teams     []FrameTeam        `bson:"teams" json:"teams"`
	Players   []FramePlayer      `bson:"players" json:"players"`
	Time      int                `bson:"time" json:"time"` // in seconds
	Timestamp time.Time          `bson:"timestamp" json:"rfc460Timestamp"`
}

type FrameTeam struct {
	TeamID     primitive.ObjectID `bson:"team_id" json:"id"`
	Gold       int                `bson:"gold" json:"totalGold"`
	Towers     int                `bson:"towers" json:"towers"`
	Dragons    int                `bson:"dragons" json:"dragons"`
	Barons     int                `bson:"barons" json:"barons"`
	Heralds    int                `bson:"heralds" json:"heralds"`
	Inhibitors int                `bson:"inhibitors" json:"inhibitors"`
	TotalKills int                `bson:"total_kills" json:"totalKills"`
}

type FramePlayer struct {
	PlayerID            primitive.ObjectID `bson:"player_id" json:"id"`
	Level               int                `bson:"level" json:"level"`
	Kills               int                `bson:"kills" json:"kills"`
	Deaths              int                `bson:"deaths" json:"deaths"`
	Assists             int                `bson:"assists" json:"assists"`
	TotalGoldEarned     int                `bson:"total_gold_earned" json:"totalGoldEarned"`
	CreepScore          int                `bson:"creep_score" json:"creepScore"`
	KillParticipation   float64            `bson:"kill_participation" json:"killParticipation"`
	ChampionDamageShare float64            `bson:"champion_damage_share" json:"championDamageShare"`
	WardsPlaced         int                `bson:"wards_placed" json:"wardsPlaced"`
	WardsDestroyed      int                `bson:"wards_destroyed" json:"wardsDestroyed"`
	AttackDamage        int                `bson:"attack_damage" json:"attackDamage"`
	AbilityPower        int                `bson:"ability_power" json:"abilityPower"`
	CriticalChance      float64            `bson:"critical_chance" json:"criticalChance"`
	AttackSpeed         int                `bson:"attack_speed" json:"attackSpeed"`
	LifeSteal           int                `bson:"life_steal" json:"lifeSteal"`
	Armor               int                `bson:"armor" json:"armor"`
	MagicResistance     int                `bson:"magic_resistance" json:"magicResistance"`
	Tenacity            float64            `bson:"tenacity" json:"tenacity"`
	Items               []int              `bson:"items" json:"items"`
	Runes               Runes              `bson:"runes" json:"runes"`
	Abilities           []string           `bson:"abilities" json:"abilities"`
}

type Runes struct {
	Main      int   `bson:"main" json:"main"`
	Secondary int   `bson:"secondary" json:"secondary"`
	Perks     []int `bson:"perks" json:"perks"`
}

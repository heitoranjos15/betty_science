package riot_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	modelsRiot "betty/science/app/league_of_legends/client"
	client "betty/science/app/league_of_legends/client/riot"
	"betty/science/app/league_of_legends/models"
)

var mockFrames = "../../test/mock_data/frames.json"
var mockPlayerFrames = "../../test/mock_data/player_frames.json"

func loadMockGamesData(path string) client.FramesResponse {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open mock data file: %v", err)
	}
	defer file.Close()

	var data client.FramesResponse
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Fatalf("Failed to decode mock data: %v", err)
	}
	return data
}

func loadMockPlayerFramesData(path string) client.PlayerFramesResponse {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open mock data file: %v", err)
	}
	defer file.Close()

	var data client.PlayerFramesResponse
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Fatalf("Failed to decode mock data: %v", err)
	}
	return data
}

func TestLoadFrames(t *testing.T) {
	mockValid := loadMockGamesData(mockFrames)
	mockPlayerValid := loadMockPlayerFramesData(mockPlayerFrames)
	gameID := primitive.NewObjectID()

	team1 := primitive.NewObjectID()
	team2 := primitive.NewObjectID()

	gameEnd, err := time.Parse(time.RFC3339, "2025-09-11T09:31:02.605Z")
	if err != nil {
		t.Fatalf("Failed to parse timestamp: %v", err)
	}

	gameScheduleTime, err := time.Parse(time.RFC3339, "2025-09-11T09:10:00.000Z")

	gameStart := gameEnd.Add(-20 * time.Minute)

	frame := models.Frame{
		GameID:    gameID,
		TimeStamp: gameEnd,
		Teams: []models.FrameTeam{
			{
				TeamID:     team1,
				Gold:       55253,
				Inhibitors: 1,
				Towers:     8,
				Barons:     1,
				TotalKills: 21,
				Dragons:    []string{"hextech", "hextech"},
			},
			{
				TeamID:     team2,
				Gold:       44482,
				Inhibitors: 0,
				Towers:     1,
				Barons:     0,
				TotalKills: 9,
				Dragons:    []string{"ocean", "cloud"},
			},
		},
		Players: []models.FramePlayer{
			{
				ExternalID:          "107492121209956865",
				Level:               14,
				Kills:               2,
				Deaths:              4,
				Assists:             1,
				TotalGoldEarned:     10559,
				CreepScore:          216,
				KillParticipation:   0.6,
				ChampionDamageShare: 0.18872078867778097,
				WardsPlaced:         15,
				WardsDestroyed:      5,
				AttackDamage:        276,
				AbilityPower:        0,
				CriticalChance:      0.0,
				AttackSpeed:         130,
				LifeSteal:           0,
				Armor:               124,
				MagicResistance:     116,
				Tenacity:            0.0,
				Items:               []int{6692, 3111, 3156, 3364, 2019, 1037, 2055},
				Abilities:           []string{"Q", "W", "E", "Q", "Q", "R", "Q", "E", "Q", "E", "R", "E", "E", "W"},
				Runes: models.Runes{
					Main:      8400,
					Secondary: 8300,
					Perks:     []int{8437, 8401, 8444, 8242, 8345, 8304, 5008, 5011},
				},
			},
			{
				ExternalID:          "98767975940969117",
				Level:               13,
				Kills:               0,
				Deaths:              6,
				Assists:             0,
				TotalGoldEarned:     8875,
				CreepScore:          187,
				KillParticipation:   0.0,
				ChampionDamageShare: 0.06831090768922408,
				WardsPlaced:         12,
				WardsDestroyed:      10,
				AttackDamage:        206,
				AbilityPower:        0,
				CriticalChance:      0.0,
				AttackSpeed:         132,
				LifeSteal:           0,
				Armor:               94,
				MagicResistance:     70,
				Tenacity:            0.0,
				Items:               []int{3364, 6610, 3111, 3071, 2055, 2055, 1036},
				Abilities:           []string{"Q", "W", "Q", "E", "Q", "R", "Q", "W", "Q", "W", "R", "W", "W"},
				Runes: models.Runes{
					Main:      8000,
					Secondary: 8300,
					Perks:     []int{8010, 9111, 9105, 8017, 8304, 8347, 5008, 5011},
				},
			},
			{
				ExternalID:          "98767975949926165",
				Level:               17,
				Kills:               0,
				Deaths:              2,
				Assists:             4,
				TotalGoldEarned:     11139,
				CreepScore:          288,
				KillParticipation:   0.8,
				ChampionDamageShare: 0.22390541849409396,
				WardsPlaced:         8,
				WardsDestroyed:      7,
				AttackDamage:        105,
				AbilityPower:        465,
				CriticalChance:      0.0,
				AttackSpeed:         143,
				LifeSteal:           0,
				Armor:               88,
				MagicResistance:     72,
				Tenacity:            0.0,
				Items:               []int{6657, 3364, 3111, 3089},
				Abilities:           []string{"E", "Q", "Q", "W", "Q", "R", "Q", "E", "Q", "E", "R", "E", "E", "W", "W", "R", "W"},
				Runes: models.Runes{
					Main:      8200,
					Secondary: 8400,
					Perks:     []int{8230, 8226, 8234, 8237, 8444, 8242, 5005, 5010, 5011},
				},
			},
			{
				ExternalID:          "102483272156027229",
				Level:               15,
				Kills:               3,
				Deaths:              2,
				Assists:             1,
				TotalGoldEarned:     12885,
				CreepScore:          297,
				KillParticipation:   0.8,
				ChampionDamageShare: 0.4242806782130524,
				WardsPlaced:         13,
				WardsDestroyed:      15,
				AttackDamage:        295,
				AbilityPower:        0,
				CriticalChance:      0.0,
				AttackSpeed:         177,
				LifeSteal:           8,
				Armor:               111,
				MagicResistance:     47,
				Tenacity:            0.0,
				Items:               []int{1055, 3078, 3363, 3047, 6676, 1037},
				Abilities:           []string{"E", "Q", "W", "Q", "Q", "R", "Q", "E", "Q", "E", "R", "E", "E", "W", "W"},
				Runes: models.Runes{
					Main:      8000,
					Secondary: 8300,
					Perks:     []int{8010, 8009, 9103, 8017, 8316, 8304, 5005, 5008, 5011},
				},
			},
			{
				ExternalID:          "106267619732229182",
				Level:               10,
				Kills:               0,
				Deaths:              6,
				Assists:             4,
				TotalGoldEarned:     6879,
				CreepScore:          35,
				KillParticipation:   0.8,
				ChampionDamageShare: 0.0947822069258486,
				WardsPlaced:         69,
				WardsDestroyed:      16,
				AttackDamage:        89,
				AbilityPower:        0,
				CriticalChance:      0.0,
				AttackSpeed:         122,
				LifeSteal:           0,
				Armor:               80,
				MagicResistance:     67,
				Tenacity:            0.0,
				Items:               []int{3364, 3107, 3111, 3067, 2055},
				Abilities:           []string{"E", "Q", "W", "W", "W", "R", "W", "E", "W", "E"},
				Runes: models.Runes{
					Main:      8300,
					Secondary: 8400,
					Perks:     []int{8360, 8306, 8345, 8347, 8473, 8242, 5007, 5010, 5011},
				},
			},
			{
				ExternalID:          "105501790021137688",
				Level:               16,
				Kills:               5,
				Deaths:              0,
				Assists:             11,
				TotalGoldEarned:     12241,
				CreepScore:          209,
				KillParticipation:   0.8,
				ChampionDamageShare: 0.29748164468325444,
				WardsPlaced:         12,
				WardsDestroyed:      9,
				AttackDamage:        138,
				AbilityPower:        307,
				CriticalChance:      0.0,
				AttackSpeed:         244,
				LifeSteal:           0,
				Armor:               135,
				MagicResistance:     79,
				Tenacity:            0.0,
				Items:               []int{1082, 6653, 3111, 3364, 8010, 2055},
				Abilities:           []string{"E", "Q", "W", "Q", "Q", "R", "Q", "E", "Q", "E", "R", "E", "E", "W", "W", "R"},
				Runes: models.Runes{
					Main:      8200,
					Secondary: 8300,
					Perks:     []int{8229, 8275, 8210, 8237, 8345, 8347, 5007, 5008, 5011},
				},
			},
			{
				ExternalID:          "107492130895812257",
				Level:               15,
				Kills:               1,
				Deaths:              2,
				Assists:             12,
				TotalGoldEarned:     10863,
				CreepScore:          198,
				KillParticipation:   0.65,
				ChampionDamageShare: 0.10876518010390435,
				WardsPlaced:         10,
				WardsDestroyed:      15,
				AttackDamage:        250,
				AbilityPower:        47,
				CriticalChance:      0.0,
				AttackSpeed:         174,
				LifeSteal:           0,
				Armor:               98,
				MagicResistance:     88,
				Tenacity:            0.0,
				Items:               []int{6610, 3111, 3071, 2021, 2055},
				Abilities:           []string{"W", "E", "Q", "W", "W", "R", "W", "E", "W", "Q", "R", "E", "E", "E", "Q"},
				Runes: models.Runes{
					Main:      8000,
					Secondary: 8300,
					Perks:     []int{8010, 9111, 9104, 8014, 8304, 8347, 5005, 5008, 5011},
				},
			},
			{
				ExternalID:          "105501715923396261",
				Level:               17,
				Kills:               3,
				Deaths:              1,
				Assists:             12,
				TotalGoldEarned:     12294,
				CreepScore:          258,
				KillParticipation:   0.75,
				ChampionDamageShare: 0.1771268411402263,
				WardsPlaced:         18,
				WardsDestroyed:      10,
				AttackDamage:        138,
				AbilityPower:        342,
				CriticalChance:      0.0,
				AttackSpeed:         131,
				LifeSteal:           0,
				Armor:               96,
				MagicResistance:     77,
				Tenacity:            0.0,
				Items:               []int{3111, 3116, 3364, 3135},
				Abilities:           []string{"Q", "E", "W", "Q", "Q", "R", "Q", "E", "Q", "E", "R", "E", "E", "W", "W", "R", "W"},
				Runes: models.Runes{
					Main:      8200,
					Secondary: 8000,
					Perks:     []int{8230, 8226, 8233, 8237, 9105, 8017, 5005, 5008, 5011},
				},
			},
			{
				ExternalID:          "109523135356383683",
				Level:               17,
				Kills:               11,
				Deaths:              0,
				Assists:             5,
				TotalGoldEarned:     17036,
				CreepScore:          328,
				KillParticipation:   0.8,
				ChampionDamageShare: 0.3309599891356028,
				WardsPlaced:         17,
				WardsDestroyed:      18,
				AttackDamage:        323,
				AbilityPower:        47,
				CriticalChance:      0.0,
				AttackSpeed:         201,
				LifeSteal:           8,
				Armor:               140,
				MagicResistance:     52,
				Tenacity:            0.0,
				Items:               []int{1055, 3508, 3363, 6675, 3071, 6673, 3174},
				Abilities:           []string{"Q", "E", "W", "Q", "Q", "R", "Q", "E", "Q", "E", "R", "E", "E", "W", "W", "R", "W"},
				Runes: models.Runes{
					Main:      8000,
					Secondary: 8300,
					Perks:     []int{8005, 9111, 9103, 8017, 8304, 8347, 5005, 5008, 5011},
				},
			},
			{
				ExternalID:          "101388913291808185",
				Level:               13,
				Kills:               0,
				Deaths:              2,
				Assists:             20,
				TotalGoldEarned:     8041,
				CreepScore:          16,
				KillParticipation:   1.0,
				ChampionDamageShare: 0.08566634493701211,
				WardsPlaced:         69,
				WardsDestroyed:      5,
				AttackDamage:        137,
				AbilityPower:        47,
				CriticalChance:      0.0,
				AttackSpeed:         138,
				LifeSteal:           0,
				Armor:               94,
				MagicResistance:     83,
				Tenacity:            0.0,
				Items:               []int{3364, 3107, 3111, 3067, 3114},
				Abilities:           []string{"Q", "E", "W", "Q", "Q", "R", "E", "E", "E", "E", "R", "Q", "Q"},
				Runes: models.Runes{
					Main:      8300,
					Secondary: 8400,
					Perks:     []int{8360, 8304, 8345, 8347, 8463, 8444, 5007, 5010, 5011},
				},
			},
		},
	}

	tests := []struct {
		name     string
		game     models.Game
		mock     func(mockClient *MockClient)
		expected modelsRiot.FrameResponse
		wantErr  bool
	}{
		{
			name: "Valid Data",
			game: models.Game{
				ID:           gameID,
				ExternalID:   "113503500599535523",
				ScheduleTime: gameScheduleTime,
				Teams: []models.GameTeam{
					{
						ID:   team1,
						Side: "blue",
					},
					{
						ID:   team2,
						Side: "red",
					},
				},
			},
			mock: func(mockClient *MockClient) {
				now, _ := time.Parse(time.RFC3339, "2025-10-21T09:30:58.492Z")
				mockClient.On("GetFrames", "113503500599535523", now).Return(mockValid, nil)
				mockClient.On("GetFrames", "113503500599535523", gameStart).Return(mockValid, client.ErrorGameFrameNoContent)
				mockClient.On("GetPlayerFrames", "113503500599535523", gameStart).Return(mockPlayerValid, nil)
			},
			expected: modelsRiot.FrameResponse{
				Frame:     frame,
				GameStart: gameStart,
				GameEnd:   gameEnd,
				Players: []modelsRiot.PlayerFrame{
					{
						Champion:   "KSante",
						TeamID:     team1,
						Side:       "blue",
						ExternalID: "107492121209956865",
						Name:       "KT PerfecT",
						Role:       "top",
					},
					{
						Champion:   "MonkeyKing",
						TeamID:     team1,
						Side:       "blue",
						ExternalID: "98767975940969117",
						Name:       "KT Cuzz",
						Role:       "jungle",
					},
					{
						Champion:   "Aurora",
						TeamID:     team1,
						Side:       "blue",
						ExternalID: "98767975949926165",
						Name:       "KT Bdd",
						Role:       "mid",
					},
					{
						Champion:   "Xayah",
						TeamID:     team1,
						Side:       "blue",
						ExternalID: "102483272156027229",
						Name:       "KT deokdam",
						Role:       "bottom",
					},
					{
						Champion:   "Rakan",
						TeamID:     team1,
						Side:       "blue",
						ExternalID: "106267619732229182",
						Name:       "KT Peter",
						Role:       "support",
					},
					{
						Champion:   "Sion",
						TeamID:     team2,
						Side:       "red",
						ExternalID: "105501790021137688",
						Name:       "BFX Clear",
						Role:       "top",
					},
					{
						Champion:   "Trundle",
						TeamID:     team2,
						Side:       "red",
						ExternalID: "107492130895812257",
						Name:       "BFX Raptor",
						Role:       "jungle",
					},
					{
						Champion:   "Morgana",
						TeamID:     team2,
						Side:       "red",
						ExternalID: "105501715923396261",
						Name:       "BFX VicLa",
						Role:       "mid",
					},
					{
						Champion:   "Sivir",
						TeamID:     team2,
						Side:       "red",
						ExternalID: "109523135356383683",
						Name:       "BFX Diable",
						Role:       "bottom",
					},
					{
						Champion:   "Alistar",
						TeamID:     team2,
						Side:       "red",
						ExternalID: "101388913291808185",
						Name:       "BFX Kellin",
						Role:       "support",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockClient)
			if tt.mock != nil {
				tt.mock(mockClient)
			}
			clientFrame := client.NewTeamFramesClient(mockClient)

			frames, err := clientFrame.LoadData(tt.game)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.expected, frames)
		})
	}
}

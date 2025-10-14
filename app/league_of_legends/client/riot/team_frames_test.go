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

func TestLoadFrames(t *testing.T) {
	mockValid := loadMockGamesData(mockFrames)
	gameID := primitive.NewObjectID()

	team1 := primitive.NewObjectID()
	team2 := primitive.NewObjectID()

	timestamp, err := time.Parse(time.RFC3339, "2025-09-11T09:30:58.492Z")
	if err != nil {
		t.Fatalf("Failed to parse timestamp: %v", err)
	}

	frame := models.Frame{
		GameID:    gameID,
		TimeStamp: timestamp,
		Teams: []models.FrameTeam{
			{
				TeamID:     team1,
				Gold:       55156,
				Inhibitors: 1,
				Towers:     8,
				Barons:     1,
				TotalKills: 21,
				Dragons:    []string{"hextech", "hextech"},
			},
			{
				TeamID:     team2,
				Gold:       44434,
				Inhibitors: 0,
				Towers:     1,
				Barons:     0,
				TotalKills: 9,
				Dragons:    []string{"ocean", "cloud"},
			},
		},
	}

	tests := []struct {
		name     string
		mock     func(mockClient *MockClient)
		expected modelsRiot.FrameResponse
		wantErr  bool
	}{
		{
			name: "Valid Data",
			mock: func(mockClient *MockClient) {
				mockClient.On("GetFrames", "113503500599535523").Return(mockValid, nil)
			},
			expected: modelsRiot.FrameResponse{
				Frame: frame,
				PlayerGamesDetails: []models.GamePlayer{
					{
						PlayerID:   primitive.NilObjectID,
						Champion:   "KSante",
						Side:       "blue",
						ExternalID: "107492121209956865",
						Name:       "KT PerfecT",
						TeamID:     team1,
						Role:       "top",
					},
					{
						PlayerID:   primitive.NilObjectID,
						Champion:   "MonkeyKing",
						Side:       "blue",
						ExternalID: "98767975940969117",
						Name:       "KT Cuzz",
						TeamID:     team1,
						Role:       "jungle",
					},
					{
						PlayerID:   primitive.NilObjectID,
						Champion:   "Aurora",
						Side:       "blue",
						ExternalID: "98767975949926165",
						Name:       "KT Bdd",
						TeamID:     team1,
						Role:       "mid",
					},
					{
						PlayerID:   primitive.NilObjectID,
						Champion:   "Xayah",
						Side:       "blue",
						ExternalID: "102483272156027229",
						Name:       "KT deokdam",
						TeamID:     team1,
						Role:       "bottom",
					},
					{
						PlayerID:   primitive.NilObjectID,
						Champion:   "Rakan",
						Side:       "blue",
						ExternalID: "106267619732229182",
						Name:       "KT Peter",
						TeamID:     team1,
						Role:       "support",
					},
					{
						PlayerID:   primitive.NilObjectID,
						Champion:   "Sion",
						Side:       "red",
						ExternalID: "105501790021137688",
						Name:       "BFX Clear",
						TeamID:     team2,
						Role:       "top",
					},
					{
						PlayerID:   primitive.NilObjectID,
						Champion:   "Trundle",
						Side:       "red",
						ExternalID: "107492130895812257",
						Name:       "BFX Raptor",
						TeamID:     team2,
						Role:       "jungle",
					},
					{
						PlayerID:   primitive.NilObjectID,
						Champion:   "Morgana",
						Side:       "red",
						ExternalID: "105501715923396261",
						Name:       "BFX VicLa",
						TeamID:     team2,
						Role:       "mid",
					},
					{
						PlayerID:   primitive.NilObjectID,
						Champion:   "Sivir",
						Side:       "red",
						ExternalID: "109523135356383683",
						Name:       "BFX Diable",
						TeamID:     team2,
						Role:       "bottom",
					},
					{
						PlayerID:   primitive.NilObjectID,
						Champion:   "Alistar",
						Side:       "red",
						ExternalID: "101388913291808185",
						Name:       "BFX Kellin",
						TeamID:     team2,
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

			game := models.Game{
				ID:         gameID,
				ExternalID: "113503500599535523",
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
			}

			frames, err := clientFrame.LoadData(game)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.expected, frames)
		})
	}
}

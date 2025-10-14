package riot_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	client "betty/science/app/league_of_legends/client/riot"
	"betty/science/app/league_of_legends/models"
)

var mockGameDetailsPath = "../../test/mock_data/game_details.json"

func loadMockGameData(path string) client.GameDetailsResponse {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open mock data file: %v", err)
	}
	defer file.Close()

	var data client.GameDetailsResponse
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Fatalf("Failed to decode mock data: %v", err)
	}
	return data
}

func TestLoadGame(t *testing.T) {
	mockValid := loadMockGameData(mockGameDetailsPath)
	matchID := primitive.NewObjectID()

	tests := []struct {
		name     string
		mock     func(mockClient *MockClient)
		expected []models.Game
		wantErr  bool
	}{
		{
			name: "Valid Data",
			expected: []models.Game{
				{
					ExternalID: "113503500599535523",
					MatchID:    matchID,
					Number:     1,
					State:      "completed",
					Teams: []models.GameTeam{
						{
							ExternalID: "99566404579461230",
							Name:       "kt Rolster",
							Side:       "blue",
						},
						{
							ExternalID: "100725845022060229",
							Name:       "BNK FEARX",
							Side:       "red",
						},
					},
				},
				{
					ExternalID: "113503500599535524",
					MatchID:    matchID,
					Number:     2,
					State:      "completed",
					Teams: []models.GameTeam{
						{
							ExternalID: "99566404579461230",
							Name:       "kt Rolster",
							Side:       "blue",
						},
						{
							ExternalID: "100725845022060229",
							Name:       "BNK FEARX",
							Side:       "red",
						},
					},
				},
				{
					ExternalID: "113503500599535525",
					MatchID:    matchID,
					Number:     3,
					State:      "completed",
					Teams: []models.GameTeam{
						{
							ExternalID: "99566404579461230",
							Name:       "kt Rolster",
							Side:       "red",
						},
						{
							ExternalID: "100725845022060229",
							Name:       "BNK FEARX",
							Side:       "blue",
						},
					},
				},
				{
					ExternalID: "113503500599535526",
					MatchID:    matchID,
					Number:     4,
					State:      "completed",
					Teams: []models.GameTeam{
						{
							ExternalID: "99566404579461230",
							Name:       "kt Rolster",
							Side:       "red",
						},
						{
							ExternalID: "100725845022060229",
							Name:       "BNK FEARX",
							Side:       "blue",
						},
					},
				},
			},
			wantErr: false,
			mock: func(mockClient *MockClient) {
				mockClient.On("GetGameDetails", "match1").Return(mockValid, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockClient)
			tt.mock(mockClient)

			clientGame := client.NewClientGame(mockClient)
			match := models.Match{ID: matchID, ExternalID: "match1"}
			games, err := clientGame.LoadData(match)

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, len(tt.expected), len(games), "Expected number of games does not match")
			assert.Equal(t, tt.expected, games, "Expected game data does not match")
			mockClient.AssertExpectations(t)
		})
	}
}

package league_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	client "betty/science/app/riot/clients/league"
)

const mockDataPath = "../../test/mock_data/schedule.json"
const mockTBDSchedulePath = "../../test/mock_data/schedule_tbd.json"

type MockClient struct {
	mock.Mock
}

func (m *MockClient) GetSchedule() (client.ScheduleResponse, error) {
	args := m.Called()

	return args.Get(0).(client.ScheduleResponse), args.Error(1)
}

func (m *MockClient) GetGameDetails(matchID string) (client.GameDetailsResponse, error) {
	args := m.Called(matchID)

	return args.Get(0).(client.GameDetailsResponse), args.Error(1)
}

func (m *MockClient) GetFrames(gameID string, startingTime time.Time) (client.FramesResponse, error) {
	args := m.Called(gameID, startingTime)

	return args.Get(0).(client.FramesResponse), args.Error(1)
}

func (m *MockClient) GetPlayerFrames(gameID string, startingTime time.Time) (client.PlayerFramesResponse, error) {
	args := m.Called(gameID, startingTime)

	return args.Get(0).(client.PlayerFramesResponse), args.Error(1)
}

func (m *MockClient) GetFirstFrame(gameID string, startingTime time.Time) (client.GameFrame, error) {
	args := m.Called(gameID, startingTime)

	return args.Get(0).(client.GameFrame), args.Error(1)
}
func loadMockData(path string) client.ScheduleResponse {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open mock data file: %v", err)
	}
	defer file.Close()

	var data client.ScheduleResponse
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Fatalf("Failed to decode mock data: %v", err)
	}
	return data
}

func TestMatchParse(t *testing.T) {
	mockValid := loadMockData(mockDataPath)

	tests := []struct {
		name     string
		mock     func(mockClient *MockClient)
		expected int
		wantErr  bool
	}{
		{
			name:     "Valid Data",
			expected: len(mockValid.Data.Schedule.Events),
			wantErr:  false,
			mock: func(mockClient *MockClient) {
				mockClient.On("GetSchedule").Return(loadMockData(mockDataPath), nil)
			},
		},
		{
			name:     "TBD Matches",
			expected: 0,
			wantErr:  false,
			mock: func(mockClient *MockClient) {
				mockClient.On("GetSchedule").Return(loadMockData(mockTBDSchedulePath), nil)
			},
		},
		{
			name:     "Empty Data",
			expected: 0,
			wantErr:  false,
			mock: func(mockClient *MockClient) {
				mockClient.On("GetSchedule").Return(client.ScheduleResponse{}, nil)
			},
		},
		{
			name:     "API Error",
			expected: 0,
			mock: func(mockClient *MockClient) {
				mockClient.On("GetSchedule").Return(client.ScheduleResponse{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockClient)

			tt.mock(mockClient)

			parser := client.NewClientMatch(mockClient)
			result, err := parser.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.expected, len(result))
		})
	}
}

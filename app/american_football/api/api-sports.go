package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SportsApiClient struct {
	APIKey string
	URL    string
}

func NewSportsApiClient(apiKey, url string) *SportsApiClient {
	return &SportsApiClient{
		APIKey: apiKey,
		URL:    url,
	}
}

func (api *SportsApiClient) get(endpoint string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", api.URL, endpoint), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.APIKey))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

func (api *SportsApiClient) GetSchedule(leagueID int, date time.Time) ([]ScheduleResponse, error) {
	fmtDate := date.Format("2006-01-02")
	url := fmt.Sprintf("games?date=%s&league=%d", fmtDate, leagueID)
	body, status, err := api.get(url)
	if err != nil {
		return []ScheduleResponse{}, err
	}
	if status != 200 {
		return []ScheduleResponse{}, fmt.Errorf("error fetching schedule: %d - %s", status, string(body))
	}
	var result DefaultResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return []ScheduleResponse{}, err
	}

	var games []ScheduleResponse
	gamesBytes, err := json.Marshal(result.Response)
	if err != nil {
		return []ScheduleResponse{}, err
	}
	if err := json.Unmarshal(gamesBytes, &games); err != nil {
		return []ScheduleResponse{}, err
	}

	return games, nil
}

func (api *SportsApiClient) GetGameTeamsStats(gameID int) ([]GameTeamsStatsResponse, error) {
	url := fmt.Sprintf("games/statistics/teams?id=%d", gameID)
	body, status, err := api.get(url)
	if err != nil {
		return []GameTeamsStatsResponse{}, err
	}
	if status != 200 {
		return []GameTeamsStatsResponse{}, fmt.Errorf("error fetching game teams stats: %d - %s", status, string(body))
	}
	var result DefaultResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return []GameTeamsStatsResponse{}, err
	}

	var teamsStats []GameTeamsStatsResponse
	teamsBytes, err := json.Marshal(result.Response)
	if err != nil {
		return []GameTeamsStatsResponse{}, err
	}
	if err := json.Unmarshal(teamsBytes, &teamsStats); err != nil {
		return []GameTeamsStatsResponse{}, err
	}

	return teamsStats, nil
}

func (api *SportsApiClient) GetGamePlayersStats(gameID int) ([]GamePlayerStatsResponse, error) {
	url := fmt.Sprintf("game/statistics/players?id=%d", gameID)
	body, status, err := api.get(url)
	if err != nil {
		return []GamePlayerStatsResponse{}, err
	}
	if status != 200 {
		return []GamePlayerStatsResponse{}, fmt.Errorf("error fetching game players stats: %d - %s", status, string(body))
	}
	var result DefaultResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return []GamePlayerStatsResponse{}, err
	}

	var playersStats []GamePlayerStatsResponse
	playersBytes, err := json.Marshal(result.Response)
	if err != nil {
		return []GamePlayerStatsResponse{}, err
	}
	if err := json.Unmarshal(playersBytes, &playersStats); err != nil {
		return []GamePlayerStatsResponse{}, err
	}

	return playersStats, nil
}

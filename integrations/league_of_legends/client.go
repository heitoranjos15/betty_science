package league_of_legends

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type LeagueOfLegendsClient struct {
	EventsURL   string
	GameLiveURL string
	APIKey      string
	Logger      *log.Logger
}

func NewLeagueOfLegendsAPI(apiKey string) *LeagueOfLegendsClient {
	return &LeagueOfLegendsClient{
		EventsURL:   "https://esports-api.lolesports.com",
		GameLiveURL: "https://feed.lolesports.com",
		APIKey:      apiKey,
		Logger:      log.New(os.Stdout, "[lol_api] ", log.LstdFlags),
	}
}

func (api *LeagueOfLegendsClient) get(url string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("x-api-key", api.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

func (api *LeagueOfLegendsClient) GetSchedule() ([]interface{}, error) {
	url := fmt.Sprintf("%s/persisted/gw/getSchedule?hl=en-US", api.EventsURL)
	body, status, err := api.get(url)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("error fetching events: %d - %s", status, string(body))
	}
	var result struct {
		Data struct {
			Schedule struct {
				Events []interface{} `json:"events"`
			} `json:"schedule"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Data.Schedule.Events, nil
}

func (api *LeagueOfLegendsClient) GetEventDetail(eventID int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/persisted/gw/getEventDetails?hl=en-US&id=%d", api.EventsURL, eventID)
	body, status, err := api.get(url)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("error fetching event data: %d - %s", status, string(body))
	}
	var result struct {
		Data struct {
			Event map[string]interface{} `json:"event"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Data.Event, nil
}

type LolMatchStartTime time.Time

func (t LolMatchStartTime) formatSeconds() time.Time {

	date := time.Now().UTC()
	seconds := date.Second()
	if seconds%10 > 0 {
		seconds = seconds - (seconds % 10)
	}
	date = date.Add(-60 * time.Second)
	return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), seconds, 0, time.UTC)
}

func (t LolMatchStartTime) ToString() string {
	return t.formatSeconds().Format(time.RFC3339) + "Z"
}

func (api *LeagueOfLegendsClient) GetLiveFrame(gameID string, startingTime LolMatchStartTime) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/livestats/v1/window/%s?startingTime=%s", api.GameLiveURL, gameID, startingTime.ToString())
	body, status, err := api.get(url)
	if err != nil {
		return nil, err
	}
	if status == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}
		return result, nil
	}
	if status == 204 {
		return nil, errors.New("game is not live or no data available yet")
	}
	return nil, fmt.Errorf("error fetching game data: %d - %s", status, string(body))
}

func (api *LeagueOfLegendsClient) GetFirstFrame(gameID string, matchStartTime LolMatchStartTime) (map[string]interface{}, error) {
	retryCount := 0
	retryLimit := 150

	for {
		if retryCount >= retryLimit {
			return nil, errors.New("exceeded maximum retries to fetch first frame")
		}

		frameData, err := api.GetLiveFrame(gameID, matchStartTime)
		if err != nil {
			return nil, err
		}

		if frameData != nil {
			if frames, ok := frameData["frames"].([]interface{}); ok && len(frames) > 0 {
				firstFrame, ok := frames[0].(map[string]interface{})
				if ok && api.isFirstFrame(firstFrame) {
					return firstFrame, nil
				}
			}

			matchStartTime = LolMatchStartTime(matchStartTime.formatSeconds().Add(5 * time.Minute))
			time.Sleep(2 * time.Second)
			retryCount++
		}
	}
}

func (api *LeagueOfLegendsClient) GetPlayerFrame(gameID string, startTime LolMatchStartTime) (map[string]interface{}, error) {
	body, status, err := api.get(fmt.Sprintf("%s/livestats/v1/details/%s?startingTime=%s", api.GameLiveURL, gameID, startTime.ToString()))
	if err != nil {
		return nil, err
	}
	if status == 200 {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}
		return result, nil
	}
	if status == 204 {
		return nil, errors.New("no player stats found for game (possibly not live yet)")
	}
	return nil, fmt.Errorf("error fetching player stats: %d - %s", status, string(body))
}

func (api *LeagueOfLegendsClient) isFirstFrame(frame map[string]interface{}) bool {
	if frame == nil {
		return false
	}
	if _, ok := frame["totalGoldEarned"].(float64); ok {
		return true
	}
	return false
}

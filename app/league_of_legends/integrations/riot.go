package integrations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"betty/science/app/league_of_legends/client/riot"
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

func (api *LeagueOfLegendsClient) GetSchedule() (riot.ScheduleResponse, error) {
	url := fmt.Sprintf("%s/persisted/gw/getSchedule?hl=en-US", api.EventsURL)
	body, status, err := api.get(url)
	if err != nil {
		return riot.ScheduleResponse{}, err
	}
	if status != 200 {
		return riot.ScheduleResponse{}, fmt.Errorf("error fetching events: %d - %s", status, string(body))
	}
	var result riot.ScheduleResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return riot.ScheduleResponse{}, err
	}
	return result, nil
}

func (api *LeagueOfLegendsClient) GetGameDetails(eventID string) (riot.GameDetailsResponse, error) {
	url := fmt.Sprintf("%s/persisted/gw/getEventDetails?hl=en-US&id=%s", api.EventsURL, eventID)
	body, status, err := api.get(url)
	if err != nil {
		return riot.GameDetailsResponse{}, err
	}
	if status != 200 {
		return riot.GameDetailsResponse{}, fmt.Errorf("error fetching game details: %d - %s", status, string(body))
	}
	var result riot.GameDetailsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return riot.GameDetailsResponse{}, err
	}
	return result, nil
}

type LolMatchStartTime time.Time

func (t LolMatchStartTime) formatSeconds(date time.Time) time.Time {
	seconds := date.Second()
	if seconds%10 > 0 {
		seconds = seconds - (seconds % 10)
	}
	return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), seconds, 0, time.UTC)
}

func (t LolMatchStartTime) ToString(date time.Time) string {
	return t.formatSeconds(date).Format(time.RFC3339)
}

func (api *LeagueOfLegendsClient) GetFrames(gameID string, startingTime time.Time) (riot.FramesResponse, error) {
	fmtStartTime := LolMatchStartTime(startingTime)
	var result riot.FramesResponse

	url := fmt.Sprintf("%s/livestats/v1/window/%s?startingTime=%s", api.GameLiveURL, gameID, fmtStartTime.ToString(startingTime))
	log.Println("Fetching frames from URL:", url)

	body, status, err := api.get(url)
	if err != nil {
		return result, err
	}

	if status == 200 {
		if err := json.Unmarshal(body, &result); err != nil {
			return result, err
		}
		return result, nil
	}

	if status == 204 {
		return result, riot.ErrorGameFrameNoContent
	}

	return result, fmt.Errorf("error fetching game data: %d - %s", status, string(body))
}

func (api *LeagueOfLegendsClient) GetFirstFrame(gameID string, startTime time.Time) (riot.GameFrame, error) {
	fmtStartTime := LolMatchStartTime(startTime)

	retryCount := 0
	retryLimit := 150

	for {
		if retryCount >= retryLimit {
			return riot.GameFrame{}, errors.New("exceeded maximum retries to fetch first frame")
		}

		frameData, err := api.GetFrames(gameID, startTime)
		if err != nil {
			return riot.GameFrame{}, err
		}

		for _, frame := range frameData.Frames {
			if api.isFirstFrame(frame) {
				return frame, nil
			}
		}

		startTime = time.Time(fmtStartTime.formatSeconds(startTime)).Add(5 * time.Minute)

		time.Sleep(2 * time.Second)
		retryCount++
	}
}

func (api *LeagueOfLegendsClient) GetPlayerFrames(gameID string, time time.Time) (riot.PlayerFramesResponse, error) {
	fmtStartTime := LolMatchStartTime(time)
	var result riot.PlayerFramesResponse
	url := fmt.Sprintf("%s/livestats/v1/details/%s?startingTime=%s", api.GameLiveURL, gameID, fmtStartTime.ToString(time))
	log.Println("Fetching player frames from URL:", url)
	body, status, err := api.get(url)
	if err != nil {
		return result, err
	}
	if status == 200 {
		if err := json.Unmarshal(body, &result); err != nil {
			return result, err
		}
		return result, nil
	}
	if status == 204 {
		return result, errors.New("no player stats found for game (possibly not live yet)")
	}
	return result, fmt.Errorf("error fetching player stats: %d - %s", status, string(body))
}

// TODO move this to riot.GameFrame
func (api *LeagueOfLegendsClient) isFirstFrame(frame riot.GameFrame) bool {
	return frame.RedTeam.TotalGold > 0 && frame.BlueTeam.TotalGold > 0
}

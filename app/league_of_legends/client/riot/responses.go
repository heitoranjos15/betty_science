package riot

type ScheduleResponse struct {
	Data struct {
		Schedule struct {
			Pages struct {
				Older string `json:"older"`
				Newer string `json:"newer"`
			} `json:"pages"`
			Events []Event `json:"events"`
		} `json:"schedule"`
	} `json:"data"`
}

type Event struct {
	StartTime string `json:"startTime"`
	State     string `json:"state"`
	Type      string `json:"type"`
	BlockName string `json:"blockName"`
	League    struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	} `json:"league"`
	Match struct {
		ID       string   `json:"id"`
		Flags    []string `json:"flags"`
		Teams    []Team   `json:"teams"`
		Strategy struct {
			Type  string `json:"type"`
			Count int    `json:"count"`
		} `json:"strategy"`
	} `json:"match"`
}

type Team struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Image  string `json:"image"`
	Result *struct {
		Outcome  *string `json:"outcome"`
		GameWins int     `json:"gameWins"`
	} `json:"result"`
	Record *struct {
		Wins   int `json:"wins"`
		Losses int `json:"losses"`
	} `json:"record"`
}

type GameDetailsResponse struct {
	Data struct {
		Event struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Tournament struct {
				ID string `json:"id"`
			} `json:"tournament"`
			League struct {
				ID    string `json:"id"`
				Slug  string `json:"slug"`
				Image string `json:"image"`
				Name  string `json:"name"`
			} `json:"league"`
			Match struct {
				Strategy struct {
					Count int `json:"count"`
				} `json:"strategy"`
				Teams []GameDetailsTeam `json:"teams"`
				Games []GamesDetails    `json:"games"`
			} `json:"match"`
			Streams []interface{} `json:"streams"`
		} `json:"event"`
	} `json:"data"`
}

type GamesDetails struct {
	Number int    `json:"number"`
	ID     string `json:"id"`
	State  string `json:"state"`
	Teams  []struct {
		ID   string `json:"id"`
		Side string `json:"side"`
	} `json:"teams"`
}

type GameDetailsTeam struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Image  string `json:"image"`
	Result struct {
		GameWins int `json:"gameWins"`
	} `json:"result"`
}

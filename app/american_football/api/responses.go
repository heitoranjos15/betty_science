package api

type DefaultResponse struct {
	Get        string            `json:"get"`
	Parameters map[string]string `json:"parameters"`
	Errors     []string          `json:"errors"`
	Results    int               `json:"results"`
	Response   any               `json:"response"`
}

type ScheduleResponse struct {
	Game   GameInfo   `json:"game"`
	League LeagueInfo `json:"league"`
	Teams  TeamsInfo  `json:"teams"`
	Scores ScoresInfo `json:"scores"`
}

type GameInfo struct {
	ID     int        `json:"id"`
	Stage  string     `json:"stage"`
	Week   string     `json:"week"`
	Date   GameDate   `json:"date"`
	Venue  VenueInfo  `json:"venue"`
	Status StatusInfo `json:"status"`
}

type GameDate struct {
	Timezone  string `json:"timezone"`
	Date      string `json:"date"`
	Time      string `json:"time"`
	Timestamp int64  `json:"timestamp"`
}

type VenueInfo struct {
	Name string `json:"name"`
	City string `json:"city"`
}

type StatusInfo struct {
	Short string  `json:"short"`
	Long  string  `json:"long"`
	Timer *string `json:"timer"`
}

type LeagueInfo struct {
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Season  string      `json:"season"`
	Logo    string      `json:"logo"`
	Country CountryInfo `json:"country"`
}

type CountryInfo struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Flag string `json:"flag"`
}

type TeamsInfo struct {
	Home TeamInfo `json:"home"`
	Away TeamInfo `json:"away"`
}

type TeamInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

type ScoresInfo struct {
	Home ScoreDetail `json:"home"`
	Away ScoreDetail `json:"away"`
}

type ScoreDetail struct {
	Quarter1 *int `json:"quarter_1"`
	Quarter2 *int `json:"quarter_2"`
	Quarter3 *int `json:"quarter_3"`
	Quarter4 *int `json:"quarter_4"`
	Overtime *int `json:"overtime"`
	Total    int  `json:"total"`
}

type GameTeamsStatsResponse struct {
	Team       TeamInfo       `json:"team"`
	Statistics TeamStatistics `json:"statistics"`
}

type TeamStatistics struct {
	FirstDowns       FirstDownsStats `json:"first_downs"`
	Plays            PlaysStats      `json:"plays"`
	Yards            YardsStats      `json:"yards"`
	Passing          PassingStats    `json:"passing"`
	Rushings         RushingsStats   `json:"rushings"`
	RedZone          RedZoneStats    `json:"red_zone"`
	Penalties        PenaltiesStats  `json:"penalties"`
	Turnovers        TurnoversStats  `json:"turnovers"`
	Posession        PosessionStats  `json:"posession"`
	Interceptions    SimpleIntStat   `json:"interceptions"`
	FumblesRecovered SimpleIntStat   `json:"fumbles_recovered"`
	Sacks            SimpleIntStat   `json:"sacks"`
	Safeties         SimpleIntStat   `json:"safeties"`
	IntTouchdowns    SimpleIntStat   `json:"int_touchdowns"`
	PointsAgainst    SimpleIntStat   `json:"points_against"`
}

type FirstDownsStats struct {
	Total                int    `json:"total"`
	Passing              int    `json:"passing"`
	Rushing              int    `json:"rushing"`
	FromPenalties        int    `json:"from_penalties"`
	ThirdDownEfficiency  string `json:"third_down_efficiency"`
	FourthDownEfficiency string `json:"fourth_down_efficiency"`
}

type PlaysStats struct {
	Total int `json:"total"`
}

type YardsStats struct {
	Total        int    `json:"total"`
	YardsPerPlay string `json:"yards_per_play"`
	TotalDrives  string `json:"total_drives"`
}

type PassingStats struct {
	Total               int    `json:"total"`
	CompAtt             string `json:"comp_att"`
	YardsPerPass        string `json:"yards_per_pass"`
	InterceptionsThrown int    `json:"interceptions_thrown"`
	SacksYardsLost      string `json:"sacks_yards_lost"`
}

type RushingsStats struct {
	Total        int    `json:"total"`
	Attempts     int    `json:"attempts"`
	YardsPerRush string `json:"yards_per_rush"`
}

type RedZoneStats struct {
	MadeAtt string `json:"made_att"`
}

type PenaltiesStats struct {
	Total string `json:"total"`
}

type TurnoversStats struct {
	Total         int `json:"total"`
	LostFumbles   int `json:"lost_fumbles"`
	Interceptions int `json:"interceptions"`
}

type PosessionStats struct {
	Total string `json:"total"`
}

type SimpleIntStat struct {
	Total int `json:"total"`
}

type GamePlayerStatsResponse struct {
	Team   TeamInfo         `json:"team"`
	Groups []TeamStatsGroup `json:"groups"`
}

type TeamStatsGroup struct {
	Name    string            `json:"name"`
	Players []TeamStatsPlayer `json:"players"`
}

type TeamStatsPlayer struct {
	Player     TeamStatsPlayerInfo   `json:"player"`
	Statistics []TeamPlayerStatistic `json:"statistics"`
}

type TeamStatsPlayerInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type TeamPlayerStatistic struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

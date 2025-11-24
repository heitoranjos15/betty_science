package league

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

type FramesResponse struct {
	EsportsGameID  string       `json:"esportsGameId"`
	EsportsMatchID string       `json:"esportsMatchId"`
	GameMetadata   GameMetadata `json:"gameMetadata"`
	Frames         []GameFrame  `json:"frames"`
}

type GameMetadata struct {
	PatchVersion     string       `json:"patchVersion"`
	BlueTeamMetadata TeamMetadata `json:"blueTeamMetadata"`
	RedTeamMetadata  TeamMetadata `json:"redTeamMetadata"`
}

type TeamMetadata struct {
	EsportsTeamID       string            `json:"esportsTeamId"`
	ParticipantMetadata []ParticipantMeta `json:"participantMetadata"`
}

type ParticipantMeta struct {
	ParticipantID   int    `json:"participantId"`
	EsportsPlayerID string `json:"esportsPlayerId"`
	SummonerName    string `json:"summonerName"`
	ChampionID      string `json:"championId"`
	Role            string `json:"role"`
}

type GameFrame struct {
	Rfc460Timestamp string    `json:"rfc460Timestamp"`
	GameState       string    `json:"gameState"`
	BlueTeam        TeamFrame `json:"blueTeam"`
	RedTeam         TeamFrame `json:"redTeam"`
}

type TeamFrame struct {
	TotalGold    int           `json:"totalGold"`
	Inhibitors   int           `json:"inhibitors"`
	Towers       int           `json:"towers"`
	Barons       int           `json:"barons"`
	TotalKills   int           `json:"totalKills"`
	Dragons      []string      `json:"dragons"`
	Participants []Participant `json:"participants"`
}

type Participant struct {
	ParticipantID int `json:"participantId"`
	TotalGold     int `json:"totalGold"`
	Level         int `json:"level"`
	Kills         int `json:"kills"`
	Deaths        int `json:"deaths"`
	Assists       int `json:"assists"`
	CreepScore    int `json:"creepScore"`
	CurrentHealth int `json:"currentHealth"`
	MaxHealth     int `json:"maxHealth"`
}

type PlayerFramesResponse struct {
	Frames []ParticipantFrame `json:"frames"`
}

type ParticipantFrame struct {
	Rfc460Timestamp string                `json:"rfc460Timestamp"`
	Participants    []ParticipantSnapshot `json:"participants"`
}

type ParticipantSnapshot struct {
	ParticipantID       int          `json:"participantId"`
	Level               int          `json:"level"`
	Kills               int          `json:"kills"`
	Deaths              int          `json:"deaths"`
	Assists             int          `json:"assists"`
	TotalGoldEarned     int          `json:"totalGoldEarned"`
	CreepScore          int          `json:"creepScore"`
	KillParticipation   float64      `json:"killParticipation"`
	ChampionDamageShare float64      `json:"championDamageShare"`
	WardsPlaced         int          `json:"wardsPlaced"`
	WardsDestroyed      int          `json:"wardsDestroyed"`
	AttackDamage        int          `json:"attackDamage"`
	AbilityPower        int          `json:"abilityPower"`
	CriticalChance      float64      `json:"criticalChance"`
	AttackSpeed         int          `json:"attackSpeed"`
	LifeSteal           int          `json:"lifeSteal"`
	Armor               int          `json:"armor"`
	MagicResistance     int          `json:"magicResistance"`
	Tenacity            float64      `json:"tenacity"`
	Items               []int        `json:"items"`
	PerkMetadata        PerkMetadata `json:"perkMetadata"`
	Abilities           []string     `json:"abilities"`
}

type PerkMetadata struct {
	StyleID    int   `json:"styleId"`
	SubStyleID int   `json:"subStyleId"`
	Perks      []int `json:"perks"`
}

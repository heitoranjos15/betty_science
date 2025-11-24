package league

import (
	"errors"
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"

	models "betty/science/app/riot"
)

var ErrorGameFrameNoContent = errors.New("game frame has no content")

type clientFrame struct {
	api api
}

func NewFramesClient(api api) *clientFrame {
	return &clientFrame{
		api: api,
	}
}

type frameCollectionResult struct {
	Frames       []GameFrame
	PlayerFrames []ParticipantFrame
	Error        error
}

type frameCollectChan struct {
	Frame       GameFrame
	PlayerFrame ParticipantFrame
}

type workerType struct {
	name     string
	startMin int
	endMin   int
}

func (c *clientFrame) collectFrameWorker(worker workerType, gameExternalID string, frameTime time.Time, channel chan<- frameCollectChan, wg *sync.WaitGroup) {
	defer wg.Done()

	timeCounter := worker.startMin
	for {
		if worker.endMin > 0 && timeCounter > worker.endMin {
			return
		}

		searchTime := frameTime.Add(time.Duration(-timeCounter) * time.Minute)

		frames, err := c.api.GetFrames(gameExternalID, searchTime)
		if err != nil {
			if errors.Is(err, ErrorGameFrameNoContent) {
				log.Printf("[%s][client-team-frame] No more frames available for game: %s", worker.name, gameExternalID)
				return
			}
		}

		playerFrames, err := c.api.GetPlayerFrames(gameExternalID, searchTime)
		if err != nil {
			log.Printf("[%s] Error fetching player frames for game: %s err: %s", worker.name, gameExternalID, err)
			return
		}

		if c.isFirstFrame(frames.Frames[0]) {
			log.Printf("[%s] Reached first frame for game: %s", worker.name, gameExternalID)
			return
		}
		log.Printf("[%s] Collected frame at %s for game: %s", worker.name, frames.Frames[len(frames.Frames)-1].Rfc460Timestamp, gameExternalID)
		channel <- frameCollectChan{
			Frame:       frames.Frames[len(frames.Frames)-1],
			PlayerFrame: playerFrames.Frames[len(playerFrames.Frames)-1],
		}

		randomDelay := time.Duration(rand.Intn(1500)+500) * time.Millisecond
		time.Sleep(randomDelay) // to avoid rate limiting

		timeCounter++
	}

}

func (c *clientFrame) isFirstFrame(frame GameFrame) bool {
	return frame.RedTeam.TotalGold == 0 && frame.BlueTeam.TotalGold == 0
}

func (c *clientFrame) collectFrames(frames []GameFrame, gameExternalID string) (frameCollectionResult, error) {
	endGameFrame := frames[len(frames)-1]

	endTime, err := time.Parse(time.RFC3339, endGameFrame.Rfc460Timestamp)
	if err != nil {
		return frameCollectionResult{}, err
	}

	playerFrameResp, err := c.api.GetPlayerFrames(gameExternalID, endTime)
	if err != nil {
		return frameCollectionResult{}, err
	}
	playerFrames := playerFrameResp.Frames

	channel := make(chan frameCollectChan, 3)
	wg := sync.WaitGroup{}

	workers := []workerType{
		{name: "buzz", startMin: 1, endMin: 15},
		{name: "crowbar jones", startMin: 15, endMin: 25},
		{name: "flick", startMin: 25, endMin: 0},
	}

	for _, worker := range workers {
		wg.Add(1)
		go c.collectFrameWorker(worker, gameExternalID, endTime, channel, &wg)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for result := range channel {
		frames = append(frames, result.Frame)
		playerFrames = append(playerFrames, result.PlayerFrame)
	}

	sort.SliceStable(frames, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, frames[i].Rfc460Timestamp)
		timeJ, _ := time.Parse(time.RFC3339, frames[j].Rfc460Timestamp)
		return timeI.Before(timeJ)
	})

	sort.SliceStable(playerFrames, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, playerFrames[i].Rfc460Timestamp)
		timeJ, _ := time.Parse(time.RFC3339, playerFrames[j].Rfc460Timestamp)
		return timeI.Before(timeJ)
	})
	return frameCollectionResult{
		Frames:       frames,
		PlayerFrames: playerFrames,
	}, nil
}

func (c *clientFrame) LoadData(game models.Game) (FrameResponse, error) {
	var resp FrameResponse

	now := time.Now().UTC().Add(-2 * time.Minute)
	lastFrame, err := c.api.GetFrames(game.ExternalID, now)
	if err != nil {
		log.Println("[client-team-frame] Error fetching game frames for game:", game.ExternalID, err)
		return resp, err
	}

	allFrames, err := c.collectFrames(lastFrame.Frames, game.ExternalID)
	if err != nil {
		return resp, err
	}

	resp.Players = c.playersDetails(game.Teams, lastFrame.GameMetadata)
	firstFrameTime, _ := time.Parse(time.RFC3339, allFrames.Frames[0].Rfc460Timestamp)
	lastFrameTime, _ := time.Parse(time.RFC3339, lastFrame.Frames[len(lastFrame.Frames)-1].Rfc460Timestamp)
	resp.GameStart = firstFrameTime
	resp.GameEnd = lastFrameTime

	frames := []models.Frame{}
	log.Printf("debug: total frames collected: %d %d", len(allFrames.Frames), len(allFrames.PlayerFrames))
	for i, f := range allFrames.Frames {
		players := ParticipantFrame{}
		if len(allFrames.PlayerFrames) <= i {
			log.Printf("warning: no player frame for frame index %d, total player frames: %d", i, len(allFrames.PlayerFrames))
			players = allFrames.PlayerFrames[len(allFrames.PlayerFrames)-1]
		} else {
			players = allFrames.PlayerFrames[i]
		}

		timestamp, _ := time.Parse(time.RFC3339, f.Rfc460Timestamp)
		frames = append(frames, models.Frame{
			GameID:    game.ID,
			Timestamp: timestamp,
			Teams:     c.frameTeams(game.Teams, f),
			Players:   c.framePlayers(game.Teams, lastFrame.GameMetadata, players),
		})
	}
	resp.Frames = frames

	winner := c.findWinner(frames[len(frames)-1])
	resp.WinnerID = winner.TeamID

	return resp, nil
}

func (c clientFrame) findWinner(frame models.Frame) models.FrameTeam {
	score := 1
	blueTeam := frame.Teams[0]
	redTeam := frame.Teams[1]

	if blueTeam.Gold > redTeam.Gold {
		score++
	}
	if blueTeam.Towers > redTeam.Towers {
		score++
	}
	if len(blueTeam.Dragons) > len(redTeam.Dragons) {
		score++
	}
	if blueTeam.Barons > redTeam.Barons {
		score++
	}

	if score >= 2 {
		return blueTeam
	}

	return redTeam
}

func (c clientFrame) framePlayers(teams []models.GameTeam, gameMetadata GameMetadata, playerFrame ParticipantFrame) []models.FramePlayer {
	var frame []models.FramePlayer
	if len(playerFrame.Participants) == 0 {
		log.Println("[client-frame] No player frame data available")
		return frame
	}

	for _, team := range teams {
		playerMeta := c.findParticipantMetadata(team.Side, gameMetadata)
		for _, pm := range playerMeta {
			playerID := pm.ParticipantID
			pf := playerFrame.Participants[playerID-1]
			frame = append(frame, models.FramePlayer{
				ExternalID:          pm.EsportsPlayerID,
				Level:               pf.Level,
				Kills:               pf.Kills,
				Deaths:              pf.Deaths,
				Assists:             pf.Assists,
				TotalGoldEarned:     pf.TotalGoldEarned,
				CreepScore:          pf.CreepScore,
				KillParticipation:   pf.KillParticipation,
				ChampionDamageShare: pf.ChampionDamageShare,
				WardsPlaced:         pf.WardsPlaced,
				WardsDestroyed:      pf.WardsDestroyed,
				AttackDamage:        pf.AttackDamage,
				AbilityPower:        pf.AbilityPower,
				CriticalChance:      pf.CriticalChance,
				AttackSpeed:         pf.AttackSpeed,
				LifeSteal:           pf.LifeSteal,
				Armor:               pf.Armor,
				MagicResistance:     pf.MagicResistance,
				Tenacity:            pf.Tenacity,
				Items:               pf.Items,
				Runes: models.Runes{
					Main:      pf.PerkMetadata.StyleID,
					Secondary: pf.PerkMetadata.SubStyleID,
					Perks:     pf.PerkMetadata.Perks,
				},
				Abilities: pf.Abilities,
			})
		}
	}
	return frame
}

func (c clientFrame) frameTeams(teams []models.GameTeam, gameMetadata GameFrame) []models.FrameTeam {
	var frame []models.FrameTeam

	for _, team := range teams {
		teamMeta := c.findTeamMetadata(team.Side, gameMetadata)
		frame = append(frame, models.FrameTeam{
			TeamID:     team.ID,
			Gold:       teamMeta.TotalGold,
			Towers:     teamMeta.Towers,
			Inhibitors: teamMeta.Inhibitors,
			Dragons:    teamMeta.Dragons,
			Barons:     teamMeta.Barons,
			TotalKills: teamMeta.TotalKills,
		})
	}

	return frame
}

func (c clientFrame) playersDetails(team []models.GameTeam, gameMetadata GameMetadata) []models.GamePlayer {
	var players []models.GamePlayer

	for _, team := range team {
		teamMeta := c.findParticipantMetadata(team.Side, gameMetadata)
		for _, player := range teamMeta {
			players = append(players, models.GamePlayer{
				ExternalID: player.EsportsPlayerID,
				Role:       player.Role,
				Name:       player.SummonerName,
				Side:       team.Side,
				Champion:   player.ChampionID,
				TeamID:     team.ID,
			})
		}
	}

	return players
}

func (c clientFrame) findTeamMetadata(side string, frame GameFrame) TeamFrame {
	if side == "blue" {
		return frame.BlueTeam
	}
	return frame.RedTeam
}

func (c clientFrame) findParticipantMetadata(side string, metadata GameMetadata) []ParticipantMeta {
	if side == "blue" {
		return metadata.BlueTeamMetadata.ParticipantMetadata
	}
	return metadata.RedTeamMetadata.ParticipantMetadata
}

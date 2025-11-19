package types

import "time"

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

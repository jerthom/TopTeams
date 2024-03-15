package models

import (
	"time"
)

// Setting up a clock to enable us to consistently test
type clock interface {
	Now() time.Time
}

type realClock struct{}

func (c realClock) Now() time.Time {
	return time.Now()
}

var c clock = realClock{}

// PlayerInput represents a player as returned by the opendota API, filtered for the fields that we care about
type PlayerInput struct {
	TeamId          int       `json:"team_id"`
	Personaname     string    `json:"personaname"`
	FullHistoryTime time.Time `json:"full_history_time"`
	CountryCode     string    `json:"country_code"`
}

// PlayerOutput represents a player with the fields that we want to output
type PlayerOutput struct {
	Personaname string `yaml:"Personaname"`
	Experience  int    `yaml:"Player Experience"`
	CountryCode string `yaml:"Country Code"`
}

// Output transforms a PlayerInput struct into PlayerOutput;
// a player's experience is calculated as the number of seconds since their FullHistoryTime
func (pin PlayerInput) Output() PlayerOutput {
	// Use clock instead of time for reliable testing
	exp := c.Now().Sub(pin.FullHistoryTime).Seconds()

	return PlayerOutput{
		Personaname: pin.Personaname,
		Experience:  int(exp),
		CountryCode: pin.CountryCode,
	}
}

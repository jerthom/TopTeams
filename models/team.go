package models

// Team represents
type Team struct {
	Name       string         `json:"name" yaml:"Team Name"`
	Id         int            `json:"team_id" yaml:"Team Id"`
	Wins       int            `json:"wins" yaml:"Wins"`
	Losses     int            `json:"losses" yaml:"Losses"`
	Rating     float32        `json:"rating" yaml:"Rating"`
	Experience int            `yaml:"Team Experience"`
	Players    []PlayerOutput `yaml:"Players"`
}

// ComputeExperience calculates a team's experience as the sum of each of its player's experience
func (t *Team) ComputeExperience() {
	var sum int
	for _, p := range t.Players {
		sum += p.Experience
	}
	t.Experience = sum
}

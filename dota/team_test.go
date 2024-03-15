package dota

import (
	"testing"
)

func TestTeam_ComputeExperience(t *testing.T) {
	cases := []struct {
		desc    string
		players []PlayerOutput
		want    int
	}{{
		desc:    "empty players",
		players: nil,
		want:    0,
	}, {
		desc:    "single player",
		players: []PlayerOutput{{Experience: 100}},
		want:    100,
	}, {
		desc:    "multiple players",
		players: []PlayerOutput{{Experience: 100}, {Experience: 200}, {Experience: 300}},
		want:    600,
	}, {
		desc:    "zero experience for all players",
		players: []PlayerOutput{{Experience: 0}, {Experience: 0}, {Experience: 0}},
		want:    0,
	}}

	for _, tc := range cases {
		t.Run(tc.desc, func(tt *testing.T) {
			// Arrange
			team := &Team{Players: tc.players}

			// Act
			team.ComputeExperience()
			got := team.Experience

			// Assert
			if tc.want != got {
				tt.Errorf("mismatched experience\n\twant: %d\n\tgot: %d", tc.want, team.Experience)
			}
		})
	}

}

package dota

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestTopTeams(t *testing.T) {
	// Restore mocked values once we are done testing
	tmpPP := proPlayers
	tmpT := teams
	tmpC := c
	defer func() {
		proPlayers = tmpPP
		teams = tmpT
		c = tmpC
	}()

	var sampleTeams = map[int]*Team{
		0: {Id: 0},
		1: {Id: 1},
		2: {Id: 2},
	}

	// Simple mock for teams
	teams = func(ids map[int]struct{}) []*Team {
		var ret []*Team
		for id, _ := range ids {
			ret = append(ret, sampleTeams[id])
		}
		return ret
	}

	// Mock our clock
	fixedTime := time.Date(2024, 3, 13, 0, 0, 0, 0, time.UTC)
	c = mockClock{fixedTime}

	cases := []struct {
		desc    string
		players []PlayerInput
		pErr    error
		n       int
		want    []*Team
		wantErr error
	}{{
		desc:    "error getting proPlayers gets surfaced",
		players: nil,
		pErr:    errors.New(""),
		n:       0,
		want:    nil,
		wantErr: errors.New("unable to get players: "),
	}, {
		desc:    "no players give nil teams",
		players: nil,
		pErr:    nil,
		n:       10,
		want:    nil,
		wantErr: nil,
	}, {
		desc:    "zero n gives no teams",
		players: []PlayerInput{{TeamId: 1, FullHistoryTime: fixedTime}},
		pErr:    nil,
		n:       0,
		want:    []*Team{},
		wantErr: nil,
	}, {
		desc:    "single player",
		players: []PlayerInput{{TeamId: 1, FullHistoryTime: fixedTime.Add(-1 * time.Second)}},
		pErr:    nil,
		n:       1,
		want:    []*Team{{Id: 1, Experience: 1, Players: []PlayerOutput{{Experience: 1}}}},
		wantErr: nil,
	}, {
		desc: "multiple players on one team",
		players: []PlayerInput{
			{TeamId: 1, FullHistoryTime: fixedTime.Add(-1 * time.Second)},
			{TeamId: 1, FullHistoryTime: fixedTime.Add(-1 * time.Second)},
		},
		pErr:    nil,
		n:       1,
		want:    []*Team{{Id: 1, Experience: 2, Players: []PlayerOutput{{Experience: 1}, {Experience: 1}}}},
		wantErr: nil,
	}, {
		desc: "n filters down multiple teams to desired size, getting top by TeamId",
		players: []PlayerInput{
			{TeamId: 1, FullHistoryTime: fixedTime.Add(-1 * time.Second)},
			{TeamId: 2, FullHistoryTime: fixedTime.Add(-1 * time.Second)},
		},
		pErr:    nil,
		n:       1,
		want:    []*Team{{Id: 1, Experience: 1, Players: []PlayerOutput{{Experience: 1}}}},
		wantErr: nil,
	}, {
		desc: "multiple teams get sorted by team experience",
		players: []PlayerInput{
			{TeamId: 1, FullHistoryTime: fixedTime.Add(-1 * time.Second)},
			{TeamId: 1, FullHistoryTime: fixedTime.Add(-2 * time.Second)},
			{TeamId: 2, FullHistoryTime: fixedTime.Add(-5 * time.Second)},
		},
		pErr: nil,
		n:    2,
		want: []*Team{
			{Id: 2, Experience: 5, Players: []PlayerOutput{{Experience: 5}}},
			{Id: 1, Experience: 3, Players: []PlayerOutput{{Experience: 1}, {Experience: 2}}},
		},
		wantErr: nil,
	}, {
		desc:    "n larger than present teams has no impact",
		players: []PlayerInput{{TeamId: 1, FullHistoryTime: fixedTime.Add(-1 * time.Second)}},
		pErr:    nil,
		n:       100,
		want:    []*Team{{Id: 1, Experience: 1, Players: []PlayerOutput{{Experience: 1}}}},
		wantErr: nil,
	}}

	for _, tc := range cases {
		t.Run(tc.desc, func(tt *testing.T) {
			// Arrange
			proPlayers = func() ([]PlayerInput, error) { return tc.players, tc.pErr }

			// Act
			got, gotErr := TopTeams(tc.n)

			// Assert
			if !reflect.DeepEqual(tc.wantErr, gotErr) {
				tt.Errorf("mismatched errors\n\twant: %v\n\tgot: %v", tc.wantErr, gotErr)
			}

			if !reflect.DeepEqual(tc.want, got) {
				tt.Errorf("mismatched team slices\n\twant: %s\n\tgot: %s", jsonMarshal(tc.want), jsonMarshal(got))
			}

		})
	}
}

func TestTeams(t *testing.T) {
	// Restore mocked values once we're done testing
	tmpT := team
	defer func() {
		team = tmpT
	}()

	var sampleTeams = map[int]*Team{
		0: {Id: 0},
		1: {Id: 1},
	}

	// Simple mock for team
	team = func(id int) (*Team, error) {
		// Use negatives to simulate any error from team()
		if id < 0 {
			return nil, errors.New("")
		}
		return sampleTeams[id], nil
	}

	cases := []struct {
		desc  string
		input map[int]struct{}
		want  []*Team
	}{{
		desc:  "empty set",
		input: map[int]struct{}{},
		want:  nil,
	}, {
		desc:  "errors from team() get skipped",
		input: map[int]struct{}{1: {}, -1: {}},
		want:  []*Team{{Id: 1}},
	}, {
		desc:  "teams with TeamId=0 get skipped",
		input: map[int]struct{}{1: {}, 0: {}},
		want:  []*Team{{Id: 1}},
	}}

	for _, tc := range cases {
		t.Run(tc.desc, func(tt *testing.T) {
			// Act
			got := teams(tc.input)

			// Assert
			if !reflect.DeepEqual(tc.want, got) {
				tt.Errorf("mismatched team slices\n\twant: %s\n\tgot: %s", jsonMarshal(tc.want), jsonMarshal(got))
			}
		})
	}

}

// This helper discards the json.Marshal error so we can marshal in our tests easily
func jsonMarshal(i interface{}) []byte {
	b, _ := json.Marshal(i)
	return b
}

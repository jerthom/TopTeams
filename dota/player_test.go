package dota

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

type mockClock struct {
	FixedTime time.Time
}

func (c mockClock) Now() time.Time {
	return c.FixedTime
}

func TestPlayerInput_Output(t *testing.T) {
	// Restore mocked values once we're done testing
	tmpC := c
	defer func() {
		c = tmpC
	}()

	// Mock our clock
	fixedTime := time.Date(2024, 3, 13, 0, 0, 0, 0, time.UTC)
	c = mockClock{fixedTime}

	cases := []struct {
		desc string
		test PlayerInput
		want PlayerOutput
	}{{
		desc: "experience calculation correct for 24 hours ago",
		test: PlayerInput{
			TeamId:          1,
			Personaname:     "A",
			FullHistoryTime: fixedTime.Add(-24 * time.Hour),
			CountryCode:     "us",
		},
		want: PlayerOutput{
			Personaname: "A",
			Experience:  86400, // 24 hours in seconds
			CountryCode: "us",
		},
	}, {
		desc: "experience calculation correct for 48 hours ago",
		test: PlayerInput{
			TeamId:          2,
			Personaname:     "B",
			FullHistoryTime: fixedTime.Add(-48 * time.Hour),
			CountryCode:     "cn",
		},
		want: PlayerOutput{
			Personaname: "B",
			Experience:  172800, // 48 hours in seconds
			CountryCode: "cn",
		},
	}, {
		desc: "experience calculation correct for zero time",
		test: PlayerInput{
			TeamId:          3,
			Personaname:     "NoHistory",
			FullHistoryTime: fixedTime,
			CountryCode:     "ru",
		},
		want: PlayerOutput{
			Personaname: "NoHistory",
			Experience:  0,
			CountryCode: "ru",
		},
	}}

	for _, tc := range cases {
		t.Run(tc.desc, func(tt *testing.T) {
			// Act
			got := tc.test.Output()

			// Assert
			if !reflect.DeepEqual(tc.want, got) {
				w, _ := json.Marshal(tc.want)
				g, _ := json.Marshal(got)
				tt.Errorf("mismatched output\n\twant: %s\n\tgot: %s", w, g)
			}
		})
	}
}

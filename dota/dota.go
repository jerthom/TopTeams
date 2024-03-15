package dota

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"topTeams/models"
)

// TopTeams provides
func TopTeams(n int) ([]*models.Team, error) {
	players, err := proPlayers()
	if err != nil {
		return nil, fmt.Errorf("unable to get players: %v", err)
	}
	teamsToPlayers := map[int][]models.PlayerOutput{}
	// A set for Team Ids, so we don't make duplicate calls for the same team
	teamIds := map[int]struct{}{}
	for _, p := range players {
		if p.TeamId == 0 {
			continue
		}
		teamIds[p.TeamId] = struct{}{}
		teamsToPlayers[p.TeamId] = append(teamsToPlayers[p.TeamId], p.Output())
	}

	ts := teams(teamIds)

	// Take top N teams, defined by TeamId
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Id < ts[j].Id
	})
	if n < len(ts) {
		ts = ts[:n]
	}

	// Add appropriate teams to each team and calculate team experience
	for _, t := range ts {
		t.Players = teamsToPlayers[t.Id]
		t.ComputeExperience()
	}

	// Sort the teams by experience
	sort.Slice(ts, func(i, j int) bool {
		return ts[i].Experience > ts[j].Experience
	})

	return ts, nil
}

var proPlayers = func() ([]models.PlayerInput, error) {
	req, err := http.NewRequest("GET", "https://api.opendota.com/api/proPlayers", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting from endpoint: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var players []models.PlayerInput
	err = json.Unmarshal(body, &players)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %v", err)
	}

	return players, nil
}

// teams concurrently fetches
var teams = func(ids map[int]struct{}) []*models.Team {
	teamChan := make(chan *models.Team, len(ids))
	var wg sync.WaitGroup
	for id := range ids {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			t, err := team(id)
			// Exclude teams we can't get
			if err != nil {
				return
			}
			teamChan <- t
		}(id)
	}

	// Wait for all goroutines to finish fetching teams
	wg.Wait()
	close(teamChan)

	var ts []*models.Team
	for t := range teamChan {
		if t.Id != 0 {
			ts = append(ts, t)
		}
	}

	return ts
}

var team = func(id int) (*models.Team, error) {
	url := "https://api.opendota.com/api/teams/" + strconv.Itoa(id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting from endpoint: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var t models.Team
	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %v", err)
	}

	return &t, nil
}

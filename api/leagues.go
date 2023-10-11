package api

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	leaguesPath = "/leagues"
	// LeagueTypeParamLeague can be passed as query `Type` parameter to return only leagues championships.
	LeagueTypeParamLeague LeagueTypeParam = "league"
	// LeagueTypeParamCup can be passed as query `Type` parameter to return only cups.
	LeagueTypeParamCup LeagueTypeParam = "cup"
)

// LeagueTypeParam represent the type of the league.
type LeagueTypeParam string

// LeaguesQueryParams represents the parameters to pass to the /leagues endpoint.
// validate tags are for the go-playground/validator.
// url tags are for google/go-querystring.
// `validate:"omitempty," url:",omitempty"`.
type LeaguesQueryParams struct {
	ID      int             `validate:"omitempty,gte=0" url:"id,omitempty"`
	Name    string          `validate:"omitempty,min=1" url:"name,omitempty"`
	Country string          `validate:"omitempty,min=1" url:"country,omitempty"`
	Code    string          `validate:"omitempty,len=2" url:"code,omitempty"`
	Season  int             `validate:"omitempty,gte=1000,lte=9999" url:"season,omitempty"`
	Team    int             `validate:"omitempty,gte=0" url:"team,omitempty"`
	Type    LeagueTypeParam `validate:"omitempty" url:"type,omitempty"`
	Current bool            `validate:"omitempty" url:"current,omitempty"`
	Search  string          `validate:"omitempty,min=3" url:"search,omitempty"`
	Last    int             `validate:"omitempty,lte=99" url:"last,omitempty"`
}

// League wraps league top objects.
type League struct {
	LeagueInfo LeagueInfo   `json:"league"`
	Country    Country      `json:"country"`
	Seasons    []SeasonInfo `json:"seasons"`
}

// LeagueInfo wraps basic information on the league.
type LeagueInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Logo string `json:"logo"`
}

// SeasonInfo wraps basic information on the league's season.
type SeasonInfo struct {
	Year     int      `json:"year"`
	Start    string   `json:"start"`
	End      string   `json:"end"`
	Current  bool     `json:"current"`
	Coverage Coverage `json:"coverage"`
}

// Coverage wraps information on the league's season data coverage.
type Coverage struct {
	Standings   bool             `json:"standings"`
	Players     bool             `json:"players"`
	TopScorers  bool             `json:"top_scorers"`
	TopAssists  bool             `json:"top_assists"`
	TopCards    bool             `json:"top_cards"`
	Injuries    bool             `json:"injuries"`
	Predictions bool             `json:"predictions"`
	Odds        bool             `json:"odds"`
	Fixtures    FixturesCoverage `json:"fixtures"`
}

// FixturesCoverage wraps specific information on fixtures coverage.
type FixturesCoverage struct {
	Events             bool `json:"events"`
	Lineups            bool `json:"lineups"`
	StatisticsFixtures bool `json:"statistics_fixtures"`
	StatisticsPlayers  bool `json:"statistics_players"`
}

// LeaguesResult wraps the api raw response as well as the list of leagues.
type LeaguesResult struct {
	*ResponseOK
	Leagues []League `json:"leagues"`
}

// Leagues is the main function to request the /leagues endpoint.
// params *LeaguesQueryParams can be passed as optional request parameters, nil is accepted if there are no parameters to provide.
func (c *Client) Leagues(ctx context.Context, params *LeaguesQueryParams) (*LeaguesResult, error) {
	logger := c.logger

	req, err := buildQuery(ctx, c, leaguesPath, params)
	if err != nil {
		logger.ErrorContext(ctx, "error while building query")

		return nil, err
	}

	apiResp, err := executeQuery(ctx, c, req)
	if err != nil {
		logger.ErrorContext(ctx, "error while executing query")

		return nil, err
	}

	var ret LeaguesResult
	ret.ResponseOK = apiResp

	leagues := []League{}

	if err := json.Unmarshal(ret.Response, &leagues); err != nil {
		logger.ErrorContext(ctx, "error while parsing response field")

		return nil, fmt.Errorf("error while parsing response field: %w", err)
	}

	ret.Leagues = leagues

	return &ret, nil
}

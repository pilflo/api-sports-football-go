package api

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	teamsInformationPath = "/teams"
)

// TeamsInformationQueryParams is a struct for wrapping teams infos endpoint query parameters.
type TeamsInformationQueryParams struct {
	ID      *int    `validate:"omitempty,gte=0" url:"id,omitempty"`
	Name    *string `validate:"omitempty,min=1" url:"name,omitempty"`
	Country *string `validate:"omitempty,min=1" url:"country,omitempty"`
	Season  *int    `validate:"omitempty,gte=1000,lte=9999" url:"season,omitempty"`
	Search  *string `validate:"omitempty,min=3" url:"search,omitempty"`
	League  *int    `validate:"omitempty,gte=0" url:"league,omitempty"`
	Code    *string `validate:"omitempty,len=3" url:"code,omitempty"`
	Venue   *int    `validate:"omitempty,gte=0" url:"venue,omitempty"`
}

// TeamInformation wraps a team top objects.
type TeamInformation struct {
	Team  Team  `json:"team"`
	Venue Venue `json:"venue"`
}

// Team wraps basic information on the team.
type Team struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Country  string `json:"country"`
	Founded  int    `json:"founded"`
	National bool   `json:"national"`
	Logo     string `json:"logo"`
}

// Venue wraps basic information on the team's venue.
type Venue struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Capacity int    `json:"capacity"`
	Surface  string `json:"surface"`
	Image    string `json:"image"`
}

// TeamsInformationResult wraps the api raw response as well as the list of teams.
type TeamsInformationResult struct {
	*ResponseOK
	Teams []TeamInformation `json:"teams"`
}

// TeamsInformation is the main function to request the /teams endpoint.
// params *TeamsInformationQueryParams can be passed as optional request parameters, nil is accepted if there are no parameters to provide.
func (c *Client) TeamsInformation(ctx context.Context, params *TeamsInformationQueryParams) (*TeamsInformationResult, error) {
	logger := c.logger

	req, err := buildQuery(ctx, c, teamsInformationPath, params)
	if err != nil {
		logger.ErrorContext(ctx, "error while building query")

		return nil, err
	}

	apiResp, err := executeQuery(ctx, c, req)
	if err != nil {
		logger.ErrorContext(ctx, "error while executing query")

		return nil, err
	}

	var ret TeamsInformationResult
	ret.ResponseOK = apiResp

	teams := []TeamInformation{}

	if err := json.Unmarshal(ret.Response, &teams); err != nil {
		logger.ErrorContext(ctx, "error while parsing response field")

		return nil, fmt.Errorf("error while parsing response field: %w", err)
	}

	ret.Teams = teams

	return &ret, nil
}

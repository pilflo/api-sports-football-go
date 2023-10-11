package api

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	countriesPath = "/countries"
)

// Country wraps basic information on the country.
type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Flag string `json:"flag"`
}

// CountriesQueryParams represents the parameters to pass to the /countries endpoint.
// validate tags are for the go-playground/validator.
// url tags are for google/go-querystring.
// `validate:"omitempty," url:",omitempty"`.
type CountriesQueryParams struct {
	Name   string `validate:"omitempty,min=1" url:"name,omitempty"`
	Code   string `validate:"omitempty,len=2" url:"code,omitempty"`
	Search string `validate:"omitempty,min=3" url:"search,omitempty"`
}

// CountriesResult wraps the api raw response as well as the list of countries.
type CountriesResult struct {
	*ResponseOK
	Countries []Country
}

// Countries is the main function to request the /countries endpoint.
// params *CountriesQueryParams can be passed as optional request parameters, nil is accepted if there are no parameters to provide.
func (c *Client) Countries(ctx context.Context, params *CountriesQueryParams) (*CountriesResult, error) {
	logger := c.logger

	req, err := buildQuery(ctx, c, countriesPath, params)
	if err != nil {
		logger.ErrorContext(ctx, "error while building query")

		return nil, err
	}

	apiResp, err := executeQuery(ctx, c, req)
	if err != nil {
		logger.ErrorContext(ctx, "error while executing query")

		return nil, err
	}

	var ret CountriesResult
	ret.ResponseOK = apiResp

	countries := []Country{}

	if err := json.Unmarshal(ret.Response, &countries); err != nil {
		logger.ErrorContext(ctx, "error while parsing response field")

		return nil, fmt.Errorf("error while parsing response field: %w", err)
	}

	ret.Countries = countries

	return &ret, nil
}

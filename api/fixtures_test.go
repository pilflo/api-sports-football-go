package api_test

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/pilflo/api-sports-football-go/api"
	"github.com/pilflo/api-sports-football-go/api/mockserver"
	"github.com/stretchr/testify/assert"
)

type fixturesTestCase struct {
	params          *api.FixturesQueryParams
	jsonFilePath    string
	responseCode    int
	expectedResults int
}

func TestFixturesOK(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]fixturesTestCase{
		"team=33,season=2021": {
			params: &api.FixturesQueryParams{
				Team:   ptr(33),
				Season: ptr(2021),
			},
			jsonFilePath:    "./test_files/fixtures_33_2021.json",
			responseCode:    http.StatusOK,
			expectedResults: 54,
		},
		"team=37,season=2023": {
			params: &api.FixturesQueryParams{
				Team:   ptr(37),
				Season: ptr(2023),
			},
			jsonFilePath:    "./test_files/fixtures_37_2023.json",
			responseCode:    http.StatusOK,
			expectedResults: 51,
		},
	}
	server := mockserver.GetServer()

	client := api.NewClient(api.SubTypeAPISports).WithCustomAPIURL(server.URL)

	for _, tc := range tests {

		queryParams := &url.Values{}
		if tc.params.Team != nil {
			queryParams.Add("team", strconv.Itoa(*tc.params.Team))
		}
		if tc.params.Season != nil {
			queryParams.Add("season", strconv.Itoa(*tc.params.Season))
		}

		mockserver.AddJSONHandler(t, mockserver.MockJSONResponse{
			Path:         "/fixtures",
			QueryParams:  queryParams,
			ResponseCode: tc.responseCode,
			FilePath:     tc.jsonFilePath,
		})

		res, err := client.Fixtures(context.Background(), tc.params)

		assert.Nil(err)
		assert.NotNil(res)
		assert.Len(res.Fixtures, tc.expectedResults)
	}
}

func TestFixturesValidationErrors(t *testing.T) {
	tests := map[string]*api.FixturesQueryParams{
		"id negative":            {ID: ptr(-1)},
		"season incorrect range": {Season: ptr(666)},
		"team negative":          {Team: ptr(-1)},
		"last too big":           {Last: ptr(100)},
	}

	client := api.NewClient(api.SubTypeAPISports)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := client.Fixtures(context.Background(), tc)
			if got != nil {
				t.Fatalf("Expected result to be nil, got %v", got)
			}

			if gotErr == nil {
				t.Fatalf("Expected result NOT to be nil, got %v, message %v", gotErr, gotErr.Error())
			}

			if gotErr != nil {
				if !reflect.DeepEqual(reflect.TypeOf(gotErr).String(), "*api.FieldValidationError") {
					t.Fatalf("Expected a validation error, got %v", reflect.TypeOf(gotErr))
				}
			}
		})
	}
}

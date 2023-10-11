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
	expectedError   error
}

func TestFixturesOK(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]fixturesTestCase{
		"fixtures,team=33,season=2021": {
			params: &api.FixturesQueryParams{
				Team:   33,
				Season: 2021,
			},
			jsonFilePath:    "./test_files/fixtures_33_2021.json",
			responseCode:    http.StatusOK,
			expectedResults: 54,
		},
		"fixtures,team=37,season=2023": {
			params: &api.FixturesQueryParams{
				Team:   37,
				Season: 2023,
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
		if tc.params.Team > 0 {
			queryParams.Add("team", strconv.Itoa(tc.params.Team))
		}
		if tc.params.Season > 1000 {
			queryParams.Add("season", strconv.Itoa(tc.params.Season))
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

func TestFixturesAPIError(t *testing.T) {
	assert := assert.New(t)
	tests := map[string]fixturesTestCase{
		"id=33,live=all": {
			params: &api.FixturesQueryParams{
				ID:   1132381,
				Live: true,
			},
			jsonFilePath:    "./test_files/fixtures_error_live_id.json",
			responseCode:    http.StatusOK,
			expectedResults: 0,
			expectedError:   &api.ResponseError{},
		},
	}
	server := mockserver.GetServer()

	client := api.NewClient(api.SubTypeAPISports).WithCustomAPIURL(server.URL)

	for _, tc := range tests {

		queryParams := &url.Values{}
		if tc.params.ID > 0 {
			queryParams.Add("id", strconv.Itoa(tc.params.ID))
		}
		if tc.params.Live {
			queryParams.Add("live", "all")
		}

		mockserver.AddJSONHandler(t, mockserver.MockJSONResponse{
			Path:         "/fixtures",
			QueryParams:  queryParams,
			ResponseCode: tc.responseCode,
			FilePath:     tc.jsonFilePath,
		})

		res, err := client.Fixtures(context.Background(), tc.params)

		assert.Nil(res)
		assert.NotNil(err)
		assert.IsType(err, tc.expectedError)
	}
}

func TestFixturesValidationErrors(t *testing.T) {
	tests := map[string]*api.FixturesQueryParams{
		"id negative":            {ID: -1},
		"season incorrect range": {Season: 666},
		"team negative":          {Team: -1},
		"last too big":           {Last: 100},
	}

	client := api.NewClient(api.SubTypeAPISports).WithCustomAPIURL("http://test.com")

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

package api_test

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/pilflo/api-sports-football-go/api"
	"github.com/pilflo/api-sports-football-go/api/mockserver"
	"github.com/stretchr/testify/assert"
)

type leaguesTestCase struct {
	params          *api.LeaguesQueryParams
	jsonFilePath    string
	responseCode    int
	expectedResults int
}

func TestLeaguesOK(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]leaguesTestCase{
		"team=33,season=2021": {
			params: &api.LeaguesQueryParams{
				Country: ptr("FR"),
			},
			jsonFilePath:    "./test_files/leagues_fr.json",
			responseCode:    http.StatusOK,
			expectedResults: 23,
		},
	}
	server := mockserver.GetServer()

	client := api.NewClient(api.SubTypeAPISports).WithCustomAPIURL(server.URL)

	for _, tc := range tests {

		queryParams := &url.Values{}
		if tc.params.Country != nil {
			queryParams.Add("country", *tc.params.Country)
		}

		mockserver.AddJSONHandler(t, mockserver.MockJSONResponse{
			Path:         "/leagues",
			QueryParams:  queryParams,
			ResponseCode: tc.responseCode,
			FilePath:     tc.jsonFilePath,
		})

		res, err := client.Leagues(context.Background(), tc.params)

		assert.Nil(err)
		assert.NotNil(res)
		assert.Len(res.Leagues, tc.expectedResults)
	}
}

func TestLeaguesValidationErrors(t *testing.T) {
	tests := map[string]*api.LeaguesQueryParams{
		"id negative":            {ID: ptr(-1)},
		"name too short":         {Name: ptr("")},
		"country too short":      {Country: ptr("")},
		"code too short":         {Code: ptr("F")},
		"code too long":          {Code: ptr("FRA")},
		"season incorrect range": {Season: ptr(666)},
		"team negative":          {Team: ptr(-1)},
		"search too short":       {Search: ptr("FR")},
		"last too big":           {Last: ptr(100)},
	}

	client := api.NewClient(api.SubTypeAPISports).WithCustomAPIURL("http://test.com")

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := client.Leagues(context.Background(), tc)
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

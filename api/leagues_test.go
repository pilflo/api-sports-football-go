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

type leaguesTestCase struct {
	params             *api.LeaguesQueryParams
	jsonFilePath       string
	responseCode       int
	expectedResults    int
	expectedAttributes map[string]any
}

func TestLeaguesOK(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]leaguesTestCase{
		"leagues,country=france": {
			params: &api.LeaguesQueryParams{
				Country: "france",
			},
			jsonFilePath:    "./test_files/leagues_fr.json",
			responseCode:    http.StatusOK,
			expectedResults: 23,
			expectedAttributes: map[string]any{
				"ID":   61,
				"Name": "Ligue 1",
				"Type": "League",
				"Logo": "https://media.api-sports.io/football/leagues/61.png",
			},
		},
		"leagues,id=39": {
			params: &api.LeaguesQueryParams{
				ID: 39,
			},
			jsonFilePath:    "./test_files/leagues_id_39.json",
			responseCode:    http.StatusOK,
			expectedResults: 1,
			expectedAttributes: map[string]any{
				"ID":   39,
				"Name": "Premier League",
				"Type": "League",
				"Logo": "https://media-4.api-sports.io/football/leagues/39.png",
			},
		},
	}
	server := mockserver.GetServer()

	client := api.NewClient(api.SubTypeAPISports).WithCustomAPIURL(server.URL)

	for _, tc := range tests {

		queryParams := &url.Values{}
		if tc.params.Country != "" {
			queryParams.Add("country", tc.params.Country)
		}

		if tc.params.ID > 0 {
			queryParams.Add("id", strconv.Itoa(tc.params.ID))
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
		league := res.Leagues[0].LeagueInfo
		for k, v := range tc.expectedAttributes {
			val := reflect.ValueOf(league)
			field := val.FieldByName(k)
			assert.True(field.IsValid())
			assert.EqualValues(v, field.Interface())
		}
	}
}

func TestLeaguesValidationErrors(t *testing.T) {
	tests := map[string]*api.LeaguesQueryParams{
		"id negative":            {ID: -1},
		"code too short":         {Code: "F"},
		"code too long":          {Code: "FRA"},
		"season incorrect range": {Season: 666},
		"team negative":          {Team: -1},
		"search too short":       {Search: "FR"},
		"last too big":           {Last: 100},
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

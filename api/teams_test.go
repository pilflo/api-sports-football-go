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

type teamsTestCase struct {
	params             *api.TeamsInformationQueryParams
	jsonFilePath       string
	responseCode       int
	expectedResults    int
	expectedAttributes map[string]any
}

func TestTeamsInfoOK(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]teamsTestCase{
		"teams,country=france": {
			params: &api.TeamsInformationQueryParams{
				Country: ptr("france"),
			},
			jsonFilePath:    "./test_files/teams_fr.json",
			responseCode:    http.StatusOK,
			expectedResults: 894,
			expectedAttributes: map[string]any{
				"ID":       2,
				"Name":     "France",
				"Code":     "FRA",
				"Country":  "France",
				"Founded":  1919,
				"National": true,
				"Logo":     "https://media-4.api-sports.io/football/teams/2.png",
			},
		},
		"teams,id=42": {
			params: &api.TeamsInformationQueryParams{
				Country: ptr("france"),
			},
			jsonFilePath:    "./test_files/teams_id_42.json",
			responseCode:    http.StatusOK,
			expectedResults: 1,
			expectedAttributes: map[string]any{
				"ID":       42,
				"Name":     "Arsenal",
				"Code":     "ARS",
				"Country":  "England",
				"Founded":  1886,
				"National": false,
				"Logo":     "https://media-4.api-sports.io/football/teams/42.png",
			},
		},
	}
	server := mockserver.GetServer()

	client := api.NewClient(api.SubTypeAPISports).WithCustomAPIURL(server.URL)

	for _, tc := range tests {

		queryParams := &url.Values{}
		if tc.params.Country != nil {
			queryParams.Add("country", *tc.params.Country)
		}

		if tc.params.ID != nil {
			queryParams.Add("id", strconv.Itoa(*tc.params.ID))
		}

		mockserver.AddJSONHandler(t, mockserver.MockJSONResponse{
			Path:         "/teams",
			QueryParams:  queryParams,
			ResponseCode: tc.responseCode,
			FilePath:     tc.jsonFilePath,
		})

		res, err := client.TeamsInformation(context.Background(), tc.params)

		assert.Nil(err)
		assert.NotNil(res)
		assert.Len(res.Teams, tc.expectedResults)
		team := res.Teams[0].Team
		for k, v := range tc.expectedAttributes {
			val := reflect.ValueOf(team)
			field := val.FieldByName(k)
			assert.True(field.IsValid())
			assert.EqualValues(v, field.Interface())
		}
	}
}

func TestTeamsValidationErrors(t *testing.T) {
	tests := map[string]*api.TeamsInformationQueryParams{
		"id negative":            {ID: ptr(-1)},
		"league negative":        {League: ptr(-1)},
		"name too short":         {Name: ptr("")},
		"country too short":      {Country: ptr("")},
		"season incorrect range": {Season: ptr(666)},
		"search too short":       {Search: ptr("FR")},
	}

	client := api.NewClient(api.SubTypeAPISports).WithCustomAPIURL("http://test.com")

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := client.TeamsInformation(context.Background(), tc)
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

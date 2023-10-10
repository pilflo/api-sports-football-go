package api_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/pilflo/api-sports-football-go/api"
	"github.com/pilflo/api-sports-football-go/api/mockserver"
	"github.com/stretchr/testify/assert"
)

func TestCountriesOK(t *testing.T) {
	assert := assert.New(t)

	server := mockserver.GetServer()

	mockResponse := mockserver.MockJSONResponse{
		Path:         "/countries",
		ResponseCode: http.StatusOK,
		FilePath:     "./test_files/countries_all.json",
	}

	mockserver.AddJSONHandler(t, mockResponse)

	client := api.NewClient(api.SubTypeAPISports).WithCustomAPIURL(server.URL)

	resp, err := client.Countries(context.Background(), nil)

	assert.Nil(err)
	assert.Len(resp.Countries, 164)
}

func TestCountriesValidationErrors(t *testing.T) {
	tests := map[string]*api.CountriesQueryParams{
		"invalid code":     {Code: ptr("FRA")},
		"search too short": {Search: ptr("FR")},
	}

	client := api.NewClient(api.SubTypeAPISports)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, gotErr := client.Countries(context.Background(), tc)
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

func TestCountries2ValidationErrors(t *testing.T) {
	assert := assert.New(t)

	// httpClient is nil because it is not supposed to be used anyway.
	client := api.NewClient(api.SubTypeAPISports)

	wrongSizeCode := "FRA"

	query := &api.CountriesQueryParams{
		Code: &wrongSizeCode,
	}

	resp, err := client.Countries(context.Background(), query)

	assert.Nil(resp)
	assert.NotNil(err)

	emptyName := ""

	query = &api.CountriesQueryParams{
		Name: &emptyName,
	}

	resp, err = client.Countries(context.Background(), query)

	assert.Nil(resp)
	assert.NotNil(err)

	tooShortSearch := "fr"

	query = &api.CountriesQueryParams{
		Search: &tooShortSearch,
	}

	resp, err = client.Countries(context.Background(), query)

	assert.Nil(resp)
	assert.NotNil(err)
}

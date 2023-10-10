package api_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/pilflo/api-sports-football-go/api/mockserver"
)

func TestClientOK(t *testing.T) {
	server := mockserver.GetServer()
	mockResponse := mockserver.MockJSONResponse{
		Path:         "/",
		ResponseCode: http.StatusOK,
		FilePath:     "./test_files/countries_all.json",
	}

	mockserver.AddJSONHandler(t, mockResponse)

	client := http.DefaultClient

	req, _ := http.NewRequest("GET", server.URL+mockResponse.Path, nil)

	res, _ := client.Do(req)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("unexpected error when reading body %s", err.Error())
	}
	fmt.Println(string(body))
}

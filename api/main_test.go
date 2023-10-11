package api_test

import (
	"os"
	"testing"

	"github.com/pilflo/api-sports-football-go/api/mockserver"
)

func TestMain(m *testing.M) {
	_ = mockserver.GetServer()

	os.Setenv("API_SPORTS_KEY", "abcdef12345")

	exitCode := m.Run()

	os.Exit(exitCode)
}

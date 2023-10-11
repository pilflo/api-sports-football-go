package main

import (
	"context"
	"log"

	sports "github.com/pilflo/api-sports-football-go/api"
)

func main() {
	client := sports.NewClient(sports.SubTypeAPISports)
	ctx := context.Background()
	resC, err := client.Countries(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("There are %v countries.", len(resC.Countries))

	resF, err := client.Fixtures(ctx, &sports.FixturesQueryParams{
		Live: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("There are %v live fixtures.", len(resF.Fixtures))

	resL, err := client.Leagues(ctx, &sports.LeaguesQueryParams{
		ID:      39,
		Current: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	league := resL.Leagues[0]
	log.Printf("There current season of %v ends on %v", league.LeagueInfo.Name, league.Seasons[0].End)

	resT, err := client.TeamsInformation(ctx, &sports.TeamsInformationQueryParams{
		ID: 42,
	})
	if err != nil {
		log.Fatal(err)
	}

	teamInfo := resT.Teams[0]
	log.Printf("%v football team was founded in %v and plays at %v", teamInfo.Team.Name, teamInfo.Team.Founded, teamInfo.Venue.Name)
}

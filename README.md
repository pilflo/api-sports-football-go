# api-sports-football-go
[![API Sports Football Go](https://github.com/pilflo/api-sports-football-go/actions/workflows/ci.yml/badge.svg)](https://github.com/pilflo/api-sports-football-go/actions/workflows/ci.yml)

A golang library for api-sports.io football API

https://www.api-football.com/documentation-v3

## API Keys

An API Key is required in order to perform API calls.  
If you have a subscription to `API-Sports`, use the environment variable `API_SPORTS_KEY`.  
If you have a subscription to `RapidAPI`, use the environment variable `RAPID_API_KEY`.  
```bash
export API_SPORTS_KEY=abdef123xxxxxxxxxxxx45ghijk
# OR
export RAPID_API_KEY=abdef123xxxxxxxxxxxx45ghijk
```


## Usage

```go
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
```

## Development

Before each pull request, make sure that all the steps (imports, format, lint, test) are successfull.  

```bash
# Requires GNU Make and Docker
make check
```

## Coverage

| ENDPOINT  | COVERAGE 
|--|--
| /timezone | ❌
| /countries | ✅
| /leagues | ✅
| /leagues/seasons | ❌
| /teams | ✅
| /teams/statistics | ❌
| /teams/seasons | ❌
| /teams/countries | ❌
| /venues | ❌
| /standings | ❌
| /fixtures/ | ✅
| /fixtures/rounds | ❌
| /fixtures/headtohead | ❌
| /fixtures/statistics | ❌
| /fixtures/events | ❌
| /fixtures/lineups | ❌
| /fixtures/players | ❌
| /injuries | ❌
| /predictions | ❌
| /coachs | ❌
| /players | ❌
| /players/seasons | ❌
| /players/squads | ❌
| /players/topscorers | ❌
| /players/topassists | ❌
| /players/topyellowcards | ❌
| /players/topredcards | ❌
| /transfers | ❌
| /trophies | ❌
| /sidelined | ❌
| /odds | ❌
| /odds/mappings | ❌
| /odds/bookmakers | ❌
| /odds/bets| ❌
| /odds/live | ❌
| /odds/live/bets | ❌

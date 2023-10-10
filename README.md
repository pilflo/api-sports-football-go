# api-sports-football-go

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

}
```

## Development

Before each pull request, make sure that all the steps (imports, format, lint, test) are successfull.  

```bash
# Requires GNU Make and Docker
make check
```

## Coverage


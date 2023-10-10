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

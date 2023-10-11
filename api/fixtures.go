package api

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FixtureStatusType represents the status of the game.
type FixtureStatusType string

// FixtureStatisticsType represents type of statistics for a fixture.
type FixtureStatisticsType string

// FixtureEventType represents type of event for a fixture.
type FixtureEventType string

// FixtureLineupsType represents type of linepus for a fixture.
type FixtureLineupsType string

const (
	fixturesPath = "/fixtures"
	// FixtureStatusTBD : Time To Be Defined.
	FixtureStatusTBD FixtureStatusType = "TBD"
	// FixtureStatusNS : Not Started.
	FixtureStatusNS FixtureStatusType = "NS"
	// FixtureStatus1H : First Half, Kick Off.
	FixtureStatus1H FixtureStatusType = "1H"
	// FixtureStatusHT : Halftime.
	FixtureStatusHT FixtureStatusType = "HT"
	// FixtureStatus2H : Second Half, 2nd Half Started.
	FixtureStatus2H FixtureStatusType = "2H"
	// FixtureStatusET : Extra Time.
	FixtureStatusET FixtureStatusType = "ET"
	// FixtureStatusP : Penalty In Progress.
	FixtureStatusP FixtureStatusType = "P"
	// FixtureStatusFT : Match Finished.
	FixtureStatusFT FixtureStatusType = "FT"
	// FixtureStatusAET : Match Finished After Extra Time.
	FixtureStatusAET FixtureStatusType = "AET"
	// FixtureStatusPEN : Match Finished After Penalty.
	FixtureStatusPEN FixtureStatusType = "PEN"
	// FixtureStatusBT : Break Time (in Extra Time).
	FixtureStatusBT FixtureStatusType = "BT"
	// FixtureStatusSUSP : Match Suspended.
	FixtureStatusSUSP FixtureStatusType = "SUSP"
	// FixtureStatusINT : Match Interrupted.
	FixtureStatusINT FixtureStatusType = "INT"
	// FixtureStatusPST : Match Postponed.
	FixtureStatusPST FixtureStatusType = "PST"
	// FixtureStatusCANC : Match Cancelled.
	FixtureStatusCANC FixtureStatusType = "CANC"
	// FixtureStatusABD : Match Abandoned.
	FixtureStatusABD FixtureStatusType = "ABD"
	// FixtureStatusAWD : Technical Loss.
	FixtureStatusAWD FixtureStatusType = "AWD"
	// FixtureStatusWO : WalkOver.
	FixtureStatusWO FixtureStatusType = "WO"
	// FixtureStatisticsTypeShotsOnGoal : "Shots on Goal".
	FixtureStatisticsTypeShotsOnGoal FixtureStatisticsType = "Shots on Goal"
	// FixtureStatisticsTypeShotsOffGoal : "Shots off Goal".
	FixtureStatisticsTypeShotsOffGoal FixtureStatisticsType = "Shots off Goal"
	// FixtureStatisticsTypeTotalShots : "Total Shots".
	FixtureStatisticsTypeTotalShots FixtureStatisticsType = "Total Shots"
	// FixtureStatisticsTypeBlockedShots : "Blocked Shots".
	FixtureStatisticsTypeBlockedShots FixtureStatisticsType = "Blocked Shots"
	// FixtureStatisticsTypeShotsInsidebox : "Shots insidebox".
	FixtureStatisticsTypeShotsInsidebox FixtureStatisticsType = "Shots insidebox"
	// FixtureStatisticsTypeShotsOutsidebox : "Shots outsidebox".
	FixtureStatisticsTypeShotsOutsidebox FixtureStatisticsType = "Shots outsidebox"
	// FixtureStatisticsTypeFouls : "Fouls".
	FixtureStatisticsTypeFouls FixtureStatisticsType = "Fouls"
	// FixtureStatisticsTypeCornerKicks : "Corner Kicks".
	FixtureStatisticsTypeCornerKicks FixtureStatisticsType = "Corner Kicks"
	// FixtureStatisticsTypeOffsides : "Offsides".
	FixtureStatisticsTypeOffsides FixtureStatisticsType = "Offsides"
	// FixtureStatisticsTypeBallPossession : "Ball Possession".
	FixtureStatisticsTypeBallPossession FixtureStatisticsType = "Ball Possession"
	// FixtureStatisticsTypeYellowCards : "Yellow Cards".
	FixtureStatisticsTypeYellowCards FixtureStatisticsType = "Yellow Cards"
	// FixtureStatisticsTypeRedCards : "Red Cards".
	FixtureStatisticsTypeRedCards FixtureStatisticsType = "Red Cards"
	// FixtureStatisticsTypeGoalkeeperSaves : "Goalkeeper Saves".
	FixtureStatisticsTypeGoalkeeperSaves FixtureStatisticsType = "Goalkeeper Saves"
	// FixtureStatisticsTypeTotalPasses : "Total passes".
	FixtureStatisticsTypeTotalPasses FixtureStatisticsType = "Total passes"
	// FixtureStatisticsTypePassesAccurate : "Passes accurate".
	FixtureStatisticsTypePassesAccurate FixtureStatisticsType = "Passes accurate"
	// FixtureStatisticsTypePassesPct : "Passes %".
	FixtureStatisticsTypePassesPct FixtureStatisticsType = "Passes %"
	// FixtureEventTypeGoal : event type goal.
	FixtureEventTypeGoal FixtureEventType = "goal"
	// FixtureEventTypeCard : event type card.
	FixtureEventTypeCard FixtureEventType = "card"
	// FixtureEventTypeSubst : event type substitution.
	FixtureEventTypeSubst FixtureEventType = "subst"
	// FixtureLineupsTypeFormation : event type formation.
	FixtureLineupsTypeFormation FixtureLineupsType = "formation"
	// FixtureLineupsTypeCoach : event type coach.
	FixtureLineupsTypeCoach FixtureLineupsType = "coach"
	// FixtureLineupsTypeStartXI : event type startxi.
	FixtureLineupsTypeStartXI FixtureLineupsType = "startxi"
	// FixtureLineupsTypeSubstitutes : event type substitutes.
	FixtureLineupsTypeSubstitutes FixtureLineupsType = "substitutes"
)

// FixturesQueryParams represents the parameters to pass to the /fixtures endpoint.
type FixturesQueryParams struct {
	ID          int
	IDs         []int
	Live        bool
	LiveLeagues []int
	Date        time.Time
	League      int
	Season      int
	Team        int
	Last        int
	Next        int
	From        time.Time
	To          time.Time
	Round       string
	Status      FixtureStatusType
	Timezone    string
}

// fixturesQueryParams represents the parameters to pass to the /fixtures endpoint.
// validate tags are for the go-playground/validator.
// url tags are for google/go-querystring.
// `validate:"omitempty," url:",omitempty"`.
type fixturesQueryParams struct {
	ID       int               `validate:"omitempty,gte=0" url:"id,omitempty"`
	IDs      string            `validate:"omitempty" url:"ids,omitempty"`
	Live     string            `validate:"omitempty" url:"live,omitempty"`
	Date     time.Time         `validate:"omitempty" url:"date,omitempty" layout:"2006-01-02"`
	League   int               `validate:"omitempty,gte=0" url:"league,omitempty"`
	Season   int               `validate:"omitempty,gte=1000,lte=9999" url:"season,omitempty"`
	Team     int               `validate:"omitempty,gte=0" url:"team,omitempty"`
	Last     int               `validate:"omitempty,gte=0,lte=99" url:"last,omitempty"`
	Next     int               `validate:"omitempty,gte=0,lte=99" url:"next,omitempty"`
	From     time.Time         `validate:"omitempty" url:"from,omitempty" layout:"2006-01-02"`
	To       time.Time         `validate:"omitempty" url:"to,omitempty" layout:"2006-01-02"`
	Round    string            `validate:"omitempty,min=1" url:"round,omitempty"`
	Status   FixtureStatusType `validate:"omitempty,min=1" url:"status,omitempty"`
	Timezone string            `validate:"omitempty,min=1" url:"timezone,omitempty"`
}

// Fixture wraps league top objects.
type Fixture struct {
	FixtureInfo FixtureInfo       `json:"fixture"`
	LeagueInfo  FixtureLeagueInfo `json:"league"`
	Teams       FixtureTeams      `json:"teams"`
	Goals       FixtureGoals      `json:"goals"`
	Score       FixtureScore      `json:"score"`
}

// FixtureInfo wraps basic information on the fixture.
type FixtureInfo struct {
	ID        int           `json:"id"`
	Referee   string        `json:"referee"`
	Timezone  string        `json:"timezone"`
	Date      time.Time     `json:"date"`
	Timestamp int           `json:"timestamp"`
	Periods   Periods       `json:"periods"`
	Venue     FixtureVenue  `json:"venue"`
	Status    FixtureStatus `json:"status"`
}

// Periods represents timestamp for first and second period.
type Periods struct {
	First  int
	Second int
}

// FixtureVenue wraps basic info about the fixture's venue.
type FixtureVenue struct {
	ID   int
	Name string
	City string
}

// FixtureStatus represents the current status of the fixture.
type FixtureStatus struct {
	Long    string
	Short   string
	Elapsed int
}

// FixtureLeagueInfo wraps basic information on the fixture's league.
type FixtureLeagueInfo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Logo    string `json:"logo"`
	Flag    string `json:"flag"`
	Season  int    `json:"season"`
	Round   string `json:"round"`
}

// FixtureTeams wraps basic information on the fixture's teams.
type FixtureTeams struct {
	Home FixtureTeam `json:"home"`
	Away FixtureTeam `json:"away"`
}

// FixtureTeam wraps basic information on a fixture's team.
type FixtureTeam struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Winner bool   `json:"winner"`
}

// FixtureScore wraps basic information on the fixture's league.
type FixtureScore struct {
	Halftime  FixtureGoals `json:"halftime"`
	Fulltime  FixtureGoals `json:"fulltime"`
	Extratime FixtureGoals `json:"extratime"`
	Penalty   FixtureGoals `json:"penalty"`
}

// FixtureGoals wraps basic information on the fixture's goals.
// Goals can be nil when related to a period that did not happen.
type FixtureGoals struct {
	Home *int `json:"home"`
	Away *int `json:"away"`
}

// FixturesResult wraps the api raw response as well as the list of leagues.
type FixturesResult struct {
	*ResponseOK
	Fixtures []Fixture `json:"fixtures"`
}

func translateParams(params *FixturesQueryParams) *fixturesQueryParams {
	ret := fixturesQueryParams{
		ID:       params.ID,
		Date:     params.Date,
		League:   params.League,
		Season:   params.Season,
		Team:     params.Team,
		Last:     params.Last,
		Next:     params.Next,
		From:     params.From,
		To:       params.To,
		Round:    params.Round,
		Status:   params.Status,
		Timezone: params.Timezone,
	}

	liveStr := "all"

	if params.Live {
		if params.LiveLeagues != nil {
			liveLeagues := params.LiveLeagues
			switch len(liveLeagues) {
			case 0:
				// Do nothing, keep Live param to 'all' and get every live match.
			case 1:
				// Keep Live param to 'all' and sort with the single league.
				// Overrides a potential League param.
				ret.League = liveLeagues[0]
			default:
				// Set Live param to id1-id2-id3... league ids.
				liveStr = arrayToString(liveLeagues, "-")
			}
		}

		ret.Live = liveStr
	}

	if params.IDs != nil {
		ids := params.IDs
		switch len(ids) {
		case 0:
			// Ignore if the array is empty.
		case 1:
			// If only a single
			ret.IDs = strconv.Itoa(ids[0])
		default:
			// Set Ids param to id1-id2-id3... fixture ids.
			ret.IDs = arrayToString(ids, "-")
		}
	}

	return &ret
}

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.ReplaceAll(fmt.Sprint(a), " ", delim), "[]")
}

// Fixtures is the main function to request the /fixtures endpoint.
// params *FixturesQueryParams can be passed as optional request parameters, nil is accepted if there are no parameters to provide.
func (c *Client) Fixtures(ctx context.Context, params *FixturesQueryParams) (*FixturesResult, error) {
	logger := c.logger

	formattedParams := translateParams(params)

	req, err := buildQuery(ctx, c, fixturesPath, formattedParams)
	if err != nil {
		logger.ErrorContext(ctx, "error while building query")

		return nil, err
	}

	apiResp, err := executeQuery(ctx, c, req)
	if err != nil {
		logger.ErrorContext(ctx, "error while executing query")

		return nil, err
	}

	var ret FixturesResult
	ret.ResponseOK = apiResp

	fixtures := []Fixture{}

	if err := json.Unmarshal(ret.Response, &fixtures); err != nil {
		logger.ErrorContext(ctx, "error while parsing response field")

		return nil, fmt.Errorf("error while parsing response field: %w", err)
	}

	ret.Fixtures = fixtures

	return &ret, nil
}

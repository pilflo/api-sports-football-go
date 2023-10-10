package api

// Config defaults to APISports.

// config is a wrapping struct for Client configuration.
type config struct {
	// The type of subscription to the service.
	// Can be APISports or RapidAPI.
	// Defaults to APISports.
	subType          SubscriptionType
	basePath         string
	apiKeyEnvVar     string
	apiKeyHTTPHeader string
}

func newConfig(t SubscriptionType) config {
	if t == SubTypeRapidAPI {
		return config{
			subType:          SubTypeRapidAPI,
			basePath:         "https://api-football-v1.p.rapidapi.com/v3",
			apiKeyEnvVar:     "RAPID_API_KEY",
			apiKeyHTTPHeader: "x-rapidapi-key",
		}
	}

	return config{
		subType:          SubTypeAPISports,
		basePath:         "https://v3.football.api-sports.io",
		apiKeyEnvVar:     "API_SPORTS_KEY",
		apiKeyHTTPHeader: "x-apisports-key",
	}
}

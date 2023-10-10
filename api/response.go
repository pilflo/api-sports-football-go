package api

import "fmt"

type apiResponseRaw struct {
	Get        string         `json:"get"`
	Parameters any            `json:"parameters"`
	Errors     any            `json:"errors"`
	Results    int            `json:"results"`
	Paging     map[string]int `json:"paging"`
	Response   any            `json:"response"`
}

// ResponseOK represents a valid response from the server.
type ResponseOK struct {
	Get        string         `json:"get"`
	Parameters map[string]any `json:"parameters"`
	Errors     map[string]any `json:"errors"`
	Results    int            `json:"results"`
	Paging     map[string]int `json:"paging"`
	Response   []byte         `json:"response"`
}

// ResponseError represents an invalid response from the server.
type ResponseError struct {
	Message string `json:"message"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("Error(s) from API : %s\n", e.Message)
}

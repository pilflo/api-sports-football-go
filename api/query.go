package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-querystring/query"
)

// baseQuery creates a basic http.Request to the API and adds headers.
func baseQuery(ctx context.Context, client *Client, queryPath string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.config.basePath+queryPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new http request: %w", err)
	}

	return req, nil
}

// buildQuery builds the request to send to the API.
// nil is accepted for params if there is no query parameters.
// path should have the format : '/path/to/endpoint'.
func buildQuery(ctx context.Context, client *Client, path string, params any) (*http.Request, error) {
	logger := client.logger

	req, err := baseQuery(ctx, client, path)
	if err != nil {
		logger.ErrorContext(ctx, "error while building base query")

		return nil, err
	}

	if reflect.ValueOf(params).IsNil() {
		return req, nil
	}

	if err = validateQueryParams(params); err != nil {
		logger.ErrorContext(ctx, "error while validating query parameters")

		return nil, err
	}

	if err = addQueryParams(req, params); err != nil {
		logger.ErrorContext(ctx, "error while adding query parameters")

		return nil, err
	}

	return req, nil
}

// validateQueryParams validates the query parameter according to the validate tags.
// no-op if params is nil.
func validateQueryParams(params any) error {
	if params == nil {
		return nil
	}

	validate := validator.New()

	if err := validate.Struct(params); err != nil {
		return newFieldValidationError(err)
	}

	return nil
}

// addQueryParams adds query parameters to the request.
// no-op if params is nil.
func addQueryParams(req *http.Request, params any) error {
	if params == nil {
		return nil
	}

	v, err := query.Values(params)
	if err != nil {
		return fmt.Errorf("failed to generate url query values: %w", err)
	}

	req.URL.RawQuery = v.Encode()

	return nil
}

// executeQuery runs the pre-built request and handles the http response.
func executeQuery(ctx context.Context, c *Client, req *http.Request) (*ResponseOK, error) {
	logger := c.logger

	httpClient := c.httpClient

	res, err := httpClient.Do(req)
	if err != nil {
		logger.ErrorContext(ctx, "failed to execute request")

		return nil, fmt.Errorf("failed to execute http request: %w", err)
	}

	logger.DebugContext(ctx, "Response Code : %v", slog.Int("status_code", res.StatusCode))

	defer func() {
		err = res.Body.Close()
		if err != nil {
			logger.ErrorContext(ctx, "error while closing response body")
		}
	}()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		logger.ErrorContext(ctx, "failed to read response body")

		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result *ResponseOK

	switch code := res.StatusCode; {
	// 200 to 399
	case code >= http.StatusOK && code < http.StatusBadRequest:
		result, err = parseResult(bytes)
		if err != nil {
			logger.ErrorContext(ctx, "error while getting api response")

			return nil, err
		}
		// 400 to 500
	case code >= http.StatusBadRequest && code <= http.StatusInternalServerError:
		logger.ErrorContext(ctx, "API responded with status code %v", slog.String("status_code", res.Status))

		apiErr, err := parseError(bytes)
		if err != nil {
			logger.ErrorContext(ctx, "error while parsing error from API")

			return nil, err
		}

		return nil, apiErr
	default:
		return nil, newUnknownHTTPCodeError(res.StatusCode)
	}

	return result, nil
}

// parseResult unmarshals a valid response to a struct.
func parseResult(res []byte) (*ResponseOK, error) {
	rawResp := apiResponseRaw{}
	if err := json.Unmarshal(res, &rawResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json raw response: %w", err)
	}

	ret := &ResponseOK{}

	// some fields like parameters, errors can be either maps or arrays depending on the context.
	switch conv := rawResp.Errors.(type) {
	case map[string]any:
		// if errors come as a map, it means that the API has returned an error.
		var errMsgBuilder strings.Builder
		for k, v := range conv {
			// there are potentially multiple error, let's concatenate them into a string.
			errMsgBuilder.WriteString(fmt.Sprintf("%s : %s\n", k, v))
		}

		return nil, &ResponseError{Message: errMsgBuilder.String()}
	default:
		ret.Errors = make(map[string]any)
	}

	switch conv := rawResp.Parameters.(type) {
	case map[string]any:
		ret.Parameters = conv
	default:
		ret.Parameters = make(map[string]any)
	}

	ret.Get = rawResp.Get
	ret.Paging = rawResp.Paging
	ret.Results = rawResp.Results
	// Response has revealed to be either an array (99% of the time) or a map (status and teamStatistics endpoints).
	// So we keep it as a byte array and leave the responsibility to results parser to make the correct conversion.
	jsonBody, err := json.Marshal(rawResp.Response)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal raw response to json bytes array: %w", err)
	}

	ret.Response = jsonBody

	return ret, nil
}

// parseError unmarshals an error response into a struct.
func parseError(res []byte) (*ResponseError, error) {
	errResp := ResponseError{}
	if err := json.Unmarshal(res, &errResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json error response: %w", err)
	}

	return &errResp, nil
}

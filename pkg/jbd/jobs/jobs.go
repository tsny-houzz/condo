package jobs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// FetchJobsOption is a type for setting optional query parameters for FetchJobs.
type FetchJobsOption func(url.Values)

// WithLimit sets the "limit" query parameter for FetchJobs.
func WithLimit(limit int) FetchJobsOption {
	return func(params url.Values) {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}
}

// WithOffset sets the "offset" query parameter for FetchJobs.
func WithOffset(offset int) FetchJobsOption {
	return func(params url.Values) {
		params.Set("offset", fmt.Sprintf("%d", offset))
	}
}

// WithTaskType sets the "taskType" query parameter for FetchJobs.
func WithTaskType(taskType string) FetchJobsOption {
	return func(params url.Values) {
		params.Set("taskType", taskType)
	}
}

// WithIssuer sets the "issuer" query parameter for FetchJobs.
func WithIssuer(issuer string) FetchJobsOption {
	return func(params url.Values) {
		params.Set("issuer", issuer)
	}
}

// FetchJobs retrieves jobs with optional query parameters.
func (jc *JobClient) FetchJobs(opts ...FetchJobsOption) (*JobsResponse, error) {
	baseURL := fmt.Sprintf("%s://%s/toolsvr/jobs", jc.Scheme, jc.Host)
	url, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %w", err)
	}

	// Apply options to set query parameters
	params := url.Query()
	for _, opt := range opts {
		opt(params)
	}
	url.RawQuery = params.Encode()
	println(url.String())

	// Create a new HTTP client with the specified timeout
	client := &http.Client{Timeout: jc.Timeout}
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var jobsResponse JobsResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobsResponse); err != nil {
		return nil, fmt.Errorf("could not decode response: %w", err)
	}

	return &jobsResponse, nil
}

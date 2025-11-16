package pullrequests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ReassignPRParams struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

type ReassignPRResult struct {
	PR         ReassignPRResultPR `json:"pr"`
	ReplacedBy string             `json:"replacedBy"`
}

type ReassignPRResultPR struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}

func (c Client) ReassignPR(ctx context.Context, params ReassignPRParams) (ReassignPRResult, error) {
	reqBodyBytes, err := json.Marshal(params)
	if err != nil {
		return ReassignPRResult{}, fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseUrl+"/pullRequest/reassign",
		bytes.NewReader(reqBodyBytes),
	)
	if err != nil {
		return ReassignPRResult{}, fmt.Errorf("error building request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return ReassignPRResult{}, fmt.Errorf("response error: %w", err)
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return ReassignPRResult{}, fmt.Errorf(
			"unsuccessful request, status code = %d, response body = %s, request body = %s",
			resp.StatusCode,
			string(body),
			string(reqBodyBytes),
		)
	}

	var response ReassignPRResult

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return ReassignPRResult{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return response, nil
}

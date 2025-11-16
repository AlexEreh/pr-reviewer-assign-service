package pullrequests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MergePRParams struct {
	PullRequestID string `json:"pull_request_id"`
}

type MergePRResult struct {
	PR MergePRResultPR `json:"pr"`
}

type MergePRResultPR struct {
	PullRequestID     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	AuthorID          string   `json:"author_id"`
	Status            string   `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	MergedAt          string   `json:"mergedAt"`
}

func (c Client) MergePR(ctx context.Context, params MergePRParams) (MergePRResult, error) {
	reqBodyBytes, err := json.Marshal(params)
	if err != nil {
		return MergePRResult{}, fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseUrl+"/pullRequest/merge",
		bytes.NewReader(reqBodyBytes),
	)
	if err != nil {
		return MergePRResult{}, fmt.Errorf("error building request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return MergePRResult{}, fmt.Errorf("response error: %w", err)
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return MergePRResult{}, fmt.Errorf(
			"unsuccessful request, status code = %d, response body = %s, request body = %s",
			resp.StatusCode,
			string(body),
			string(reqBodyBytes),
		)
	}

	var response MergePRResult

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return MergePRResult{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return response, nil
}

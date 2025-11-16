package users

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetReviewPRsParams struct {
	UserID string `json:"-"`
}

type GetReviewPRsResult struct {
	UserID       string                 `json:"user_id"`
	PullRequests []GetReviewPRsResultPR `json:"pull_requests"`
}

type GetReviewPRsResultPR struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

func (c Client) GetReviewPRs(
	ctx context.Context,
	params GetReviewPRsParams,
) (GetReviewPRsResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseUrl+"/users/getReview", nil)
	if err != nil {
		return GetReviewPRsResult{}, fmt.Errorf("error building request: %w", err)
	}

	q := req.URL.Query()
	q.Add("user_id", params.UserID)
	req.URL.RawQuery = q.Encode()

	resp, err := c.c.Do(req)
	if err != nil {
		return GetReviewPRsResult{}, fmt.Errorf("response error: %w", err)
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return GetReviewPRsResult{}, fmt.Errorf(
			"unsuccessful request: status %d, body: %s",
			resp.StatusCode,
			string(body),
		)
	}

	var response GetReviewPRsResult

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return GetReviewPRsResult{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return response, nil
}

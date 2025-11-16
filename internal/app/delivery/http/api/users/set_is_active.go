package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SetIsActiveParams struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveResultUser struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveResult struct {
	User SetIsActiveResultUser `json:"user"`
}

func (c Client) SetIsActive(
	ctx context.Context,
	params SetIsActiveParams,
) (SetIsActiveResult, error) {
	reqBodyBytes, err := json.Marshal(params)
	if err != nil {
		return SetIsActiveResult{}, fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseUrl+"/users/setIsActive",
		bytes.NewReader(reqBodyBytes),
	)
	if err != nil {
		return SetIsActiveResult{}, fmt.Errorf("error building request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return SetIsActiveResult{}, fmt.Errorf("response error: %w", err)
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return SetIsActiveResult{}, fmt.Errorf(
			"unsuccessful request: status %d, body: %s",
			resp.StatusCode,
			string(body),
		)
	}

	var response SetIsActiveResult

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return SetIsActiveResult{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return response, nil
}

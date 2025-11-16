package teams

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AddTeamParams struct {
	TeamName string              `json:"team_name"`
	Members  []AddTeamParamsUser `json:"members"`
}

type AddTeamParamsUser struct {
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type AddTeamResult struct {
	Team AddTeamResultTeam `json:"team"`
}

type AddTeamResultTeam struct {
	TeamName string              `json:"team_name"`
	Members  []AddTeamResultUser `json:"members"`
}

type AddTeamResultUser struct {
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (c Client) AddTeam(ctx context.Context, params AddTeamParams) (AddTeamResult, error) {
	reqBodyBytes, err := json.Marshal(params)
	if err != nil {
		return AddTeamResult{}, fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseUrl+"/teams/add",
		bytes.NewReader(reqBodyBytes),
	)
	if err != nil {
		return AddTeamResult{}, fmt.Errorf("error building request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return AddTeamResult{}, fmt.Errorf("response error: %w", err)
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return AddTeamResult{}, fmt.Errorf(
			"unsuccessful request: status %d, body: %s",
			resp.StatusCode,
			string(body),
		)
	}

	var response AddTeamResult

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&response)
	if err != nil {
		return AddTeamResult{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return response, nil
}

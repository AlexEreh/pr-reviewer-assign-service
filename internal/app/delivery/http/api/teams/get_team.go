package teams

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetTeamParams struct {
	TeamName string
}

type GetTeamResult struct {
	TeamName string              `json:"team_name"`
	Members  []GetTeamResultUser `json:"members"`
}

type GetTeamResultUser struct {
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (c Client) GetTeam(ctx context.Context, params GetTeamParams) (GetTeamResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseUrl+"/teams/get", http.NoBody)
	if err != nil {
		return GetTeamResult{}, fmt.Errorf("error building request: %w", err)
	}

	q := req.URL.Query()
	q.Add("team_name", params.TeamName)
	req.URL.RawQuery = q.Encode()

	resp, err := c.c.Do(req)
	if err != nil {
		return GetTeamResult{}, fmt.Errorf("response error: %w", err)
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return GetTeamResult{}, fmt.Errorf(
			"unsuccessful request: status %d, body: %s",
			resp.StatusCode,
			string(body),
		)
	}

	var response GetTeamResult

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&response)
	if err != nil {
		return GetTeamResult{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return response, nil
}

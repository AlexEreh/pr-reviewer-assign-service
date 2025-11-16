package health

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ReadyZParams struct{}

type ReadyZResult struct {
	OK bool `json:"ok"`
}

type readyz interface {
	ReadyZ(ctx context.Context, params ReadyZParams) (ReadyZResult, error)
}

func (c Client) ReadyZ(ctx context.Context, _ ReadyZParams) (ReadyZResult, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseUrl+"/health/readyz",
		http.NoBody,
	)
	if err != nil {
		return ReadyZResult{}, fmt.Errorf("error building request: %w", err)
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return ReadyZResult{}, fmt.Errorf("response error: %w", err)
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode >= http.StatusBadRequest {
		return ReadyZResult{}, errors.New("unsuccessful request")
	}

	var response ReadyZResult

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&response)
	if err != nil {
		return ReadyZResult{}, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return response, nil
}

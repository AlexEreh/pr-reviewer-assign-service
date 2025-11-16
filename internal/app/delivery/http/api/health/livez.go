package health

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type LiveZParams struct{}

type LiveZResult struct{}

type livez interface {
	LiveZ(ctx context.Context, params LiveZParams) (LiveZResult, error)
}

func (c Client) LiveZ(ctx context.Context, _ LiveZParams) (LiveZResult, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		c.baseUrl+"/health/livez",
		http.NoBody,
	)
	if err != nil {
		return LiveZResult{}, fmt.Errorf("error building request: %w", err)
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return LiveZResult{}, fmt.Errorf("response error: %w", err)
	}

	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode >= http.StatusBadRequest {
		return LiveZResult{}, errors.New("unsuccessful request")
	}

	return LiveZResult{}, nil
}

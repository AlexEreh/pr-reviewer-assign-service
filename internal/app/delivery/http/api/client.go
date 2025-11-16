package api

import (
	"net/http"

	"pr-reviewer-assign-service/internal/app/delivery/http/api/health"
	"pr-reviewer-assign-service/internal/app/delivery/http/api/pullrequests"
	"pr-reviewer-assign-service/internal/app/delivery/http/api/statistics"
	"pr-reviewer-assign-service/internal/app/delivery/http/api/teams"
	"pr-reviewer-assign-service/internal/app/delivery/http/api/users"
)

type Client struct {
	healthClient       health.Client
	pullRequestsClient pullrequests.Client
	statisticsClient   statistics.Client
	teamsClient        teams.Client
	usersClient        users.Client
}

func NewClient(c *http.Client, baseUrl string) Client {
	return Client{
		healthClient:       health.NewClient(c, baseUrl),
		pullRequestsClient: pullrequests.NewClient(c, baseUrl),
		statisticsClient:   statistics.NewClient(c, baseUrl),
		teamsClient:        teams.NewClient(c, baseUrl),
		usersClient:        users.NewClient(c, baseUrl),
	}
}

func (c Client) Health() health.Client {
	return c.healthClient
}

func (c Client) PR() pullrequests.Client {
	return c.pullRequestsClient
}

func (c Client) Statistics() statistics.Client {
	return c.statisticsClient
}

func (c Client) Teams() teams.Client {
	return c.teamsClient
}

func (c Client) Users() users.Client {
	return c.usersClient
}

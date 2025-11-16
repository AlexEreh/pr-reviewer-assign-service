package impl

import (
	"github.com/knadh/koanf/v2"

	"pr-reviewer-assign-service/internal/app/delivery/http/impl/health"
	"pr-reviewer-assign-service/internal/app/delivery/http/impl/pullrequests"
	"pr-reviewer-assign-service/internal/app/delivery/http/impl/statistics"
	"pr-reviewer-assign-service/internal/app/delivery/http/impl/teams"
	"pr-reviewer-assign-service/internal/app/delivery/http/impl/users"
	"pr-reviewer-assign-service/internal/app/domain/usecase"
	"pr-reviewer-assign-service/pkg/http"
)

type API struct {
	useCase *usecase.UseCase
	server  *http.Server

	healthHandler       *health.Handler
	statisticsHandler   *statistics.Handler
	pullRequestsHandler *pullrequests.Handler
	teamsHandler        *teams.Handler
	usersHandler        *users.Handler

	errorMiddleware  *ErrorMiddleware
	loggerMiddleware *LogMiddleware
}

// NewAPI
//
//	@title						PR Reviewer Assignment Service (Test Task, Fall 2025)
//	@version					1.0
//	@description				This is a sample server.
//	@termsOfService				http://swagger.io/terms/
//
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@contact.email				support@swagger.io
//
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@host						127.0.0.1:8080
//
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func NewAPI(cfg *koanf.Koanf, useCase *usecase.UseCase, server *http.Server) *API {
	api := &API{
		useCase:             useCase,
		server:              server,
		healthHandler:       health.NewHandler(useCase),
		statisticsHandler:   statistics.NewHandler(useCase),
		pullRequestsHandler: pullrequests.NewHandler(useCase),
		teamsHandler:        teams.NewHandler(useCase),
		usersHandler:        users.NewHandler(useCase),
		errorMiddleware:     NewErrorMiddleware(),
		loggerMiddleware:    NewLogMiddleware(cfg),
	}

	return api
}

func (a *API) Init() {
	teamsGroup := a.server.Group("/teams", a.loggerMiddleware.Call, a.errorMiddleware.Call)
	teamsGroup.Post("/add", a.teamsHandler.AddTeam)
	teamsGroup.Get("/get", a.teamsHandler.GetTeam)

	usersGroup := a.server.Group("/users", a.loggerMiddleware.Call, a.errorMiddleware.Call)
	usersGroup.Post("/setIsActive", a.usersHandler.SetIsActive)
	usersGroup.Get("/getReview", a.usersHandler.GetReview)

	pullRequestsGroup := a.server.Group(
		"/pullRequest",
		a.loggerMiddleware.Call,
		a.errorMiddleware.Call,
	)
	pullRequestsGroup.Post("/create", a.pullRequestsHandler.CreatePR)
	pullRequestsGroup.Post("/merge", a.pullRequestsHandler.MergePR)
	pullRequestsGroup.Post("/reassign", a.pullRequestsHandler.ReassignPR)

	healthGroup := a.server.Group("/health", a.loggerMiddleware.Call, a.errorMiddleware.Call)
	healthGroup.Get("/livez", a.healthHandler.LiveZ)
	healthGroup.Get("/readyz", a.healthHandler.ReadyZ)

	statisticsGroup := a.server.Group(
		"/statistics",
		a.loggerMiddleware.Call,
		a.errorMiddleware.Call,
	)
	statisticsGroup.Get("/get", a.statisticsHandler.GetStatistics)
}

package app

import (
	"context"

	"github.com/knadh/koanf/v2"

	"pr-reviewer-assign-service/internal/app/data/postgres"
	http2 "pr-reviewer-assign-service/internal/app/delivery/http/impl"
	"pr-reviewer-assign-service/internal/app/domain/usecase"
	"pr-reviewer-assign-service/migrations"
	"pr-reviewer-assign-service/pkg/app"
	"pr-reviewer-assign-service/pkg/db"
	"pr-reviewer-assign-service/pkg/http"
	"pr-reviewer-assign-service/pkg/log"
	"pr-reviewer-assign-service/pkg/txman"
)

func Run(ctx context.Context, cfg *koanf.Koanf) []error {
	return app.Run(ctx, func(app *app.App) error {
		log.Init(cfg)

		server := http.Init(app, cfg.Cut("server"))
		database, err := db.Init(app, cfg.Cut("database"))
		if err != nil {
			return err
		}

		err = migrations.MigrateUp(app, database)
		if err != nil {
			return err
		}

		man := txman.New(database)

		repo := postgres.NewRepository(man)
		uc := usecase.New(repo, man)
		api := http2.NewAPI(cfg, uc, server)

		api.Init()

		return nil
	})
}

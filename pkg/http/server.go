package http

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/knadh/koanf/v2"

	"pr-reviewer-assign-service/pkg/app"
)

type Server = fiber.App

func Init(ctx *app.App, cfg *koanf.Koanf) *Server {
	server := fiber.New(fiber.Config{ //nolint:exhaustruct // Похоже, тут без ручек не обойтись
		StreamRequestBody:  true,
		UnescapePath:       true,
		ProxyHeader:        "X-Forwarded-For",
		EnableIPValidation: true,
	})

	ctx.AfterInit(func() error {
		ctx.Go(func() error {
			return server.Listen(fmt.Sprintf("%s:%d", cfg.String("host"), cfg.Int("port")))
		})
		ctx.OnComplete(func() error {
			return server.Shutdown()
		})

		return nil
	})

	return server
}

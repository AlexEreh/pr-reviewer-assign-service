package impl

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"

	goerrors "errors"

	"pr-reviewer-assign-service/pkg/errors"
	httperr "pr-reviewer-assign-service/pkg/errors/http"
	jsonerr "pr-reviewer-assign-service/pkg/errors/json"
	"pr-reviewer-assign-service/pkg/log"
)

type ErrorMiddleware struct{}

func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

func (mw *ErrorMiddleware) Call(c *fiber.Ctx) (err error) {
	err = c.Next()

	if err == nil {
		return nil
	}

	var eerr *errors.Error
	if goerrors.As(err, &eerr) {
		statusCode, ok := httperr.GetStatus(eerr)
		if !ok {
			statusCode = http.StatusInternalServerError
		}

		c.Status(statusCode)

		body := jsonerr.Marshal(eerr, func(e *jsonerr.MarshalerConfig) {
			e.IsPrivateKey = func(key string) bool {
				return key == "StackTrace"
			}
		}, httperr.JSONPrivate)

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		_, err = c.Write(body)
		if err != nil {
			return err
		}

		return nil
	}

	var fiberErr *fiber.Error
	if goerrors.As(err, &fiberErr) {
		c.Status(fiberErr.Code)

		var e error

		switch {
		case fiberErr.Code >= fiber.StatusBadRequest && fiberErr.Code < fiber.StatusInternalServerError:
			e = fiberErr
		default:
			e = errors.Wrap(fiberErr, errors.InternalError, httperr.WithStatus(fiberErr.Code))
		}

		body := jsonerr.Marshal(e, func(e *jsonerr.MarshalerConfig) {
			e.IsPrivateKey = func(key string) bool {
				return key == "StackTrace"
			}
		}, httperr.JSONPrivate)

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		_, err = c.Write(body)
		if err != nil {
			return err
		}

		return nil
	}

	internalError := errors.Wrap(err, errors.InternalError)

	c.Status(http.StatusInternalServerError)
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	_, err = c.Write(jsonerr.Marshal(internalError))
	if err != nil {
		return err
	}

	return nil
}

type LogMiddleware struct {
	cfg *koanf.Koanf
}

func NewLogMiddleware(cfg *koanf.Koanf) *LogMiddleware {
	return &LogMiddleware{
		cfg: cfg,
	}
}

func (mw *LogMiddleware) Call(ctx *fiber.Ctx) error {
	err := ctx.Next()

	var eerr *errors.Error

	if err != nil && goerrors.As(err, &eerr) {
		status, ok := httperr.GetStatus(err)
		if !ok {
			return err
		}

		if status >= fiber.StatusInternalServerError {
			log.LoggerFromCtx(context.Background()).Error(
				"operation error",
				zap.Error(eerr),
				zap.String("stacktrace", eerr.StackTrace().String()),
			)
		}

		return err
	}
	if err != nil {
		log.LoggerFromCtx(context.Background()).Error(
			"operation error",
			zap.Error(err),
		)

		return err
	}

	return nil
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	goerrors "errors"

	"pr-reviewer-assign-service/internal/app"
	"pr-reviewer-assign-service/internal/app/config"
	"pr-reviewer-assign-service/pkg/errors"
)

func main() {
	cfgPath := flag.String("config", "config/local.yml", "Path to config storage.")
	flag.Parse()

	cfg, err := config.FromFile(*cfgPath)
	if err != nil {
		log.Println(err)

		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	errs := app.Run(ctx, cfg)
	if len(errs) != 0 {
		fmt.Printf("error:\n\t%s\n", errorsToString(errs)) //nolint:forbidigo

		return
	}
}

func errorsToString(errs []error) string {
	sb := strings.Builder{}

	for index, err := range errs {
		if index != 0 {
			sb.WriteString("\t\n")
		}

		sb.WriteString(errToString(err))
	}

	return sb.String()
}

func errToString(err error) string {
	var typedErr *errors.Error
	if goerrors.As(err, &typedErr) {
		return fmt.Sprintf("%+v", typedErr.Unwrap())
	}

	return fmt.Sprintf("%+v", err)
}

package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"

	"github.com/mtpereira/deck/deck"
	"github.com/mtpereira/deck/web"
)

func main() {
	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	err := run(ctx, log)
	if err != nil {
		log.Error(error.Error(err))
	}
}

func run(ctx context.Context, log *slog.Logger) error {
	cfg := struct {
		APIHost string `conf:"default:127.0.0.1:9000"`
	}{}
	prefix := "DECK"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("error parsing config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("error outputting conf: %w", err)
	}
	log.Info("startup", "configs", out)

	log.Info("startup", "status", "initalising api")
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	ds := deck.NewStore(log)
	mux := web.NewMux(log, ds)

	api := http.Server{
		Handler:  mux,
		Addr:     cfg.APIHost,
		ErrorLog: slog.NewLogLogger(slog.NewTextHandler(os.Stdout, nil), slog.LevelError),
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info("startup", "status", "api listening", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Info("shutdown", "status", "initiated", "signal", sig)
		defer log.Info("shutdown", "status", "complete")

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			api.Close()
			return fmt.Errorf("could not shutdown gracefully: %w", err)
		}
	}

	return nil
}

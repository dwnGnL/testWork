package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dwnGnL/testWork/internal/config"
	"github.com/dwnGnL/testWork/internal/worker"
	"github.com/dwnGnL/testWork/lib/goerrors"
	"golang.org/x/sync/errgroup"
)

func StartWorker(ctx context.Context, cfg *config.Config) error {
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()
	s, err := buildService(ctx, cfg)
	if err != nil {
		return fmt.Errorf("build service err:%w", err)
	}
	w, err := worker.New(&cfg.Consumer)
	if err != nil {
		return fmt.Errorf("worker.New err: %w", err)
	}

	err = w.StartWorker(s)
	if err != nil {
		return fmt.Errorf("w.StartWorker err: %w", err)
	}
	var group errgroup.Group

	group.Go(func() error {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		goerrors.Log().Debug("wait for Ctrl-C")
		<-sigCh
		goerrors.Log().Debug("Ctrl-C signal")
		cancelCtx()

		w.StopWorker()

		return nil
	})

	if err := group.Wait(); err != nil {
		goerrors.Log().WithError(err).Error("Stopping service with error")
	}
	return nil
}

package auth

import (
	"context"
	"fmt"
	"github.com/mefourr/tgdevob/authentication/internal/controller/grpc/generator"
	"github.com/mefourr/tgdevob/authentication/internal/usecase"
	"github.com/mefourr/tgdevob/authentication/pkg/logger"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	srv  *grpc.Server
	port int
}

func New(port int) *App {
	srv := grpc.NewServer()

	generator.Register(srv, usecase.New())

	return &App{srv: srv, port: port}
}

func (a *App) MustRun(ctx context.Context) {
	if err := a.run(ctx); err != nil {
		panic(err)
	}
}

func (a *App) run(ctx context.Context) error {
	// TODO: add controller logger attribute to slog

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to listen:", "err", err)
		return err
	}

	slog.InfoContext(ctx, "server listening at", "addr", lis.Addr())

	if err := a.srv.Serve(lis); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "failed to serve:", "err", err)
		return err
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) {
	slog.InfoContext(ctx, "shutting down grpc server")
	a.srv.GracefulStop()
	slog.InfoContext(ctx, "grpc server stopped")
}

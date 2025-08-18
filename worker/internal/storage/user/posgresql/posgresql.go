package posgresql

import (
	"context"
	"fmt"
	"gihub.com/mefourr/tgdevob/worker/config"
	"gihub.com/mefourr/tgdevob/worker/internal/utils"
	"gihub.com/mefourr/tgdevob/worker/pkg/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func New(ctx context.Context, conf config.StorageConfig) (pool *pgxpool.Pool, err error) {
	hp := net.JoinHostPort(conf.Hostname, conf.Port)
	dsn := fmt.Sprintf("postgresql://%s:%s@%s/%s", conf.Username, conf.Password, hp, conf.Database)
	slog.InfoContext(ctx, "try to connect to psql by", "dsn", dsn)

	err = utils.TryConn(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}
		if err = pool.Ping(ctx); err != nil {
			slog.ErrorContext(logging.ErrorCtx(ctx, err), "Unable to ping database")
			return err
		}

		return nil
	}, conf)

	if err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "failed to TryConn")
		return nil, err
	}

	return
}

package utils

import (
	"context"
	"github.com/mefourr/tgdevob/telegram-api/pkg/logging"
	"log/slog"
	"os"
)

func CleanStorage(ctx context.Context, filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		slog.ErrorContext(logging.ErrorCtx(ctx, err), "Error: "+err.Error())
		return err
	}
	return nil
}

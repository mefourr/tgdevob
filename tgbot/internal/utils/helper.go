package utils

import (
	"context"
	"github.com/mefourr/tgdevob/tgbot/pkg/logger"
	"log/slog"
	"os"
)

func CleanStorage(ctx context.Context, filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Error: "+err.Error())
		return err
	}
	return nil
}

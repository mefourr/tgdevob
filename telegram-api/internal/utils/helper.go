package utils

import (
	"context"
	"log/slog"
	"os"
)

func CleanStorage(ctx context.Context, filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		slog.ErrorContext(ErrorCtx(ctx, err), "Error: "+err.Error())
		return err
	}
	return nil
}

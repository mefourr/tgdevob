package logger

import (
	"context"
	"errors"
	"log/slog"
	"os"
)

type HandlerMiddleware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddleware {
	return &HandlerMiddleware{next: next}
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	if _, ok := ctx.Value(key).(logCtx); ok {
	}
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithAttrs(attrs)} // не забыть обернуть, но осторожно
}

func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithGroup(name)} // не забыть обернуть, но осторожно
}

type logCtx struct {
}

type keyType int

const key = keyType(0)

// -----------------------------------------------

type errorWithLogCtx struct {
	next error
	ctx  logCtx
}

func (e *errorWithLogCtx) Error() string {
	return e.next.Error()
}

func WrapError(ctx context.Context, err error) error {
	c := logCtx{}
	if x, ok := ctx.Value(key).(logCtx); ok {
		c = x
	}
	return &errorWithLogCtx{
		next: err,
		ctx:  c,
	}
}

func ErrorCtx(ctx context.Context, err error) context.Context {
	var e *errorWithLogCtx
	if errors.As(err, &e) {
		return context.WithValue(ctx, key, e.ctx)
	}
	return ctx
}

// -----------------------------------------------

func Init() context.Context {
	handler := slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	handler = NewHandlerMiddleware(handler)
	slog.SetDefault(slog.New(handler))
	return context.Background()
}

//func Handler(ctx context.Context, userID int) {
//	ctx = WithLogUserID(ctx, userID)
//	slog.InfoContext(ctx, "Handler started")
//	phone, _ := GetPhoneByID(ctx)
//	ctx = WithLogPhone(ctx, phone)
//	err := SendSMS(ctx)
//	if err != nil {
//		slog.ErrorContext(ErrorCtx(ctx, err), "Error: "+err.Error())
//		return
//	}
//	slog.InfoContext(ctx, "controller done")
//}

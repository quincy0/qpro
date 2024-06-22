package rescue

import (
	"context"
	"github.com/quincy0/qpro/qLog"

	"go.uber.org/zap"
)

func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		qLog.Error("recover failed", zap.Any("p", p), zap.Stack("stack"))
	}
}

func RecoverCtx(ctx context.Context, cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}
	if p := recover(); p != nil {
		qLog.TraceError(ctx, "recover failed", zap.Any("p", p), zap.Stack("stack"))
	}
}

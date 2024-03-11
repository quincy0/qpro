package qTrace

import (
	"context"

	"go.opentelemetry.io/otel"
)

func New(ctx context.Context) context.Context {
	provider := otel.GetTracerProvider()
	tr := provider.Tracer("dh-scheduler-otel")
	tCTX, _ := tr.Start(ctx, "cron-job")
	return tCTX
}

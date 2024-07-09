package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"go.opentelemetry.io/otel/trace"
)

type Connection interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

func spanWithQuery(ctx context.Context, trace trace.Tracer, query string) (context.Context, trace.Span) {
	ctx, span := trace.Start(ctx, "db:query")
	span.SetAttributes(semconv.DBStatementKey.String(query))
	return ctx, span
}

package bunzap

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Field names
const (
	OperationFieldName     = "operation"
	OperationTimeFieldName = "operation_time_ms"
)

// QueryHook defines the
// structure of our query hook
// it implements the bun.QueryHook
// interface
type QueryHook struct {
	bun.QueryHook

	logger   *zap.Logger
	slowTime time.Duration
}

// QueryHookOptions defines the
// available options for a new
// query hook.
type QueryHookOptions struct {
	Logger   *zap.Logger
	SlowTime time.Duration
}

// NewQueryHook returns a new query hook for use with
// uptrace/bun.
func NewQueryHook(options QueryHookOptions) QueryHook {
	return QueryHook{
		logger:   options.Logger,
		slowTime: options.SlowTime,
	}
}

func (qh QueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (qh QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	queryTime := time.Since(event.StartTime)
	fields := []zapcore.Field{
		zap.String(OperationFieldName, event.Operation()),
		zap.Int64(OperationTimeFieldName, queryTime.Milliseconds()),
	}

	// Errors will always be logged
	if event.Err != nil {
		fields = append(fields, zap.Error(event.Err))
		qh.logger.Error(event.Query, fields...)
		return
	}

	// Queries over a slow time duration
	// will be logged as debug
	if queryTime > qh.slowTime {
		qh.logger.Debug(event.Query, fields...)
	}
}

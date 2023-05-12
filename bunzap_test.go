package bunzap_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexlast/bunzap"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestQueryHookError(t *testing.T) {
	core, obs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	defer logger.Sync()

	qh := bunzap.NewQueryHook(bunzap.QueryHookOptions{
		Logger: logger,
	})

	event := &bun.QueryEvent{
		StartTime: time.Now(),
		Query:     "SELECT * FROM users WHERE id = $1",
		Err:       errors.New("database error"),
	}

	qh.AfterQuery(context.Background(), event)
	assert.Equal(t, 1, obs.Len())

	logs := obs.All()
	assert.Equal(t, event.Query, logs[0].Message)
	assert.Equal(t, zapcore.ErrorLevel, logs[0].Level)
	assert.Len(t, logs[0].Context, 3)
	assert.Equal(t, []zap.Field{
		{
			Key:    bunzap.OperationFieldName,
			Type:   zapcore.StringType,
			String: event.Operation(),
		},
		{
			Key:     bunzap.OperationTimeFieldName,
			Type:    zapcore.Int64Type,
			Integer: 0,
		},
		{
			Key:       "error",
			Type:      zapcore.ErrorType,
			Interface: event.Err,
		},
	}, logs[0].Context)
}

func TestQueryHookDebug(t *testing.T) {
	core, obs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	defer logger.Sync()

	qh := bunzap.NewQueryHook(bunzap.QueryHookOptions{
		Logger: logger,
	})

	event := &bun.QueryEvent{
		StartTime: time.Now(),
		Query:     "SELECT * FROM users WHERE id = $1",
	}

	qh.AfterQuery(context.Background(), event)
	assert.Equal(t, 1, obs.Len())

	logs := obs.All()
	assert.Equal(t, event.Query, logs[0].Message)
	assert.Equal(t, zapcore.DebugLevel, logs[0].Level)
	assert.Len(t, logs[0].Context, 2)
	assert.Equal(t, []zap.Field{
		{
			Key:    bunzap.OperationFieldName,
			Type:   zapcore.StringType,
			String: event.Operation(),
		},
		{
			Key:     bunzap.OperationTimeFieldName,
			Type:    zapcore.Int64Type,
			Integer: 0,
		},
	}, logs[0].Context)
}

func TestQueryHookFast(t *testing.T) {
	core, obs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	defer logger.Sync()

	qh := bunzap.NewQueryHook(bunzap.QueryHookOptions{
		SlowTime: 200 * time.Millisecond,
		Logger:   logger,
	})

	event := &bun.QueryEvent{
		StartTime: time.Now(),
		Query:     "SELECT * FROM users WHERE id = $1",
	}

	qh.AfterQuery(context.Background(), event)
	assert.Equal(t, 0, obs.Len())
}

func TestQueryHookSlow(t *testing.T) {
	core, obs := observer.New(zap.DebugLevel)
	logger := zap.New(core)

	defer logger.Sync()

	qh := bunzap.NewQueryHook(bunzap.QueryHookOptions{
		Logger:   logger,
		SlowTime: 200 * time.Millisecond,
	})

	event := &bun.QueryEvent{
		StartTime: time.Now().Add(-300 * time.Millisecond),
		Query:     "SELECT * FROM users WHERE id = $1",
	}

	qh.AfterQuery(context.Background(), event)
	assert.Equal(t, 1, obs.Len())

	logs := obs.All()
	assert.Equal(t, event.Query, logs[0].Message)
	assert.Equal(t, zapcore.DebugLevel, logs[0].Level)
	assert.Len(t, logs[0].Context, 2)
	assert.Equal(t, []zap.Field{
		{
			Key:    bunzap.OperationFieldName,
			Type:   zapcore.StringType,
			String: event.Operation(),
		},
		{
			Key:     bunzap.OperationTimeFieldName,
			Type:    zapcore.Int64Type,
			Integer: 300,
		},
	}, logs[0].Context)
}

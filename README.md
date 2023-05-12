# bunzap
A query hook for [uptrace/bun](https://github.com/uptrace/bun) that logs with [uber-go/zap](https://github.com/uber-go/zap).

```bash
$ go get github.com/alexlast/bunzap
```

All errors will be logged at error level with the hook enabled, everything else will be logged as debug. If `SlowTime` is defined, only operations taking longer than the defined duration will be logged.

## Usage
```go
db.AddQueryHook(bunzap.NewQueryHook(bunzap.QueryHookOptions{
    Logger:       options.Logger,
    SlowDuration: 200 * time.Millisecond, // Omit to log all operations as debug
}))
```

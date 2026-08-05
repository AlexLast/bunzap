

# bunzap
Un hook de consulta para [uptrace/bun](https://github.com/uptrace/bun) que registra con [uber-go/zap](https://github.com/uber-go/zap).

```bash
$ go get github.com/alexlast/bunzap
```

Con el hook habilitado, todos los errores se registrarán en el nivel de error; todo lo demás se registrará como debug. Si se define `SlowDuration`, solo se registrarán las operaciones que duren más que la duración definida.

## Uso
```go
db := bun.NewDB()
logger, err := zap.NewProduction()
db.AddQueryHook(bunzap.NewQueryHook(bunzap.QueryHookOptions{
    Logger:       logger,
    SlowDuration: 200 * time.Millisecond, // Omit to log all operations as debug
}))
```

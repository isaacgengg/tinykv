# tinykv notes

## Day 1 - 2026/05/01

**Built:** a race-free in-memory KV store with GET/PUT/DELETE over HTTP.

### Concurrency
- Go is multithreaded. In Node.js, concurrency is cooperative (`await`). In Go, concurrency is preemptive.
- Race conditions are very real — two handlers can execute simultaneously on different threads. Always use a mutex.

### HTTP routing
`http.HandleFunc` is the Go equivalent of Express routing:
```go
// Go
http.HandleFunc("GET /kv/{key}", func(w http.ResponseWriter, r *http.Request) {})

// Node
app.get('/kv/:key', (req, res) => {})
```

### Error handling
No `try/catch` in Go. Functions signal failure by returning an `error` as their last return value:
```go
value, err := someFunction()
if err != nil {
    // something went wrong
}

// If you only care about the error, not the result:
_, err := someFunction()
```

### defer
When acquiring a resource, immediately defer its release. Runs when the function exits no matter what:
```go
mu.Lock()
defer mu.Unlock()

defer r.Body.Close()
```

### Zero values
Go variables are always initialized — no `undefined` surprises:
- `sync.RWMutex` is ready to use without `= sync.RWMutex{}`
- `int` → `0`, `string` → `""`, `bool` → `false`, pointers → `nil`

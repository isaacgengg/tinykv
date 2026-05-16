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

## Day 2 - 2026/05/15
**Built:** Tests using httptest, Store struct, WAL persistence that survives restarts
### WAL
Write ahead log, stores the requests before actually doing the requests and updating the map so that upon server restarts, you can "replay" and not lose data.
Isolate the real wal.log from a temporary test one otherwise tests will fail/affect the real wal.log file
### MUX vs Direct Handler Call
The mux is what extracts stuff from the URL and stores it on a request. If you call the handler directly, it will bypass the mux so the PathValue will be "". Always use mux.ServeHTTP(w, req) in tests.
### Pointer vs Value Receiver
pointer receiver (*Store) gives you the real struct so that mutations are shared. If you used a value receiver (Store), it gives a copy so changes are lost.

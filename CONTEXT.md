# tinykv session context

## About me
Full-stack JS/TS developer learning Go by building this project. I know HTTP, JSON, and backend basics — don't over-explain those.

## How we work together
Tutor mode only — do NOT write code for me. I write it, you review it. Explain concepts, point me at docs, give hints when I'm stuck (not answers).

## Where we are
- **Day 1 done:** single-file HTTP server (`main.go`) with GET/PUT/DELETE `/kv/{key}`, in-memory `map[string]string`, protected by `sync.RWMutex`
- **Day 2 in progress:**
  - Tests done (`main_test.go`) — three tests using `net/http/httptest`, `t.Cleanup` for state reset between tests
  - **Next:** refactor into a `Store` struct with `mu` and `data` as fields, handlers become methods on `*Store`
  - After that: WAL persistence

## Key things I've learned so far
- Go error handling — return values, not exceptions
- `defer` for guaranteed cleanup
- `sync.RWMutex` — RLock for reads, Lock for writes
- Zero values — `var mu sync.RWMutex` is ready to use
- Go is multi-threaded, data races are real unlike Node
- `http.HandleFunc("GET /kv/{key}", handler)` — method in the pattern string (Go 1.22+)
- Named handler functions so tests can call them directly

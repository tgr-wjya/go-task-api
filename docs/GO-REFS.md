# Go — Practical Reference

> Written for someone coming from TypeScript + Bun + Hono. You already know REST, routing, middleware, and service/handler separation. This doc covers what's different, what's unfamiliar, and what you actually need to build.

---

## Module System

Go's equivalent of `package.json` + `bun.lock`.

```bash
go mod init github.com/yourname/projectname   # creates go.mod (like bun init)
go get github.com/some/package                # install a dependency
go mod tidy                                   # clean unused deps (like bun install --clean)
```

`go.mod` — the manifest. `go.sum` — the lockfile. Don't hand-edit either.

Running your project:

```bash
go run main.go          # like bun run main.ts
go run .                # runs the package in current dir
go build -o myapp .     # compiles to a binary
go test ./...           # run all tests recursively
```

---

## Types

Go is statically typed. No `any` by default. Types are declared explicitly or inferred via `:=`.

```go
// explicit
var name string = "mercury"
var mass float64 = 0.055

// inferred (preferred inside functions)
name := "mercury"
mass := 0.055
active := true

// zero values — Go initializes everything
var count int       // 0
var label string    // ""
var flag bool       // false
var score float64   // 0.0
```

### Basic types

| Go type | TS equivalent |
|---|---|
| `string` | `string` |
| `int` | `number` (integer) |
| `float64` | `number` (decimal) |
| `bool` | `boolean` |
| `[]string` | `string[]` |
| `map[string]int` | `Record<string, number>` |
| `interface{}` | `any` (avoid) |
| `any` | `any` (Go 1.18+, alias for `interface{}`) |

### Type declarations

```go
// named type
type Direction string

// using it
var heading Direction = "north"
```

---

## Structs

Go's equivalent of a TypeScript `interface` or `type` with data. Structs are value types — they hold data, not references (unless you use a pointer, covered later).

```go
type Planet struct {
    Name   string
    Moons  int
    Radius float64
}

// create a struct literal
earth := Planet{
    Name:   "Earth",
    Moons:  1,
    Radius: 6371.0,
}

// access fields
fmt.Println(earth.Name)   // "Earth"
earth.Moons = 2           // mutation is allowed
```

### Embedded structs

```go
type Timestamp struct {
    CreatedAt string
    UpdatedAt string
}

type Star struct {
    Name string
    Timestamp               // embedded — Star gets CreatedAt and UpdatedAt fields
}

sun := Star{Name: "Sol"}
fmt.Println(sun.CreatedAt) // accessible directly
```

### JSON struct tags

This is how you control JSON serialization. The backtick annotations are metadata — not magic, just a string the `encoding/json` package reads via reflection.

```go
type Planet struct {
    Name   string  `json:"name"`
    Moons  int     `json:"moons"`
    Hidden string  `json:"-"`          // omitted from JSON entirely
    Notes  string  `json:"notes,omitempty"` // omitted if empty/zero value
}
```

If you skip the tag, Go uses the field name as-is — `Name` becomes `"Name"` in JSON, capital N. Always add tags.

---

## Functions

```go
// basic function
func greet(name string) string {
    return "hello, " + name
}

// multiple return values — this is Go's primary error handling pattern
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// calling it
result, err := divide(10, 2)
if err != nil {
    // handle error
}
```

`nil` is Go's `null`/`undefined`. For errors: `nil` = no error. For pointers, interfaces, slices, maps: `nil` = empty/unset.

### Methods — functions on structs

```go
type Counter struct {
    Value int
}

// method with value receiver — gets a copy
func (c Counter) Current() int {
    return c.Value
}

// method with pointer receiver — can mutate the original
func (c *Counter) Increment() {
    c.Value++
}

c := Counter{Value: 0}
c.Increment()
fmt.Println(c.Current()) // 1
```

Rule of thumb: if the method needs to modify the struct, use `*Counter`. If it's read-only, `Counter` is fine. Be consistent — if any method uses a pointer receiver, use it for all.

---

## Error Handling

No try/catch. Errors are values returned from functions. You check them explicitly.

```go
file, err := os.Open("data.txt")
if err != nil {
    return err   // propagate up
}
// use file
```

### Creating errors

```go
import "errors"
import "fmt"

// simple
err := errors.New("something went wrong")

// formatted
err := fmt.Errorf("planet %s not found", name)
```

### Custom error types

```go
type NotFoundError struct {
    Resource string
    ID       string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s with id %s not found", e.Resource, e.ID)
}

// use it
return nil, &NotFoundError{Resource: "planet", ID: id}

// check type
var notFound *NotFoundError
if errors.As(err, &notFound) {
    // it's a NotFoundError
}
```

This is your `TaskNotFound` class equivalent. `errors.As` is `instanceof`.

---

## Slices

Go doesn't have arrays the way you'd think. You work with slices — a view over an underlying array.

```go
// create
planets := []string{"mercury", "venus", "earth"}

// append (returns a new slice — assign it back)
planets = append(planets, "mars")

// length
len(planets)   // 4

// iterate
for i, p := range planets {
    fmt.Println(i, p)
}

// iterate values only
for _, p := range planets {
    fmt.Println(p)
}

// slice of structs
type Planet struct{ Name string }
all := []Planet{}
all = append(all, Planet{Name: "earth"})
```

### Filter (no .filter() — use a loop)

```go
filtered := []Planet{}
for _, p := range planets {
    if p.Name != "earth" {
        filtered = append(filtered, p)
    }
}
```

### Find (no .find() — use a loop)

```go
func findByName(planets []Planet, name string) (Planet, bool) {
    for _, p := range planets {
        if p.Name == name {
            return p, true
        }
    }
    return Planet{}, false
}

p, ok := findByName(planets, "mars")
if !ok {
    // not found
}
```

`(value, bool)` is Go's idiomatic "maybe found" pattern. `ok` is the convention name — you'll see it everywhere.

---

## Maps

Go's `Record<K, V>` / `Map<K, V>`.

```go
// declare and initialize (always use make, or it's nil and will panic)
store := make(map[string]Planet)

// set
store["mars"] = Planet{Name: "mars"}

// get — returns value + existence bool
p, ok := store["mars"]
if !ok {
    // key doesn't exist
}

// delete
delete(store, "mars")

// iterate
for id, planet := range store {
    fmt.Println(id, planet.Name)
}

// length
len(store)
```

A `nil` map will panic on write. Always `make(map[K]V)` or initialize with a literal:

```go
store := map[string]Planet{
    "earth": {Name: "earth"},
}
```

---

## Pointers

Go has pointers. You will touch them, but not painfully.

```go
name := "mercury"
ptr := &name       // & = "address of" — ptr is *string
*ptr = "venus"     // * = "dereference" — modify the value at the address
fmt.Println(name)  // "venus"
```

You use pointers when:
- A function needs to modify a value (method receivers)
- You want to express "this might be nil" (optional field)
- You're passing large structs and don't want to copy them

```go
// optional field — pointer to express "might not be set"
type Config struct {
    Port    int
    Timeout *int   // nil = "not set"
}
```

You'll mostly see `*SomeStruct` in method receivers and custom error types. Don't overthink it — if the compiler complains, you'll know.

---

## Interfaces

Go interfaces are implicit — no `implements` keyword. If a type has all the methods an interface requires, it satisfies the interface automatically.

```go
type Describer interface {
    Describe() string
}

type Planet struct{ Name string }

func (p Planet) Describe() string {
    return "planet: " + p.Name
}

// Planet satisfies Describer — no declaration needed
var d Describer = Planet{Name: "earth"}
fmt.Println(d.Describe())
```

This is how you write testable services in Go — define an interface, write a concrete implementation, inject either into handlers.

```go
type PlanetStore interface {
    GetAll() []Planet
    GetByID(id string) (Planet, error)
    Add(p Planet) Planet
    Remove(id string) error
}
```

Your handler takes `PlanetStore` — doesn't care if it's in-memory or Postgres.

---

## `encoding/json`

Your `c.json()` and `c.req.valid("json")` equivalent.

### Decode (unmarshal) — request body → struct

```go
import "encoding/json"

var p Planet
err := json.NewDecoder(r.Body).Decode(&p)   // r.Body is io.Reader
if err != nil {
    // malformed JSON
}
```

### Encode (marshal) — struct → response

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(p)
```

Or if you need the bytes first:

```go
data, err := json.Marshal(p)   // returns []byte
```

### Into a map (dynamic shape)

```go
var body map[string]any
json.NewDecoder(r.Body).Decode(&body)
```

---

## `net/http` — The Stdlib Server

No framework. This is what Hono sits on top of, Go-style.

### Handler signature

Every handler is a function with this exact signature:

```go
func(w http.ResponseWriter, r *http.Request)
```

- `w` — write your response to this (`ResponseWriter`)
- `r` — read the request from this (`*Request`)

This is your `(c Context)` in Hono. Everything lives in `w` and `r`.

### Start a server

```go
mux := http.NewServeMux()

mux.HandleFunc("GET /greet", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"hello": "world"})
})

http.ListenAndServe(":3000", mux)
```

Go 1.22+ supports method+path patterns directly: `"GET /greet"`. Before 1.22, you'd check `r.Method` manually.

### Path parameters (Go 1.22+)

```go
mux.HandleFunc("GET /planets/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")   // equivalent of c.req.param("id")
    // use id
})
```

### Reading request body

```go
var body struct {
    Name  string `json:"name"`
    Moons int    `json:"moons"`
}

if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
    http.Error(w, "bad request", http.StatusBadRequest)
    return
}
```

### Writing responses

```go
// status + body
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusNotFound)        // must come before writing body
json.NewEncoder(w).Encode(map[string]string{"error": "not found"})

// shorthand for plain text error
http.Error(w, "not found", http.StatusNotFound)
```

Status code constants:

```go
http.StatusOK           // 200
http.StatusCreated      // 201
http.StatusBadRequest   // 400
http.StatusNotFound     // 404
http.StatusInternalServerError // 500
```

### Middleware

A function that takes a handler and returns a handler.

```go
func logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("%s %s\n", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)   // call the next handler
    })
}

// apply it
mux.Handle("/", logger(myHandler))
```

This is your `app.use()` equivalent — wrap handlers, call `next`.

---

## Chi Router

Drop-in on top of `net/http`. Same handler signature, adds routing ergonomics.

```bash
go get github.com/go-chi/chi/v5
```

```go
import "github.com/go-chi/chi/v5"
import "github.com/go-chi/chi/v5/middleware"

r := chi.NewRouter()

r.Use(middleware.Logger)           // global middleware
r.Use(middleware.Recoverer)        // catches panics

r.Get("/orbits", getAllHandler)
r.Post("/orbits", createHandler)
r.Get("/orbits/{id}", getOneHandler)
r.Patch("/orbits/{id}", updateHandler)
r.Delete("/orbits/{id}", deleteHandler)

// path param
id := chi.URLParam(r, "id")        // instead of r.PathValue("id")

// group routes with shared prefix
r.Route("/systems", func(r chi.Router) {
    r.Get("/", listSystems)
    r.Post("/", createSystem)
})

http.ListenAndServe(":3000", r)
```

---

## Testing

Go's stdlib test runner. No separate install.

```go
// file must end in _test.go
// function must start with Test

import "testing"

func TestAdd(t *testing.T) {
    result := add(2, 3)
    if result != 5 {
        t.Errorf("expected 5, got %d", result)
    }
}
```

### HTTP handler testing

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestGetOrbit(t *testing.T) {
    req := httptest.NewRequest("GET", "/orbits/123", nil)
    w := httptest.NewRecorder()

    getOrbitHandler(w, req)

    res := w.Result()
    if res.StatusCode != http.StatusOK {
        t.Errorf("expected 200, got %d", res.StatusCode)
    }
}
```

`httptest.NewRecorder()` is your `app.request()` / `app.handle()` equivalent. You call the handler directly without starting a real server.

### Table-driven tests (idiomatic Go)

```go
func TestDivide(t *testing.T) {
    cases := []struct {
        name     string
        a, b     float64
        want     float64
        wantErr  bool
    }{
        {"normal division", 10, 2, 5, false},
        {"divide by zero", 10, 0, 0, true},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            got, err := divide(tc.a, tc.b)
            if (err != nil) != tc.wantErr {
                t.Errorf("unexpected error: %v", err)
            }
            if got != tc.want {
                t.Errorf("want %v, got %v", tc.want, got)
            }
        })
    }
}
```

This is Go's version of `it.each`. One test function, multiple cases — each gets its own name and failure message.

---

## Project Structure

The pattern that matches what you already know:

```
myproject/
  main.go             ← entry point, wires everything together
  go.mod
  go.sum
  internal/
    models/
      planet.go       ← struct definitions, types
    store/
      memory.go       ← in-memory map, implements the interface
    service/
      planet.go       ← business logic
    handler/
      planet.go       ← HTTP handlers, calls service
```

`internal/` means the package can only be imported by code within the same module — equivalent to not exporting from a file.

Packages are directories. Everything in the same directory is the same package. Import paths are module name + directory path:

```go
import "github.com/yourname/myproject/internal/handler"
```

---

## Package Visibility

Capital letter = exported (public). Lowercase = unexported (private). That's the entire rule.

```go
type Planet struct {   // exported — other packages can use this
    Name string        // exported field
    mass float64       // unexported field — only accessible within this package
}

func GetAll() {}       // exported
func filter() {}       // unexported
```

No `export` keyword. No `public`/`private`. Just the first letter.

---

## Concurrency (brief — you'll need this eventually)

Go's concurrency model is part of the language, not a library. You don't need it for basic CRUD but you'll hit it when multiple requests modify shared state (like your in-memory map).

```go
import "sync"

type SafeStore struct {
    mu    sync.RWMutex
    items map[string]Planet
}

func (s *SafeStore) Get(id string) (Planet, bool) {
    s.mu.RLock()         // multiple readers allowed simultaneously
    defer s.mu.RUnlock() // always unlock — defer runs when function returns
    p, ok := s.items[id]
    return p, ok
}

func (s *SafeStore) Set(id string, p Planet) {
    s.mu.Lock()          // exclusive write lock
    defer s.mu.Unlock()
    s.items[id] = p
}
```

`sync.RWMutex` — read/write mutex. Multiple concurrent reads are fine. Writes are exclusive. Without this, concurrent requests to your in-memory map will cause a race condition (Go will tell you with `-race` flag).

---

## Common Patterns You'll Use Immediately

### Helper to write JSON responses (you'll make this yourself)

```go
func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
```

Call it everywhere instead of repeating the three lines.

### Dependency injection into handlers

```go
type Handler struct {
    service PlanetService
}

func NewHandler(s PlanetService) *Handler {
    return &Handler{service: s}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
    planets := h.service.GetAll()
    writeJSON(w, http.StatusOK, planets)
}
```

`h.service` is your injected dependency. This is how you wire things in `main.go` and keep handlers testable.

### Wiring in main.go

```go
func main() {
    store := store.NewMemoryStore()
    svc := service.NewPlanetService(store)
    h := handler.NewHandler(svc)

    r := chi.NewRouter()
    r.Get("/planets", h.GetAll)
    r.Post("/planets", h.Create)

    http.ListenAndServe(":3000", r)
}
```

---

## What Go Doesn't Have (That You Might Look For)

| You want | Go equivalent |
|---|---|
| `?.` optional chaining | Check `nil` explicitly |
| `??` nullish coalescing | `if x == nil { x = defaultValue }` |
| `async/await` | goroutines + channels (don't need for basic CRUD) |
| Union types | Interface or separate structs |
| Generics (limited) | Go 1.18+ has them, stdlib doesn't lean on them heavily |
| `Array.map()` | Write a loop |
| `Array.filter()` | Write a loop |
| Exceptions / try-catch | `(value, error)` return pattern |
| `nodemon` / `--watch` | `air` — `go install github.com/air-verse/air@latest` |

---

## Quick Reference Card

| Concept | Syntax |
|---|---|
| Declare + assign | `x := value` |
| Declare only | `var x Type` |
| Struct literal | `Planet{Name: "earth", Moons: 1}` |
| Pointer to struct | `&Planet{...}` |
| Dereference | `*ptr` |
| Error check | `if err != nil { ... }` |
| Slice append | `s = append(s, item)` |
| Map get | `v, ok := m[key]` |
| Map set | `m[key] = value` |
| Map delete | `delete(m, key)` |
| Range over slice | `for i, v := range s` |
| Range over map | `for k, v := range m` |
| Method on struct | `func (s *Struct) Method() {}` |
| Exported | Capital first letter |
| Unexported | lowercase first letter |
# Go Guardrails

Rules that apply when writing Go code in this repository.

## Rules

### Handle all errors
Never discard errors with `_`. Check and handle every returned error.

❌ Bad:
```go
data, _ := os.ReadFile("config.yaml")
```

✅ Good:
```go
data, err := os.ReadFile("config.yaml")
if err != nil {
    return fmt.Errorf("reading config: %w", err)
}
```

### Use `errors.Is` and `errors.As`
Use `errors.Is` and `errors.As` instead of direct comparison or type assertion on errors.

❌ Bad:
```go
if err == sql.ErrNoRows {
    return nil
}
```

✅ Good:
```go
if errors.Is(err, sql.ErrNoRows) {
    return nil
}
```

### Keep interfaces small
Define interfaces with one or two methods. Accept interfaces, return concrete types.

❌ Bad:
```go
type Store interface {
    Get(id string) (Item, error)
    List() ([]Item, error)
    Create(item Item) error
    Update(item Item) error
    Delete(id string) error
    Count() (int, error)
}
```

✅ Good:
```go
type Reader interface {
    Get(id string) (Item, error)
}

type Writer interface {
    Create(item Item) error
}
```

### Early return
Return early on errors instead of nesting logic inside `else` blocks.

❌ Bad:
```go
func process(data []byte) error {
    if len(data) > 0 {
        result, err := parse(data)
        if err == nil {
            save(result)
            return nil
        } else {
            return err
        }
    } else {
        return errors.New("empty data")
    }
}
```

✅ Good:
```go
func process(data []byte) error {
    if len(data) == 0 {
        return errors.New("empty data")
    }
    result, err := parse(data)
    if err != nil {
        return err
    }
    save(result)
    return nil
}
```

### Pass `context.Context` as the first parameter
Always pass `context.Context` as the first parameter to functions that do I/O or may be cancelled.

❌ Bad:
```go
func fetchUser(id string) (*User, error) {
    return db.Query("SELECT * FROM users WHERE id = $1", id)
}
```

✅ Good:
```go
func fetchUser(ctx context.Context, id string) (*User, error) {
    return db.QueryContext(ctx, "SELECT * FROM users WHERE id = $1", id)
}
```

### No `init()` abuse
Avoid `init()` for anything beyond simple variable registration. Use explicit initialization instead.

❌ Bad:
```go
var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("postgres", os.Getenv("DB_URL"))
    if err != nil {
        log.Fatal(err)
    }
}
```

✅ Good:
```go
func newDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("opening database: %w", err)
    }
    return db, nil
}
```

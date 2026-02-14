# Logging output example.

## Error context grouped by places where it was added.
```json
{
  "time": "2026-02-14T18:31:43.766419+03:00",
  "level": "INFO",
  "source": {
    "function": "main.main",
    "file": "/Users/d.cheremisov/Sources/mine/errors/internal/example/main.go",
    "line": 47
  },
  "msg": "hello world",
  "err": "ask to do something: failed to do something",
  "@err": {
    "CTX": {
      "@location": "/Users/d.cheremisov/Sources/mine/errors/internal/example/main.go:27",
      "pi": 3.141592653589793,
      "e": 2.718281828459045
    },
    "WRAP: ask to do something": {
      "@location": "/Users/d.cheremisov/Sources/mine/errors/internal/example/main.go:24",
      "insert-locations": true
    },
    "NEW: failed to do something": {
      "@location": "/Users/d.cheremisov/Sources/mine/errors/internal/example/main.go:19",
      "int-value": 13,
      "string-value": "world"
    }
  },
  "name": "value"
}
```

## Flat structure, error context is shown all together.
```json
{
  "time": "2026-02-14T18:31:43.766713+03:00",
  "level": "INFO",
  "source": {
    "function": "main.main",
    "file": "/Users/d.cheremisov/Sources/mine/errors/internal/example/main.go",
    "line": 60
  },
  "msg": "hello world",
  "err": "ask to do something: failed to do something",
  "@err": {
    "pi": 3.141592653589793,
    "e": 2.718281828459045,
    "insert-locations": true,
    "int-value": 13,
    "string-value": "world",
    "@locations": {
      "/Users/d.cheremisov/Sources/mine/errors/internal/example/main.go:27": "CTX",
      "/Users/d.cheremisov/Sources/mine/errors/internal/example/main.go:24": "WRAP: ask to do something",
      "/Users/d.cheremisov/Sources/mine/errors/internal/example/main.go:19": "NEW: failed to do something"
    }
  },
  "name": "value"
}
```

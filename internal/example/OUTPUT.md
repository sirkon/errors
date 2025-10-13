# Пример выхлопа логирования.

## Log with all context from errors shown all together.
```json
{
  "time": "2025-10-03T02:56:10.326348+03:00",
  "level": "INFO",
  "source": {
    "function": "main.LogFlat",
    "file": "/Users/d.cheremisov/Sources/work/errors/internal/example/main.go",
    "line": 100
  },
  "msg": "logging test with flat output",
  "int-value": 12,
  "string-value": "hello",
  "err": "ask to do something: failed to do something",
  "@err": {
    "pi": 3.141592653589793,
    "e": 2.718281828459045,
    "insert-locations": true,
    "int-value": 13,
    "string-value": "world",
    "@locations": {
      "/Users/d.cheremisov/Sources/work/errors/internal/example/main.go:49": "CTX",
      "/Users/d.cheremisov/Sources/work/errors/internal/example/main.go:46": "WRAP: ask to do something",
      "/Users/d.cheremisov/Sources/work/errors/internal/example/main.go:41": "NEW: failed to do something"
    }
  },
  "x": 1.5,
  "err-naked": "naked error message"
}
```

## Log with error context grouped by the places it was added.
```json
{
  "time": "2025-10-03T02:56:10.326488+03:00",
  "level": "INFO",
  "source": {
    "function": "main.LogGrouped",
    "file": "/Users/d.cheremisov/Sources/work/errors/internal/example/main.go",
    "line": 106
  },
  "msg": "logging test with error context grouped by the places it was added",
  "grouped-structure": true,
  "err": "ask to do something: failed to do something",
  "@err": {
    "CTX": {
      "@location": "/Users/d.cheremisov/Sources/work/errors/internal/example/main.go:49",
      "pi": 3.141592653589793,
      "e": 2.718281828459045
    },
    "WRAP: ask to do something": {
      "@location": "/Users/d.cheremisov/Sources/work/errors/internal/example/main.go:46",
      "insert-locations": true
    },
    "NEW: failed to do something": {
      "@location": "/Users/d.cheremisov/Sources/work/errors/internal/example/main.go:41",
      "int-value": 13,
      "string-value": "world"
    }
  },
  "greeting": "Hello, World!"
}
```

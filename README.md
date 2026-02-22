# errors

A Go library for working with errors that supports structured context.

## Getting started

Add the dependency with:

```shell
go get github.com/sirkon/errors@latest
```

See a minimal setup and usage example in [internal/example/main.go](internal/example/main.go).

## Features and goals

- Errors are considered to be processes, not values. This means you should not compare two errors as you would
  not try to compare two processes and do not use `errors.New` and `errors.Newf` to create sentinel errors. There are
  `errors.NewSentinel` and `errors.NewSentinelf` for this.
- Almost drop-in replacement for the standard errors package.
- Avoid the inconsistency in the standard library where you use `errors.New` but `fmt.Errorf`.
- Real first class error wrapping support with `errors.Wrap` and `errors.Wrapf`.
- Structured context support, so you don’t have to log the same error at multiple stages just to add details — simply
  attach context to the error and the extra data will be rendered by default.
- Optional inclusion of the file:line location where the error was created/handled with `errors.InsertLocations()`.

## Usage examples

```go
// You can create errors and attach context to it.
return errors.New("unexpected name").Str("expected", expected).Str("actual", name)

// You can wrap and add context.
if err != nil {
return errors.Wrap(err, "do something").Int("int", intVal).Str("str", strVal)
}

// Sometimes text annotation doesn't make a sense, but some context info does.
if err != nil {
return errors.Just(err).Int("int", intVal)
}
```

## Performance

The [benchmark](./bench_test.go) produced the following numbers:

| Operation                                                   | ns/op  | B/op | Allocs/op             |
|-------------------------------------------------------------|--------|------|-----------------------|
| Wrap.                                                       | 126.7  | 528  | 6                     |
| fmt.Errorf("…: %w")                                         | 252.3  | 296  | 9                     |
| Wrap with short context.                                    | 130.6  | 528  | 6                     |
| fmt.Errorf with text formatting matching that short context | 306.9  | 512  | 8 (strange, why -1 ?) |
| errors.Wrap with large context                              | 602.2  | 3898 | 11                    |
| fmt.Errorf with large text formatting                       | 1214.0 | 2769 | 18                    |

As you see, this library has both faster and doesn't degrade at the scale.

Now, pipeline benchmarking. We get an error, we annotate it, we log it. Four cases in here:

1. Development mode logging with slog.
2. Production mode logging with slog.
3. We log at every step of stdlib's error processing to mimic this library' behavior.
4. We just put everything in text format to deliver the context.

Here outputs look like

Dev.

```json
{
  "time": "2026-02-21T23:24:48.261254+03:00",
  "level": "ERROR",
  "msg": "log error with tree structured context",
  "err": "check error: this is an error",
  "@err": {
    "NEW: this is an error": {
      "bytes": "AQID",
      "text-bytes": "Hello World!"
    },
    "WRAP: check error": {
      "count": 333,
      "is-wrap-layer": true
    },
    "CTX": {
      "pi": 3.141592653589793,
      "e": 2.718281828459045
    }
  }
}
```

Prod.

```json
{
  "time": "2026-02-21T23:24:48.261709+03:00",
  "level": "ERROR",
  "msg": "log error with flat structured context",
  "err": "check error: this is an error",
  "@err": {
    "bytes": "AQID",
    "text-bytes": "Hello World!",
    "count": 333,
    "is-wrap-layer": true,
    "pi": 3.141592653589793,
    "e": 2.718281828459045
  }
}
```

fmt.Errorf with the logging at every step

```json lines
{
  "time": "2026-02-21T23:34:53.317258+03:00",
  "level": "ERROR",
  "msg": "failed to do something 1",
  "err": "this is an error",
  "bytes": "AQID",
  "text-bytes": "Hello World!"
}
{
  "time": "2026-02-21T23:34:53.317293+03:00",
  "level": "ERROR",
  "msg": "failed to check error",
  "err": "this is an error",
  "count": 333,
  "is-wrap-layer": true
}
{
  "time": "2026-02-21T23:34:53.317299+03:00",
  "level": "ERROR",
  "msg": "got an error",
  "err": "check error: this is an error",
  "pi": 3.141592653589793,
  "e": 2.718281828459045
}
```

fmt.Errorf with text format.

```json
{
  "time": "2026-02-21T23:29:56.447285+03:00",
  "level": "ERROR",
  "msg": "failed to do something 1",
  "err": "context pi[3.141592653589793] e[2.718281828459045]: check error count[333] is-wrap-layer[true]: this is an error bytes[[1 2 3]] text-bytes[Hello World!]"
}
```

We disabled location logging in both logger and this library (which makes it at every New, Wrap and Just).

| Test                            | ns/op Apple M4Pro | ns/op Intel 12700K on Linux | ns/op AMD Ryzen 7 5700X on Linux |
|---------------------------------|-------------------|-----------------------------|----------------------------------|
| Tree                            | 3119              | 4498                        | 5763                             |
| Flat                            | 2928              | 4005                        | 5454                             |
| fmt.Errorf and multiple logging | 7037              | 4542                        | 11731                            |
| fmt.Errorf and text format      | 2611              | 1888                        | 4624                             |

We see how fast Intel on dumb things: AVX/whatever SIMD makes wonders for text formatting and syscalls are relatively cheap with them.
Apple branch prediction is out of this world but Darwin is not on par with Linux.
And AMD … well, it is known they are not single core champions.

It is safe to say context extraction in Tree and Flat deconstructions will benefit from broader branch prediction
units on EPYCs and Zeons in the manner close to M4.

Anyway, fmt.Errorf with text format will be no-go in any half-sane environment due to observability it fails
to deliver or making it too expensive. And Flat context which is the best for observability goals is still faster
than multiple logging (which is an antipattern and just harder to do).

## Appendix.

This is how full Dev output looks like:

```json
{
  "time": "2026-02-21T23:58:48.542593+03:00",
  "level": "ERROR",
  "msg": "log error with tree structured context",
  "err": "check error: this is an error",
  "@err": {
    "NEW: this is an error": {
      "@location": "/Users/d.cheremisov/Sources/mine/errors/internal/example/example.go:16",
      "bytes": "AQID",
      "text-bytes": "Hello World!"
    },
    "WRAP: check error": {
      "@location": "/Users/d.cheremisov/Sources/mine/errors/internal/example/example.go:19",
      "count": 333,
      "is-wrap-layer": true
    },
    "CTX": {
      "@location": "/Users/d.cheremisov/Sources/mine/errors/internal/example/example.go:22",
      "pi": 3.141592653589793,
      "e": 2.718281828459045
    }
  }
}
```




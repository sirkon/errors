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

The [benchmark](./bench_test.go) produced the following numbers on Apple M4Pro:

| Operation                                                   | ns/op  | B/op | Allocs/op             |
|-------------------------------------------------------------|--------|------|-----------------------|
| Wrap.                                                       | 126.7  | 528  | 6                     |
| fmt.Errorf("…: %w")                                         | 252.3  | 296  | 9                     |
| Wrap with short context.                                    | 130.6  | 528  | 6                     |
| fmt.Errorf with text formatting matching that short context | 306.9  | 512  | 8 (strange, why -1 ?) |
| errors.Wrap with large context                              | 602.2  | 3898 | 11                    |
| fmt.Errorf with large text formatting                       | 1214.0 | 2769 | 18                    |

Intel 12700K

| Operation                                                   | ns/op  | B/op | Allocs/op |
|-------------------------------------------------------------|--------|------|-----------|
| Wrap.                                                       | 210.1  | 528  | 6         |
| fmt.Errorf("…: %w")                                         | 441.8  | 296  | 9         |
| Wrap with short context.                                    | 229.1  | 528  | 6         |
| fmt.Errorf with text formatting matching that short context | 553.4  | 512  | 8         |
| errors.Wrap with large context                              | 921.6  | 3896 | 11        |
| fmt.Errorf with large text formatting                       | 2087.0 | 2768 | 18        |

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

| Test                                       | ns/op Apple M4Pro | ns/op Intel 12700K on Linux |
|--------------------------------------------|-------------------|-----------------------------|
| Tree                                       | 3119              | 2866                        |
| Flat                                       | 2887              | 2650                        |
| fmt.Errorf and multiple logging            | 7037              | 4674                        |
| fmt.Errorf and text format                 | 2611              | 1755                        |
| Error context assemble and formatting cost | 1190              | 2056                        |
| Log write cost. Basically, syscall cost    | 1731              | 540.2                       |

As you can see, Intel is a lot dumber and spends almost twice more time on slog formatting and this became a culprit
point: a construction and then the rendition of more complex structs in slog costs almost a microsecond over dummy
text.

I actually work around it in my [sirkon/blog](https://github.com/sirkon/blog) and structured errors package
[sirkon/blog/beer](https://github.com/sirkon/blog/beer), where I got rid of formatting altogether using binary logs.
So, my `beer.Error` costs a bit more than `errors.Error` (less than `fmt.Errorf` anyway) and avoid formatting at all
by just pushing bytes. Where I have structured view when needed and this doesn't cost an arm:

| Test                         | ns/op Intel 12700K on Linux |
|------------------------------|-----------------------------|
| blog beer.Error with context | 690.6                       |
| blog text error              | 1262                        |
| slog text error              | 1772                        |

And with then the beer.Error's related output will look like

```
2026-03-05T14:24:44 INFO (/home/emacs/Sources/mine/blog/internal/playground/main.go:31) test
  ├─ text: Hello world!
  ├─ time: 2026-03-05 14:24:44.387840771 +0300 MSK
  ├─ math-constants
  │  ├─ pi: 3.141592653589793
  │  └─ e: 2.718281828459045
  ├─ duration: 16.455µs
  ├─ err: EOF
  ├─ words: [I am waiting for the spring]
  └─ err-with-ctx
     ├─ text: top check: another check: check error: EOF
     └─ @context
        ├─ WRAP: check error
        │  ├─ @location: /home/emacs/Sources/mine/blog/internal/playground/main.go:27
        │  └─ tag: tag value
        ├─ CTX
        │  ├─ @location: /home/emacs/Sources/mine/blog/internal/playground/main.go:28
        │  └─ key: 12
        └─ WRAP: another check
           ├─ @location: /home/emacs/Sources/mine/blog/internal/playground/main.go:29
           └─ bool: true
```

for this code

```go
start := time.Now()
err = beer.Wrap(io.EOF, "check error").Str("tag", "tag value")
err = beer.Just(err).Int("key", 12)
err = beer.Wrap(err, "another check").Bool("bool", true)
err = fmt.Errorf("top check: %w", err)
log.Info(context.Background(),
    "test",
    blog.Str("text", "Hello world!"),
    blog.Time("time", start),
    blog.Group("math-constants",
        blog.Flt64("pi", math.Pi),
        blog.Flt64("e", math.E),
    ),
    blog.Duration("duration", time.Since(start)),
    blog.Err(io.EOF),
    blog.Strs("words", []string{"I", "am", "waiting", "for", "the", "spring"}),
    blog.Error("err-with-ctx", err),
)
```

But here with text logs we have a classical tradeoff, where richer stuff may be too expensive.

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




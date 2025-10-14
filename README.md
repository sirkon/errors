# errors

A Go library for working with errors that supports structured context.

## Getting started

Add the dependency with:

```shell
go get github.com/sirkon/errors@latest
```

See a minimal setup and usage example in [internal/example/main.go](internal/example/main.go).

## Features and goals

- A drop-in replacement for the standard errors package.
- Avoid the inconsistency in the standard library where you use `errors.New` but `fmt.Errorf`.
- Built-in support for the most common error-wrapping pattern: "annotation: %w".
- Structured context support, so you don’t have to log the same error at multiple layers just to add details — simply attach context to the error and the extra data will be rendered by default.
- Optional inclusion of the file:line location where the error was created/handled. See [loc.go](./loc.go).
- A dedicated type for defining constant errors.

## Usage examples

```go
if err != nil {
    return errors.Wrap(err, "do something").Int("int", intVal).Str("str", strVal)
}
```

```go
const (
    Err1 errors.Const = "error 1"
    Err2 = "error 2"
    …
)

…

io.EOF = errors.New("new io.EOF") // This compiles and works.
Err1 = errors.New("new error 1") // This does not compile.
```

```go
if err != nil {
    return errors.New("new error").Loc(0)
}
```

## Performance

The [benchmark](./bench_test.go) produced the following numbers:

| Benchmark Name                           | ns/op | B/op   | Allocs/op |
|------------------------------------------|-------|--------|-----------|
| errors.Wrap                              | 188.7 | 760    | 14        |
| fmt.Errorf without text formatting       | 240.2 | 296    | 9         |
| errors.Wrapf                             | 364.6 | 1200   | 20        |
| fmt.Errorf with text formatting          | 358.8 | 576    | 12        |
| errors.Wrap with 4 context values        | 303.8 | 1248 B | 18        |
| errors.Wrap with large context           | 1024  | 4634   | 49        |
| fmt.Errorf with large context            | 1495  | 2769   | 18        |


As expected, performance is roughly on par thanks to less reliance on reflection, even in the `Wrapf` case.
It seems that `%w` itself is relatively heavy. Fewer allocations per operation do not always help.

One reason for the higher number of allocations and memory usage is the chosen design: every operation creates
a new `error.Error` object. The approach can be changed to reuse existing objects, which yields noticeably
better performance:

| Benchmark Name                       | ns/op | B/op | Allocs/op |
|--------------------------------------|-------|------|-----------|
| errors.Wrap                          | 101.5 | 312  | 6         |
| fmt.Errorf without text formatting   | 256.1 | 296  | 9         |
| errors.Wrapf                         | 307.4 | 976  | 14        |
| fmt.Errorf with text formatting      | 350.1 | 576  | 12        |

This result comes from minimal changes to the library logic. There are additional easy ways to reduce allocations
(which would effectively bring the library back to its prototype state, lol). These changes are not part of the
current codebase, so the first table should be considered authoritative.

> In any case, `errors.Wrapf` and `errors.Newf` should have no fewer allocations than `Errorf`, because they can do
> everything `Errorf("annotation: %w", err)` does plus some extra work that costs something.
>
> Memory usage will also be higher because we store more information.

### Performance when saving error locations

With saving of error locations enabled, performance drops by roughly 6–7x. This is about 1600 ns/op for `errors.Wrap`
and about 6900 ns/op for `errors.Wrap` with a long context.

## About adding structured context

Currently, the pattern is "<create/wrap error>" followed by a chain of context additions:

```go
errors.New("error").Int("int-val", intVal).Str("str-val", strVal)
```

We considered the following alternatives:

### Like in zap

```go
return errors.Wrap(err, "do something", errors.Int("int-val", intVal), errors.Str("str-val"))
```

This approach was dismissed because it makes it awkward to provide formatted error messages.

### Like in zerolog

In one of the prototypes — to be precise, in a prototype of a prototype — we tried a zerolog-like approach:

```go
return errors.Flt64("x", x).Flt64("y", y).Errorf("do %v", action)
```

This required keeping both functions and methods with the same names:

```go
func Flt64(name string, value float64) *Context {
    return NewContext().Flt64(name, value)
}

…

func (ctx *Context) Flt64(name, value) *Context {
    ctx.values = append(ctx.values, ctxTouple{name, value})
    return ctx
}
```

It worked, but the following issues emerged:

- Such expressions are harder to read, which matters given how often they appear. Because human perception is better when starting from the main thing, going down to details if needed. It is the opposite with
  this approach: start from error context (details) and only add the essence at the very end.
- The package context `gitlab.example.com/common/errors` became very large, making IDE autocompletion cumbersome.
  We shouldn’t be fighting the tool while we work.

### Summary comparison of approaches

|                     | Current                                                              | Zap                                                          | Zerolog                                                                                                                                                                                                                     |
|---------------------|----------------------------------------------------------------------|--------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Readability         | High.                                                                | High.                                                        | Significantly worse than the others.<br/>Because it reverses the order of "general → specific".<br/>Instead of "error, error details" you get<br/>"error details, error".<br/>Gets very bad when there is a lot of context. |
| Formatting support  | Yes.                                                                 | No.                                                          | Yes.                                                                                                                                                                                                                        |
| Library footprint   | Low.                                                                 | Significant.<br/>You need to maintain many helper functions. | Slightly worse than Zap.                                                                                                                                                                                                    |
| Context composition | Somewhat complicated.<br/>Requires a special object.                 | Somewhat complicated.<br/>Requires a special object.         | Supported.                                                                                                                                                                                                                  |
| Performance         | Lowest among the methods.<br/>About 2× slower than the Zerolog-like. | Medium.                                                      | Mirrors the usual Zerolog vs Zap comparison (for loggers).                                                                                                                                                                  |

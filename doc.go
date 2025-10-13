// Package errors is a drop-in replacement for the standard errors package that provides a consistent API
// and adds capabilities to make error messages more informative.
//
// Features
//
//   - Add structured context to errors in addition to plain text.
//   - Attach the location where an error was created/handled (file:line).
//   - Define constant error values via a dedicated type.
//
// # Rationale
//
// The standard [fmt.Errorf] becomes insufficient when you need to attach rich context describing
// the conditions that led to an error. Once you add more than a couple of values, the quality of the
// error message tends to suffer:
//
//   - The separation between general and specific information is lost; everything gets mixed together,
//     which makes the message harder to read.
//
// A common workaround is to log at every point where an error is observed, but this has drawbacks:
//
//   - You have to pass a logger into places that otherwise wouldn't need it.
//   - Logs become bloated, as the same error (with different annotations) is logged multiple times.
//   - Methodologically, each log line combines error + context + system-wide metadata, which hurts the signal-to-noise
//     ratio.
//
// This library lets you attach context to the error itself. As a result, logging preserves the
// separation of general and specific information without overloading the logs. The final output
// quality can be higher than with logging at every stage (depending on the logger implementation).
package errors

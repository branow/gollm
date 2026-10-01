# Working on gollm

Instructions for a coding agent changing this repository. `README.md` says what the
module does and why it is shaped this way; read it before writing code.

## Stack

Go, one module: `github.com/branow/gollm`. The root package is named `llm`, so import
it aliased: `llm "github.com/branow/gollm"`. No CLI, no main package, no framework
dependencies.

Dependencies are **pinned**. Never run a bare `go get`, never upgrade a pin as a side
effect. `go mod tidy` is fine once your code has imports. If you genuinely need a new
dependency, stop and ask rather than adding one — the whole point of this module is
that a caller can take it without taking a dependency graph.

```sh
go build ./... && go vet ./... && go test -race ./...
gofmt -l .        # must print nothing
```

## Scope

This is a provider seam: one prompt plus one JSON schema in, one structured answer
out, plus middleware that wraps a `Client` and returns a `Client`. Chat history, tool
calling, streaming, agent loops and config-file parsing are the caller's business and
do not belong here. A new provider is a new package implementing `llm.Client`; a new
cross-cutting concern is a new middleware.

The public API is a contract with other repositories. A breaking change to `Client`,
`Config`, `Request`, `Response` or `Error` is a release decision, not an
implementation detail — stop and report it rather than making it.

## Non-negotiables

- **No deployment-specific identifiers.** No real hostnames, model deployments,
  employer or product references — not in code, docs, tests, or fixtures. Use
  `example.internal`.
- **Secrets never enter a log line, an error message, or a test fixture.** An API key
  arrives on `Config` at construction time and goes nowhere else. `WithCache` stores
  the response only, never the prompt, because a prompt may carry data the caller
  would not want on disk.
- **No emoji.** Not in code, comments, docs, or commit messages.

## Code

- Errors are structured values, never formatted strings. Every error leaving this
  module is an `*llm.Error` with a `Class`, because the caller's retry and exit-code
  decisions depend on the class and string-matching a provider's prose is not a
  decision procedure.
- Rule lists are **data walked by an engine**, never an if/else ladder. Provider
  registries and error-classification tables are rows; a new case is a new row.
- Dependencies flow one way: sub-packages import the root, the root imports no
  sub-package, and sibling providers never import each other.
- Doc-comment every exported symbol: one line on what it is, plus *why* where the code
  would otherwise look arbitrary. Comments explain why, never restate the code. No
  commented-out code, no committed TODOs.
- `gofmt` clean, imports ordered, lines at most 100 columns.

## Tests

- Tests live beside the code as `*_test.go`. Whoever writes the code writes its tests.
- Table-driven. **No network and no real model call in any test.** Providers are
  exercised against an in-process HTTP server or a fake command; `fake_test.go` holds
  the in-memory `Client` the middleware tests use.
- Test the decision, not the implementation. Several tests exist to stop a past
  mistake being reintroduced — a structured-output field being renamed, a retry
  honouring `Retry-After`, a cache key missing a field. If one of those fails, the fix
  is almost certainly in the code, not the test.

## Commits

The repo builds and passes tests at **every** commit. One commit is one meaningful
step, landed with the tests that prove it.

Conventional Commits: `feat:`, `fix:`, `test:`, `refactor:`, `docs:`, `chore:`.
Lowercase first word, no body, no attribution or co-author line.

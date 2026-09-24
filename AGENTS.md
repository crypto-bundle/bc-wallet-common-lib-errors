# Development Guidelines — bc-wallet-common-lib-errors

> **Repository:** `github.com/crypto-bundle/bc-wallet-common-lib-errors`
> **Maintainer:** [@gudron (Alex V Kotelnikov)](https://github.com/gudron)
> **License:** [MIT NON-AI](./LICENSE)

## 1. Purpose & Scope

This is a Go library providing structured error formatting with multiple strategies:
- **SimplyFormatted** — basic error messages with prefix/scope
- **Scoped** — scope-based categorization (e.g., blockchain, wallet, network)
- **Valued** — value-based fields (KindDetails, KindScope, KindCode, KindPublicCode)
- **CodeContains** — search-by-error-code wrapper (depends on `bc-wallet-common-lib-tinyerrors`)

The library wraps `fmtService` implementations via the `ErrorFormatterService` interface for testability and polymorphism.

## 2. Repository Rules

| Rule | Description |
|---|---|
| Read-only by default | AI agents operate in read-first mode unless explicitly asked to modify |
| Package path | All source lives under `pkg/errformatter/`. Do not move or rename this directory. |
| Module name | Must remain `github.com/crypto-bundle/bc-wallet-common-lib-errors` |
| No new top-level dirs | Keep `README.md`, `CHANGELOG.md`, `Makefile`, `.golangci.yml`, `go.mod`, `LICENSE` at root |
| Tests alongside code | Test files must use `_test.go` suffix next to the source file they test |

## 3. Tech Stack

| Component | Value | Notes |
|---|---|---|
| Language | Go 1.23 | Current `go.mod` directive; verify before bumping |
| Dependencies | `github.com/crypto-bundle/bc-wallet-common-lib-tinyerrors v0.0.3` | Only external dep; stdlib elsewhere |
| Linter | golangci-lint | ~100+ linters enabled in `.golangci.yml` |
| Build | Makefile | Single `lint` target; runs `golangci-lint run --config .golangci.yml -v ./pkg/errformatter/` |
| Module Mode | vendor (`vendor/` directory present) | Use `go mod vendor` after dependency changes |

### Dependency Policy

Only `tinyerrors` is permitted as an external dependency (enforced by `deguard` linter rule allowing `$gostd` and `github.com/crypto-bundle/`). Adding third-party packages requires explicit justification and approval.

## 4. Code Style

### Naming Conventions

- **Structs**: PascalCase (e.g., `ErrValuedWithDefaults`, `FmtService`)
- **Interfaces**: Suffix `-er` or noun form (e.g., `ErrorFormatterService`, `fmtService`)
- **Functions**: camelCase for internal, PascalCase for exported (e.g., `newErrValuedWithDefaults`, `NewErrValuedWithDefaults`)
- **Constants**: ALL_CAPS with underscore separator (e.g., `PrefixDefault`, `KeyScope`)

### Struct Tags

Use JSON tags consistently for serialization contexts. Example:
```go
type ErrValuedWithDefaults struct {
    Message string   `json:"message"`
    // ...
}
```

### Imports

Order enforced by `gci` linter:
1. Standard library
2. Alias imports (if any)
3. Local module (`github.com/crypto-bundle/`)
4. Third-party (none currently)
5. Default / dot / blank

Never sort manually; rely on `goimports` / `gci`.

## 5. Architecture

### Interfaces

```go
// ErrorFormatterService — main public interface combining value formatting,
code attachment, and multi-method convenience wrappers.
type ErrorFormatterService interface {
    errorCodeContainable   // ErrorWithCode, ErrWithCode, NewErrorWithCode,
                           // ErrorGetCode, ErrGetCode, ErrorCodeIsOneOf, ErrCodeIsOneOf
    ErrorNoWrap(err error) error
    ErrNoWrap(err error) error
    ErrorOnly(err error, details ...string) error
    Error(err error, details ...string) error
    Errorf(err error, format string, args ...interface{}) error
    NewError(details ...string) error
    NewErrorf(format string, args ...interface{}) error
}

// ErrorFormatterBuilder — creates configured formatter instances.
type ErrorFormatterBuilder interface {
    MakeScoped(scope string) ErrorFormatterService
    MakeValued(values ...Value) ErrorFormatterService
    MakeSimply() ErrorFormatterService
    Make() ErrorFormatterService
}
```

### Service Implementations

| Internal Type | Purpose | Constructor |
|---|---|---|
| `service` | Basic default formatter (delegates to stdlib functions) | `NewErrorBasicFormatter()` |
| `serviceScoped` | Scope-prefixed errors (e.g., `[wallet] ...`) | `NewScopedErrorFormatter(scope)` |
| `*serviceValued` | Value-based fields (KindDetails, KindScope, KindCode, KindPublicCode) | via builder / direct construction |
| `serviceValuedWithDefaults` | Like Valued but always carries default values | `NewErrValuedWithDefaults(defaultValues...)` |
| *(code contains wrapper)* | Wraps inner service; searches for matching error codes | via builder |

All implement `ErrorFormatterService`.

### Value System

The `Value` struct carries typed metadata on errors:

```go
type Value struct { num Kind; any any }
type Kind uint  // KindEmpty, KindDetails, KindScope, KindCode, KindPublicCode
```

Each `Value` has a `Bits` bitmask (Set/Clear/Toggle/Has) for efficient state tracking.

### Builder Pattern

`builder.go` provides construction helpers:
```go
NewErrValuedWithDefaults(...).
    WithScope(...).
    WithCode(...).
    Build()
```

### Delegation Design

Formatters may embed other formatters or delegate to helper functions (`ValuedErrorOnly`, `ScopedErrorOnly`, etc.). The delegation chain ensures that `ErrorWithCode`, `ErrorNoWrap`, etc. all route through the appropriate formatter-specific logic.

## 6. API Conventions

### Function Pairs

Exported functions come in pairs for convenience:
- `Error*()` / `Err*()` — full error object vs short alias
- `NewError()` / `NewErrorf()` — constructor variants

Both return compatible types; aliases exist for brevity.

### Deprecation Policy

Mark deprecated APIs with a doc comment starting with `Deprecated:` and add a corresponding entry to `CHANGELOG.md`. Keep deprecated symbols for at least one minor version before removal.

## 7. Error Wrapping Rules

### Pseudo-Wrap vs Real Wrap

Most formatters produce formatted strings via `fmtService` methods — these are **pseudo-wraps** (text concatenation, not `errors.Wrap`). The underlying wrapped error may or may not be stored internally.

### ErrorNoWrap Contract

All formatters MUST implement `ErrorNoWrap()` which returns the formatted message WITHOUT wrapping any inner error. This allows callers who want clean output to bypass the chain.

### wrapcheck Configuration

The linter's `wrapcheck` ignores these sigs (defined in `.golangci.yml`):
- Custom methods: `ErrorWithCode`, `ErrWithCode`, `ErrorGetCode`, `ErrGetCode`, `ErrorNoWrap`, `ErrNoWrap`, `Errorf`, `ErrorOnly`, `Error`, `NewError`, `NewErrorf`

Do not remove these from `ignoreSigs` without understanding the wrapping strategy.

## 8. License Constraints

**MIT NON-AI** license applies. Key restrictions:
- No use for AI/ML model training
- No incorporation into datasets used for training artificial intelligence systems
- Commercial use allowed under standard MIT terms if AI-training restriction is honored

Include the LICENSE file unmodified in all distributions. Do not alter the license text.

## 9. Testing

- Use standard `testing` package (no external test frameworks)
- Table-driven subtests preferred (`t.Run("case-name", func(t *testing.T) {...})`)
- Test files live next to source: `service_valued_with_defaults.go` → `service_valued_with_defaults_test_*.go`
- Run tests: `go test ./pkg/errformatter/...`
- Linter excludes test files (`tests: false` in `.golangci.yml`)

### Coverage Expectations

Every formatter type should have:
- Constructor validation (nil checks, invalid input handling)
- Method output verification (formatted strings match expected patterns)
- Edge cases (empty fields, unicode input, deeply nested wraps)

## 10. CI / Quality Gates

### golangci-lint Highlights

~100 linters enabled across these presets:
- **bugs**: `staticcheck`, `errcheck`, `gas`
- **format**: `gofmt`, `gocritic`, `misspell`
- **complexity**: `cyclop`, `funlen`, `gocognit`
- **error**: `errorlint`, `err113`, `wrapcheck`
- **style**: `revive`, `varnamelen`, `godot`

Run locally: `make lint`

Important settings:
- `modules-download-mode: readonly` — vendor dir must be up-to-date
- `timeout: 5m` — allow for large projects
- `tests: false` — lint does NOT analyze test files

## 11. Examples

This repository does NOT contain example applications. Example implementations (HTTP server, gRPC, game engine) live in the parent/dependency repository `bc-wallet-common-lib-tinyerrors`.

To reference working usage patterns, clone `tinyerrors` separately and examine its examples directory.

## 12. Contributing

### How to Propose Changes

1. Fork `bc-wallet-common-lib-errors`
2. Create feature branch: `feat/<short-description>`
3. Ensure `make lint` passes and `go test ./pkg/errformatter/...` succeeds
4. Update `CHANGELOG.md` with entry under `Unreleased`
5. Submit PR with clear description of change and rationale

### Maintainer Checklist

- [ ] Linting passes (`make lint` → exit 0)
- [ ] All tests pass (`go test ./pkg/errformatter/...`)
- [ ] No new dependencies added (or justified)
- [ ] CHANGELOG.md updated
- [ ] README.md updated if public API changed
- [ ] License headers preserved on modified files

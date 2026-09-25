/*
Package errformatter implements structured error formatting with multiple
strategies.

It wraps fmtService implementations via the ErrorFormatterService interface,
allowing formatters to be chained and composed without tight coupling.

Interfaces

	ErrorFormatterBuilder — constructs configured formatter instances from
	  parameters such as scope string or typed values.

	ErrorFormatterService — combines value formatting, code attachment, and
	  multi-method convenience wrappers for common error-shaping scenarios.

Constants

	KindEmpty, KindDetails, KindScope, KindCode, KindPublicCode — recognized
	  Value kinds; MaxKindValue is an internal sentinel.

Types

	Bits      — uint8 bitmask with Set / Clear / Toggle / Has helpers.
	ErrorFormatterService — the primary formatter contract.  Three factory
	  functions return ready-built instances:

	    NewErrorFormatter()          — default / simply-formatted
	    NewScopedErrorFormatter(s)   — scope-prefixed
	    NewValuesErrorFormatter(vs…) — value-based

	  All other methods on ErrorFormatterService follow a consistent pattern:
	  ErrorNoWrap / ErrNoWrap, ErrorOnly / ErrOnly, Error / Err, Errorf,
	  NewError, NewErrorf, plus code-lookup helpers ErrorCodeIsOneOf /
	  ErrCodeIsOneOf, ErrorWithCode / ErrWithCode, ErrorGetCode / ErrGetCode.

	Kind       — uint alias distinguishing the five Kind variants above.
	Value      — carries typed metadata on errors; populated according to its
	             Kind discriminator.

	ErrorScoped — exported facade wrapping scopedError, returned by the
	              ScopedError* family of helpers.

# Usage recommendation

For structured error handling, use service-component formatters via
ErrorFormatterService methods rather than standalone helper functions.
This gives consistent formatting and allows dependency injection of
the formatter in tests.

Create a separate formatter instance for each distinct service scope
so that every layer prefixes errors with its own responsibility tag.
Reuse the same instance across all methods in that component.

	// wallet/service/wallet.go — package-level scoped formatter
	var walletFmt = errfmt.NewScopedErrorFormatter("wallet")

	func Transfer(ctx context.Context, from, to string, amount int64) error {
		walletFmt.ErrorOnly(db.ErrRecordNotFound, "source account not found")
		// -> "[wallet] source account not found"
	}

Choose a formatter strategy based on your needs:

	// Simply formatted — quick ad-hoc errors without metadata
	svc := errfmt.NewErrorBasicFormatter()

	// Scoped — domain-layer errors where scope prefix matters ([wallet], [network], ...)
	svc := errfmt.NewScopedErrorFormatter("wallet")

	// Valued (UNIVERSAL) — Error Scope + Code + Detail + Public Code in one type
	svc := errfmt.NewValuesErrorFormatter(vs...)

NewValuesErrorFormatter is recommended as the default when in doubt — it
is the single universal formatter carrying all four value types (details,
scope, code, public code) and inspectable at runtime via ValuedErrorGetCode.

# Agent instructions

When modifying or writing code in this package follow these rules:

**Formatter selection.** Prefer NewValuedErrorFormatter unless simplicity is required.
It carries KindDetails + KindScope + KindCode + KindPublicCode simultaneously — no
strategy-switching mid-stack is ever needed.

**One instance per scope.** Create a package-level formatter var for each responsibility
boundary (wallet, network, blockchain). Do not pass scope strings through every call site;
reuse the same configured instance across all methods in that component.

**Runtime inspection.** Use ValuedErrorGetCode(err) to extract error codes at call sites.
Errors produced by Valued* helpers satisfy errors.As(*valuedError) for rich introspection.

**Re-wrap over re-create.** When wrapping an existing *valuedError prefer MultiValuedErrorOnly
so internal Bits state is preserved instead of losing metadata on a fresh allocation.

**Pseudo-wrap vs real-wrap.** Most formatters build formatted strings via fmtService methods;
they do NOT use errors.Wrap — they produce "pseudo-wraps" (text concatenation before %w).
Use ErrorNoWrap when callers need clean output without chain accumulation.

**Godoc convention.** Every exported top-level function must carry a single-sentence godoc
comment starting with the function name. Match existing style — terse, imperative, no trailing
periods on first line. Do not document receiver methods or internal services.

**Nolint discipline.** Preserve existing //nolint directives unless their rule has changed.
Do not add blanket nolint overrides — justify each new directive inline with the affected linter(s).

**No factory additions.** All creation goes through three constructors:
NewErrorBasicFormatter(), NewScopedErrorFormatter(scope), NewValuesErrorFormatter(values…).
Adding new factories or builders requires justification and approval.

Internal Formatters (unexported)

These structs implement ErrorFormatterService but are not exposed directly.
Callers interact with them exclusively through the factory functions
or the ErrorFormatterBuilder interface.

	formatterBuilder            — builds scoped, valued, simply, or default
	                              formatters via MakeScoped / MakeValued /
	                              MakeSimply / Make.

	service                     — basic formatter delegating to standard library
	                              error creation functions.

	serviceScoped               — scope-prefixed errors (e.g. "[wallet] ...").

	serviceValuedWithDefaults   — extends serviceValued by carrying default
	                              Values across every call unless overridden.

All internal types satisfy ErrorFormatterService so they may be freely mixed
and nested behind the builder or the convenience functions listed below.

# Convenience Functions

Each strategy provides a family of top-level helpers that bypass the
interface entirely.  The most common entry points are:

	// Simply formatted (default formatter):
	err := errfmt.NewError("something went wrong")
	err := errfmt.Error(oldErr, "extra context")

	// Scope-based:
	err := errfmt.NewScopedError("wallet", "failed to sign tx")
	err := errfmt.ScopedError(oldErr, "network", "request timed out")

	// Value-based:
	v := errfmt.NewValue(errfmt.KindDetails, []string{"field"})
	err := errfmt.ValuedErrorOnly(oldErr, v)

Use the module-level documentation for the full list:

	go doc github.com/crypto-bundle/bc-wallet-common-lib-errors
*/
package errformatter

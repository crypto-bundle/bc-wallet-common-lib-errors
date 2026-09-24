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

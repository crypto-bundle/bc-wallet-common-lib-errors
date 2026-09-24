/*
Package errformatter provides structured error formatting for the bc-wallet-common-lib-errors library.

It offers four complementary strategies for formatting errors:

  - SimplyFormatted — basic error messages with an optional prefix and scope label
  - Scoped        — scope-based categorization (e.g., "wallet", "network", "blockchain")
  - Valued        — value-based fields that attach typed metadata (KindDetails, KindScope,
    KindCode, KindPublicCode) to each formatted error
  - CodeContains  — search-by-error-code wrapper that delegates to an inner formatter while
    checking whether a given error carries a matching code

All strategies implement the ErrorFormatterService interface defined in the
pkg/errformatter subpackage. Use the ErrorFormatterBuilder to construct a configured
formatter at runtime:

	builder := NewErrorBuilder()
	svc := builder.MakeScoped("wallet")

go doc github.com/crypto-bundle/bc-wallet-common-lib-errors/pkg/errformatter
→ detailed API reference for all types, interfaces, and constructors.
*/
package errformatter // import "github.com/crypto-bundle/bc-wallet-common-lib-errors"

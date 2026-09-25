# bc-wallet-common-lib-errors

## Description

Structured error formatting library for Go with multiple strategies: SimplyFormatted, Scoped, Valued, and CodeContains.

## Usage

### Quick start — pick your formatter

| Strategy | Constructor | When to use |
|----------|-------------|-------------|
| SimplyFormatted | `errfmt.NewErrorBasicFormatter()` | Ad-hoc errors with no metadata |
| Scoped | `errfmt.NewScopedErrorFormatter("wallet")` | Domain-layer errors where scope prefix matters (`[wallet]`, `[network]`, ...) |
| Valued (**Universal**, default) | `errfmt.NewValuedErrorFormatter(values...)` | **Single universal solution** — provides Error Scope, Error Code, Error Detail, AND Error Public Code simultaneously. No need to switch strategies mid-stack. |

> **Recommendation:** prefer `NewValuedErrorFormatter` as the default when in doubt. It is a **universal** error-formatting solution — it simultaneously provides **Error Scope**, **Error Code**, **Error Detail**, and **Error Public Code**. Runtime inspection via `ValuedErrorGetCode(err)` and `errors.As(*valuedError)`.

> **Note:** `NewValuesErrorFormatter` is deprecated; use `NewValuedErrorFormatter` instead (they are functionally identical).

> **Pattern:** create one configured instance per service-component scope and reuse it across all methods in that package. This avoids passing scope strings through every call site and keeps tests swappable.

### Example — scoped formatter per service

```go
// wallet/service/wallet.go
package wallet

var svc = errfmt.NewScopedErrorFormatter("wallet")

func Transfer(from, to string, amount int64) error {
    // wrapped error with "[wallet]" prefix
    return svc.ErrorOnly(db.ErrRecordNotFound, "source account not found")
}

func GetBalance(addr string) (*big.Int, error) {
    return nil, svc.NewErrorf("address %s has no balance record", addr)
}
```

### Example — valued formatter with code inspection

```go
package blockchain

var svc = errfmt.NewValuedErrorFormatter(
    errfmt.NewValue(errfmt.KindScope, "blockchain"),
)

func ValidateTx(tx *Tx) error {
    if tx.Fee < minFee {
        return svc.ValuedError(nil, []errfmt.Value{
            errfmt.NewValue(errfmt.KindCode, int(FeeTooLow)),
            errfmt.NewValue(errfmt.KindDetails, []string{"insufficient fee"}),
        })
    }
    return nil
}

// Inspect at call site
if code := errfmt.ValuedErrorGetCode(err); code == int(FeeTooLow) {
    // handle specific case
}
```

### Example — valued formatter with code operations

The valued formatter embeds `errorCodeContainable`, which provides seven methods for attaching and extracting error codes:

```go
package payment

var svc = errfmt.NewValuedErrorFormatter(
    errfmt.NewValue(errfmt.KindScope, "payment"),
)

// Attach a code to an existing error via ErrorWithCode():
return svc.ErrorWithCode(db.ErrNotFound, int(CodesAccountLocked))

// Create a fresh error with text + embedded code (no wrapping):
return svc.NewErrorWithCode("transaction rejected", int(CodesTxRejected))

// Extract the code at the call site:
if code := svc.ErrorGetCode(err); code == int(CodesAccountLocked) {
    // account is locked — prompt re-authentication
}

// Check against multiple known codes at once:
codes := []int{int(CodesLocked), int(CodesFrozen)}
if matched, found := svc.ErrorCodeIsOneOf(err, codes...); found {
    log.Warnf("matched code %d — blocking operation", matched)
}
```

## Contributors

* Author and maintainer - [@gudron (Alex V Kotelnikov)](https://github.com/gudron)

## Licence

**bc-wallet-common-lib-errors** is licensed under the [MIT NON-AI](./LICENSE) License.
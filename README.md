# bc-wallet-common-lib-errors

## Description

Structured error formatting library for Go with multiple strategies: SimplyFormatted, Scoped, Valued, and CodeContains.

## Usage

### Quick start — pick your formatter

| Strategy | Constructor | When to use |
|----------|-------------|-------------|
| SimplyFormatted | `errfmt.NewErrorBasicFormatter()` | Ad-hoc errors with no metadata |
| Scoped | `errfmt.NewScopedErrorFormatter("wallet")` | Domain-layer errors where scope prefix matters (`[wallet]`, `[network]`, ...) |
| Valued (**Universal**, default) | `errfmt.NewValuesErrorFormatter(values...)` | **Single universal solution** — provides Error Scope, Error Code, Error Detail, AND Error Public Code simultaneously. No need to switch strategies mid-stack. |

> **Recommendation:** prefer `NewValuesErrorFormatter` as the default when in doubt. It is a **universal** error-formatting solution — it simultaneously provides **Error Scope**, **Error Code**, **Error Detail**, and **Error Public Code**. Runtime inspection via `ValuedErrorGetCode(err)` and `errors.As(*valuedError)`.

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

var svc = errfmt.NewValuesErrorFormatter(
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

## Contributors

* Author and maintainer - [@gudron (Alex V Kotelnikov)](https://github.com/gudron)

## Licence

**bc-wallet-common-lib-errors** is licensed under the [MIT NON-AI](./LICENSE) License.
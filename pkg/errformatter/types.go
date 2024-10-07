/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

package errformatter

import "fmt"

type Bits uint8

const (
	ValueDetailsIsSet Bits = 1 << iota
	ValueScopeIsSet
	ValueCodeIsSet
	ValuePublicCodeIsSet
)

func (b *Bits) Set(flag Bits) {
	*b |= flag
}

func (b *Bits) Clear(flag Bits) {
	*b &^= flag
}

func (b *Bits) Toggle(flag Bits) {
	*b ^= flag
}

func (b *Bits) Has(flag Bits) bool {
	return *b&flag != 0
}

// A Value can represent any Go value, but unlike type any,
// it can represent most small values without an allocation.
// The zero Value corresponds to nil.
type Value struct {
	_ [0]func() // disallow ==
	// num holds the value for Kinds KindScope, KindCode
	num Kind
	// If any is of type Kind, then the value is in num as described above.
	any any
}

func NewValue(kind Kind, value any) Value {
	//nolint:exhaustruct //it's ok - field _ disallow struct comparison
	return Value{
		num: kind,
		any: value,
	}
}

func (v Value) Kind() Kind {
	return v.num
}

func (v Value) GetCode() int {
	if g, w := v.Kind(), KindCode; g != w {
		panic(fmt.Sprintf("Value kind is %s, not %s", g, w))
	}

	return v.getCode()
}

func (v Value) getCode() int {
	if code, ok := v.any.(int); ok {
		return code
	}

	return ValueCodeMissing
}

func (v Value) GetPublicCode() int {
	if g, w := v.Kind(), KindPublicCode; g != w {
		panic(fmt.Sprintf("Value kind is %s, not %s", g, w))
	}

	return v.getPublicCode()
}

func (v Value) getPublicCode() int {
	if code, ok := v.any.(int); ok {
		return code
	}

	return -1
}

func (v Value) GetDetails() []string {
	if g, w := v.Kind(), KindDetails; g != w {
		panic(fmt.Sprintf("Value kind is %s, not %s", g, w))
	}

	return v.getDetails()
}

func (v Value) getDetails() []string {
	if details, ok := v.any.([]string); ok {
		return details
	}

	return nil
}

func (v Value) SetDetails(details ...string) []string {
	if g, w := v.Kind(), KindDetails; g != w {
		v.num = KindDetails
		v.any = details
	}

	return details
}

func (v Value) MergeDetails(details ...string) []string {
	if g, w := v.Kind(), KindDetails; g != w {
		panic(fmt.Sprintf("Value kind is %s, not %s", g, w))
	}

	return v.addDetails(details...)
}

func (v Value) mergeDetails(newDetails ...string) []string {
	currentDetails, ok := v.any.([]string)
	if !ok {
		v.any = newDetails

		return newDetails
	}

	existMap := make(map[string]struct{}, len(currentDetails))
	merged := append(currentDetails, newDetails...)
	result := make([]string, 0, len(currentDetails))

	for i := range merged {
		detailValue := merged[i]

		_, isExists := existMap[detailValue]
		if isExists {
			continue
		}

		existMap[detailValue] = struct{}{}
		result = append(result, detailValue)
	}

	v.any = result

	return result
}

func (v Value) AddDetails(details ...string) []string {
	if g, w := v.Kind(), KindDetails; g != w {
		panic(fmt.Sprintf("Value kind is %s, not %s", g, w))
	}

	return v.addDetails(details...)
}

func (v Value) addDetails(newDetails ...string) []string {
	if currentDetails, ok := v.any.([]string); ok {
		merged := append(currentDetails, newDetails...)
		v.any = merged

		return merged
	}

	v.any = newDetails

	return nil
}

func (v Value) GetScope() string {
	if g, w := v.Kind(), KindScope; g != w {
		panic(fmt.Sprintf("Value kind is %s, not %s", g, w))
	}

	return v.getScope()
}

func (v Value) getScope() string {
	if scope, ok := v.any.(string); ok {
		return scope
	}

	return ""
}

func (v Value) SetScope(scopeName string) Value {
	if g, w := v.Kind(), KindScope; g != w {
		panic(fmt.Sprintf("Value kind is %s, not %s", g, w))
	}

	return v.setScope(scopeName)
}

func (v Value) setScope(scopeName string) Value {
	v.any = scopeName
	v.num = KindScope

	return v
}

// Kind is the kind of [Value].
type Kind uint

// The following list is sorted alphabetically, but it's also important that
// KindAny is 0 so that a zero Value represents nil.

const (
	KindEmpty Kind = iota
	KindDetails
	KindScope
	KindCode
	KindPublicCode
	// MaxKindValue - used as size of array of Value. !!!PLZ do not touch this constant.
	// This constant must be last in order of Kind constants.
	// Usage example in `valuedError` struct...
	MaxKindValue = iota - 1

	KinaEmptyName      = "kind_empty"
	KindDetailsName    = "kind_details"
	KindScopeName      = "kind_scope"
	KindCodeName       = "kind_code"
	KindPublicCodeName = "kind_public_code"
)

func (k Kind) String() string {
	switch k {
	case KindEmpty:
		return KinaEmptyName
	case KindDetails:
		return KindDetailsName
	case KindScope:
		return KindScopeName
	case KindCode:
		return KindCodeName
	case KindPublicCode:
		return KindPublicCodeName
	default:
		return KinaEmptyName
	}
}

func (k Kind) Bits() Bits {
	switch k {
	case KindEmpty:
		return 0
	case KindDetails:
		return ValueDetailsIsSet
	case KindScope:
		return ValueScopeIsSet
	case KindCode:
		return ValueCodeIsSet
	case KindPublicCode:
		return ValuePublicCodeIsSet
	default:
		return 0
	}
}

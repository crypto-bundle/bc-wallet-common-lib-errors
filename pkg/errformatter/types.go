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

type Bits uint8

const (
	ValueDetailsIsSet Bits = 1 << iota
	ValueScopeIsSet
	ValueCodeIsSet
	ValuePublicCodeIsSet
)

func (b Bits) Set(flag Bits) Bits    { return b | flag }
func (b Bits) Clear(flag Bits) Bits  { return b &^ flag }
func (b Bits) Toggle(flag Bits) Bits { return b ^ flag }
func (b Bits) Has(flag Bits) bool    { return b&flag != 0 }

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

// Kind is the kind of a [Value].
type Kind uint

// The following list is sorted alphabetically, but it's also important that
// KindAny is 0 so that a zero Value represents nil.

const (
	KindEmpty Kind = iota
	KindDetails
	KindScope
	KindCode
	KindPublicCode
)

var kindStrings = []string{
	"Emtpy",
	"Details",
	"Scope",
	"Code",
	"PublicCode",
}

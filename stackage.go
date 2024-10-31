package schemax

/*
stackage.go contains links to go-stackage for convenience, and
implements some basic methods for RFC 4512 source types such as
QuotedDescrList and unofficial types such as collection.
*/

import (
	"github.com/JesseCoretta/go-stackage"
)

var stackageList func(...int) stackage.Stack = stackage.List
var stackageBasic func(...int) stackage.Stack = stackage.Basic

func newCollection(name string) (c collection) {
	return collection(
		stackageList().
			NoPadding(true).
			SetID(name).
			SetCategory(`collection`).
			SetDelimiter(rune(10)).
			Mutex())
}

func newQDescrList(name string) (qdlist QuotedDescriptorList) {
	return QuotedDescriptorList(
		stackageList().
			SetID(name).
			Encap(`'`).
			SetCategory(`qdescrlist`).
			SetDelimiter(' ').
			Paren(true).
			Mutex())
}

func newQStringList(name string) (qstr QuotedStringList) {
	return QuotedStringList(stackageList().
		SetID(name).
		Encap(`'`).
		SetCategory(`qstrlist`).
		SetDelimiter(' ').
		Paren(true).
		Mutex())
}

func newRuleIDList(name string) RuleIDList {
	return RuleIDList(stackageList().
		SetID(name).
		SetAuxiliary(map[string]any{
			`options`:  newOpts(),
			`stringer`: nil,
		}).
		SetCategory(`ruleids`).
		SetDelimiter(`_null_`).
		Paren(true).
		Mutex())
}

func newOIDList(name string) oIDList {
	return oIDList(stackageList().
		SetID(name).
		SetAuxiliary(map[string]any{
			`options`:  newOpts(),
			`stringer`: nil,
		}).
		SetCategory(`oidlist`).
		SetDelimiter('$').
		Paren(true).
		Mutex())
}

func newExtensions(o ...Option) Extensions {
	opts := newOpts()
	for i := 0; i < len(o); i++ {
		opts.Shift(o[i])
	}

	return Extensions(stackageList().
		SetID(`extensions`).
		SetAuxiliary(map[string]any{
			`options`:  opts,
			`stringer`: nil,
		}).
		SetCategory(`extensions`).
		SetDelimiter(' ').
		Mutex())
}

func (r ObjectClasses) cast() stackage.Stack {
	return stackage.Stack(r)
}

func (r QuotedStringList) cast() stackage.Stack {
	return stackage.Stack(r)
}

func (r Extensions) cast() stackage.Stack {
	return stackage.Stack(r)
}

func (r AttributeTypes) cast() stackage.Stack {
	return stackage.Stack(r)
}

func (r MatchingRules) cast() stackage.Stack {
	return stackage.Stack(r)
}

func (r MatchingRuleUses) cast() stackage.Stack {
	return stackage.Stack(r)
}

func (r NameForms) cast() stackage.Stack {
	return stackage.Stack(r)
}

func (r LDAPSyntaxes) cast() stackage.Stack {
	return stackage.Stack(r)
}

func (r DITStructureRules) cast() stackage.Stack {
	return stackage.Stack(r)
}

func (r DITContentRules) cast() stackage.Stack {
	return stackage.Stack(r)
}

func (r QuotedDescriptorList) cast() stackage.Stack {
	return stackage.Stack(r)
}

func (r collection) cast() stackage.Stack {
	return stackage.Stack(r)
}

func (r Schema) cast() stackage.Stack {
	return stackage.Stack(r)
}

/*
Index returns the instance of string found within the receiver stack
instance at index N.  If no instance is found at the index specified,
a zero string instance is returned.
*/
func (r QuotedStringList) Index(idx int) string {
	return r.index(idx)
}

func (r QuotedStringList) index(idx int) (v string) {
	slice, _ := stackage.Stack(r).Index(idx)
	if str, ok := slice.(string); ok {
		v = str
	}

	return
}

/*
Index returns the instance of string found within the receiver stack
instance at index N.  If no instance is found at the index specified,
a zero string instance is returned.
*/
func (r Extensions) Index(idx int) Extension {
	return r.index(idx)
}

func (r Extensions) index(idx int) (v Extension) {
	slice, _ := stackage.Stack(r).Index(idx)
	if extn, ok := slice.(Extension); ok {
		v = extn
	}

	return
}

func (r QuotedStringList) String() string {
	return r.cast().String()
}

/*
List returns an instance of []string containing values derived from
the receiver instance.
*/
func (r QuotedStringList) List() (list []string) {
	for i := 0; i < r.len(); i++ {
		list = append(list, r.index(i))
	}

	return
}

/*
Contains returns a Boolean value indicative of the presence of the indicated
key within the receiver instance.
*/
func (r QuotedStringList) Contains(val string) bool {
	return r.contains(val)
}

func (r QuotedStringList) contains(val string) bool {
	for i := 0; i < r.len(); i++ {
		v := r.index(i)
		nsv := repAll(v, ` `, ``)
		vnsv := repAll(val, ` `, ``)
		if eq(nsv, vnsv) {
			return true
		}
	}

	return false
}

/*
Len returns the current integer length of the receiver instance.
*/
func (r QuotedStringList) Len() int {
	return r.len()
}

func (r QuotedStringList) len() int {
	return r.cast().Len()
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r QuotedStringList) IsZero() bool {
	return r.cast().IsZero()
}

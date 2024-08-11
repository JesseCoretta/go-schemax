package schemax

import (
	"math/rand"
	"strconv"
	"strings"
	"unicode"

	"github.com/JesseCoretta/go-schemax/internal/rfc2079"
	"github.com/JesseCoretta/go-schemax/internal/rfc2307"
	"github.com/JesseCoretta/go-schemax/internal/rfc2589"
	"github.com/JesseCoretta/go-schemax/internal/rfc2798"
	"github.com/JesseCoretta/go-schemax/internal/rfc3045"
	"github.com/JesseCoretta/go-schemax/internal/rfc3671"
	"github.com/JesseCoretta/go-schemax/internal/rfc3672"
	"github.com/JesseCoretta/go-schemax/internal/rfc4403"
	"github.com/JesseCoretta/go-schemax/internal/rfc4512"
	"github.com/JesseCoretta/go-schemax/internal/rfc4517"
	"github.com/JesseCoretta/go-schemax/internal/rfc4519"
	"github.com/JesseCoretta/go-schemax/internal/rfc4523"
	"github.com/JesseCoretta/go-schemax/internal/rfc4524"
	"github.com/JesseCoretta/go-schemax/internal/rfc4530"
	"github.com/JesseCoretta/go-schemax/internal/rfc5020"
	"github.com/JesseCoretta/go-schemax/internal/x501"

	"github.com/JesseCoretta/go-shifty"
)

var (
	atoi   func(string) (int, error)           = strconv.Atoi
	itoa   func(int) string                    = strconv.Itoa
	repAll func(string, string, string) string = strings.ReplaceAll
	eq     func(string, string) bool           = strings.EqualFold
	split  func(string, string) []string       = strings.Split
	join   func([]string, string) string       = strings.Join
	hasPfx func(string, string) bool           = strings.HasPrefix
	hasSfx func(string, string) bool           = strings.HasSuffix
	lc     func(string) string                 = strings.ToLower
	uc     func(string) string                 = strings.ToUpper
	trim   func(string, string) string         = strings.Trim
	trimL  func(string, string) string         = strings.TrimLeft
	trimR  func(string, string) string         = strings.TrimRight
	trimS  func(string) string                 = strings.TrimSpace
)

var (
	isDigit func(rune) bool = unicode.IsDigit
	isAlpha func(rune) bool = unicode.IsLetter
)

func newStringBuilder() strings.Builder {
	return strings.Builder{}
}

func newOpts() Options {
	return Options(shifty.New(shifty.Uint16))
}

func (r Options) cast() shifty.BitValue {
	return shifty.BitValue(r)
}

/*
Unshift assigns the variadic bit values of x to the receiver instance.

Variadic input can include non-negative int or explicit [Option] constant slices.

This is a fluent method.
*/
func (r Options) Shift(x ...any) Options {
	r.cast().Shift(x...)
	return r
}

/*
Unshift unassigns the variadic bit values of x to the receiver instance.

Variadic input can include non-negative int or explicit [Option] constant slices.

This is a fluent method.
*/
func (r Options) Unshift(x ...any) Options {
	r.cast().Unshift(x...)
	return r
}

/*
Positive returns a Boolean value indicative of whether option x has
been set within the receiver instance.

Input value x may be a non-negative int or explicit [Option] constants.
*/
func (r Options) Positive(x any) bool {
	return r.cast().Positive(x)
}

/*
randString is only used to generate a random string for the purpose of
temporary file/directory creation in path-based parser unit tests (see
schema_test.go).
*/
func randString(num int) string {
	bts := make([]byte, num)

	chars := `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`
	rand.Read(bts)

	for idx, bite := range bts {
		bts[idx] = chars[bite%byte(len(chars))]
	}
	return string(bts)
}

func uitoa(x any) (s string) {
	switch tv := x.(type) {
	case uint:
		s = strconv.FormatUint(uint64(tv), 10)
	case uint64:
		s = strconv.FormatUint(tv, 10)
	}

	return
}

func atoui(x string) (ui uint, ok bool) {
	_ui, err := strconv.ParseUint(x, 10, 32)

	if 0 <= _ui && _ui <= uint64(^uint(0)) {
		if err == nil {
			ok = true
			ui = uint(_ui)
		}
	}

	return
}

/*
these private variables contain sourced textual definitions from myriad
RFC implementations found within the internal package directory.  These
are used anytime a Load<Type> method is called.
*/
var (
	rfc2307Macros map[string]string = rfc2307.Macros

	x501AttributeTypes    x501.AttributeTypeDefinitions    = x501.AllAttributeTypes
	rfc2079AttributeTypes rfc2079.AttributeTypeDefinitions = rfc2079.AllAttributeTypes
	rfc2307AttributeTypes rfc2307.AttributeTypeDefinitions = rfc2307.AllAttributeTypes
	rfc2589AttributeTypes rfc2589.AttributeTypeDefinitions = rfc2589.AllAttributeTypes
	rfc2798AttributeTypes rfc2798.AttributeTypeDefinitions = rfc2798.AllAttributeTypes
	rfc3045AttributeTypes rfc3045.AttributeTypeDefinitions = rfc3045.AllAttributeTypes
	rfc3671AttributeTypes rfc3671.AttributeTypeDefinitions = rfc3671.AllAttributeTypes
	rfc3672AttributeTypes rfc3672.AttributeTypeDefinitions = rfc3672.AllAttributeTypes
	rfc4403AttributeTypes rfc4403.AttributeTypeDefinitions = rfc4403.AllAttributeTypes
	rfc4512AttributeTypes rfc4512.AttributeTypeDefinitions = rfc4512.AllAttributeTypes
	rfc4519AttributeTypes rfc4519.AttributeTypeDefinitions = rfc4519.AllAttributeTypes
	rfc4523AttributeTypes rfc4523.AttributeTypeDefinitions = rfc4523.AllAttributeTypes
	rfc4524AttributeTypes rfc4524.AttributeTypeDefinitions = rfc4524.AllAttributeTypes
	rfc4530AttributeTypes rfc4530.AttributeTypeDefinitions = rfc4530.AllAttributeTypes
	rfc5020AttributeTypes rfc5020.AttributeTypeDefinitions = rfc5020.AllAttributeTypes

	rfc2079ObjectClasses rfc2079.ObjectClassDefinitions = rfc2079.AllObjectClasses
	rfc2307ObjectClasses rfc2307.ObjectClassDefinitions = rfc2307.AllObjectClasses
	rfc2589ObjectClasses rfc2589.ObjectClassDefinitions = rfc2589.AllObjectClasses
	rfc2798ObjectClasses rfc2798.ObjectClassDefinitions = rfc2798.AllObjectClasses
	rfc3671ObjectClasses rfc3671.ObjectClassDefinitions = rfc3671.AllObjectClasses
	rfc3672ObjectClasses rfc3672.ObjectClassDefinitions = rfc3672.AllObjectClasses
	rfc4403ObjectClasses rfc4403.ObjectClassDefinitions = rfc4403.AllObjectClasses
	rfc4512ObjectClasses rfc4512.ObjectClassDefinitions = rfc4512.AllObjectClasses
	rfc4519ObjectClasses rfc4519.ObjectClassDefinitions = rfc4519.AllObjectClasses
	rfc4523ObjectClasses rfc4523.ObjectClassDefinitions = rfc4523.AllObjectClasses
	rfc4524ObjectClasses rfc4524.ObjectClassDefinitions = rfc4524.AllObjectClasses

	rfc4517LDAPSyntaxes rfc4517.LDAPSyntaxDefinitions = rfc4517.AllLDAPSyntaxes
	rfc2307LDAPSyntaxes rfc2307.LDAPSyntaxDefinitions = rfc2307.AllLDAPSyntaxes
	rfc4523LDAPSyntaxes rfc4523.LDAPSyntaxDefinitions = rfc4523.AllLDAPSyntaxes
	rfc4530LDAPSyntaxes rfc4530.LDAPSyntaxDefinitions = rfc4530.AllLDAPSyntaxes

	rfc2307MatchingRules rfc2307.MatchingRuleDefinitions = rfc2307.AllMatchingRules
	rfc4517MatchingRules rfc4517.MatchingRuleDefinitions = rfc4517.AllMatchingRules
	rfc4523MatchingRules rfc4523.MatchingRuleDefinitions = rfc4523.AllMatchingRules
	rfc4530MatchingRules rfc4530.MatchingRuleDefinitions = rfc4530.AllMatchingRules

	rfc4403NameForms rfc4403.NameFormDefinitions = rfc4403.AllNameForms

	rfc4403DITStructureRules rfc4403.DITStructureRuleDefinitions = rfc4403.AllDITStructureRules
)

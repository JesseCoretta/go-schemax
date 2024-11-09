package schemax

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

/*
This example demonstrates the means for marshaling an instance of
[AttributeType] from a map[string]any instance.
*/
func ExampleAttributeType_Marshal() {
	m := map[string]any{
		`NAME`:                 `exampleAttributeType`,
		`DESC`:                 `This is an example`,
		`NUMERICOID`:           `1.3.6.1.4.1.56521.999.12.34.56`,
		`COLLECTIVE`:           `FALSE`,
		`NO-USER-MODIFICATION`: `TRUE`,
		`SINGLE-VALUE`:         `TRUE`,
		`OBSOLETE`:             `FALSE`,
		`EQUALITY`:             `integerMatch`,
		`ORDERING`:             `integerOrderingMatch`,
		`SUBSTR`:               `caseIgnoreSubstringsMatch`,
		`SYNTAX`:               `1.3.6.1.4.1.1466.115.121.1.27`,
		`USAGE`:                `dSAOperation`,
		`X-ORIGIN`:             `RFCXXXX`,
	}

	var def AttributeType = mySchema.NewAttributeType()
	if err := def.Marshal(m); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s\n", def)
	// Output: ( 1.3.6.1.4.1.56521.999.12.34.56
	//     NAME 'exampleAttributeType'
	//     DESC 'This is an example'
	//     EQUALITY integerMatch
	//     SUBSTR caseIgnoreSubstringsMatch
	//     ORDERING integerOrderingMatch
	//     SYNTAX 1.3.6.1.4.1.1466.115.121.1.27
	//     SINGLE-VALUE
	//     NO-USER-MODIFICATION
	//     USAGE dSAOperation
	//     X-ORIGIN 'RFCXXXX' )
}

/*
This example demonstrates the means of gathering references to every
superior [AttributeType] in the relevant super type chain.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_SuperChain() {
	cn := mySchema.AttributeTypes().Get(`nisMapName`)
	chain := cn.SuperChain()

	for i := 0; i < chain.Len(); i++ {
		fmt.Println(chain.Index(i).OID())
	}

	// Output: name
}

/*
This example demonstrates the means of accessing all subordinate type
instances of the receiver instance.

In essence, this method is the opposite of the [AttributeType.SuperType]
method and may return zero (0) or more [AttributeType] instances within
the return [AttributeTypes] instance.
*/
func ExampleAttributeType_SubTypes() {
	def := mySchema.AttributeTypes().Get(`name`)
	fmt.Printf("%d subordinate types found", def.SubTypes().Len())
	// Output: 15 subordinate types found
}

/*
This example demonstrates a compliancy check of the "name" [AttributeType].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_Compliant() {
	name := mySchema.AttributeTypes().Get(`name`)
	fmt.Println(name.Compliant())
	// Output: true
}

/*
This example demonstrates a compliancy check of the "name" [AttributeType].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeTypes_Compliant() {
	attrs := mySchema.AttributeTypes()
	fmt.Println(attrs.Compliant())
	// Output: true
}

/*
This example demonstrates the means for accessing all [AttributeType]
instances which bear the specified `X-ORIGIN` extension value.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeTypes_XOrigin() {
	defs := mySchema.AttributeTypes()
	matches := defs.XOrigin(`RFC4512`)
	fmt.Printf("Matched %d of %d %s\n", matches.Len(), defs.Len(), defs.Type())
	// Output: Matched 24 of 318 attributeTypes
}

/*
This example demonstrates determining the USAGE of an [AttributeType].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_Usage() {
	ctime := mySchema.AttributeTypes().Get(`createTimestamp`)
	fmt.Println(ctime.Usage())
	// Output: directoryOperation
}

/*
This example demonstrates the means of walking the super type chain to
determine the effective [LDAPSyntax] instance held by an [AttributeType]
instance, whether direct or indirect.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_EffectiveSyntax() {
	at := mySchema.AttributeTypes().Get(`roleOccupant`)
	fmt.Println(at.EffectiveSyntax().Description())
	// Output: DN
}

/*
This example demonstrates the means of walking the super type chain to
determine the effective EQUALITY [MatchingRule] instance held by an
[AttributeType] instance, whether direct or indirect.
*/
func ExampleAttributeType_EffectiveEquality() {
	at := mySchema.AttributeTypes().Get(`registeredAddress`)
	fmt.Println(at.EffectiveEquality().OID())
	// Output: caseIgnoreListMatch
}

/*
This example demonstrates the means of walking the super type chain
to determine the effective SUBSTR [MatchingRule] instance held by an
[AttributeType] instance, whether direct or indirect.
*/
func ExampleAttributeType_EffectiveSubstring() {
	at := mySchema.AttributeTypes().Get(`registeredAddress`)
	fmt.Println(at.EffectiveSubstring().OID())
	// Output: caseIgnoreListSubstringsMatch
}

/*
This example demonstrates the means of walking the super type chain
to determine the effective ORDERING [MatchingRule] instance held by
an [AttributeType] instance, whether direct or indirect.
*/
func ExampleAttributeType_EffectiveOrdering() {
	at := mySchema.AttributeTypes().Get(`createTimestamp`)
	fmt.Println(at.EffectiveOrdering().OID())
	// Output: generalizedTimeOrderingMatch
}

/*
This example demonstrates a conventional means of checking a given
value under the terms of a specific [AttributeType]'s assigned
[ValueQualifier].

Naturally this example is overly simplified, with support extended
for nil value states purely for educational purposes only.  A real
life implementation would likely be more stringent.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_QualifyValue_withSet() {
	// Let's use "hasSubordinates" due to its common
	// use within multiple popular directory products.
	hS := mySchema.AttributeTypes().Get(`hasSubordinates`)
	if hS.IsZero() {
		fmt.Println("hasSubordinates not found!")
		return
	}

	// hasSubordinates is a BOOLEAN type, so let's make
	// a very "forgiving" handler for such values, both
	// for string and (Go) bool/*bool values.
	hS.SetValueQualifier(func(x any) error {
		var err error

		switch tv := x.(type) {
		case string:
			switch strings.ToLower(tv) {
			case `true`, `false`, `undefined`, ``:
				// OK: Match all valid string values in one shot
				// in a manner compliant with the caseIgnoreMatch
				// equality matching rule.
			default:
				// BAD: No other string value is applicable here.
				err = ErrInvalidSyntax
			}
		case bool, *bool, nil:
			// OK: Guaranteed to be valid, with a nil instance
			// equivalent to the LDAP "Undefined" BOOLEAN state.
		default:
			// BAD: no other type is applicable here.
			err = ErrInvalidType
		}

		return err
	})

	// Let's subject our newly-assigned SyntaxQualifier to
	// a series of valid values of supported types.
	for _, possibleValue := range []any{
		`True`,
		false,
		`False`,
		true,
		`FALSE`,
		`fALse`,
		``,
		`undefineD`,
		nil,
	} {
		// None of these should return errors.
		if err := hS.QualifyValue(possibleValue); err != nil {
			fmt.Println(err)
			return
		}
	}

	// Let's pass a known bogus value just to
	// make sure this thing is indeed working.
	err := hS.QualifyValue(`falsch`) // Entschuldigung, kein deutscher support :(
	fmt.Println(err)
	// Output: Value does not meet the prescribed syntax qualifications
}

/*
This example demonstrates a means of parsing a raw definition into a new
instance of [AttributeType].

Note: this example assumes a legitimate schema variable is defined in place
of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_Parse() {
	nattr := mySchema.NewAttributeType()

	// feed the parser a subtly bogus definition ...
	err := nattr.Parse(`( 1.3.6.1.4.1.56521.999.14.56.1
		NAME 'fakeAttribute'
		DESC 'It\'s not real'
		SINGLE-VALUE
		COLLECTIVE
		X-ORIGIN 'YOUR FACE'
	)`)

	fmt.Println(err)
	// Output: AttributeType is both single-valued and collective; aborting (1.3.6.1.4.1.56521.999.14.56.1)
}

/*
This example demonstrates the preferred means of initializing a new instance
of [AttributeType].  This strategy will automatically associate the receiver
instance of [Schema] with the return value.

Note: this example assumes a legitimate schema variable is defined in place
of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleSchema_NewAttributeType() {
	nattr := mySchema.NewAttributeType()
	nattr.SetNumericOID(`1.3.6.1.4.1.56521.999.14.56.1`)
	fmt.Println(nattr.NumericOID())
	// Output: 1.3.6.1.4.1.56521.999.14.56.1
}

/*
This example demonstrates an alternative to [Schema.NewAttributeType].
The return value must be manually configured and must also be manually
associated with the relevant [Schema] instance.  Use of this function
is only meaningful when dealing with multiple [Schema] instances.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleNewAttributeType() {
	// lookup and get the Directory String syntax
	dStr := mySchema.LDAPSyntaxes().Get(`1.3.6.1.4.1.1466.115.121.1.15`)
	if dStr.IsZero() {
		return
	}

	// lookup and get the caseIgnoreMatch equality matching rule
	cIM := mySchema.MatchingRules().Get(`caseIgnoreMatch`)
	if cIM.IsZero() {
		return
	}

	// prepare new var instance
	var def AttributeType = NewAttributeType()

	// set values in fluent form
	def.SetSchema(mySchema).
		SetNumericOID(`1.3.6.1.4.1.56521.999.5`).
		SetName(`cb`).
		SetDescription(`Celestial Body`).
		SetSyntax(dStr).
		SetMinimumUpperBounds(64).
		SetEquality(cIM).
		SetSingleValue().
		SetExtension(`X-ORIGIN`, `NOWHERE`).
		SetStringer() // use default stringer

	fmt.Printf("%s", def)
	// Output: ( 1.3.6.1.4.1.56521.999.5
	//     NAME 'cb'
	//     DESC 'Celestial Body'
	//     EQUALITY caseIgnoreMatch
	//     SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{64}
	//     SINGLE-VALUE
	//     X-ORIGIN 'NOWHERE' )
}

/*
This example demonstrates the replacement process of an [AttributeType]
instance within an instance of [AttributeTypes].

For reasons of oversight, we've added a custom extension X-WARNING to
remind users and admin alike of the modification.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_Replace() {

	// Lets create a new attributeType: coolattr
	attr := mySchema.NewAttributeType()
	goodraw := `( 1.3.6.1.4.1.56521.999.14.56.1 NAME 'coolattr' SUP cn )`
	if err := attr.Parse(goodraw); err != nil {
		fmt.Println(err)
		return
	}

	// Parsing says its valid, so let's push this
	// new type into our official type stack.
	mySchema.AttributeTypes().Push(attr)

	// Oh no! We realized we used the wrong supertype.
	// We wanted name, not cn :(

	// Retrieve the type
	attr = mySchema.AttributeTypes().Get(`coolattr`)

	// Craft a near identical type instance, changing
	// the supertype to name. Also, for good measure,
	// lets make a note of this modification using
	// an "X-WARNING" extension...
	nattr := mySchema.NewAttributeType().
		SetName(attr.Name()).
		SetNumericOID(attr.NumericOID()).
		SetSuperType(`name`).
		SetExtension(`X-WARNING`, `MODIFIED`). // optional
		SetStringer()

	// Replace attr with nattr, while preserving its pointer
	// address so that references within stacks do not fail.
	attr.Replace(nattr)

	// call the new one (just to be sure)
	fmt.Println(mySchema.AttributeTypes().Get(`coolattr`))
	// Output: ( 1.3.6.1.4.1.56521.999.14.56.1
	//     NAME 'coolattr'
	//     SUP name
	//     X-WARNING 'MODIFIED' )
}

/*
This example demonstrates the assignment of arbitrary data to an instance
of [AttributeType].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_SetData() {
	// The value can be any type, but we'll
	// use a string for simplicity.
	documentation := `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`

	// Obtain the target attribute type to bear
	// the assigned value.
	drink := mySchema.AttributeTypes().Get(`drink`)

	// Set it.
	drink.SetData(documentation)

	// Get it and compare to the original.
	equal := documentation == drink.Data().(string)

	fmt.Printf("Values are equal: %t", equal)
	// Output: Values are equal: true
}

func ExampleAttributeType_QualifySyntax() {

	// Obtain the syntax of interest
	dstr := mySchema.LDAPSyntaxes().Get(`directoryString`)

	// Assign a new syntax qualifier to our syntax
	// to perform a naïve assessment of x in order
	// to determine whether it is UTF8.
	dstr.SetSyntaxQualifier(func(x any) (err error) {

		// Type assert, allowing string or
		// byte values to be processed.
		switch tv := x.(type) {
		case string:
			if !utf8.ValidString(tv) {
				err = ErrInvalidSyntax
			}
		case []byte:
			if !utf8.ValidString(string(tv)) {
				err = ErrInvalidSyntax
			}
		default:
			err = ErrInvalidType
		}

		return
	})

	// Check an attribute that is known to use the above syntax
	cn := mySchema.AttributeTypes().Get(`2.5.4.3`) // or "cn"

	// Test a value against the qualifier function
	ok := cn.QualifySyntax(`Coolie McLoach`) == nil

	fmt.Printf("Syntax ok: %t", ok)
	// Output: Syntax ok: true
}

/*
This example demonstrates the means of performing a substring match
assertion between two values by way of an [AssertionMatcher] closure
assigned to the relevant [MatchingRule] instance in use by one or
more [AttributeType] instances.

For this example, we'll use the [regexp] package for brevity.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_SubstringAssertion() {

	// Obtain the syntax of interest
	cism := mySchema.MatchingRules().Get(`caseIgnoreSubstringsMatch`)

	// Assign a new assertion matcher to our matching rule to
	// allow caseless substring matching between two values.
	cism.SetAssertionMatcher(func(val, substr any) (err error) {

		// Type assert x, allowing string or
		// byte values to be processed.
		var value, substring string
		switch tv := val.(type) {
		case string:
			value = tv
		case []byte:
			value = string(tv)
		default:
			err = ErrInvalidType
			return
		}

		// Now type assert y similarly.
		switch tv := substr.(type) {
		case string:
			substring = tv
		case []byte:
			substring = string(tv)
		default:
			err = ErrInvalidType
			return
		}

		// create expression, altering wildcards
		// to conform to regexp.
		pat := strings.ReplaceAll(substring, "*", ".*")

		// Compile the expression.
		re, _ := regexp.Compile("(?i)" + pat)
		if match := re.MatchString(value); !match {
			err = ErrNoMatch
		}

		return
	})

	// Check an attribute that is known to use the above
	// matching rule.
	cn := mySchema.AttributeTypes().Get(`2.5.4.3`) // or "cn"

	// Compare two values via the SubstringAssertion method.
	// In the context of an assertion check via LDAP, the
	// first value (Kenny) could represent a value within
	// the database being compared, while the second value
	// (k*NN*) is the substring statement input by the user,
	// ostensibly within an LDAP Search Filter.
	ok := cn.SubstringAssertion(`Kenny`, `k*NN*`) == nil

	fmt.Printf("Values match: %t", ok)
	// Output: Values match: true
}

/*
This example demonstrates the means of performing an ordering match
assertion between two values by way of an [AssertionMatcher] closure
assigned to the relevant [MatchingRule] instance in use by one or
more [AttributeType] instances.

For this example, we'll be comparing two string-based timestamps in
"YYYYMMDDHHmmss" timestamp format. The values are marshaled into
proper [time.Time] instances and then compared ordering-wise.

The first input value is the higher order value, while the second value
is the lower order value. A comparison error returned indicates that the
first value is NOT greater or equal to the second.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_OrderingAssertion() {

	// Obtain the syntax of interest
	cism := mySchema.MatchingRules().Get(`generalizedTimeOrderingMatch`)

	// Assign a new assertion matcher to our matching rule to allow
	// time ordering matching (e.g.: is timeA later than timeB).
	cism.SetAssertionMatcher(func(after, before any) (err error) {

		// Type assert x, allowing string or
		// byte values to be processed.
		var A, B string
		switch tv := after.(type) {
		case string:
			A = tv
		case []byte:
			A = string(tv)
		default:
			err = ErrInvalidType
			return
		}

		// Now type assert y similarly.
		switch tv := before.(type) {
		case string:
			B = tv
		case []byte:
			B = string(tv)
		default:
			err = ErrInvalidType
			return
		}

		format := `20060102150405`

		After, _ := time.Parse(format, A)
		Before, _ := time.Parse(format, B)

		if !(After.After(Before) || (After.Equal(Before))) {
			err = ErrNoMatch
		}

		return
	})

	// Check an attribute that is known to use the above matching rule.
	modTime := mySchema.AttributeTypes().Get(`modifyTimestamp`)

	// Compare two values via the SubstringAssertion method.
	// In the context of an assertion check via LDAP, the
	// first value (Kenny) could represent a value within
	// the database being compared, while the second value
	// (k*NN*) is the substring statement input by the user,
	// ostensibly within an LDAP Search Filter.
	timeA := `20150107145309`
	timeB := `20090417110844`
	ok := modTime.OrderingAssertion(timeA, timeB) == nil

	fmt.Printf("Values match: %t", ok)
	// Output: Values match: true
}

/*
This example demonstrates the means of performing an equality match
assertion between two values by way of an [AssertionMatcher] closure
assigned to the relevant [MatchingRule] instance in use by one or
more [AttributeType] instances.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_EqualityAssertion() {

	// Obtain the syntax of interest
	cim := mySchema.MatchingRules().Get(`caseIgnoreMatch`)

	// Assign a new assertion matcher to our matching rule to
	// allow caseless equality matching between two values.
	cim.SetAssertionMatcher(func(x, y any) (err error) {

		// Type assert x, allowing string or
		// byte values to be processed.
		var X, Y string
		switch tv := x.(type) {
		case string:
			X = tv
		case []byte:
			X = string(tv)
		default:
			err = ErrInvalidType
			return
		}

		// Now type assert y similarly.
		switch tv := y.(type) {
		case string:
			Y = tv
		case []byte:
			Y = string(tv)
		default:
			err = ErrInvalidType
			return
		}

		if !strings.EqualFold(X, Y) {
			err = ErrNoMatch
		}

		return
	})

	// Check an attribute that is known to use the above
	// matching rule.
	cn := mySchema.AttributeTypes().Get(`2.5.4.3`) // or "cn"

	// Compare two values via the EqualityAssertion method.
	ok := cn.EqualityAssertion(`kenny`, `Kenny`) == nil

	fmt.Printf("Values match: %t", ok)
	// Output: Values match: true
}

/*
This example demonstrates the creation of an [Inventory] instance based
upon the current contents of an [AttributeTypes] stack instance.  Use
of an [Inventory] instance is convenient in cases where a receiver of
schema information may not be able to directly receive working stack
instances and requires a more portable and generalized type.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeTypes_Inventory() {
	at := mySchema.AttributeTypes().Inventory()
	fmt.Println(at[`2.5.4.3`][0])
	// Output: cn
}

func ExampleAttributeTypes_Type() {
	at := mySchema.AttributeTypes()
	fmt.Println(at.Type())
	// Output: attributeTypes
}

func ExampleAttributeType_Type() {
	var def AttributeType
	fmt.Println(def.Type())
	// Output: attributeType
}

/*
This example demonstrates the means of transferring an [AttributeType]
into an instance of [DefinitionMap].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_Map() {
	def := mySchema.AttributeTypes().Get(`cn`)
	fmt.Println(def.Map()[`NUMERICOID`][0]) // risky, just for simplicity
	// Output: 2.5.4.3
}

/*
This example demonstrates use of the [AttributeTypes.Maps] method, which
produces slices of [DefinitionMap] instances containing [AttributeType]
derived values

Here, we (quite recklessly) call index three (3) and reference index zero
(0) of its `SYNTAX` key to obtain the relevant [LDAPSyntax] OID string value.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeTypes_Maps() {
	defs := mySchema.AttributeTypes().Maps()
	fmt.Println(defs[3][`SYNTAX`][0]) // risky, just for simplicity
	// Output: 1.3.6.1.4.1.1466.115.121.1.24
}

/*
This example demonstrates a means of checking whether a particular instance
of [AttributeType] is present within an instance of [AttributeTypes].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeTypes_Contains() {
	attrs := mySchema.AttributeTypes()
	fmt.Println(attrs.Contains(`cn`)) // or "2.5.4.3"
	// Output: true
}

/*
This example demonstrates a means of determining whether an [AttributeType]
instance is known by the numeric OID or descriptor input.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_IsIdentifiedAs() {
	surprise := mySchema.AttributeTypes().Get(`0.9.2342.19200300.100.1.5`)
	knownBy := surprise.IsIdentifiedAs(`drink`)
	fmt.Printf("Definition is known by 'drink': %t", knownBy)
	// Output: Definition is known by 'drink': true
}

func ExampleAttributeType_NoUserModification() {
	modTime := mySchema.AttributeTypes().Get(`modifyTimestamp`)
	fmt.Printf("Definition is immutable: %t", modTime.NoUserModification())
	// Output: Definition is immutable: true
}

func ExampleAttributeType_Obsolete() {
	modTime := mySchema.AttributeTypes().Get(`modifyTimestamp`)
	fmt.Printf("Definition is obsolete: %t", modTime.Obsolete())
	// Output: Definition is obsolete: false
}

func ExampleAttributeType_Names() {
	cn := mySchema.AttributeTypes().Get(`2.5.4.3`)
	fmt.Println(cn.Names())
	// Output: ( 'cn' 'commonName' )
}

/*
This example demonstrates a means of accessing the underlying [Extensions]
stack instance within the receiver instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_Extensions() {
	cn := mySchema.AttributeTypes().Get(`cn`)
	fmt.Println(cn.Extensions())
	// Output: X-ORIGIN 'RFC4519'
}

func ExampleAttributeType_Description() {
	cn := mySchema.AttributeTypes().Get(`cn`)
	fmt.Println(cn.Description())
	// Output: RFC4519: common name(s) for which the entity is known by
}

/*
This example demonstrates use of the [AttributeType.SetStringer] method
to impose a custom [Stringer] closure over the default instance.

Naturally the end-user would opt for a more useful stringer, such as one
that produces singular CSV rows per instance.

To avoid impacting other unit tests, we reset the default stringer
via the [AttributeType.SetStringer] method again, with no arguments.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_SetStringer() {
	cn := mySchema.AttributeTypes().Get(`cn`)
	cn.SetStringer(func() string {
		return "This useless message brought to you by a dumb stringer"
	})

	msg := fmt.Sprintln(cn)
	cn.SetStringer() // return it to its previous state if need be ...

	fmt.Printf("Original: %s\nOld: %s", cn, msg)
	// Output: Original: ( 2.5.4.3
	//     NAME ( 'cn' 'commonName' )
	//     DESC 'RFC4519: common name(s) for which the entity is known by'
	//     SUP name
	//     X-ORIGIN 'RFC4519' )
	// Old: This useless message brought to you by a dumb stringer
}

/*
This example demonstrates use of the [AttributeTypes.SetStringer] method
to impose a custom [Stringer] closure upon all stack members.

Naturally the end-user would opt for a more useful stringer, such as one
that produces a CSV file containing all [AttributeType] instances.

To avoid impacting other unit tests, we reset the default stringer
via the [AttributeTypes.SetStringer] method again, with no arguments.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeTypes_SetStringer() {
	attrs := mySchema.AttributeTypes()
	attrs.SetStringer(func() string {
		return "" // make a null stringer
	})

	output := attrs.String()
	attrs.SetStringer() // return to default

	fmt.Println(output)
	// Output:
}

/*
This example demonstrates the assignment of a minimum upper bounds value,
meant to declare the maximum limit for a value of this [AttributeType].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_SetMinimumUpperBounds() {
	// First we'll craft a fake attribute
	raw := `( 1.3.6.1.4.1.56521.999.14.56.5
		NAME 'coolattr'
		SYNTAX 1.3.6.1.4.1.1466.115.121.1.26 )`

	var attr AttributeType = mySchema.NewAttributeType()
	if err := attr.Parse(raw); err != nil {
		fmt.Println(err)
		return
	}

	// Oh no! We forgot to specify the min. upper bounds!
	// No worries, it can be done after the fact:
	attr.SetMinimumUpperBounds(128)

	fmt.Println(attr.MinimumUpperBounds())
	// Output: 128
}

/*
This example demonstrates the assignment of an [LDAPSyntax] instance to
an [AttributeType] instance during assembly.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_SetSyntax() {
	// First we'll craft a fake attribute
	raw := `( 1.3.6.1.4.1.56521.999.14.56.1
                NAME 'coolattr'
		EQUALITY caseIgnoreMatch )`

	var attr AttributeType = mySchema.NewAttributeType()
	if err := attr.Parse(raw); err != nil {
		fmt.Println(err)
		return
	}

	// Oh no! We forgot to specify the desired syntax!
	// No worries, it can be done after the fact:
	attr.SetSyntax(`1.3.6.1.4.1.1466.115.121.1.26`)

	fmt.Println(attr.Syntax().Description())
	// Output: IA5 String
}

/*
This example demonstrates the assignment of an EQUALITY [MatchingRule]
instance to an [AttributeType] instance during assembly.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_SetEquality() {
	// First we'll craft a fake attribute
	raw := `( 1.3.6.1.4.1.56521.999.14.56.1
                NAME 'coolattr'
		SYNTAX 1.3.6.1.4.1.1466.115.121.1.26 )`

	var attr AttributeType = mySchema.NewAttributeType()
	if err := attr.Parse(raw); err != nil {
		fmt.Println(err)
		return
	}

	// Oh no! We forgot to specify the desired equality matching
	// rule! No worries, it can be done after the fact:
	attr.SetEquality(`caseIgnoreMatch`)

	fmt.Println(attr.Equality().NumericOID())
	// Output: 2.5.13.2
}

/*
This example demonstrates the assignment of a SUBSTR [MatchingRule]
instance to an [AttributeType] instance during assembly.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_SetSubstring() {
	// First we'll craft a fake attribute
	raw := `( 1.3.6.1.4.1.56521.999.14.56.1
                NAME 'coolattr'
                SYNTAX 1.3.6.1.4.1.1466.115.121.1.26 )`

	var attr AttributeType = mySchema.NewAttributeType()
	if err := attr.Parse(raw); err != nil {
		fmt.Println(err)
		return
	}

	// Oh no! We forgot to specify the desired substring matching
	// rule! No worries, it can be done after the fact:
	attr.SetSubstring(`caseIgnoreSubstringsMatch`)

	fmt.Println(attr.Substring().NumericOID())
	// Output: 2.5.13.4
}

/*
This example demonstrates the assignment of an ORDERING [MatchingRule]
instance to an [AttributeType] instance during assembly.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_SetOrdering() {
	// First we'll craft a fake attribute
	raw := `( 1.3.6.1.4.1.56521.999.14.56.1
                NAME 'coolattr'
                SYNTAX 1.3.6.1.4.1.1466.115.121.1.27 )`

	var attr AttributeType = mySchema.NewAttributeType()
	if err := attr.Parse(raw); err != nil {
		fmt.Println(err)
		return
	}

	// Oh no! We forgot to specify the desired ordering matching
	// rule! No worries, it can be done after the fact:
	attr.SetOrdering(`integerOrderingMatch`)

	fmt.Println(attr.Ordering().NumericOID())
	// Output: 2.5.13.15
}

/*
This example demonstrates accessing the minimum upper bounds of an instance
of [AttributeType], if set.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleAttributeType_MinimumUpperBounds() {
	ip := mySchema.AttributeTypes().Get(`ipHostNumber`)
	fmt.Println(ip.MinimumUpperBounds())
	// Output: 128
}

/*
Do stupid things to make schemax panic, gain additional
coverage in the process.
*/
func TestAttributeType_codecov(t *testing.T) {
	attr := AttributeType{}
	if err := attr.Parse(`(garbage)`); err == nil {
		t.Errorf("%s failed: expected error, got nothing", t.Name())
		return
	}
	attr = AttributeType{&attributeType{OID: `1.2.3.4.5`}}
	if err := attr.Parse(`(garbage)`); err == nil {
		t.Errorf("%s failed: expected error, got nothing", t.Name())
		return
	}
	attr = mySchema.NewAttributeType()
	goodraw := `( 1.3.6.1.4.1.56521.999.14.56.1 NAME 'coolattr' SUP cn )`
	if err := attr.Parse(goodraw); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	var oo *attributeType
	zz := AttributeType{oo}

	_ = zz.Compliant()
	_ = zz.Map()
	_ = zz.QualifyValue(``)
	_ = zz.QualifySyntax(``)
	_ = zz.EqualityAssertion(nil, rune(33))
	_ = zz.SubstringAssertion(nil, rune(33))
	_ = zz.OrderingAssertion(nil, rune(33))

	oo = new(attributeType)

	oo.OID = ``
	attr.replace(AttributeType{oo})
	oo.OID = `freakz`
	attr.replace(AttributeType{oo})

	zz = AttributeType{oo}
	zz.SetSchema(mySchema)
	_ = zz.Compliant()
	zz.setOID(`1.2.3.4.5.6.7`)
	zz.macro()

	goodraw = `( 1.3.6.1.4.1.56521.999.14.56.1 NAME 'coolattr' )`
	if err := zz.Parse(goodraw); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	_ = zz.QualifySyntax(nil)
	_ = zz.QualifyValue(nil)
	_ = zz.EqualityAssertion(nil, rune(33))
	_ = zz.SubstringAssertion(nil, rune(33))
	_ = zz.OrderingAssertion(nil, rune(33))
	_ = zz.SetSubstring(mySchema.MatchingRules().Get(`caseIgnoreSubstringsMatch`))
	_ = zz.SetOrdering(mySchema.MatchingRules().Get(`generalizedTimeOrderingMatch`))
	_ = zz.SetEquality(mySchema.MatchingRules().Get(`caseIgnoreMatch`))

	zz.setDescription(`02`)
	zz.setDescription(`aa02'`)
	zz.setDescription(`'aa02`)
	zz.setDescription(`'aa02'`)
	zz.attributeType.setMinimumUpperBounds(uint(64))

	attrs := mySchema.AttributeTypes()
	attrs.cast().NoPadding(true)
	attrs.oIDsStringerStd()
	attrs.cast().NoPadding(false)

	goodraw = `( 1.3.6.1.4.1.56521.999.14.56.1 NAME 'coolattr' SUP createTimestamp )`
	if err := zz.Parse(goodraw); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	zz.EffectiveEquality()
	zz.EffectiveOrdering()
	zz.EffectiveSubstring()

	var ats AttributeTypes
	_ = ats.Maps()
	ats.oIDsStringerPretty(1)
	ats.canPush()
	ats.canPush(nil)
	ats.canPush(rune(0))
	ats.Push(nil)
	ats.Push(rune(0))

	goodraw = `( 1.3.6.1.4.1.56521.999.14.56.1 NAME 'coolattr' SUP cn )`
	if err := zz.Parse(goodraw); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}
	zz.superChain()
	zz.setSuperType(mySchema.AttributeTypes().Get(`cn`))
	zz.SetCollective()
	zz.SetCollective()
	zz.SetObsolete()
	zz.SetObsolete()
	zz.SetNoUserModification()
	zz.SetUsage(`directoryOperation`)
	zz.Usage()
	zz.SetUsage(`distributedOperation`)
	zz.Usage()
	zz.SetUsage(`dSAOperation`)
	zz.Usage()
	zz.SetUsage(`userApplication`)
	zz.Usage()
	zz.SetUsage(`userApplications`)
	zz.Usage()
	zz.SetUsage(uint(0))
	zz.SetUsage(0)
	zz.SetUsage(1)
	zz.SetUsage(2)
	zz.SetUsage(3)

	var bmr collection = newCollection(``)
	bmr.cast().Push(AttributeType{&attributeType{OID: `1.2.3.4.5`, Collective: true, Single: true}})
	AttributeTypes(bmr).Compliant()

	badls := LDAPSyntax{&lDAPSyntax{OID: `fdsdfs`, Desc: "bad syntax"}}
	badeq := MatchingRule{&matchingRule{OID: `1.3.6.1.4.1.56521.999.96.1`, Desc: "bad equality"}}
	badord := MatchingRule{&matchingRule{OID: ``, Desc: "bad ordering"}}
	badss := MatchingRule{&matchingRule{OID: `1.3.6.1.4.1.56521.999.96.3`, Desc: "bad substring"}}

	badsch := NewEmptySchema()
	zz.attributeType.schema = badsch
	badsch.LDAPSyntaxes().cast().Push(badls)
	badsup := AttributeType{&attributeType{
		OID:    `....`,
		schema: badsch,
	}}
	badsch.AttributeTypes().cast().Push(badsup)

	zz.attributeType.SuperType = badsup
	zz.Compliant()

	badsch.LDAPSyntaxes().cast().Push(badls)
	zz.attributeType.Syntax = badls
	zz.Compliant()

	badsch.MatchingRules().cast().Push(badeq)
	badsch.MatchingRules().cast().Push(badss)
	badsch.MatchingRules().cast().Push(badord)

	zz = AttributeType{&attributeType{
		OID:       `1.3.6.1.4.1.56521.999.831.5`,
		schema:    badsch,
		Equality:  badeq,
		Substring: badss,
		Ordering:  badord,
	}}
	zz.Compliant()
}

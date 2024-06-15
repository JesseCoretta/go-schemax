package schemax

import (
	"fmt"
	"testing"
)

/*
This example demonstrates the process of parsing a raw string-based
matchingRule definition into a proper instance of [MatchingRule].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRule_Parse() {
	// Craft and push an assembled (and fictional)
	// LDAPSyntax instance into our schema.
	mySchema.LDAPSyntaxes().Push(mySchema.NewLDAPSyntax().
		SetNumericOID(`1.3.6.1.4.1.56521.999.5`).
		SetDescription(`Frequency`).
		SetExtension(`X-ORIGIN`, `NOWHERE`).
		SetExtension(`X-NOT-HUMAN-READABLE`, `TRUE`).
		SetStringer())

	var raw string = `( 1.3.6.1.4.1.56521.999.88.5 NAME 'frequencyMatch' SYNTAX 1.3.6.1.4.1.56521.999.5 X-ORIGIN 'NOWHERE' )`
	var def MatchingRule = mySchema.NewMatchingRule()
	if err := def.Parse(raw); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(def.SetStringer())
	// Output: ( 1.3.6.1.4.1.56521.999.88.5
	//     NAME 'frequencyMatch'
	//     SYNTAX 1.3.6.1.4.1.56521.999.5
	//     X-ORIGIN 'NOWHERE' )
}

/*
This example demonstrates the creation of a new [MatchingRule]
instance for manual assembly in a fluent manner.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleNewMatchingRule() {
	// Craft and push an assembled (and fictional)
	// LDAPSyntax instance into our schema.
	mySchema.LDAPSyntaxes().Push(mySchema.NewLDAPSyntax().
		SetNumericOID(`1.3.6.1.4.1.56521.999.5`).
		SetDescription(`Frequency`).
		SetExtension(`X-ORIGIN`, `NOWHERE`).
		SetStringer())

	// set values in fluent form
	def := NewMatchingRule().
		SetSchema(mySchema).
		SetNumericOID(`1.3.6.1.4.1.56521.999.88.5`).
		SetName(`frequencyMatch`).
		SetSyntax(`Frequency`).
		SetExtension(`X-ORIGIN`, `NOWHERE`).
		SetStringer()

	fmt.Printf("%s", def)
	// Output: ( 1.3.6.1.4.1.56521.999.88.5
	//     NAME 'frequencyMatch'
	//     SYNTAX 1.3.6.1.4.1.56521.999.5
	//     X-ORIGIN 'NOWHERE' )
}

/*
This example demonstrates the assignment of an [LDAPSyntax] instance
to a [MatchingRule].
*/
func ExampleMatchingRule_SetSyntax() {

	// Integer syntax
	syn := mySchema.LDAPSyntaxes().Get(`integer`)

	// set values in fluent form
	def := NewMatchingRule().
		SetSchema(mySchema).
		SetNumericOID(`1.3.6.1.4.1.56521.999.81.3`).
		SetName(`salaryMatch`).
		SetSyntax(syn).
		SetExtension(`X-ORIGIN`, `NOWHERE`).
		SetStringer()

	fmt.Println(def.Syntax().NumericOID())
	// Output: 1.3.6.1.4.1.1466.115.121.1.27
}

/*
This example demonstrates accessing the principal name of the receiver
instance.
*/
func ExampleMatchingRule_Name() {
	im := mySchema.MatchingRules().Get(`2.5.13.14`)
	fmt.Println(im.Name())
	// Output: integerMatch
}

/*
This example demonstrates accessing the numeric OID of the receiver
instance.
*/
func ExampleMatchingRule_NumericOID() {
	im := mySchema.MatchingRules().Get(`integerMatch`)
	fmt.Println(im.NumericOID())
	// Output: 2.5.13.14
}

/*
This example demonstrates accessing the OID -- whether it is the principal
name or numeric OID -- of the receiver instance.
*/
func ExampleMatchingRule_OID() {
	im := mySchema.MatchingRules().Get(`2.5.13.14`)
	fmt.Println(im.OID())
	// Output: integerMatch
}

/*
This example demonstrates the creation of a new [MatchingRule]
instance which will be replaced in memory by another. This change
will be recognized in any and all stacks in which the replaced
[MatchingRule] resides.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRule_Replace() {

	// Craft and push an assembled (and fictional)
	// LDAPSyntax instance into our schema.
	mySchema.LDAPSyntaxes().Push(mySchema.NewLDAPSyntax().
		SetNumericOID(`1.3.6.1.4.1.56521.999.5`).
		SetDescription(`Frequency`).
		SetExtension(`X-ORIGIN`, `NOWHERE`).
		SetStringer())

	// Here is our bad version
	orig := NewMatchingRule().
		SetSchema(mySchema).
		SetNumericOID(`1.3.6.1.4.1.56521.999.88.5`).
		SetName(`freakwencyMatch`).
		SetSyntax(`Frequency`).
		SetExtension(`X-OERIGIN`, `NOWHERE`).
		SetStringer()

	// Here is our good version
	good := NewMatchingRule().
		SetSchema(mySchema).
		SetNumericOID(`1.3.6.1.4.1.56521.999.88.5`).
		SetName(`frequencyMatch`).
		SetSyntax(`Frequency`).
		SetExtension(`X-ORIGIN`, `NOWHERE`).
		SetStringer()

	// Swap orig for good, but while preserving
	// the same pointer address to keep our
	// references valid.
	orig.Replace(good)

	fmt.Printf("%s", orig)
	// Output: ( 1.3.6.1.4.1.56521.999.88.5
	//     NAME 'frequencyMatch'
	//     SYNTAX 1.3.6.1.4.1.56521.999.5
	//     X-ORIGIN 'NOWHERE' )
}

/*
This example demonstrates the creation of a new [MatchingRule]
instance for manual assembly as an OBSOLETE instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRule_SetObsolete() {

	// Craft and push an assembled (and fictional)
	// LDAPSyntax instance into our schema.
	mySchema.LDAPSyntaxes().Push(mySchema.NewLDAPSyntax().
		SetNumericOID(`1.3.6.1.4.1.56521.999.5`).
		SetDescription(`Frequency`).
		SetExtension(`X-ORIGIN`, `NOWHERE`).
		SetStringer())

	// set values in fluent form
	def := NewMatchingRule().
		SetSchema(mySchema).
		SetNumericOID(`1.3.6.1.4.1.56521.999.88.5`).
		SetName(`frequencyMatch`).
		SetSyntax(`Frequency`).
		SetObsolete().
		SetExtension(`X-ORIGIN`, `NOWHERE`).
		SetStringer()

	fmt.Printf("%s", def)
	// Output: ( 1.3.6.1.4.1.56521.999.88.5
	//     NAME 'frequencyMatch'
	//     OBSOLETE
	//     SYNTAX 1.3.6.1.4.1.56521.999.5
	//     X-ORIGIN 'NOWHERE' )
}

func ExampleMatchingRule_Obsolete() {
	def := mySchema.MatchingRules().Get(`caseExactMatch`)
	fmt.Println(def.Obsolete())
	// Output: false
}

/*
This example demonstrates instant compliance checks for all [LDAPSyntax]
instances present within an instance of [MatchingRules].
*/
func ExampleMatchingRules_Compliant() {
	mrs := mySchema.MatchingRules()
	fmt.Printf("All %d %s are compliant: %t", mrs.Len(), mrs.Type(), mrs.Compliant())
	// Output: All 44 matchingRules are compliant: true
}

func ExampleMatchingRule_Data() {
	mr := mySchema.MatchingRules().Get(`caseIgnoreMatch`)

	// Let's pretend img ([]uint8) represents
	// some JPEG data (e.g.: a diagram)
	var img []uint8 = []uint8{0x1, 0x2, 0x3, 0x4}
	mr.SetData(img)

	got := mr.Data().([]uint8)

	fmt.Printf("%T, Len:%d", got, len(got))
	// Output: []uint8, Len:4
}

func ExampleMatchingRule_SetData() {
	mr := mySchema.MatchingRules().Get(`caseIgnoreMatch`)

	// Let's pretend img ([]uint8) represents
	// some JPEG data (e.g.: a diagram)
	var img []uint8 = []uint8{0x1, 0x2, 0x3, 0x4}
	mr.SetData(img)

	got := mr.Data().([]uint8)

	fmt.Printf("%T, Len:%d", got, len(got))
	// Output: []uint8, Len:4
}

/*
This example demonstrates use of the [MatchingRules.Type] method to determine
the type of stack defined within the receiver. This is mainly useful in cases
where multiple stacks are being iterated in [Definitions] interface contexts
and is more efficient when compared to manual type assertion.
*/
func ExampleMatchingRules_Type() {
	mrs := mySchema.MatchingRules()
	fmt.Printf("We have %d %s", mrs.Len(), mrs.Type())
	// Output: We have 44 matchingRules
}

/*
This example demonstrates the means of accessing the integer length of
an [MatchingRules] stack instance.
*/
func ExampleMatchingRules_Len() {
	mrs := mySchema.MatchingRules()
	fmt.Printf("We have %d %s", mrs.Len(), mrs.Type())
	// Output: We have 44 matchingRules
}

/*
This example demonstrates the means of accessing a specific slice value
within an instance of [MatchingRules] by way of its associated integer
index.
*/
func ExampleMatchingRules_Index() {
	slice := mySchema.MatchingRules().Index(3)
	fmt.Println(slice)
	// Output: ( 2.5.13.2
	//     NAME 'caseIgnoreMatch'
	//     SYNTAX 1.3.6.1.4.1.1466.115.121.1.15
	//     X-ORIGIN 'RFC4517' )
}

/*
This example demonstrates assigning and verifying a descriptive text value
within a new (and incomplete) [MatchingRule] instance.
*/
func ExampleMatchingRule_SetDescription() {
	var def MatchingRule = NewMatchingRule()
	def.SetDescription(`Important Notes`)
	fmt.Println(def.Description())
	// Output: Important Notes
}

func ExampleMatchingRule_Map() {
	def := mySchema.MatchingRules().Get(`caseIgnoreMatch`)
	fmt.Println(def.Map()[`SYNTAX`][0]) // risky, just for simplicity
	// Output: 1.3.6.1.4.1.1466.115.121.1.15
}

/*
This example demonstrates use of the [MatchingRules.Maps] method, which
produces slices of [DefinitionMap] instances born of the [MatchingRules]
stack in which they reside.  We (quite recklessly) call index three (3)
and reference index zero (0) of its `SYNTAX` key to obtain the relevant
[LDAPSyntax] OID string value.
*/
func ExampleMatchingRules_Maps() {
	defs := mySchema.MatchingRules().Maps()
	fmt.Println(defs[3][`SYNTAX`][0]) // risky, just for simplicity
	// Output: 1.3.6.1.4.1.1466.115.121.1.15
}

/*
Do stupid things to make schemax panic, gain additional
coverage in the process.
*/
func TestMatchingRule_codecov(t *testing.T) {
	_ = mySchema.MatchingRules().SetStringer().Contains(``)
	mySchema.MatchingRules().Push(rune(10))
	mySchema.MatchingRules().IsZero()
	mySchema.MatchingRules().String()
	cim := mySchema.MatchingRules().Get(`caseIgnoreMatch`)
	mySchema.MatchingRules().canPush()
	mySchema.MatchingRules().canPush(``, ``, ``, ``, cim)
	mySchema.MatchingRules().canPush(cim, cim)
	bmr := newCollection(``)
	MatchingRules(bmr.cast()).Push(NewMatchingRule().SetSchema(mySchema))
	MatchingRules(bmr.cast()).Push(NewMatchingRule().SetSchema(mySchema).SetNumericOID(`1.2.3.4.5`))
	bmr.cast().Push(NewMatchingRule().SetSchema(mySchema))
	bmr.cast().Push(NewMatchingRule().SetSchema(mySchema).SetNumericOID(`1.2.3.4.5`))

	MatchingRules(bmr).Push(NewMatchingRule().SetSchema(mySchema))
	MatchingRules(bmr).Push(NewMatchingRule().SetSchema(mySchema).SetNumericOID(`1.2.3.4.5`))
	MatchingRules(bmr).Compliant()

	var def MatchingRule

	_ = def.String()
	_ = def.SetStringer()
	_ = def.Description()
	_ = def.Name()
	_ = def.Names()
	_ = def.Extensions()
	_ = def.Syntax()
	_ = def.Schema()
	_ = def.Map()
	_ = def.Compliant()
	_ = def.macro()
	_ = def.Obsolete()

	def.setOID(`4.3.2.1`)
	var raw string = `( 2.5.13.2 NAME 'caseIgnoreMatch' SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4517' )`
	if err := def.Parse(raw); err != ErrNilReceiver {
		t.Errorf("%s failed: expected ErrNilReceiver, got %v", t.Name(), err)
		return
	}

	def = NewMatchingRule()
	def.SetDescription(`'a`)
	def.SetDescription(`'Unnecessary quoted value to be overwritten'`)

	if err := def.Parse(raw); err != ErrNilSchemaRef {
		t.Errorf("%s failed: expected ErrNilSchemaRef, got %v", t.Name(), err)
		return
	}

	// Try again. Properly.
	def.SetSchema(mySchema)
	if def.Schema().IsZero() {
		t.Errorf("%s failed: no schema reference!", t.Name())
		return
	}
	def.setStringer(func() string {
		return "blarg"
	})

	def.SetData(`fake`)
	def.SetData(nil)
	def.Data()

	if err := def.Parse(raw); err != nil {
		t.Errorf("%s failed: expected success, got %v", t.Name(), err)
		return
	}
	_ = def.macro()
	def.setOID(`2.5.13.2`)

	var def2 MatchingRule
	_ = def2.Replace(def) // will fail

}

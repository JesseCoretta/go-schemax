package schemax

import (
	"fmt"
	"testing"
)

/*
This example demonstrates a compliancy check of the "account" [ObjectClass].
*/
func ExampleDITContentRule_Compliant() {
	dc := mySchema.DITContentRules().Index(0)
	fmt.Println(dc.Compliant())
	// Output: true
}

/*
This example demonstrates a compliancy check of all [ObjectClasses] members.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITContentRules_Compliant() {
	rules := mySchema.DITContentRules()
	fmt.Println(rules.Compliant())
	// Output: true
}

/*
This example demonstrates use of the [ObjectClass.SetStringer] method
to impose a custom [Stringer] closure over the default instance.

Naturally the end-user would opt for a more useful stringer, such as one
that produces singular CSV rows per instance.

To avoid impacting other unit tests, we reset the default stringer
via the [ObjectClass.SetStringer] method again, with no arguments.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITContentRule_SetStringer() {
	opers := mySchema.DITContentRules().Get(`rootArcContent`)
	opers.SetStringer(func() string {
		return "This useless message brought to you by a dumb stringer"
	})

	msg := fmt.Sprintln(opers)
	opers.SetStringer() // return it to its previous state if need be ...

	fmt.Printf("Original: %s\nOld: %s", opers, msg)
	// Output: Original: ( 1.3.6.1.4.1.56521.101.2.5.2
	//     NAME 'rootArcContent'
	//     DESC 'root arc entry content rule'
	//     AUX ( iSORegistration
	//         $ iTUTRegistration
	//         $ jointISOITUTRegistration )
	//     MUST ( aSN1Notation
	//          $ iRI
	//          $ identifier
	//          $ n
	//          $ unicodeValue )
	//     MAY ( additionalUnicodeValue
	//         $ nameAndNumberForm
	//         $ secondaryIdentifier
	//         $ standardizedNameForm )
	//     NOT ( dotNotation
	//         $ registrationRange
	//         $ registrationStatus
	//         $ supArc
	//         $ topArc )
	//     X-ORIGIN 'draft-coretta-oiddir-schema; unofficial supplement'
	//     X-WARNING 'UNOFFICIAL' )
	// Old: This useless message brought to you by a dumb stringer
}

/*
This example demonstrates use of the [ObjectClasses.SetStringer] method
to impose a custom [Stringer] closure upon all stack members.

Naturally the end-user would opt for a more useful stringer, such as one
that produces a CSV file containing all [ObjectClass] instances.

To avoid impacting other unit tests, we reset the default stringer
via the [ObjectClasses.SetStringer] method again, with no arguments.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITContentRules_SetStringer() {
	attrs := mySchema.DITContentRules()
	attrs.SetStringer(func() string {
		return "" // make a null stringer
	})

	output := attrs.String()
	attrs.SetStringer() // return to default

	fmt.Println(output)
	// Output:
}

/*
This example demonstrates the assignment of arbitrary data to an instance
of [ObjectClass].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITContentRule_SetData() {
	// The value can be any type, but we'll
	// use a string for simplicity.
	documentation := `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`

	// Obtain the target attribute type to bear
	// the assigned value.
	dvc := mySchema.DITContentRules().Index(0)

	// Set it.
	dvc.SetData(documentation)

	// Get it and compare to the original.
	equal := documentation == dvc.Data().(string)

	fmt.Printf("Values are equal: %t", equal)
	// Output: Values are equal: true
}

func ExampleDITContentRule_IsIdentifiedAs() {
	oc := mySchema.DITContentRules().Get(`rootArcContent`)
	fmt.Println(oc.IsIdentifiedAs(`1.3.6.1.4.1.56521.101.2.5.2`))
	// Output: true
}

/*
 */
func ExampleDITContentRule_Replace() {
	gon := mySchema.DITContentRules().Index(0)

	// Craft a near identical groupOfNames instance,
	// save for the one change we intend to make.
	ngon := mySchema.NewDITContentRule().
		SetName(gon.Name()).
		SetNumericOID(gon.NumericOID()).
		SetDescription(gon.Description()).
		SetMust(`cn`).
		SetMay(`member`, `businessCategory`, `seeAlso`, `owner`, `ou`, `o`, `description`).
		SetNot(`drink`).
		SetExtension(`X-ORIGIN`, `NOWHERE`).
		SetExtension(`X-WARNING`, `MODIFIED`). // optional
		SetStringer()

	// Replace gon with ngon, while preserving its pointer
	// address so that references within stacks do not fail.
	gon.Replace(ngon)

	// call the new one (just to be sure)
	fmt.Println(mySchema.DITContentRules().Index(0))
	// Output: ( 1.3.6.1.4.1.56521.101.2.5.2
	//     NAME 'rootArcContent'
	//     DESC 'root arc entry content rule'
	//     MUST cn
	//     MAY ( member
	//         $ businessCategory
	//         $ seeAlso
	//         $ owner
	//         $ ou
	//         $ o
	//         $ description )
	//     NOT drink
	//     X-ORIGIN 'NOWHERE'
	//     X-WARNING 'MODIFIED' )
}

func ExampleDITContentRule_SetObsolete() {
	fake := NewDITContentRule().
		SetNumericOID(`1.3.6.1.4.1.56521.999.108.4`).
		SetName(`obsoleteClass`).
		SetObsolete()

	fmt.Println(fake.Obsolete())
	// Output: true
}

/*
This example demonstrates a means of checking whether a particular instance
of [ObjectClass] is present within an instance of [ObjectClasses].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITContentRules_Contains() {
	rules := mySchema.DITContentRules()
	fmt.Println(rules.Contains(`rootArcContent`)) // or "2.5.6.0"
	// Output: true
}

func ExampleDITContentRules_Inventory() {
	dc := mySchema.DITContentRules().Inventory()
	fmt.Println(dc[`1.3.6.1.4.1.56521.101.2.5.2`][0])
	// Output: rootArcContent
}

func ExampleDITContentRules_Type() {
	oc := mySchema.DITContentRules()
	fmt.Println(oc.Type())
	// Output: dITContentRules
}

func ExampleDITContentRule_Type() {
	var def DITContentRule
	fmt.Println(def.Type())
	// Output: dITContentRule
}

func ExampleDITContentRule_Map() {
	def := mySchema.DITContentRules().Index(0)
	fmt.Println(def.Map()[`NUMERICOID`][0]) // risky, just for simplicity
	// Output: 1.3.6.1.4.1.56521.101.2.5.2
}

/*
This example demonstrates use of the [ObjectClasses.Maps] method, which
produces slices of [DefinitionMap] instances born of the [ObjectClasses]
stack in which they reside.  We (quite recklessly) call index three (3)
and reference index zero (0) of its `SYNTAX` key to obtain the relevant
[LDAPSyntax] OID string value.
*/
func ExampleDITContentRules_Maps() {
	defs := mySchema.DITContentRules().Maps()
	fmt.Println(defs[0][`NUMERICOID`][0]) // risky, just for simplicity
	// Output: 1.3.6.1.4.1.56521.101.2.5.2
}

/*
This example demonstrates the manual (non-parsed) assembly of a new
[ObjectClass] instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleNewDITContentRule() {
	oc := NewDITContentRule() // Important! Initializes internal stacks

	// Conveniently input values in fluent form ...
	oc.SetSchema(mySchema).
		SetName(`engineeringPersonnel`).
		SetDescription(`EP-46: Engineering employee`).
		SetNumericOID(`0.9.2342.19200300.100.4.5`).
		SetMust(`uid`).
		SetMay(`sn`, `cn`, `l`, `st`, `c`, `co`).
		SetNot(`drink`).
		SetExtension(`X-ORIGIN`, `NOWHERE`).
		SetStringer() // use default stringer

	fmt.Println(oc)
	// Output: ( 0.9.2342.19200300.100.4.5
	//     NAME 'engineeringPersonnel'
	//     DESC 'EP-46: Engineering employee'
	//     MUST uid
	//     MAY ( sn
	//         $ cn
	//         $ l
	//         $ st
	//         $ c
	//         $ co )
	//     NOT drink
	//     X-ORIGIN 'NOWHERE' )
}

/*
Do stupid things to make schemax panic, gain additional
coverage in the process.
*/
func TestDITContentRule_codecov(t *testing.T) {
	_ = mySchema.DITContentRules().SetStringer().Contains(``)
	mySchema.DITContentRules().Push(rune(10))
	mySchema.DITContentRules().IsZero()
	_ = mySchema.DITContentRules().String()
	cim := mySchema.DITContentRules().Get(`account`)
	mySchema.DITContentRules().canPush()
	mySchema.DITContentRules().canPush(``, ``, ``, ``, cim)
	mySchema.DITContentRules().canPush(cim, cim)
	bmr := newCollection(``)
	bma := newCollection(``)
	ObjectClasses(bmr.cast()).Push(NewDITContentRule().SetSchema(mySchema))
	ObjectClasses(bmr.cast()).Push(NewDITContentRule().SetSchema(mySchema).SetNumericOID(`1.2.3.4.5`))
	bmr.cast().Push(NewDITContentRule().SetSchema(mySchema))
	bmr.cast().Push(NewDITContentRule().SetSchema(mySchema).SetNumericOID(`1.2.3.4.5`))
	var bad ObjectClass
	bmr.cast().Push(bad)

	ObjectClasses(bmr.cast()).oIDsStringerPretty(0)
	ObjectClasses(bmr.cast()).oIDsStringerStd()
	ObjectClasses(bmr.cast()).canPush()
	ObjectClasses(bmr.cast()).canPush(`things`)
	var ocs ObjectClasses
	ocs.oIDsStringerPretty(0)
	ocs.oIDsStringerStd()
	ocs.canPush(`forks`)
	ocs.Push(NewDITContentRule().SetSchema(mySchema))
	bmr.cast().Push(AttributeType{&attributeType{OID: `1.2.3.4.5`, Collective: true, Single: true}})
	bma.cast().Push(AttributeType{&attributeType{OID: ``, Collective: true, Single: true}})
	xoc := ObjectClass{&objectClass{
		Must: AttributeTypes(bmr),
	}}
	yoc := ObjectClass{&objectClass{
		May: AttributeTypes(bma),
	}}

	xoc.Compliant()
	yoc.Compliant()

	ocs.Push(bad)

	ObjectClasses(bmr).Push(NewDITContentRule().SetSchema(mySchema))
	ObjectClasses(bmr).Push(NewDITContentRule().SetSchema(mySchema).SetNumericOID(`1.2.3.4.5`))
	ObjectClasses(bmr).Compliant()
	mySchema.DITContentRules().Compliant()

	var def DITContentRule

	_ = def.String()
	_ = def.SetStringer()
	_ = def.Description()
	_ = def.Name()
	_ = def.Names()
	_ = def.Extensions()
	_ = def.Must()
	_ = def.May()
	_ = def.Not()
	_ = def.Schema()
	_ = def.Map()
	_ = def.Compliant()
	_ = def.macro()
	_ = def.Obsolete()

	def.setOID(`4.3.2.1`)
	var raw string = `( 1.3.6.1.4.1.56521.101.2.5.2
                NAME 'rootArcContent'
                DESC 'root arc entry content rule'
                MUST ( aSN1Notation
                     $ iRI
                     $ identifier
                     $ n
                     $ unicodeValue )
                MAY ( additionalUnicodeValue
                    $ nameAndNumberForm
                    $ secondaryIdentifier
                    $ standardizedNameForm )
                NOT ( dotNotation
                    $ registrationRange
                    $ registrationStatus
                    $ supArc
                    $ topArc )
                X-ORIGIN 'draft-coretta-oiddir-schema; unofficial supplement'
		X-WARNING 'UNOFFICIAL' )`

	if err := def.Parse(raw); err != ErrNilReceiver {
		t.Errorf("%s failed: expected ErrNilReceiver, got %v", t.Name(), err)
		return
	}

	def = NewDITContentRule()
	def.SetDescription(`'a`)
	def.SetDescription(`'Unnecessary quoted value to be overwritten'`)

	oo := new(dITContentRule)
	oo.OID = mySchema.ObjectClasses().Get(`device`)
	def.replace(DITContentRule{oo})

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

	def.SetAux(mySchema.ObjectClasses().Get(`dcObject`))
	def.SetAux(`dcObject`)
	def.SetAux(rune(8))
	def.SetMust(mySchema.AttributeTypes().Get(`cn`))
	def.SetMust(rune(11))
	def.SetMay(mySchema.AttributeTypes().Get(`sn`))
	def.SetMay(rune(11))
	def.SetNot(rune(8))
	def.SetNot(`l`)
	def.SetNot(mySchema.AttributeTypes().Get(`l`))
	def.Map()
	mySchema.DITContentRules().canPush(DITContentRule{}, DITContentRule{new(dITContentRule)})
	if err := def.Parse(raw); err != nil {
		t.Errorf("%s failed: expected success, got %v", t.Name(), err)
		return
	}
	_ = def.macro()
	def.setOID(`2.5.6.2`)

	def.SetData(`fake`)
	def.SetData(nil)
	def.Data()

	auxs := NewObjectClassOIDList()
	auxs.Push(mySchema.ObjectClasses().Get(`dcObject`))
	dcrs := NewDITContentRules()
	dcrs.oIDsStringer()
	dcrs.Push(mySchema.DITContentRules().Index(0))
	dcrs.oIDsStringer()
	dcrs.Push(DITContentRule{&dITContentRule{
		OID: mySchema.ObjectClasses().Get(`2.5.6.2`),
		Aux: auxs,
	}})
	dcrs.oIDsStringer()
	dcrs.cast().NoPadding(true)
	dcrs.oIDsStringer()

	def.Map()

	var def2 DITContentRule
	_ = def2.Replace(def) // will fail

}

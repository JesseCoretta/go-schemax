package schemax

import (
	"fmt"
	"testing"
)

/*
This example demonstrates the creation of a new instance of [NameForm].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleSchema_NewNameForm() {
	var def NameForm = mySchema.NewNameForm()

	// set values in fluent form
	def.SetSchema(mySchema).
		SetNumericOID(`1.3.6.1.4.1.56521.999.97.7`).
		SetName(`deviceNameForm`).
		SetOC(mySchema.ObjectClasses().Get(`device`)).
		SetMust(`cn`).
		SetStringer()

	fmt.Printf("%s", def)
	// Output: ( 1.3.6.1.4.1.56521.999.97.7
	//     NAME 'deviceNameForm'
	//     OC device
	//     MUST cn )
}

/*
This example demonstrates a compliancy check of the "account" [ObjectClass].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleNameForm_Compliant() {
	dnaf := mySchema.NameForms().Get(`dotNotationArcForm`)
	fmt.Println(dnaf.Compliant())
	// Output: true
}

/*
This example demonstrates a compliancy check of all [NameForms] members.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleNameForms_Compliant() {
	forms := mySchema.NameForms()
	fmt.Println(forms.Compliant())
	// Output: true
}

func ExampleNameForm_IsIdentifiedAs() {
	nf := mySchema.NameForms().Get(`nRootArcForm`)
	fmt.Println(nf.IsIdentifiedAs(`1.3.6.1.4.1.56521.101.2.7.1`))
	// Output: true
}

/*
This example demonstrates a common (and most unfortunate) modification to an
OFFICIAL [ObjectClass] definition -- "groupOfNames", found within Section 3.5
of RFC 4519.

The design of this particular class is widely considered to be inconvenient
due to its mandate that at least one (1) instance of the "member" [AttributeType]
(from Section 2.17 of RFC 4519).

As such, this has forced many LDAP architects to literally modify this [ObjectClass]
definition within the given directory schema, moving the "member" [AttributeType]
from the MUST clause to the MAY clause.

For reasons of oversight, we've added the RFC of origin as an X-ORIGIN extension, and
a custom extension X-WARNING to remind users and admin alike that we've resorted to
this risky trick.
*/
func ExampleNameForm_Replace() {
	// Obtain the groupOfNames (gon) ObjectClass so
	// we can copy some of its values.
	gon := mySchema.NameForms().Get(`nRootArcForm`)

	// Craft a near identical groupOfNames instance,
	// save for the one change we intend to make.
	ngon := mySchema.NewNameForm().
		SetName(gon.Name()).
		SetDescription("new root arc name form for a number form RDN").
		SetNumericOID(gon.NumericOID()).
		SetOC(`rootArc`).
		SetMust(`n`).
		SetExtension(`X-ORIGIN`, `draft-coretta-oiddir-schema`).
		SetExtension(`X-WARNING`, `MODIFIED`).
		SetStringer()

	// Replace gon with ngon, while preserving its pointer
	// address so that references within stacks do not fail.
	gon.Replace(ngon)

	// call the new one (just to be sure)
	fmt.Println(mySchema.NameForms().Get(`nRootArcForm`))
	// Output: ( 1.3.6.1.4.1.56521.101.2.7.1
	//     NAME 'nRootArcForm'
	//     DESC 'new root arc name form for a number form RDN'
	//     OC rootArc
	//     MUST n
	//     X-ORIGIN 'draft-coretta-oiddir-schema'
	//     X-WARNING 'MODIFIED' )
}

/*
This example demonstrates use of the [NameForm.SetStringer] method
to impose a custom [Stringer] closure over the default instance.

Naturally the end-user would opt for a more useful stringer, such as one
that produces singular CSV rows per instance.

To avoid impacting other unit tests, we reset the default stringer
via the [NameForm.SetStringer] method again, with no arguments.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleNameForm_SetStringer() {
	dnaf := mySchema.NameForms().Get(`dotNotationArcForm`)
	dnaf.SetStringer(func() string {
		return "This useless message brought to you by a dumb stringer"
	})

	msg := fmt.Sprintln(dnaf)
	dnaf.SetStringer() // return it to its previous state if need be ...

	fmt.Printf("Original: %s\nOld: %s", dnaf, msg)
	// Output: Original: ( 1.3.6.1.4.1.56521.101.2.7.3
	//     NAME 'dotNotationArcForm'
	//     DESC 'arc name form for a numeric OID RDN'
	//     OC arc
	//     MUST dotNotation
	//     X-ORIGIN 'draft-coretta-oiddir-schema' )
	// Old: This useless message brought to you by a dumb stringer
}

/*
This example demonstrates use of the [NameForms.SetStringer] method
to impose a custom [Stringer] closure upon all stack members.

Naturally the end-user would opt for a more useful stringer, such as one
that produces a CSV file containing all [NameForm] instances.

To avoid impacting other unit tests, we reset the default stringer
via the [NameForms.SetStringer] method again, with no arguments.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleNameForms_SetStringer() {
	nfs := mySchema.NameForms()
	nfs.SetStringer(func() string {
		return "" // make a null stringer
	})

	output := nfs.String()
	nfs.SetStringer() // return to default

	fmt.Println(output)
	// Output:
}

/*
This example demonstrates a means of parsing a raw definition into a new
instance of [NameForm].

Note: this example assumes a legitimate schema variable is defined in place
of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleNameForm_Parse() {
	nf := mySchema.NewNameForm()

	// feed the parser a subtly bogus definition ...
	err := nf.Parse(`( 1.3.6.1.4.1.56521.999.14.56.1
                NAME 'fakeForm'
                DESC 'It\'s not real'
		OC device
		MUST cn
                X-ORIGIN 'YOUR FACE'
        )`)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(nf.Must())
	// Output: cn
}

/*
This example demonstrates the act of pushing, or appending, a new instance
of [NameForm] into a new [NameForms] stack instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleNameForms_Push() {
	nf := mySchema.NameForms().Get(`nArcForm`)
	myNFs := NewNameForms()
	myNFs.Push(nf)
	fmt.Println(myNFs.Len())
	// Output: 1
}

/*
This example demonstrates the assignment of arbitrary data to an instance
of [NameForm].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleNameForm_SetData() {
	// The value can be any type, but we'll
	// use a string for simplicity.
	documentation := `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`

	// Obtain the target attribute type to bear
	// the assigned value.
	nform := mySchema.NameForms().Get(`nArcForm`)

	// Set it.
	nform.SetData(documentation)

	// Get it and compare to the original.
	equal := documentation == nform.Data().(string)

	fmt.Printf("Values are equal: %t", equal)
	// Output: Values are equal: true
}

func ExampleNameForm_SetObsolete() {
	fake := NewNameForm().
		SetNumericOID(`1.3.6.1.4.1.56521.999.108.4`).
		SetName(`obsoleteForm`).
		SetObsolete()

	fmt.Println(fake.Obsolete())
	// Output: true
}

/*
This example demonstrates a means of checking whether a particular instance
of [NameForm] is present within an instance of [NameForms].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleNameForms_Contains() {
	classes := mySchema.NameForms()
	fmt.Println(classes.Contains(`nArcForm`)) // or "1.3.6.1.4.1.56521.101.2.7.2"
	// Output: true
}

func ExampleNameForms_Inventory() {
	oc := mySchema.NameForms().Inventory()
	fmt.Println(oc[`1.3.6.1.4.1.56521.101.2.7.2`][0])
	// Output: nArcForm
}

func ExampleNameForms_Type() {
	oc := mySchema.NameForms()
	fmt.Println(oc.Type())
	// Output: nameForms
}

func ExampleNameForm_Type() {
	var def NameForm
	fmt.Println(def.Type())
	// Output: nameForm
}

func ExampleNameForm_Map() {
	def := mySchema.NameForms().Get(`nRootArcForm`)
	fmt.Println(def.Map()[`NUMERICOID`][0]) // risky, just for simplicity
	// Output: 1.3.6.1.4.1.56521.101.2.7.1
}

/*
This example demonstrates use of the [NameForms.Maps] method, which
produces slices of [DefinitionMap] instances born of the [NameForms]
stack in which they reside.
*/
func ExampleNameForms_Maps() {
	defs := mySchema.NameForms().Maps()
	fmt.Println(defs[0][`NUMERICOID`][0]) // risky, just for simplicity
	// Output: 1.3.6.1.4.1.56521.101.2.7.1
}

/*
Do stupid things to make schemax panic, gain additional
coverage in the process.
*/
func TestNameForm_codecov(t *testing.T) {
	_ = mySchema.NameForms().SetStringer().Contains(``)
	mySchema.NameForms().Push(rune(10))
	mySchema.NameForms().IsZero()
	_ = mySchema.NameForms().String()
	cim := mySchema.NameForms().Get(`account`)
	mySchema.NameForms().canPush()
	mySchema.NameForms().canPush(``, ``, ``, ``, cim)
	mySchema.NameForms().canPush(cim, cim)
	bmr := newCollection(``)
	bma := newCollection(``)
	NameForms(bmr.cast()).Push(NewNameForm().SetSchema(mySchema))
	NameForms(bmr.cast()).Push(NewNameForm().SetSchema(mySchema).SetNumericOID(`1.2.3.4.5`))
	bmr.cast().Push(NewNameForm().SetSchema(mySchema))
	bmr.cast().Push(NewNameForm().SetSchema(mySchema).SetNumericOID(`1.2.3.4.5`))
	var bad NameForm
	bmr.cast().Push(bad)

	NameForms(bmr.cast()).canPush()
	NameForms(bmr.cast()).canPush(`things`)
	var ocs NameForms
	ocs.canPush(`forks`)
	ocs.Push(NewNameForm().SetSchema(mySchema))
	bmr.cast().Push(AttributeType{&attributeType{OID: `1.2.3.4.5`, Collective: true, Single: true}})
	bma.cast().Push(AttributeType{&attributeType{OID: ``, Collective: true, Single: true}})
	xoc := NameForm{&nameForm{
		Must: AttributeTypes(bmr),
	}}
	yoc := NameForm{&nameForm{
		May: AttributeTypes(bma),
	}}

	xoc.Compliant()
	yoc.Compliant()

	ocs.Push(bad)

	NameForms(bmr).Push(NewNameForm().SetSchema(mySchema))
	NameForms(bmr).Push(NewNameForm().SetSchema(mySchema).SetNumericOID(`1.2.3.4.5`))
	NameForms(bmr).Compliant()
	mySchema.NameForms().Compliant()

	var def NameForm

	_ = def.String()
	_ = def.SetStringer()
	_ = def.Description()
	_ = def.Name()
	_ = def.Names()
	_ = def.Extensions()
	_ = def.Must()
	_ = def.May()
	_ = def.Schema()
	_ = def.Map()
	_ = def.Compliant()
	_ = def.macro()
	_ = def.Obsolete()

	def.setOID(`4.3.2.1`)
	var raw string = `( 2.999.6.11 NAME 'fakeForm' OC device MUST cn MAY ( seeAlso $ ou $ l $ description ) X-ORIGIN 'NOWHERE' )`
	if err := def.Parse(raw); err != ErrNilReceiver {
		t.Errorf("%s failed: expected ErrNilReceiver, got %v", t.Name(), err)
		return
	}

	def = NewNameForm()
	def.SetDescription(`'a`)
	def.SetDescription(`'Unnecessary quoted value to be overwritten'`)

	oo := new(nameForm)
	oo.OID = `freakz`
	def.replace(NameForm{oo})

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

	def.SetMust(mySchema.AttributeTypes().Get(`cn`))
	def.SetMust(rune(11))
	def.SetMay(mySchema.AttributeTypes().Get(`cn`))
	def.SetMay(rune(11))
	mySchema.NameForms().canPush(NameForm{}, NameForm{new(nameForm)})

	if err := def.Parse(raw); err != nil {
		t.Errorf("%s failed: expected success, got %v", t.Name(), err)
		return
	}
	_ = def.macro()
	def.setOID(`2.5.13.2`)

	def.SetData(`fake`)
	def.SetData(nil)
	def.Data()

	var def2 NameForm
	_ = def2.Replace(def) // will fail

}

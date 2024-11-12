package schemax

import (
	"fmt"
	"testing"
)

/*
This example demonstrates the means for marshaling an instance of
[DITStructureRule] from a map[string]any instance.
*/
func ExampleDITStructureRule_Marshal() {
	m := map[string]any{
		`NAME`:     `exampleRule`,
		`DESC`:     `This is an example`,
		`RULEID`:   1,
		`OBSOLETE`: `FALSE`,
		`FORM`:     `domainNameForm`,
		`X-ORIGIN`: `RFCXXXX`,
	}

	var def DITStructureRule = mySchema.NewDITStructureRule()
	if err := def.Marshal(m); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s\n", def)
	// Output: ( 1
	//     NAME 'exampleRule'
	//     DESC 'This is an example'
	//     FORM domainNameForm
	//     X-ORIGIN 'RFCXXXX' )
}

/*
This example demonstrates an analysis of a distinguished name to determine
whether it honors the receiver instance of [DITStructureRule].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_Govern() {
	dn := `dc=example,dc=com` // flattened context (1 comma)

	// create a new DIT structure rule to leverage
	// RFC 2377's 'domainNameForm' definition
	dcdsr := mySchema.NewDITStructureRule().
		SetRuleID(13).
		SetName(`domainStructureRule`).
		SetForm(mySchema.NameForms().Get(`domainNameForm`)).
		SetStringer()

	err := dcdsr.Govern(dn, 1) // flat int
	fmt.Println(err)
	// Output: <nil>
}

/*
This example demonstrates the creation of a [DITStructureRule].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleNewDITStructureRule() {
	// First create a name form that requires an
	// RDN of uid=<val>, or (optionally) an RDN
	// of uid=<val>+gidNumber=<val>
	perForm := mySchema.NewNameForm().
		SetNumericOID(`1.3.6.1.4.1.56521.999.16.7`).
		SetName(`fictionalPersonForm`).
		SetDescription(`generalized person name form`).
		SetOC(`person`).
		SetMust(`uid`).
		SetMay(`gidnumber`).
		SetStringer()

	// Create the structure rule and assign the
	// new nameform
	ds := mySchema.NewDITStructureRule().
		SetRuleID(0).
		SetName(`fictionalPersonStructure`).
		SetDescription(`person structure rule`).
		SetForm(perForm).
		SetStringer()

	fmt.Println(ds)
	// Output: ( 0
	//     NAME 'fictionalPersonStructure'
	//     DESC 'person structure rule'
	//     FORM fictionalPersonForm )
}

/*
This example demonstrates the means of accessing the STRUCTURAL [ObjectClass]
instance held by the [NameForm] instance assigned to the [DITStructureRule]
instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_NamedObjectClass() {
	ds := mySchema.DITStructureRules().Get(1) // Integer Identifier #1
	noc := ds.NamedObjectClass()
	fmt.Println(noc.OID())
	// Output: uddiBusinessEntity
}

/*
This example demonstrates the means for checking to see if the receiver
is in an error condition.
*/
func ExampleDITStructureRule_E() {
	def := mySchema.NewDITStructureRule()
	def.SetRuleID("Z") // bogus
	if err := def.E(); err != nil {
		fmt.Println(err)
	}
	// Output: Invalid integer identifier for structure rule
}

/*
This example demonstrates the means for resolving an error condition.
*/
func ExampleDITStructureRule_E_clearError() {
	def := mySchema.NewDITStructureRule()
	def.SetRuleID(`X`) // bogus

	// We realized our mistake.
	def.SetRuleID(30) // valid
	def.SetForm(`domainNameForm`)

	// But when we check again, the error is still there.
	if err := def.E(); err != nil {
		// handle error
	}

	// We must clear the error with a
	// passing compliance check.
	if def.Compliant(); def.E() == nil {
		fmt.Println("Error has been resolved")
	}
	// Output: Error has been resolved

	return
}

/*
This example demonstrates a compliancy check of a [DITStructureRule]
instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_Compliant() {
	// grab our dITStructureRule bearing an ID of zero (0)
	ds := mySchema.DITStructureRules().Get(20) // or "rootArcStructure"
	fmt.Println(ds.Compliant())
	// Output: true
}

/*
This example demonstrates a compliancy check of a [DITStructureRules]
instance.

Generally use of this method is unnecessary due to stringent checks
imposed upon submitted [DITStructureRule] instance during the "Push"
process.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRules_Compliant() {
	// grab our dITStructureRules collection
	defs := mySchema.DITStructureRules()
	fmt.Println(defs.Compliant())
	// Output: true
}

/*
This example demonstrates the means for accessing all [DITStructureRules]
instances which bear the specified `X-ORIGIN` extension value.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRules_XOrigin() {
	defs := mySchema.DITStructureRules()
	matches := defs.XOrigin(`RFC4403`)
	fmt.Printf("Matched %d %s\n", matches.Len(), defs.Type())
	// Output: Matched 10 dITStructureRules
}

/*
This example demonstrates accessing the string type name of a [DITStructureRule]
definition.  This is mainly used as a low-cost alternative to type assertion
when dealing with [Definition] interface type instances.
*/
func ExampleDITStructureRule_Type() {
	var def DITStructureRule
	fmt.Println(def.Type())
	// Output: dITStructureRule
}

/*
This example demonstrates accessing the string type name of a [DITStructureRules]
definition.  This is mainly used as a low-cost alternative to type assertion
when dealing with [Definitions] interface type instances.
*/
func ExampleDITStructureRules_Type() {
	var defs DITStructureRules
	fmt.Println(defs.Type())
	// Output: dITStructureRules
}

/*
This example demonstrates the means of checking whether an instance of
[DITStructureRule] is of a nil state.
*/
func ExampleDITStructureRule_IsZero() {
	var def DITStructureRule
	fmt.Println(def.IsZero())
	// Output: true
}

/*
This example demonstrates the means of checking whether an instance of
[DITStructureRules] is of a nil state.
*/
func ExampleDITStructureRules_IsZero() {
	var defs DITStructureRules
	fmt.Println(defs.IsZero())
	// Output: true
}

/*
This example demonstrates the means of accessing the string form of the
principal name OR rule ID of a [DITStructureRule] instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_ID() {
	def := mySchema.DITStructureRules().Get(20)
	fmt.Println(def.ID())
	// Output: rootArcStructure
}

/*
This example demonstrates the means of accessing the unsigned rule ID
held by an instance of [DITStructureRule].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_RuleID() {
	def := mySchema.DITStructureRules().Get(`arcStructure`) // or 11, or "11"
	fmt.Println(def.RuleID())
	// Output: 0
}

/*
This example demonstrates the means of accessing the principal name of
an instance of [DITStructureRule].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_Name() {
	def := mySchema.DITStructureRules().Get(20)
	fmt.Println(def.Name())
	// Output: rootArcStructure
}

/*
This example demonstrates the means of accessing the [QuotedDescriptorList]
instance containing the name(s) by which a [DITStructureRule] is known.

Use of this method will encapsulate the value(s), per § 4.1 of RFC 4512,
using single quotes (SQUOTE ('), ASCII %x27).

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_Names() {
	def := mySchema.DITStructureRules().Get(1) // or "1"
	fmt.Println(def.Names())
	// Output: 'uddiBusinessEntityStructureRule'
}

/*
This example demonstrates the futility of the [DITStructureRule.NumericOID]
method. Numeric OIDs are not used to identify instances of [DITStructureRule].

The [DITStructureRule.NumericOID] method only exists to satisfy Go's interface
signature requirements with regards to the [Definition] interface type.
*/
func ExampleDITStructureRule_NumericOID() {
	// access a known valid dITStructureRule
	def := mySchema.DITStructureRules().Get(2) // or 'dotNotArcStructure'
	fmt.Println(def.NumericOID())
	// Output:
}

/*
This example demonstrates the means of accessing any and all immediate
superior [DITStructureRule] instances for the receiver instance. The
super chain of rules is NOT traversed indefinitely.

The [DITStructureRule.NumericOID] method only exists to satisfy Go's interface
signature requirements with regards to the [Definition] interface type.
*/
func ExampleDITStructureRule_SuperRules() {
	def := mySchema.DITStructureRules().Get(2)
	fmt.Println(def.SuperRules())
	// Output: 1
}

/*
This example demonstrates the means of accessing all subordinate rule
instances of the receiver instance.

In essence, this method is the opposite of the [DITStructureRule.SuperRules]
method and may return zero (0) or more [DITStructureRule] instances within
the return [DITStructureRules] instance.
*/
func ExampleDITStructureRule_SubRules() {
	def := mySchema.DITStructureRules().Get(5)
	fmt.Printf("%d subordinate rule found", def.SubRules().Len())
	// Output: 1 subordinate rule found
}

/*
This example demonstrates the means of calling the Nth [DITStructureRule]
slice instance from the [DITStructureRules] collection instance in which
it resides.

This method should not be confused with [DITStructureRules.Get], which
deals with unsigned rule IDs and names of definitions -- not indices.
*/
func ExampleDITStructureRules_Index() {
	defs := mySchema.DITStructureRules()
	fmt.Println(defs.Index(0).Name())
	// Output: uddiBusinessEntityStructureRule
}

/*
This example demonstrates the means of transforming an instance of
[DITStructureRule] into an instance of [DefinitionMap] for simplified
use.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_Map() {
	def := mySchema.DITStructureRules().Get(20)
	fmt.Println(def.Map()[`NAME`][0])
	// Output: rootArcStructure
}

/*
This example demonstrates the assignment of arbitrary data to an instance
of [AttributeType].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_SetData() {
	nf := mySchema.NewNameForm().
		SetNumericOID(`1.3.6.1.4.1.56521.999.16.7`).
		SetName(`personForm`).
		SetDescription(`generalized person name form`).
		SetOC(`person`).
		SetMust(`uid`).
		SetMay(`gidnumber`).
		SetStringer()

	//mySchema.NameForms().Push(nf)

	// Create the structure rule and assign the
	// new nameform
	ds := mySchema.NewDITStructureRule().
		SetRuleID(0).
		SetName(`personStructure`).
		SetDescription(`person structure rule`).
		SetForm(nf).
		SetExtension(`X-ORIGIN`, `NOWHERE`).
		SetStringer()

	// The value can be any type, but we'll
	// use a string for simplicity.
	documentation := `Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`

	// Set it.
	ds.SetData(documentation)

	// Get it and compare to the original.
	equal := documentation == ds.Data().(string)

	fmt.Printf("Values are equal: %t", equal)
	// Output: Values are equal: true
}

func ExampleDITStructureRule_IsIdentifiedAs() {
	def := mySchema.DITStructureRules().Get(20)
	fmt.Println(def.IsIdentifiedAs(`rootArcStructure`))
	// Output: true
}

/*
This example demonstrates use of the [DITStructureRule.SetStringer] method
to impose a custom [Stringer] closure over the default instance.

Naturally the end-user would opt for a more useful stringer, such as one
that produces singular CSV rows per instance.

To avoid impacting other unit tests, we reset the default stringer
via the [DITStructureRule.SetStringer] method again, with no arguments.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_SetStringer() {
	ds := mySchema.DITStructureRules().Get(2)
	ds.SetStringer(func() string {
		return "This useless message brought to you by a dumb stringer"
	})

	msg := fmt.Sprintln(ds)
	ds.SetStringer() // return it to its previous state

	fmt.Printf("Original: %s\nOld: %s", ds, msg)
	// Original: ( 2
	//     NAME 'dotNotArcStructure'
	//     DESC 'structure rule for two dimensional arc entries; FOR DEMONSTRATION USE ONLY'
	//     FORM dotNotationArcForm
	//     SUP 0
	//     X-ORIGIN 'draft-coretta-oiddir-schema; unofficial supplement' )
	// Old: This useless message brought to you by a dumb stringer
}

/*
This example demonstrates use of the [DITStructureRules.SetStringer] method
to impose a custom [Stringer] closure upon all stack members.

Naturally the end-user would opt for a more useful stringer, such as one
that produces a CSV file containing all [AttributeType] instances.

To avoid impacting other unit tests, we reset the default stringer
via the [DITStructureRules.SetStringer] method again, with no arguments.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRules_SetStringer() {
	defs := mySchema.DITStructureRules()
	defs.SetStringer(func() string {
		return "" // make a null stringer
	})

	output := defs.String()
	defs.SetStringer() // return to default

	fmt.Println(output)
	// Output:
}

/*
This example demonstrates use of the [DITStructureRule.SetObsolete]
method to impose obsolescence upon the receiver instance.

Note: we craft a throw-away instance so as to avoid impacting other
unit tests as a result of declaring obsolescence upon an otherwise
valid instance.
*/
func ExampleDITStructureRule_SetObsolete() {
	defs := NewDITStructureRule()
	defs.SetName(`throwAwayStructure`)
	defs.SetRuleID(10)
	defs.SetObsolete()
	fmt.Printf("%s is obsolete: %t", defs.Name(), defs.Obsolete())
	// Output: throwAwayStructure is obsolete: true
}

/*
This example demonstrates the means of manually assigning a superior
[DITStructureRule] to the receiver instance, thereby rendering it
subordinate it context.

Note: we craft a throw-away instance just for the sake of simplicity.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_SetSuperRule() {
	super := mySchema.DITStructureRules().Get(2)
	defs := mySchema.NewDITStructureRule()
	defs.SetName(`throwAwaySubStructure`)
	defs.SetRuleID(10)
	defs.SetSuperRule(super)
	fmt.Printf("Superior rule is %s", defs.SuperRules().Index(0).ID())
	// Output: Superior rule is uddiContactStructureRule
}

/*
This example demonstrates a means of parsing a raw definition into a new
instance of [AttributeType].

Note: this example assumes a legitimate schema variable is defined in place
of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_Parse() {
	def := mySchema.NewDITStructureRule()

	err := def.Parse(`( 31
                NAME 'fakeStructureRule'
		DESC 'fake structure rule'
		FORM dotNotationArcForm
                X-ORIGIN 'NOWHERE' )`)

	fmt.Println(err)
	// Output: <nil>
}

/*
This example demonstrates the replacement process of a [DITStructureRule]
instance within an instance of [DITStructureRules].

For reasons of oversight, we've added a custom extension X-WARNING to
remind users and admin alike of the modification.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_Replace() {
	def := mySchema.NewDITStructureRule()

	err := def.Parse(`( 31
                NAME 'fakeStructureRule'
                DESC 'fake structure rule'
                FORM dotNotationArcForm
                X-ORIGIN 'NOWHERE' )`)

	if err != nil {
		fmt.Println(err)
		return
	}

	def2 := mySchema.NewDITStructureRule()
	err = def2.Parse(`( 31
                NAME 'fakeStructureRule'
                DESC 'fake structure rule updated'
                FORM dotNotationArcForm
                X-ORIGIN 'ANYWHERE' )`)

	if err != nil {
		fmt.Println(err)
		return
	}

	def.Replace(def2)
	fmt.Println(def.Description())
	// Output: fake structure rule updated
}

/*
This example demonstrates the creation of an [Inventory] instance based
upon the current contents of a [DITStructureRules] instance.

Note: this example assumes a legitimate schema variable is defined in place
of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRules_Inventory() {
	defs := mySchema.DITStructureRules()
	inv := defs.Inventory()
	fmt.Println(inv[`2`][0])
	// Output: uddiContactStructureRule
}

/*
This example demonstrates a means of checking whether a particular instance
of [DITStructureRule] is present within an instance of [DITStructureRules].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRules_Contains() {
	fmt.Println(mySchema.DITStructureRules().Contains(`rootArcStructure`))
	// Output: true
}

/*
This example demonstrates use of the [DITStructureRules.Maps] method, which
produces slices of [DefinitionMap] instances containing [DITStructureRule]
derived values.

Here, we (quite recklessly) call index three (3) and reference index zero
(0) of its `NAME` key to obtain the principal name of the definition.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRules_Maps() {
	defs := mySchema.DITStructureRules().Maps()
	fmt.Println(defs[2][`NAME`][0]) // risky, just for simplicity
	// Output: uddiAddressStructureRule
}

/*
This example demonstrates the parsing of a bogus definition, which results
in the return of an error by the underlying ANTLR parser.

Note: this example assumes a legitimate schema variable is defined in place
of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleDITStructureRule_Parse_bogus() {
	def := mySchema.NewDITStructureRule()

	// feed the parser a subtly bogus definition ...
	err := def.Parse(`( -1
                NAME 'fakeStructureRule'
		DESC 'fake structure rule'
		FORM dotNotationArcForm
                X-ORIGIN 'YOUR FACE' )`)

	fmt.Println(err)
	// Output: Inconsistent antlr4512.DITStructureRule parse results or bad content
}

func TestDITStructureRule_Govern(t *testing.T) {
	// create a new DIT structure rule to leverage
	// RFC 2377's 'domainNameForm' definition
	dcdsr := mySchema.NewDITStructureRule().
		SetRuleID(13).
		SetName(`domainStructureRule`).
		SetForm(mySchema.NameForms().Get(`domainNameForm`)).
		SetStringer()
	mySchema.DITStructureRules().Push(dcdsr)

	ounf := mySchema.NewNameForm().
		SetNumericOID(`1.3.6.1.4.1.56521.999.55.11.33`).
		SetName(`ouNameForm`).
		SetOC(`organizationalUnit`).
		SetMust(`ou`).
		SetStringer()
	mySchema.NameForms().Push(ounf)

	oudsr := mySchema.NewDITStructureRule().
		SetRuleID(14).
		SetName(`ouStructureRule`).
		SetForm(mySchema.NameForms().Get(`ouNameForm`)).
		SetSuperRule(13, `self`).
		SetStringer()

	mySchema.DITStructureRules().Push(oudsr)

	for idx, strukt := range []struct {
		DN string
		L  int
		ID int
	}{
		{`dc=example,dc=com`, 1, 13},
		{`o=example`, 1, 13},
		{`ou=People,dc=example,dc=com`, 1, 14},
		{`ou=People,dc=example,dc=com`, -1, 13},
		{`ou=Employees,ou=People,dc=example,dc=com`, 1, 14},
		{`x=People,dc=example,dc=com`, 1, 14},
		{`ou=People+ou=Employees,dc=example,dc=com`, 1, 14},
		{`ou=People+cn=Employees,dc=example,dc=com`, 1, 14},
		{`ou=People+ou=Employees,dc=example,dc=com`, 1, 14},
		{`ou=Employees,ou=People,dc=example,dc=com`, 1, 13},
	} {
		rule := mySchema.DITStructureRules().Get(strukt.ID)
		even := idx%2 == 0

		if err := rule.Govern(strukt.DN, strukt.L); err != nil {
			if even {
				t.Errorf("%s[%d] failed: %v (%v)", t.Name(), idx, strukt.DN, err)
				return
			}
		} else {
			if !even {
				t.Errorf("%s[%d] failed: expected error, got nothing", t.Name(), idx)
				return
			}
		}
	}
}

/*
Do stupid things to make schemax panic, gain additional
coverage in the process.
*/
func TestDITStructureRule_codecov(t *testing.T) {
	var rs DITStructureRules
	rs.iDsStringer()
	rs.canPush(rune(33))
	rs = NewDITStructureRules()
	rs.canPush(rs)

	var r DITStructureRule
	r.Map()
	r.NumericOID()
	r.Compliant()
	r.Replace(DITStructureRule{})
	r.ID()
	r.Name()
	r.Names()
	r.Form()
	r.RuleID()
	r.setOID(``)
	r.macro()
	r.Govern(`bogusdn`)
	r.Parse(`crap`)
	r.IsIdentifiedAs(`nothing`)
	r.Replace(DITStructureRule{&dITStructureRule{}})
	r.replace(DITStructureRule{&dITStructureRule{ID: uint(55)}})

	_r := new(dITStructureRule)
	_r.setDescription(`'s`)
	_r.setRuleID(uint64(4))
	_r.setRuleID(uint(4))
	_r.setRuleID(4)
	_r.setRuleID(`4`)

	z := DITStructureRule{_r}
	z.Parse(`crap`)
	z.Map()

	_r.setDescription(`'texts'`)
	_r.setDescription(`s'`)
	_r.schema = mySchema
	_r.setForm(`nArcForm`)
	_r.setSuperRule(float64(333))
	_r.setSuperRule(`rootArcStructure`)
	_r.setSuperRule(mySchema.DITStructureRules().Get(`rootArcStructure`))
	r = DITStructureRule{_r}
	r.ID()
	r.Replace(DITStructureRule{&dITStructureRule{ID: uint(55), schema: mySchema}})

	bmr := newCollection(``)
	DITStructureRules(bmr).iDsStringer()
	DITStructureRules(bmr).canPush()
	bmr.cast().Push(DITStructureRule{&dITStructureRule{ID: 4, Form: NameForm{&nameForm{OID: `1.2.3`}}}})
	DITStructureRules(bmr).iDsStringer()
	bmr.cast().Push(DITStructureRule{&dITStructureRule{ID: 5, Form: NameForm{&nameForm{OID: `1.2.3`}}}})
	DITStructureRules(bmr).Compliant()
	mySchema.DITStructureRules().Get(-1)
	rs.Push(mySchema.DITStructureRules().Get(1))
	rs.Push(mySchema.DITStructureRules().Get(1))
	rs.Push(mySchema.AttributeTypes().Get(`cn`))
	DITStructureRules(bmr).Push(DITStructureRule{&dITStructureRule{ID: 5, Form: NameForm{&nameForm{OID: `1.2.3`}}}})
	bmr.cast().NoPadding(true)
	DITStructureRules(bmr).iDsStringer()
}

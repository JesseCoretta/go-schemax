package schemax

import (
	"fmt"
	"testing"
)

/*
This example demonstrates manual assembly of a new [MatchingRuleUse]
instance. Note this is provided for demonstration purposes only and
in context does not perform anything useful.

In general it is not necessary for end-users to manually define
this kind of instance.  Instances of this type are normally created
by automated processes when new [AttributeType] definitions are created
or introduced which make use of a given [MatchingRule] instance.
*/
func ExampleNewMatchingRuleUse() {
	var def MatchingRuleUse = NewMatchingRuleUse().SetSchema(mySchema)

	def.SetNumericOID(`2.5.13.16`).
		SetName(`fakeBitStringMatch`).
		SetExtension(`X-ORIGIN`, `NOWHERE`)

	for _, apl := range []AttributeType{
		mySchema.AttributeTypes().Get(`cn`),
		mySchema.AttributeTypes().Get(`sn`),
		mySchema.AttributeTypes().Get(`l`),
	} {
		def.SetApplies(apl)
	}

	// We're done and ready, set the stringer
	def.SetStringer()

	fmt.Printf("%s", def)
	// Output: ( 2.5.13.16
	//     NAME 'fakeBitStringMatch'
	//     APPLIES ( cn
	//             $ sn
	//             $ l )
	//     X-ORIGIN 'NOWHERE' )
}

/*
This example demonstrates the futility of attempting to parse a raw
string-based matchingRuleUse definition into a proper instance of
[MatchingRuleUse], as these definitions are auto-generated by the
DSA governed by the relevant schema; not parsed from input.
*/
func ExampleMatchingRuleUse_Parse() {

	var raw string = `( 2.5.13.21
		NAME 'tellyMatch'
		APPLIES sponsorTelephoneNumber
		X-ORIGIN 'NOWHERE' )`

	var def MatchingRuleUse
	fmt.Println(def.Parse(raw))
	// Output: Parsing is not applicable to a MatchingRuleUse
}

/*
This example demonstrates a means of initializing a new instance of
[MatchingRuleUse], with the receiver instance of [Schema] automatically
registered as its origin.

Generally this is not needed, and is shown here merely for coverage
purposes. [MatchingRuleUse] instances are expected to be generated
automatically by the DSA(s) governed by the relevant [Schema], and
never parsed through user input.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleSchema_NewMatchingRuleUse() {
	var def MatchingRuleUse = mySchema.NewMatchingRuleUse()
	fmt.Println(def.Type())
	// Output: matchingRuleUse
}

/*
This example demonstrates accessing the principal name of the receiver
instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUse_Name() {
	im := mySchema.MatchingRuleUses().Get(`2.5.13.14`)
	fmt.Println(im.Name())
	// Output: integerMatch
}

/*
This example demonstrates accessing the numeric OID of the receiver
instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUse_NumericOID() {
	im := mySchema.MatchingRuleUses().Get(`integerMatch`)
	fmt.Println(im.NumericOID())
	// Output: 2.5.13.14
}

/*
This example demonstrates accessing the OID -- whether it is the principal
name or numeric OID -- of the receiver instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUse_OID() {
	im := mySchema.MatchingRuleUses().Get(`2.5.13.14`)
	fmt.Println(im.OID())
	// Output: integerMatch
}

//func ExampleMatchingRuleUse_Map() {
//        def := mySchema.MatchingRuleUses().Get(`2.5.13.14`)
//        fmt.Println(def.Map()[`SYNTAX`][0]) // risky, just for simplicity
// Output: 1.3.6.1.4.1.1466.115.121.1.15
//}

/*
This example demonstrates use of the [MatchingRuleUses.Maps] method, which
produces slices of [DefinitionMap] instances born of the [MatchingRuleUses]
stack in which they reside.  We (quite recklessly) call index three (3)
and reference index zero (0) of its `SYNTAX` key to obtain the relevant
[LDAPSyntax] OID string value.
*/
//func ExampleMatchingRuleUses_Maps() {
//        defs := mySchema.MatchingRuleUses().Maps()
//        fmt.Println(defs[3][`SYNTAX`][0]) // risky, just for simplicity
// Output: 1.3.6.1.4.1.1466.115.121.1.15
//}

/*
This example demonstrates the means for accessing all [MatchingRuleUse]
instances which bear the specified `X-ORIGIN` extension value.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUses_XOrigin() {
	defs := mySchema.MatchingRuleUses()
	matches := defs.XOrigin(`Bogus RFC`)
	fmt.Printf("Matched %d of %d %s\n", matches.Len(), defs.Len(), defs.Type())
	// Output: Matched 0 of 32 matchingRuleUses
}

/*
This example demonstrates use of the [MatchingRuleUses.Type] method to determine
the type of stack defined within the receiver. This is mainly useful in cases
where multiple stacks are being iterated in [Definitions] interface contexts
and is more efficient when compared to manual type assertion.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUses_Type() {
	mrs := mySchema.MatchingRuleUses()
	fmt.Printf("We have %d %s", mrs.Len(), mrs.Type())
	// Output: We have 32 matchingRuleUses
}

/*
This example demonstrates the means of accessing the integer length of
a [MatchingRuleUses] stack instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUses_Len() {
	mrs := mySchema.MatchingRuleUses()
	fmt.Printf("We have %d %s", mrs.Len(), mrs.Type())
	// Output: We have 32 matchingRuleUses
}

/*
This example demonstrates the means of accessing a specific slice value
within an instance of [MatchingRuleUses] by way of its associated integer
index.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUses_Index() {
	slice := mySchema.MatchingRuleUses().Index(3)
	fmt.Println(slice)
	// Output: ( 2.5.13.28
	//     NAME 'generalizedTimeOrderingMatch'
	//     APPLIES ( createTimestamp
	//             $ modifyTimestamp
	//             $ subschemaTimestamp
	//             $ registrationCreated
	//             $ registrationModified
	//             $ currentAuthorityStartTimestamp
	//             $ firstAuthorityStartTimestamp
	//             $ firstAuthorityEndTimestamp
	//             $ sponsorStartTimestamp
	//             $ sponsorEndTimestamp )
	//     X-ORIGIN 'RFC4517' )
}

func ExampleMatchingRuleUse_Data() {
	mr := mySchema.MatchingRuleUses().Get(`caseIgnoreMatch`)

	// Let's pretend img ([]uint8) represents
	// some JPEG data (e.g.: a diagram)
	var img []uint8 = []uint8{0x1, 0x2, 0x3, 0x4}
	mr.SetData(img)

	got := mr.Data().([]uint8)

	fmt.Printf("%T, Len:%d", got, len(got))
	// Output: []uint8, Len:4
}

func ExampleMatchingRuleUse_SetData() {
	mr := mySchema.MatchingRuleUses().Get(`caseIgnoreMatch`)

	// Let's pretend img ([]uint8) represents
	// some JPEG data (e.g.: a diagram)
	var img []uint8 = []uint8{0x1, 0x2, 0x3, 0x4}
	mr.SetData(img)

	got := mr.Data().([]uint8)

	fmt.Printf("%T, Len:%d", got, len(got))
	// Output: []uint8, Len:4
}

func ExampleMatchingRuleUse_Obsolete() {
	def := mySchema.MatchingRuleUses().Get(`caseExactMatch`)
	fmt.Println(def.Obsolete())
	// Output: false
}

/*
This example demonstrates instant compliance checks for all [MatchingRuleUse]
instances present within an instance of [MatchingRuleUses].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUses_Compliant() {
	mus := mySchema.MatchingRuleUses()
	fmt.Printf("All %d %s are compliant: %t", mus.Len(), mus.Type(), mus.Compliant())
	// Output: All 32 matchingRuleUses are compliant: true
}

/*
This example demonstrates use of the [MatchingRuleUse.SetStringer] method
to impose a custom [Stringer] closure over the default instance.

Naturally the end-user would opt for a more useful stringer, such as one
that produces singular CSV rows per instance.

To avoid impacting other unit tests, we reset the default stringer
via the [MatchingRuleUse.SetStringer] method again, with no arguments.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUse_SetStringer() {
	cim := mySchema.MatchingRuleUses().Get(`caseIgnoreMatch`)
	cim.SetStringer(func() string {
		return "This useless message brought to you by a dumb stringer"
	})

	msg := fmt.Sprintln(cim)
	cim.SetStringer() // return it to its previous state if need be ...

	fmt.Printf("Original: %s\nOld: %s", cim, msg)
	// Output: Original: ( 2.5.13.2
	//     NAME 'caseIgnoreMatch'
	//     APPLIES ( carLicense
	//             $ departmentNumber
	//             $ displayName
	//             $ employeeNumber
	//             $ employeeType
	//             $ preferredLanguage
	//             $ name
	//             $ businessCategory
	//             $ description
	//             $ destinationIndicator
	//             $ dnQualifier
	//             $ houseIdentifier
	//             $ physicalDeliveryOfficeName
	//             $ postalCode
	//             $ postOfficeBox
	//             $ serialNumber
	//             $ street
	//             $ uid
	//             $ buildingName
	//             $ co
	//             $ documentIdentifier
	//             $ documentLocation
	//             $ documentPublisher
	//             $ documentTitle
	//             $ documentVersion
	//             $ drink
	//             $ host
	//             $ info
	//             $ organizationalStatus
	//             $ personalTitle
	//             $ roomNumber
	//             $ uniqueIdentifier
	//             $ userClass
	//             $ uddiBusinessKey
	//             $ uddiOperator
	//             $ uddiName
	//             $ uddiDescription
	//             $ uddiDiscoveryURLs
	//             $ uddiUseType
	//             $ uddiPersonName
	//             $ uddiPhone
	//             $ uddiEMail
	//             $ uddiSortCode
	//             $ uddiTModelKey
	//             $ uddiAddressLine
	//             $ uddiIdentifierBag
	//             $ uddiCategoryBag
	//             $ uddiKeyedReference
	//             $ uddiServiceKey
	//             $ uddiBindingKey
	//             $ uddiAccessPoint
	//             $ uddiHostingRedirector
	//             $ uddiInstanceDescription
	//             $ uddiInstanceParms
	//             $ uddiOverviewDescription
	//             $ uddiOverviewURL
	//             $ uddiFromKey
	//             $ uddiToKey
	//             $ uddiUUID
	//             $ uddiLang
	//             $ uddiv3BusinessKey
	//             $ uddiv3ServiceKey
	//             $ uddiv3BindingKey
	//             $ uddiv3TModelKey
	//             $ uddiv3NodeId
	//             $ uddiv3SubscriptionKey
	//             $ uddiv3SubscriptionFilter
	//             $ uddiv3NotificationInterval
	//             $ uddiv3EntityKey
	//             $ aSN1Notation )
	//     X-ORIGIN 'RFC4517' )
	// Old: This useless message brought to you by a dumb stringer
}

/*
This example demonstrates use of the [MatchingRuleUses.SetStringer] method
to impose a custom [Stringer] closure upon all stack members.

Naturally the end-user would opt for a more useful stringer, such as one
that produces a CSV file containing all [MatchingRuleUse] instances.

To avoid impacting other unit tests, we reset the default stringer
via the [MatchingRuleUses.SetStringer] method again, with no arguments.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUses_SetStringer() {
	mrs := mySchema.MatchingRuleUses()
	mrs.SetStringer(func() string {
		return "" // make a null stringer
	})

	output := mrs.String()
	mrs.SetStringer() // return to default

	fmt.Println(output)
	// Output:
}

/*
This example demonstrates use of the [MatchingRuleUses.Maps] method, which
produces slices of [DefinitionMap] instances containing [MatchingRuleUse]
derived values

Here, we (quite recklessly) call index three (3) and reference index zero
(0) of its `SYNTAX` key to obtain the relevant [LDAPSyntax] OID string value.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUses_Maps() {
	defs := mySchema.MatchingRuleUses().Maps()
	fmt.Println(defs[3][`APPLIES`][0]) // risky, just for simplicity
	// Output: createTimestamp
}

/*
This example demonstrates the means of transferring a [MatchingRuleUse]
into an instance of [DefinitionMap].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUse_Map() {
	def := mySchema.MatchingRuleUses().Get(`caseIgnoreMatch`)
	fmt.Println(def.Map()[`NUMERICOID`][0]) // risky, just for simplicity
	// Output: 2.5.13.2
}

/*
This example demonstrates the creation of an [Inventory] instance based
upon the current contents of a [MatchingRuleUses] stack instance.  Use
of an [Inventory] instance is convenient in cases where a receiver of
schema information may not be able to directly receive working stack
instances and requires a more portable and generalized type.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUses_Inventory() {
	at := mySchema.MatchingRuleUses().Inventory()
	fmt.Println(at[`2.5.13.2`][0])
	// Output: caseIgnoreMatch
}

func ExampleMatchingRuleUses_IsZero() {
	var mus MatchingRuleUses
	fmt.Println(mus.IsZero())
	// Output: true
}

/*
This example demonstrates a means of accessing the underlying [Extensions]
stack instance within the receiver instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUse_Extensions() {
	cim := mySchema.MatchingRuleUses().Get(`caseIgnoreMatch`)
	fmt.Println(cim.Extensions())
	// Output: X-ORIGIN 'RFC4517'
}

/*
This example demonstrates accessing the description clause within the
receiver instance. Most [MatchingRuleUse] instances do not have any
descriptive text set, thus like others this example produces no value.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUse_Description() {
	cim := mySchema.MatchingRuleUses().Get(`caseIgnoreMatch`)
	fmt.Println(cim.Description())
	// Output:
}

func ExampleMatchingRuleUse_SetDescription() {
	cim := mySchema.MatchingRuleUses().Get(`caseIgnoreMatch`)
	cim.SetDescription("Caseless string match")
	fmt.Println(cim.Description())
	// Output: Caseless string match
}

func ExampleMatchingRuleUse_Names() {
	cim := mySchema.MatchingRuleUses().Get(`2.5.13.2`)
	fmt.Println(cim.Names())
	// Output: 'caseIgnoreMatch'
}

/*
This example demonstrates calling the 3rd index of the [MatchingRuleUses]
stack of our [Schema] and performing a compliancy checking on the slice.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUse_Compliant() {
	mu := mySchema.MatchingRuleUses().Index(3)
	fmt.Println(mu.Compliant())
	// Output: true
}

/*
This example demonstrates the string representation process for an instance
of [MatchingRuleUse].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUse_String() {
	mu := mySchema.MatchingRuleUses().Get(`2.5.13.2`)
	fmt.Println(mu)
	// Output: ( 2.5.13.2
	//     NAME 'caseIgnoreMatch'
	//     APPLIES ( carLicense
	//             $ departmentNumber
	//             $ displayName
	//             $ employeeNumber
	//             $ employeeType
	//             $ preferredLanguage
	//             $ name
	//             $ businessCategory
	//             $ description
	//             $ destinationIndicator
	//             $ dnQualifier
	//             $ houseIdentifier
	//             $ physicalDeliveryOfficeName
	//             $ postalCode
	//             $ postOfficeBox
	//             $ serialNumber
	//             $ street
	//             $ uid
	//             $ buildingName
	//             $ co
	//             $ documentIdentifier
	//             $ documentLocation
	//             $ documentPublisher
	//             $ documentTitle
	//             $ documentVersion
	//             $ drink
	//             $ host
	//             $ info
	//             $ organizationalStatus
	//             $ personalTitle
	//             $ roomNumber
	//             $ uniqueIdentifier
	//             $ userClass
	//             $ uddiBusinessKey
	//             $ uddiOperator
	//             $ uddiName
	//             $ uddiDescription
	//             $ uddiDiscoveryURLs
	//             $ uddiUseType
	//             $ uddiPersonName
	//             $ uddiPhone
	//             $ uddiEMail
	//             $ uddiSortCode
	//             $ uddiTModelKey
	//             $ uddiAddressLine
	//             $ uddiIdentifierBag
	//             $ uddiCategoryBag
	//             $ uddiKeyedReference
	//             $ uddiServiceKey
	//             $ uddiBindingKey
	//             $ uddiAccessPoint
	//             $ uddiHostingRedirector
	//             $ uddiInstanceDescription
	//             $ uddiInstanceParms
	//             $ uddiOverviewDescription
	//             $ uddiOverviewURL
	//             $ uddiFromKey
	//             $ uddiToKey
	//             $ uddiUUID
	//             $ uddiLang
	//             $ uddiv3BusinessKey
	//             $ uddiv3ServiceKey
	//             $ uddiv3BindingKey
	//             $ uddiv3TModelKey
	//             $ uddiv3NodeId
	//             $ uddiv3SubscriptionKey
	//             $ uddiv3SubscriptionFilter
	//             $ uddiv3NotificationInterval
	//             $ uddiv3EntityKey
	//             $ aSN1Notation )
	//     X-ORIGIN 'RFC4517' )
}

/*
This example demonstrates the act of pushing, or appending, a new instance
of [MatchingRuleUse] into a new [MatchingRuleUses] stack instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUses_Push() {
	mu := mySchema.MatchingRuleUses().Get(`caseIgnoreMatch`)
	myMus := NewMatchingRuleUses()
	myMus.Push(mu)
	fmt.Println(myMus.Len())
	// Output: 1
}

/*
This example demonstrates a means of checking whether a particular instance
of [MatchingRuleUse] is present within an instance of [MatchingRuleUses].

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUses_Contains() {
	mus := mySchema.MatchingRuleUses()
	fmt.Println(mus.Contains(`caseIgnoreMatch`)) // or "2.5.13.2"
	// Output: true
}

/*
This example demonstrates the creation of a new [MatchingRuleUse]
instance for manual assembly as an OBSOLETE instance.

Note: this example assumes a legitimate schema variable is defined
in place of the fictional "mySchema" var shown here for simplicity.
*/
func ExampleMatchingRuleUse_SetObsolete() {
	var def MatchingRuleUse = NewMatchingRuleUse()
	def.SetObsolete()
	fmt.Printf("Is obsolete: %t", def.Obsolete())
	// Output: Is obsolete: true
}

/*
Do stupid things to make schemax panic, gain additional
coverage in the process.
*/
func TestMatchingRuleUse_codecov(t *testing.T) {
	_ = mySchema.MatchingRuleUses().SetStringer().Contains(``)
	mySchema.MatchingRuleUses().Push(rune(10))
	mySchema.MatchingRuleUses().IsZero()
	_ = mySchema.MatchingRuleUses().String()
	cim := mySchema.MatchingRuleUses().Get(`caseIgnoreMatch`)
	mySchema.MatchingRuleUses().canPush()
	mySchema.MatchingRuleUses().canPush(``, ``, ``, ``, cim)
	mySchema.MatchingRuleUses().canPush(cim, cim)
	bmr := newCollection(``)
	MatchingRuleUses(bmr.cast()).Push(NewMatchingRuleUse().SetSchema(mySchema))
	MatchingRuleUses(bmr.cast()).Push(NewMatchingRuleUse().SetSchema(mySchema).SetNumericOID(`1.2.3.4.5`))
	bmr.cast().Push(NewMatchingRuleUse().SetSchema(mySchema))
	bmr.cast().Push(NewMatchingRuleUse().SetSchema(mySchema).SetNumericOID(`1.2.3.4.5`))

	MatchingRuleUses(bmr).Push(NewMatchingRuleUse().SetSchema(mySchema))
	MatchingRuleUses(bmr).Push(NewMatchingRuleUse().SetSchema(mySchema).SetNumericOID(`1.2.3.4.5`))
	MatchingRuleUses(bmr).Compliant()

	var def MatchingRuleUse

	_ = def.String()
	_ = def.SetStringer()
	_ = def.Description()
	_ = def.Name()
	_ = def.Names()
	_ = def.Extensions()
	_ = def.Applies()
	_ = def.Schema()
	_ = def.Map()
	_ = def.Compliant()
	_ = def.macro()
	_ = def.Obsolete()

	def.setOID(`4.3.2.1`)
	var raw string = `( 2.5.13.2 NAME 'caseIgnoreMatch' APPLIES cn X-ORIGIN 'RFC4517' )`
	if err := def.Parse(raw); err == nil {
		t.Errorf("%s failed: expected parsing incompatibility error, got nothing", t.Name())
		return
	}

	def = NewMatchingRuleUse()
	def.SetDescription(`'a`)
	def.SetDescription(`'Unnecessary quoted value to be overwritten'`)

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

	_ = def.macro()
	def.setOID(`2.5.13.2`)

	var def2 MatchingRuleUse
	_ = def2.Replace(def) // will fail

	xx := mySchema.MatchingRuleUses().Get(`caseExactMatch`)
	yy := mySchema.MatchingRuleUses().Get(`caseIgnoreMatch`)
	yy.Replace(xx)

	var oo *matchingRuleUse = new(matchingRuleUse)
	var mru MatchingRuleUse = MatchingRuleUse{oo}

	name := mySchema.AttributeTypes().Get(`name`)

	_ = mySchema.MatchingRuleUses().String()
	_ = mySchema.MatchingRuleUses().push(MatchingRuleUse{})

	mr := MatchingRule{}
	_, _ = mr.makeMatchingRuleUse()
	mySchema.updateMatchingRuleUses(AttributeTypes{})

	oo.replace(MatchingRuleUse{&matchingRuleUse{schema: mySchema}})
	oo.setApplies(`cn`)
	oo.setApplies(rune(33), `cn`, nil, MatchingRuleUse{}, AttributeType{}, name)

	oo.OID = mySchema.MatchingRules().Get(`caseIgnoreMatch`)
	mru.replace(MatchingRuleUse{oo})
	mru.SetSchema(mySchema)
	_ = mru.Compliant()
	mru.setOID(`1.2.3.4.5.6.7`)
	mru.macro()
	mru.SetStringer()
	mru.stringer()

}

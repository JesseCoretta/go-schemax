package schemax

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	//"runtime/pprof"
	"testing"

	"github.com/JesseCoretta/go-antlr4512"
)

var mySchema Schema

// For release testing only
//func TestNewSchema(t *testing.T) {
//        f, err := os.Create("cpu.prof")
//        if err != nil {
//            t.Fatal(err)
//        }
//        defer f.Close()
//
//        pprof.StartCPUProfile(f)
//        defer pprof.StopCPUProfile()
//
//	mySchema := NewSchema()
//	fmt.Printf("%d types parsed", mySchema.Counters().AT)
//	// Output: 214 types parsed
//}

/*
This example demonstrates the so-called "Quick Start Schema" initialization.
The [NewSchema] function imports all built-in definitions instantly, allowing
the user to start their activities with no fuss.
*/

func ExampleNewBasicSchema() {
	mySchema := NewBasicSchema()
	fmt.Printf("%d syntaxes parsed", mySchema.Counters().LS)
	// Output: 67 syntaxes parsed
}

func ExampleSchema_Options() {
	opts := mySchema.Options()
	opts.Shift(AllowOverride)
	fmt.Println(opts.Positive(AllowOverride))
	// Output: true
}

/*
This example demonstrates a convenient means of adding any [Definition]
qualifier, such as an [ObjectClass] instance, to the schema. This is
useful when manual type verification is cumbersome.
*/
func ExampleSchema_Push() {
	def := mySchema.NewObjectClass()
	def.SetNumericOID(`1.3.6.1.4.1.56521.999.1232`)
	def.SetKind(AuxiliaryKind)
	def.SetName(`exampleClass`)

	mySchema.Push(def)
	get := mySchema.ObjectClasses().Get(`exampleClass`)

	fmt.Printf("%s\n", get.OID())
	// Output: exampleClass
}

func ExampleSchema_Exists() {
	var def Definition = mySchema.LDAPSyntaxes().Get(`integer`)
	fmt.Println(mySchema.Exists(def))
	// Output: true
}

func ExampleSchema_Replace_objectClass() {

	gon := mySchema.ObjectClasses().Get(`groupOfNames`)
	ngon := mySchema.NewObjectClass().
		SetNumericOID(gon.NumericOID()).
		SetName(gon.Name()).
		SetDescription(gon.Description()).
		SetKind(gon.Kind()).
		SetSuperClass(`top`).
		SetMust(`cn`).
		SetMay(`member`,
			`businessCategory`,
			`seeAlso`,
			`owner`,
			`ou`,
			`o`,
			`description`).
		SetExtension(`X-ORIGIN`, `RFC4519`).
		SetExtension(`X-WARNING`, `MODIFIED`). // optional
		SetStringer()

	mySchema.Options().Unshift(AllowOverride)
	mySchema.Replace(ngon)
	mySchema.Options().Shift(AllowOverride)
	mySchema.Replace(ngon)

	fmt.Println(mySchema.ObjectClasses().Get(`groupOfNames`))
	// Output: ( 2.5.6.9
	//     NAME 'groupOfNames'
	//     SUP top
	//     STRUCTURAL
	//     MUST cn
	//     MAY ( member
	//         $ businessCategory
	//         $ seeAlso
	//         $ owner
	//         $ ou
	//         $ o
	//         $ description )
	//     X-ORIGIN 'RFC4519'
	//     X-WARNING 'MODIFIED' )
}

/*
This example demonstrates refreshing the [MatchingRuleUses] collection
within the receiver instance of [Schema]. The result of this operation
is influence by any new [AttributeType] instances that have been added
since the last refresh.
*/
func ExampleSchema_UpdateMatchingRuleUses() {
	mySchema.UpdateMatchingRuleUses()
	fmt.Printf("%d matchingRuleUses present", mySchema.Counters().MU)
	// Output: 32 matchingRuleUses present
}

/*
This example demonstrates obtaining a non thread-safe [Counters] instance,
which outlines the number of [Definition] instances in categorical fashion.
*/
func ExampleSchema_Counters() {
	fmt.Printf("%d types present", mySchema.Counters().AT)
	// Output: 319 types present
}

/*
This example demonstrates accessing the [Schema] instance's distinguished
name, if set.
*/
func ExampleSchema_DN() {
	fmt.Println(mySchema.DN())
	// Output: cn=schema
}

/*
This example demonstrates specifying a non-standard distinguished name
for use by the [Schema] instance.
*/
func ExampleSchema_SetDN() {
	mySchema := NewEmptySchema()
	mySchema.SetDN(`cn=subschema`)

	fmt.Println(mySchema.DN())
	// Output: cn=subschema
}

func TestLoadSyntaxes(t *testing.T) {
	want := 67
	if got := mySchema.LDAPSyntaxes().Len(); got != want {
		t.Errorf("%s failed: want '%d' ldapSyntaxes, got '%d'",
			t.Name(), want, got)
		return
	}
}

func TestLoadMatchingRules(t *testing.T) {
	want := 44
	if got := mySchema.MatchingRules().Len(); got != want {
		t.Errorf("%s failed: want '%d' matchingRules, got '%d'",
			t.Name(), want, got)
		return
	}
}

func TestLoadAttributeTypes(t *testing.T) {
	want := 317 // includes supplementals and dcodSchema
	if got := mySchema.AttributeTypes().Len(); got != want {
		t.Errorf("%s failed: want '%d' attributeTypes, got '%d'",
			t.Name(), want, got)
		return
	}
}

func TestLoadObjectClasses(t *testing.T) {
	want := 80
	if got := mySchema.ObjectClasses().Len(); got != want {
		t.Errorf("%s failed: want '%d' objectClasses, got '%d'",
			t.Name(), want, got)
		return
	}
}

func TestSchema_codecov(t *testing.T) {
	// hopelessly try to overwrite each definition
	// with itself.  This is only necessary for
	// coverage reasons related to the Schema type;
	// actual replacements are already covered.
	for _, def := range []Definition{
		mySchema.LDAPSyntaxes().Index(0),
		mySchema.MatchingRules().Index(0),
		mySchema.AttributeTypes().Index(0),
		mySchema.MatchingRuleUses().Index(0),
		mySchema.ObjectClasses().Index(0),
		mySchema.DITContentRules().Index(0),
		mySchema.NameForms().Index(0),
		mySchema.DITStructureRules().Index(0),
	} {
		mySchema.Replace(def)
		mySchema.Exists(def)
	}

	mySchema.Replace(LDAPSyntax{})
	mySchema.Replace(ObjectClass{})
	mySchema.Replace(DITStructureRule{})
	mySchema.Replace(DITContentRule{})
	mySchema.Replace(NameForm{})
	mySchema.Replace(AttributeType{})
	mySchema.Replace(MatchingRule{})
	mySchema.Replace(MatchingRuleUse{})
	_ = mySchema.ParseDirectory(`_fnfjnds`)
	_, _ = mySchema.marshalOC(antlr4512.ObjectClass{Name: []string{`uiw`}})
	_, _ = mySchema.marshalDC(antlr4512.DITContentRule{Name: []string{`uiw`}})
	_, _ = mySchema.marshalAT(antlr4512.AttributeType{Name: []string{`uiw`}})
	_, _ = mySchema.marshalMR(antlr4512.MatchingRule{Name: []string{`uiw`}})
	_, _ = mySchema.marshalMU(antlr4512.MatchingRuleUse{Name: []string{`uiw`}})
	_, _ = mySchema.marshalNF(antlr4512.NameForm{Name: []string{`uiw`}})
	_, _ = mySchema.marshalLS(antlr4512.LDAPSyntax{Desc: `descriptive text`})
	_, _ = mySchema.marshalOC(antlr4512.ObjectClass{OID: `2.5.6.11`, Name: []string{`uiw`}})
	_, _ = mySchema.marshalAT(antlr4512.AttributeType{OID: `2.5.6.11`, Name: []string{`uiw`}})
	_, _ = mySchema.marshalDC(antlr4512.DITContentRule{OID: `2.5.6.11`, Name: []string{`uiw`}})
	_, _ = mySchema.marshalNF(antlr4512.NameForm{OID: `2.5.6.11`, Name: []string{`uiw`}})
	_, _ = mySchema.marshalAT(antlr4512.AttributeType{OID: `2.5.6.100011`, Name: []string{`uiw`}})
	_, _ = mySchema.marshalMR(antlr4512.MatchingRule{OID: `2.5.6.100011`, Name: []string{`uiw`}})
	_, _ = mySchema.marshalMU(antlr4512.MatchingRuleUse{OID: `2.5.6.100011`, Name: []string{`uiw`}})
	_, _ = mySchema.marshalDC(antlr4512.DITContentRule{OID: `2.5.6.100011`, Name: []string{`uiw`}})
	_, _ = mySchema.marshalOC(antlr4512.ObjectClass{OID: `2.5.6.100011`, Name: []string{`uiw`}})
	_, _ = mySchema.marshalMR(antlr4512.MatchingRule{OID: `2.5.6.100011`, Name: []string{`uiw`}})
	_, _ = mySchema.marshalMU(antlr4512.MatchingRuleUse{OID: `2.5.6.100011`, Name: []string{`uiw`}})
	_, _ = mySchema.marshalMR(antlr4512.MatchingRule{OID: `INTEGERT`, Name: []string{`uiw`}})
	_, _ = mySchema.marshalMR(antlr4512.MatchingRule{
		OID:    `1.3.4.5.6`,
		Syntax: `'INTEGER'`,
		Name:   []string{`uiw`},
	})
	_, _ = mySchema.marshalMU(antlr4512.MatchingRuleUse{OID: `INTEGERT`, Name: []string{`uiw`}})
	//var err error

	_ = mySchema.incorporateLS(antlr4512.LDAPSyntaxes{antlr4512.LDAPSyntax{
		OID:  `1.3.6.1.4.1.1466.115.121.1.52`,
		Desc: `uiw`,
	}})

	_ = mySchema.incorporateDS(antlr4512.DITStructureRules{antlr4512.DITStructureRule{
		ID:   `-7`,
		Form: `badForm`,
	}})

	_ = mySchema.incorporateDS(antlr4512.DITStructureRules{antlr4512.DITStructureRule{
		ID:   `7`,
		Form: `badForm`,
	}})

	_ = mySchema.incorporateDS(antlr4512.DITStructureRules{antlr4512.DITStructureRule{
		ID:         `7`,
		Form:       `nArcForm`,
		SuperRules: []string{`111`, `88`},
	}})

	_ = mySchema.incorporateNF(antlr4512.NameForms{antlr4512.NameForm{
		OID:  `1.2.3.4.54.65`,
		OC:   `2.5.6.11`,
		Must: []string{`h`},
	}})

	_ = mySchema.incorporateNF(antlr4512.NameForms{antlr4512.NameForm{
		OID:  `1.2.3.4.54.65`,
		OC:   `2.5.6.11`,
		Must: []string{`c`},
		May:  []string{`h`},
	}})

	_ = mySchema.incorporateDS(antlr4512.DITStructureRules{antlr4512.DITStructureRule{
		ID:         `7`,
		SuperRules: []string{`100`, `11`},
	}})

	_ = mySchema.incorporateAT(antlr4512.AttributeTypes{antlr4512.AttributeType{
		OID:  `2.5.4.3`,
		Name: []string{`uiw`},
	}})

	_ = mySchema.incorporateAT(antlr4512.AttributeTypes{antlr4512.AttributeType{
		OID:    `2.5.4.1999`,
		Name:   []string{`uiw`},
		Syntax: `1.3.6.1.4.1.56521.999.72.1`,
	}})

	_ = mySchema.incorporateAT(antlr4512.AttributeTypes{antlr4512.AttributeType{
		OID:      `2.5.4.1999`,
		Name:     []string{`uiw`},
		Equality: `1.3.6.1.4.1.56521.999.72.1`,
	}})

	_ = mySchema.incorporateAT(antlr4512.AttributeTypes{antlr4512.AttributeType{
		OID:       `2.5.4.1999`,
		Name:      []string{`uiw`},
		SuperType: `blarg`,
	}})

	_ = mySchema.incorporateOC(antlr4512.ObjectClasses{antlr4512.ObjectClass{
		OID:  `2.5.6.11`,
		Name: []string{`uiw`},
		Must: []string{`cn`, ``},
	}})

	_ = mySchema.incorporateMR(antlr4512.MatchingRules{antlr4512.MatchingRule{
		OID:    ``,
		Macro:  []string{`nisSchema`, `1.3.6.1.1.1.999`},
		Name:   []string{`uiw`},
		Syntax: `1.3.6.1.4.1.56521.999.72.1`,
	}})

	_ = mySchema.incorporateMR(antlr4512.MatchingRules{antlr4512.MatchingRule{
		///OID: `2.5.6.11`,
		Name:   []string{`uiw`},
		Syntax: `INTEGER`,
	}})

	_ = mySchema.incorporateMR(antlr4512.MatchingRules{antlr4512.MatchingRule{
		//OID: `1.3.6.1.4.1.1466.109.114.1`,
		//Macro: []string{`nisSchema`,`1`},
		Name:   []string{`uiw`},
		Syntax: `1.3.6.1.4.1.1466.115.121.1.11`,
	}})

	_ = mySchema.incorporateMR(antlr4512.MatchingRules{antlr4512.MatchingRule{
		OID:    `2.5.6.11`,
		Name:   []string{`uiw`},
		Syntax: `NTEGEF`,
	}})

	_ = mySchema.incorporateMU(antlr4512.MatchingRuleUses{antlr4512.MatchingRuleUse{
		OID:     `2.5.13.13`,
		Name:    []string{`uiw`},
		Applies: []string{`cn`, `dfjskalf`},
	}})

	_ = mySchema.incorporateMU(antlr4512.MatchingRuleUses{antlr4512.MatchingRuleUse{
		OID:     `2.5.13.13`,
		Name:    []string{`uiw`},
		Applies: []string{`cn`},
		Extensions: map[string][]string{
			`X-ORIGIN`: {`NOWHERE`},
		},
	}})

	_ = mySchema.incorporateOC(antlr4512.ObjectClasses{antlr4512.ObjectClass{
		OID:  `2.5.6.11`,
		Name: []string{`uiw`},
		May:  []string{`cn`, ``},
	}})

	_ = mySchema.incorporateOC(antlr4512.ObjectClasses{antlr4512.ObjectClass{
		OID:          `2.5.6.11`,
		Name:         []string{`uiw`},
		SuperClasses: []string{`p`, `account`},
	}})

	_ = mySchema.incorporateOC(antlr4512.ObjectClasses{antlr4512.ObjectClass{
		OID:          `2.5.6.1111111`,
		Name:         []string{`uiw`},
		SuperClasses: []string{`account`},
		Kind:         `JUNK`,
	}})

	_ = mySchema.incorporateOC(antlr4512.ObjectClasses{antlr4512.ObjectClass{
		OID:          `2.5.6.1111111`,
		Name:         []string{`uiw`},
		SuperClasses: []string{`account`},
		Must:         []string{`xx`},
	}})
	_ = mySchema.incorporateOC(antlr4512.ObjectClasses{antlr4512.ObjectClass{
		OID:          `2.5.6.1111111`,
		Name:         []string{`uiw`},
		SuperClasses: []string{`account`},
		Must:         []string{`cn`},
		May:          []string{`xx`},
	}})

	_ = mySchema.incorporateNF(antlr4512.NameForms{antlr4512.NameForm{
		OID:  `2.5.6.1111111`,
		Name: []string{`uiw`},
		Must: []string{`xx`},
	}})

	_ = mySchema.incorporateNF(antlr4512.NameForms{antlr4512.NameForm{
		OID:  `2.5.6.1111111`,
		Name: []string{`uiw`},
		Must: []string{`cn`},
		May:  []string{`xx`},
	}})

	_ = mySchema.incorporateNF(antlr4512.NameForms{antlr4512.NameForm{
		OID:  `2.5.6.11`,
		Name: []string{`uiw`},
		Must: []string{`cn`, ``},
	}})

	_ = mySchema.incorporateNF(antlr4512.NameForms{antlr4512.NameForm{
		OID:  `2.5.6.11`,
		Name: []string{`uiw`},
		Must: []string{`cn`},
		May:  []string{``},
	}})

	_ = mySchema.incorporateNF(antlr4512.NameForms{antlr4512.NameForm{
		OID:  `2.5.6.11`,
		Name: []string{`uiw`},
		Must: []string{`cn`},
		May:  []string{`_l`},
	}})

	_ = mySchema.incorporateDC(antlr4512.DITContentRules{antlr4512.DITContentRule{
		OID:  `2.5.6.11`,
		Name: []string{`uiw`},
		Must: []string{`cn`, ``},
	}})

	_ = mySchema.incorporateDC(antlr4512.DITContentRules{antlr4512.DITContentRule{
		OID:  `2.5.6.11`,
		Name: []string{`uiw`},
		Must: []string{`cn`},
		May:  []string{``},
	}})

	_ = mySchema.incorporateDC(antlr4512.DITContentRules{antlr4512.DITContentRule{
		OID:  `2.5.6.11`,
		Name: []string{`uiw`},
		Must: []string{`cn`},
		May:  []string{`l`},
		Not:  []string{``},
	}})

	_ = mySchema.incorporateDC(antlr4512.DITContentRules{antlr4512.DITContentRule{
		OID:  `2.5.6.11`,
		Name: []string{`uiw`},
		Must: []string{`cn`},
		May:  []string{`l`},
		Not:  []string{`drink`},
		Aux:  []string{``},
	}})

	_ = mySchema.ParseNameForm(mySchema.NameForms().Index(0).String())
	_ = mySchema.ParseDITStructureRule(mySchema.DITStructureRules().Index(0).String())
	_ = mySchema.ParseDITContentRule(mySchema.DITContentRules().Index(0).String())
	_ = mySchema.ParseMatchingRuleUse(mySchema.MatchingRuleUses().Index(0).String())
	_a := &attributeType{}
	_aa := antlr4512.AttributeType{
		Usage: `distributedOperation`,
	}
	_a.marshalUsage(_aa)
}

/*
this test generates a temporary dir/file and uses its newly written contents
to parse into our schema. It is not practical to make true Examples for these
methods :(
*/
func TestSchema_ParseFileAndDirectory(t *testing.T) {
	// Create a temp file
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}
	defer os.RemoveAll(tempDir)

	fileName := `test.schema`
	tempFile := filepath.Join(tempDir, fileName)

	text := []byte("attributeType ( 1.3.6.1.4.1.56521.999.87.1 NAME 'fakeAttribute' )")
	if err = ioutil.WriteFile(tempFile, text, 0644); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	if err = mySchema.ParseFile(tempFile); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	if err = mySchema.ParseDirectory(tempDir); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	// We generate a fake (non-existent) file/dir name rather
	// than hard-coding something bogus that could be abused...
	bogusName := randString(8)

	_ = mySchema.ParseFile(bogusName + `.schema`)
	_ = mySchema.ParseDirectory(bogusName)
}

func TestLoads_codecov(t *testing.T) {
	coolSchema := NewEmptySchema()
	coolSchema.LoadRFC4517Syntaxes()
	coolSchema.LoadRFC2307Syntaxes()
	coolSchema.LoadRFC4523Syntaxes()
	coolSchema.LoadRFC4530Syntaxes()

	coolSchema.LoadRFC4517MatchingRules()
	coolSchema.LoadRFC2307MatchingRules()
	coolSchema.LoadRFC4523MatchingRules()
	coolSchema.LoadRFC4530MatchingRules()

	coolSchema.LoadX501AttributeTypes()
	coolSchema.LoadRFC4512AttributeTypes()
	coolSchema.LoadRFC2079AttributeTypes()
	coolSchema.LoadRFC2589AttributeTypes()
	coolSchema.LoadRFC2798AttributeTypes()
	coolSchema.LoadRFC3045AttributeTypes()
	coolSchema.LoadRFC3672AttributeTypes()
	coolSchema.LoadRFC4519AttributeTypes()
	coolSchema.LoadRFC2307AttributeTypes()
	coolSchema.LoadRFC3671AttributeTypes()
	coolSchema.LoadRFC4403AttributeTypes()
	coolSchema.LoadRFC4523AttributeTypes()
	coolSchema.LoadRFC4524AttributeTypes()
	coolSchema.LoadRFC4530AttributeTypes()
	coolSchema.LoadRFC5020AttributeTypes()

	coolSchema.LoadRFC4512ObjectClasses()
	coolSchema.LoadRFC2079ObjectClasses()
	coolSchema.LoadRFC2589ObjectClasses()
	coolSchema.LoadRFC2798ObjectClasses()
	coolSchema.LoadRFC2307ObjectClasses()
	coolSchema.LoadRFC4512ObjectClasses()
	coolSchema.LoadRFC3671ObjectClasses()
	coolSchema.LoadRFC3672ObjectClasses()
	coolSchema.LoadRFC4403ObjectClasses()
	coolSchema.LoadRFC4519ObjectClasses()
	coolSchema.LoadRFC4523ObjectClasses()
	coolSchema.LoadRFC4524ObjectClasses()

	coolSchema.LoadRFC2377NameForms()
	coolSchema.LoadRFC4403NameForms()

	coolSchema.LoadRFC4403DITStructureRules()

}

// supplemental attributeTypes not sourced from an official doc, but
// are useful in UTs, et al.
//
// Currently empty.
var suplATs []string = []string{}

func init() {
	// Prepare our UT/Example reference schema
	mySchema = NewEmptySchema(
		AllowOverride,
		SortExtensions,
		SortLists,
		HangingIndents,
		AllowReindexedStructureRules,
	)

	funks := []func() error{
		mySchema.LoadLDAPSyntaxes,      // import all built-in Syntaxes
		mySchema.LoadMatchingRules,     // import all built-in Matching Rules
		mySchema.LoadAttributeTypes,    // import all built-in Attribute Types
		mySchema.LoadObjectClasses,     // import all built-in Object Classes
		mySchema.LoadNameForms,         // import all built-in Name Forms
		mySchema.LoadDITStructureRules, // import all built-in DIT Structure Rules
	}

	var err error
	for i := 0; i < len(funks) && err == nil; i++ {
		if err = funks[i](); err != nil {
			fmt.Printf("I:%d\n", i)
			panic(err)
		}
	}

	// load some supplemental attributeTypes used for special cases
	for _, at := range suplATs {
		if err := mySchema.ParseAttributeType(at); err != nil {
			panic(err)
		}
	}

	// load dcodSchema for more meaningful examples of certain
	// lesser-used definition types. This is not a built-in, as
	// it is an unapproved I-D used here strictly for robust UTs.
	if err := mySchema.ParseRaw(dcodSchema); err != nil {
		panic(err)
	}

	// Refresh our matching rules
	mySchema.updateMatchingRuleUses(mySchema.AttributeTypes())
}

/*
dcodSchema contains definitions per draft-coretta-oiddir-schema.  These
are not considered "built-ins" as they are not importable by users. They
are present here only for testing and examples.
*/
var dcodSchema []byte = []byte(`# Author: Jesse Coretta (February 29, 2024)
# Based on: https://datatracker.ietf.org/doc/html/draft-coretta-oiddir-schema-01
# NOT FOR PRODUCTION USE !
##########################
#
# 2.3.  Attribute Types
#
#    The following subsections detail one hundred three (103) attribute
#    types created for use within implementations of this specification.
#
#    Please note that a great many of these attribute type definitions are
#    sub types of attribute types defined in the following Standards, and
#    as such are considered dependencies:
#
#      - [RFC2079] for URI support
#      - [RFC4519] for so-called "core" schema elements
#      - [RFC4524] for Cosine schema elements
#
#    If the nature of a particular directory implementation precludes the
#    use of sub typed attributes, this specification may not be practical
#    for adoption.
#
# 2.3.1.  'n'
#
#    The 'n' attribute type allows the assignment of an unsigned integer
#    used to represent the Number Form of a registration per ITU-T Rec.
#    [X.680].
#
#    A practical ABNF production, labeled 'number', is defined in Section
#    1.4 of [RFC4512].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.1
           NAME ( 'n' 'numberForm' )
           DESC 'X.680, cl. 32.3: NumberForm'
           EQUALITY integerMatch
           ORDERING integerOrderingMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.27
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "56521", "0"
#
# 2.3.2.  'dotNotation'
#
#    The 'dotNotation' attribute type allows the assignment of one (1) OID
#    to any non root registration.
#
#    A practical ABNF production, labeled 'numericoid', is defined within
#    Section 1.4 of [RFC4512].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.2
           NAME 'dotNotation'
           DESC 'X.680: OID in dotted notation'
           EQUALITY objectIdentifierMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.38
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#     Examples: "1.3.6.1", "2.999"
#
# 2.3.3.  'iRI'
#
#    The 'iRI' attribute type allows the assignment of one (1) or more
#    OID-IRI values to a registration.
#
#    A practical ABNF production for this attribute type, as derived from
#    clause 34.3 of ITU-T Rec. [X.680], is as follows:
#
#      subArcId   = SOLIDUS arcId [ subArcId ]
#      arcId      = SOLIDUS ( intUV / nintUV )
#      nintUV     = 1*iunreserved ; non-integer unicode value
#      intUV      = number        ; integer unicode value
#
#      SOLIDUS    = "%x2f"        ; "/"
#
#    The 'number' ABNF production originates in Section 1.4 of [RFC4512].
#    The 'iunreserved' ABNR production originates within Section 2.2 of
#    [RFC3987].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.3
           NAME 'iRI'
           DESC 'X.680, cl. 34: OID-IRI'
           EQUALITY octetStringMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.40
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "/ITU-T", "/ISO/Identified-Organization", "/ASN.1"
#
# 2.3.4.  'aSN1Notation'
#
#    The 'aSN1Notation' attribute type allows the assignment of an OID in
#    ASN.1, or ITU-T Rec. [X.680] ObjectIdentifierValue, notation.
#
#    A practical ABNF production for this attribute type is as follows:
#
#      asn1notation = LCURL forms RCURL
#      forms        = nanf *( SPACE nanf )
#
#      SPACE        = "%x20"   ; " "
#      LCURL        = "%x7b"   ; "{"
#      RCURL        = "%x7d"   ; "}"
#
#    The 'nanf' ABNF originates in Section 2.3.19.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.4
           NAME 'aSN1Notation'
           DESC 'X.680, cl. 32.3: ObjectIdentifierValue'
           EQUALITY caseIgnoreMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.15
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "{itu-t(0)}", "{iso(1) identified-organization(3)}"
#
# 2.3.5.  'unicodeValue'
#
#    The 'unicodeValue' attribute type allows the assignment of a single
#    non-integer Unicode label to a registration, per ITU-T Rec. [X.660].
#
#    A practical ABNF production for this attribute type is as follows:
#
#      uval = 1*iunreserved
#
#    The ABNF production 'iunreserved' is defined in Section 2.2 of
#    [RFC3987].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.5
           NAME 'unicodeValue'
           DESC 'X.660, cl. 7.5: non-integer Unicode label'
           EQUALITY octetStringMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.40
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "ITU-T", "Identified-Organization"
#
# 2.3.6.  'additionalUnicodeValue'
#
#    The 'additionalUnicodeValue' attribute type allows for the assignment
#    of one (1) or more additional non-integer Unicode labels, per clause
#    3.5.2 of ITU-T Rec. [X.660], to a registration.
#
#    A practical ABNF production for this attribute type is defined within
#    Section 2.3.5.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.6
           NAME 'additionalUnicodeValue'
           DESC 'X.660, cl. 3.5.2: additional non-integer Unicode labels'
           EQUALITY octetStringMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.40
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.3.7.  'identifier'
#
#    The 'identifier' attribute type allows for the assignment of one (1)
#    non-numeric, non-Unicode identifier, or nameForm, to a registration.
#
#    Per clause 12.3 of ITU-T Rec. [X.680]:
#
#       An "identifier" shall consist of an arbitrary number (one or more)
#       of letters, digits, and hyphens.  The initial character shall be a
#       lower-case letter.  A hyphen shall not be the last character.  A
#       hyphen shall not be immediately followed by another hyphen.
#
#    As a practical ABNF production, the above clause translates as
#    follows:
#
#       identifier  = leadkeychar *keychar
#       leadkeychar = LOWER
#       keychar     = *( [ HYPHEN ] ( ALPHA / DIGIT ) )
#
#       ALPHA       = ( LOWER / UPPER ) ; a-z / A-Z
#       UPPER       = "%x41-%x5a"       ; A-Z
#       LOWER       = "%x61-%x7a"       ; a-z
#       DIGIT       = "%x30-%x39"       ; 0-9
#       HYPHEN      = "%x2d"            ; "-"
#
#    The attribute type 'name', as defined in Section 2.18 of [RFC4519],
#    is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.7
           NAME ( 'identifier' 'nameForm' )
           DESC 'X.680, cl. 12.3: Identifier'
           SUP name
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "itu-t", "iso"
#
# 2.3.8.  'secondaryIdentifier'
#
#    The 'secondaryIdentifier' attribute type allows the assignment of
#    one (1) or more additional secondary non-numeric non-Unicode values,
#    per clause 3.5.1 of ITU-T Rec. [X.660], to a registration.
#
#    A practical ABNF production for this attribute type is defined within
#    Section 2.3.7.
#
#    The attribute type 'name', as defined in Section 2.18 of [RFC4519],
#    is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.8
           NAME 'secondaryIdentifier'
           DESC 'X.660, cl. 3.5.1: Additional Identifiers'
           SUP name
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "enterprises", "ccitt"
#
# 2.3.9.  'registrationInformation'
#
#    The 'registrationInformation' attribute type allows the OPTIONAL
#    assignment of octet-based values intended for extended information
#    relating to the registration in question.
#
#    The 'OCTET' ABNF production defined in Section 1.4 of [RFC4512] is
#    applicable in any non-zero quantity or combination, as no defined
#    syntax or standard exists to govern this type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.9
           NAME 'registrationInformation'
           DESC 'Extended octet-based registration data'
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.40
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "Acme, Co., Research & Development, Copyright (c) 2024"
#
#
# 2.3.10.  'registrationURI'
#
#    The 'registrationURI' attribute type allows for the assignment of one
#    (1) or more URI values, with optional labels, to a registration.
#
#    The attribute type 'labeledURI', as defined in [RFC2079], is a super
#    type of this attribute type.
#
#    A practical ABNF production for this attribute type is defined within
#    Appendix A of [RFC3986].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.10
           NAME 'registrationURI'
           DESC 'Uniform Resource Identifier for a registration'
           SUP labeledURI
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "http://example.com Example", "ftp://example.com"
#
# 2.3.11.  'registrationCreated'
#
#    The 'registrationCreated' attribute type allows for the assignment
#    of a generalized timestamp indicating the date and time at which a
#    registration was, or will be, created or officiated.
#
#    A practical ABNF production for this attribute type is defined within
#    Section 3.3.13 of [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.11
           NAME 'registrationCreated'
           DESC 'Generalized timestamp for a registration creation'
           EQUALITY generalizedTimeMatch
           ORDERING generalizedTimeOrderingMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.24
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "19800229114853Z", "20130109033116Z"
#
# 2.3.12.  'registrationModified'
#
#    The 'registrationModified' attribute type allows for the assignment
#    of one (1) or more generalized timestamp values indicating the dates
#    and times of all applied updates to a registration.
#
#    A practical ABNF production for this attribute type is defined within
#    Section 3.3.13 of [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.12
           NAME 'registrationModified'
           DESC 'Generalized timestamps for registration modifications'
           EQUALITY generalizedTimeMatch
           ORDERING generalizedTimeOrderingMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.24
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "19951231115959Z", "20130109033116Z"
#
# 2.3.13.  'registrationRange'
#
#    The 'registrationRange' attribute type allows for the expression of
#    an OID sibling allocation range, such as "100" to indicate 'up to
#    100', or "-1" to indicate 'to infinity'.
#
#    A practical ABNF production, labeled 'number', is defined in Section
#    1.4 of [RFC4512].  This satisfies the unsigned form of instances of
#    this type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.13
           NAME 'registrationRange'
           DESC 'Numerical registration range expression'
           EQUALITY integerMatch
           ORDERING integerOrderingMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.27
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "-1", "1999", "1000000"
#
# 2.3.14.  'registrationStatus'
#
#    The 'registrationStatus' attribute type allows for the assignment of
#    status information meant to define the state of the registration.
#
#    A practical ABNF production for the super type, 'description', is
#    found within Section 3.3.6 of [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.14
           NAME 'registrationStatus'
           DESC 'Current registration status'
           SUP description
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "OBSOLETE", "DEALLOCATED", "RESERVED"
#
# 2.3.15.  'registrationClassification'
#
#    The 'registrationClassification' attribute type allows a registration
#    to bear an informal classification value, thereby describing an OID's
#    context or category.
#
#    A practical ABNF production for the super type, 'description', can be
#    found within Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.15
           NAME 'registrationClassification'
           DESC 'Known registration classification'
           SUP description
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "Standard", "Individual", "ASN.1 Modules"
#
# 2.3.16.  'isLeafNode'
#
#    The 'isLeafNode' attribute type allows for the assignment of a single
#    Boolean value indicative of whether a registration can be a parent to
#    any subordinate registrations.
#
#    A practical ABNF production for this attribute type is found within
#    Section 3.3.3 of [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.16
           NAME 'isLeafNode'
           DESC 'Whether a registration may allocate sub arcs'
           EQUALITY booleanMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.7
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "TRUE", "FALSE", or Undefined (implies "FALSE")
#
# 2.3.17.  'isFrozen'
#
#    The 'isFrozen' attribute type allows for the assignment of a single
#    Boolean value indicative of whether a registration can be a parent
#    to any further subordinate registrations beyond those that already
#    exist at present.
#
#    A practical ABNF production for this attribute type is found within
#    Section 3.3.3 of [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.17
           NAME 'isFrozen'
           DESC 'Whether additional sub arcs are permitted'
           EQUALITY booleanMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.7
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "TRUE", "FALSE", or Undefined (implies "FALSE")
#
# 2.3.18.  'standardizedNameForm'
#
#    The 'standardizedNameForm' attribute type allows for the assignment
#    of one (1) or more Standardized NameForm values, per clauses A.2 and
#    A.3 of ITU-T Rec. [X.660], to a registration only if its 'identifier'
#    value is considered standardized.
#
#    A practical ABNF production for this attribute type is as follows:
#
#      stdnf = LCURL nonfs RCURL
#      nonfs = nonf *( SPACE nonf )
#      nonf  = ( identifier / number ) ; name OR number
#
#      SPACE = "%x20"      ; " "
#      LCURL = "%x7b"      ; "{"
#      RCURL = "%x7d"      ; "}"
#
#    The 'identifier' ABNF originates in Section 2.3.7.  The 'number' ABNF
#    production can be found within Section 1.4 of [RFC4512].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.18
           NAME 'standardizedNameForm'
           DESC 'X.660, cl. A.2-A-3: Standardized Identifier or NameForm'
           EQUALITY caseExactMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.15
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "{itu-t}", "{0 0 d}"
#
# 2.3.19.  'nameAndNumberForm'
#
#    The 'nameAndNumberForm' attribute type allows for the assignment of
#    an ITU-T Rec. [X.680] NameAndNumberForm value to a registration.
#
#    A practical ABNF production for this attribute type is as follows:
#
#      nonf       = ( nanf  / number )        ; nanf OR numberForm
#      nanf       = name LPAREN number RPAREN ; name AND numberForm
#      name       = identifier                ; name form
#
#      LPAREN     = "%x28"                    ; "("
#      RPAREN     = "%x29"                    ; ")"
#
#    The 'identifier' ABNF production can be found in Section 2.3.7. The
#    'number' ABNF production is defined in Section 1.4 of [RFC4512].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.19
           NAME 'nameAndNumberForm'
           DESC 'X.680, cl. 32.3: NameAndNumberForm'
           EQUALITY caseExactMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.15
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "private(4)", "itu-t(0)", "56521"
#
# 2.3.20.  'longArc'
#
#    The 'longArc' attribute type allows for the assignment of one (1) or
#    more so-called "Long Arc" well-known identifiers to a registration.
#
#    A practical ABNF production for this attribute type is as follows:
#
#      longArc    = SOLIDUS ( intUV / nintUV )
#      nintUV     = 1*iunreserved ; non-integer unicode value
#      intUV      = number        ; integer unicode value
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.20
           NAME 'longArc'
           DESC 'X.660, cl. A.7: Long Arc'
           EQUALITY octetStringMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.40
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "/Example", "/Ejemplo"
#
# 2.3.21.  'supArc'
#
#    The 'supArc' attribute type allows for the assignment of an LDAP DN
#    value to any non-root registration, thereby identifying the DN of the
#    immediate superior (parent) registration.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.21
           NAME 'supArc'
           DESC 'Immediate superior registration DN'
           SUP distinguishedName
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "n=1,n=2,ou=Registrations,o=rA"
#
# 2.3.22.  'c-supArc'
#
#    The 'c-supArc' attribute type is the manifestation of the 'supArc'
#    attribute type defined in Section 2.3.21 with Collective Attribute
#    [RFC3671] support.
#
#    This attribute type should only be used in directory environments
#    which actively support and require [RFC3671] capabilities.
#
#    This attribute type MUST NOT be present for entries that also bear
#    a 'supArc' attribute type instance.  The value MUST be singular.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.22
           NAME 'c-supArc'
           DESC 'Collective immediate superior registration DN'
           SUP distinguishedName
           COLLECTIVE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "n=1,n=2,ou=Registrations,o=rA"
#
# 2.3.23.  'topArc'
#
#    The 'topArc' attribute type allows for the assignment of an LDAP DN
#    value to any non-root registration, thereby identifying the superior
#    root registration.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.23
           NAME 'topArc'
           DESC 'Superior root registration DN'
           SUP distinguishedName
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "n=2,ou=Registrations,o=rA"
#
# 2.3.24.  'c-topArc'
#
#    The 'c-topArc' attribute type is the manifestation of the 'topArc'
#    attribute type defined in Section 2.3.23 with Collective Attribute
#    [RFC3671] support.
#
#    This attribute type should only be used in directory environments
#    which actively support and require [RFC3671] capabilities.
#
#    This attribute type MUST NOT be present for entries that also bear
#    a 'topArc' attribute type instance.  The value MUST be singular.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.24
           NAME 'c-topArc'
           DESC 'Collective superior root registration DN'
           SUP distinguishedName
           COLLECTIVE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "n=2,ou=Registrations,o=rA"
#
# 2.3.25.  'subArc'
#
#    The 'subArc' attribute type allows for the assignment of one (1) or
#    more LDAP DN values to a registration as a manifest of subordinate
#    registrations residing exactly one (1) logical level below, if any.
#
#    In robust implementations of this ID that cover most (or all) of the
#    OID Tree, use of this attribute type will require careful, long-term
#    planning.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.25
           NAME 'subArc'
           DESC 'Subordinate registration DN'
           SUP distinguishedName
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "n=1,n=6,n=3,n=1,ou=Registrations,o=rA"
#
# 2.3.26.  'leftArc'
#
#    The 'leftArc' attribute type allows for the assignment of a DN value
#    to a registration, referencing a registration positioned to the left
#    side of the bearer in terms of (lesser) 'n' magnitude.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.26
           NAME 'leftArc'
           DESC 'Nearest antecedent registration DN'
           SUP distinguishedName
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "n=5,n=2,ou=Registrations,o=rA"
#
# 2.3.27.  'minArc'
#
#    The 'minArc' attribute type allows for the assignment of a DN value
#    to a registration.  The value SHOULD reference the entry which bears
#    the lowest 'n' value of any of its siblings.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.27
           NAME 'minArc'
           DESC 'First or left-most sibling registration DN'
           SUP distinguishedName
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "n=0,n=2,ou=Registrations,o=rA"
#
# 2.3.28.  'c-minArc'
#
#    The 'c-minArc' attribute type is the manifestation of the attribute
#    type 'minArc', defined in Section 2.3.27 with Collective Attribute
#    [RFC3671] support.
#
#    This attribute type should only be used in directory environments
#    which actively support and require [RFC3671] capabilities.
#
#    This attribute type MUST NOT be present for entries that also bear
#    a 'minArc' attribute type instance.  The value MUST be singular.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.28
           NAME 'c-minArc'
           DESC 'Collective first or left-most sibling registration DN'
           SUP distinguishedName
           COLLECTIVE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "n=0,n=2,ou=Registrations,o=rA"
#
# 2.3.29.  'rightArc'
#
#    The 'rightArc' attribute type allows for the assignment of a DN value
#    to a registration, referencing a registration positioned to the right
#    side of the bearer in terms of (greater) 'n' magnitude.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.29
           NAME 'rightArc'
           DESC 'Nearest subsequent registration DN'
           SINGLE-VALUE
           SUP distinguishedName
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "n=2,n=2,ou=Registrations,o=rA"
#
# 2.3.30.  'maxArc'
#
#    The 'maxArc' attribute type allows for the assignment of a DN value
#    to a registration.  The value SHOULD reference the entry which bears
#    the highest 'n' value of any of its siblings.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.30
           NAME 'maxArc'
           DESC 'Final or right-most sibling registration DN'
           SUP distinguishedName
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "n=999,n=2,ou=Registrations,o=rA"
#
# 2.3.31.  'c-maxArc'
#
#    The 'c-maxArc' attribute type is the manifestation of the attribute
#    type 'maxArc', defined in Section 2.3.30 with Collective Attribute
#    [RFC3671] support.
#
#    This attribute type should only be used in directory environments
#    which actively support and require [RFC3671] capabilities.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    This attribute type MUST NOT be present for entries that also bear
#    a 'maxArc' attribute type instance.  The value MUST be singular.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.31
           NAME 'c-maxArc'
           DESC 'Collective final or right-most sibling registration DN'
           SUP distinguishedName
           COLLECTIVE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "n=999,n=2,ou=Registrations,o=rA"
#
# 2.3.32.  'discloseTo'
#
#    The 'discloseTo' attribute type allows for the assignment of one (1)
#    or more LDAP DN values to a registration, each of which reference an
#    identity that is authorized to access one (1) or more confidential
#    registrations in read-only fashion.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.32
           NAME 'discloseTo'
           DESC 'Authorized registration reader'
           SUP distinguishedName
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "cn=AuthorizedReader,ou=Groups,o=rA"
#
# 2.3.33.  'c-discloseTo'
#
#    The 'c-discloseTo' attribute type is the COLLECTIVE manifestation of
#    the attribute type 'discloseTo', defined in Section 2.3.32.
#
#    This attribute type should only be used in directory environments
#    which actively support and require [RFC3671] capabilities.
#
#    This attribute type MAY be present for entries that also bear a
#    'discloseTo' attribute type instance.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.33
           NAME 'c-discloseTo'
           DESC 'Collective authorized registration reader'
           SUP distinguishedName
           COLLECTIVE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "cn=ClearanceLevel4,ou=Groups,o=rA"
#
# 2.3.34.  'registrantID'
#
#    The 'registrantID' attribute type is intended to allow for singular
#    assignment of a UUID, GUID or some other auto-generated value to a
#    registrant entry.
#
#    No specific syntax is implied for values of this type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.34
           NAME 'registrantID'
           DESC 'Unambiguous identifier assigned to an authority'
           EQUALITY octetStringMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.40
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "rfc4519", "65116e61-cc02-4c50-bde7-5bdaf4e973e4"
#
# 2.3.35.  'currentAuthority'
#
#    The 'currentAuthority' attribute type allows for the assignment of
#    one (1) or more DN values to a registration.
#
#    The value(s) of this attribute type are meant to refer to distinct
#    entries that contain current registrant authority information for
#    the registration to which it is linked.
#
#    This attribute type is only required if registrant information is not
#    stored within a given registration directly.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.35
           NAME 'currentAuthority'
           DESC 'DN for a registrant serving as the current authority'
           SUP distinguishedName
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "registrantID=XYZ,ou=Registrants,o=rA"
#
# 2.3.36.  'c-currentAuthority'
#
#    The 'c-currentAuthority' attribute type is the manifestation of the
#    'currentAuthority' attribute type, defined in Section 2.3.35 with
#    Collective Attribute [RFC3671] support.
#
#    This attribute type should only be used in directory environments
#    which actively support and require [RFC3671] capabilities.
#
#    This attribute type MAY be present for entries that also bear
#    a 'currentAuthority' attribute type instance.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.36
           NAME 'c-currentAuthority'
           DESC 'Collective DN for a current authority'
           SUP distinguishedName
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "registrantID=XYZ,ou=Registrants,o=rA"
#
# 2.3.37.  'currentAuthorityStartTimestamp'
#
#    The 'currentAuthorityStartTimestamp' attribute type allows for the
#    assignment of a generalized timestamp value to a current registration
#    authority.  The value is indicative of the date and time at which the
#    current registration authority was, or will be, officiated.
#
#    A practical ABNF production for this attribute type is defined within
#    Section 3.3.13 of [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.37
           NAME 'currentAuthorityStartTimestamp'
           DESC 'Registration authority commencement timestamp'
           EQUALITY generalizedTimeMatch
           ORDERING generalizedTimeOrderingMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.24
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "19951231115959Z"
#
# 2.3.38.  'currentAuthorityCommonName'
#
#    The 'currentAuthorityCommonName' attribute type allows for the
#    assignment of a common name to a current authority entry.
#
#    The attribute type 'cn', as defined in Section 2.3 of [RFC4519],
#    is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.38
           NAME 'currentAuthorityCommonName'
           DESC 'Common Name for a current authority'
           SUP cn
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "Jesse Coretta", "Jane Smith"
#
# 2.3.39.  'currentAuthorityCountryCode'
#
#    The 'currentAuthorityCountryCode' attribute type allows for the
#    assignment of a country code to a current authority entry.
#
#    The attribute type 'c', as defined in Section 2.2 of [RFC4519],
#    is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.4 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.39
           NAME 'currentAuthorityCountryCode'
           DESC 'Country Code for a current authority'
           SUP c
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "US", "CA"
#
# 2.3.40.  'currentAuthorityCountryName'
#
#    The 'currentAuthorityCountryName' attribute type allows for the
#    assignment of a country name to a current authority entry.
#
#    The attribute type 'co', as defined in Section 2.4 of [RFC4519],
#    is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.40
           NAME 'currentAuthorityCountryName'
           DESC 'Country name for a current authority'
           SUP co
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "United States", "Canada"
#
# 2.3.41.  'currentAuthorityEmail'
#
#    The 'currentAuthorityEmail' attribute type allows for the assignment
#    of an email address to the current registration authority entry.
#
#    The attribute type 'mail', as defined in Section 2.16 of [RFC4524],
#    is a super type of this attribute type.
#
#    A practical ABNF production can be found in Section 3.2 of [RFC4517],
#    labeled 'IA5String'.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.41
           NAME 'currentAuthorityEmail'
           DESC 'Email address for a current authority'
           SUP mail
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "jesse.coretta@icloud.com"
#
# 2.3.42.  'currentAuthorityFax'
#
#    The 'currentAuthorityFax' attribute type allows for the assignment
#    of a facsimile telephone number to a current authority entry.
#
#    The attribute type 'facsimileTelephoneNumber', as defined in Section
#    2.10 of [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production can be found within Section 3.3.11 of
#    [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.42
           NAME 'currentAuthorityFax'
           DESC 'Facsimile telephone number for a current authority'
           SUP facsimileTelephoneNumber
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "+11234567890"
#
# 2.3.43.  'currentAuthorityLocality'
#
#    The 'currentAuthorityLocality' attribute type allows for a locality
#    name to be assigned to a current authority entry.
#
#    The attribute type 'l', as defined in Section 2.16 of [RFC4519], is
#    a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.43
           NAME 'currentAuthorityLocality'
           DESC 'Locality name for a current authority'
           SUP l
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "Palm Springs", "Anna Maria Island"
#
# 2.3.44.  'currentAuthorityMobile'
#
#    The 'currentAuthorityMobile' attribute type allows for the assignment
#    of a mobile telephone number to a current authority entry.
#
#    The attribute type 'mobile', as defined in Section 2.18 of [RFC4524],
#    is a super type of this attribute type.
#
#    A practical ABNF production can be found in Section 3.2 of [RFC4517],
#    labeled 'PrintableString'.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.44
           NAME 'currentAuthorityMobile'
           DESC 'Mobile telephone number for a current authority'
           SUP mobile
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "+11234567890"
#
# 2.3.45.  'currentAuthorityOrg'
#
#    The 'currentAuthorityOrg' attribute type allows for the assignment of
#    an organization name to a current authority entry.
#
#    The attribute type 'o', as defined in Section 2.19 of [RFC4519], is
#    a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.45
           NAME 'currentAuthorityOrg'
           DESC 'Organization name for a current authority'
           SUP o
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "Acme, Co."
#
# 2.3.46.  'currentAuthorityPOBox'
#
#    The 'currentAuthorityPOBox' attribute type allows for the assignment
#    of a post office box number to a current authority entry.
#
#    The attribute type 'postOfficeBox', as defined in Section 2.25 of
#    [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.46
           NAME 'currentAuthorityPOBox'
           DESC 'Post office box number for a current authority'
           SUP postOfficeBox
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "555", "475"
#
# 2.3.47.  'currentAuthorityPostalAddress'
#
#    The 'currentAuthorityPostalAddress' attribute type allows for the
#    assignment of a complete postal address to a current authority entry.
#    This single attribute may be used instead of other individual address
#    component attribute types, but will require field parsing on part of
#    the client.
#
#    The attribute type 'postalAddress', as defined in Section 2.23 of
#    [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.28 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.47
           NAME 'currentAuthorityPostalAddress'
           DESC 'Full postal address for a current authority'
           SUP postalAddress
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "1 Fake St$Anytown$CA$12345$US"
#
# 2.3.48.  'currentAuthorityPostalCode'
#
#    The 'currentAuthorityPostalCode' attribute type allows for a postal
#    code to be assigned to a current authority entry.
#
#    The attribute type 'postalCode', as defined in Section 2.23 of
#    [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.48
           NAME 'currentAuthorityPostalCode'
           DESC 'Postal code for a current authority'
           SUP postalCode
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#         Examples: "92262", "34216"
#
# 2.3.49.  'currentAuthorityState'
#
#    The 'currentAuthorityState' attribute type allows for a state or
#    province name to be assigned to a current authority entry.
#
#    The attribute type 'st', as defined in Section 2.33 of [RFC4519],
#    is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.49
           NAME 'currentAuthorityState'
           DESC 'State or province name for a current authority'
           SUP st
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "California", "North Dakota"
#
# 2.3.50.  'currentAuthorityStreet'
#
#    The 'currentAuthorityStreet' attribute type allows for the assignment
#    of a street name and number to a current authority entry.
#
#    The attribute type 'street', as defined in Section 2.34 of [RFC4519],
#    is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.50
           NAME 'currentAuthorityStreet'
           DESC 'Street name and number for a current authority'
           SUP street
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "1 Fake Street"
#
# 2.3.51.  'currentAuthorityTelephone'
#
#    The 'currentAuthorityTelephone' attribute type allows for a telephone
#    number to be assigned to a current authority entry.
#
#    The attribute type 'telephoneNumber', as defined in Section 2.35 of
#    [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production can be found in Section 3.2 of [RFC4517],
#    labeled 'PrintableString'.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.51
           NAME 'currentAuthorityTelephone'
           DESC 'Telephone number assigned to a current authority'
           SUP telephoneNumber
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "+11234567890"
#
# 2.3.52.  'currentAuthorityTitle'
#
#    The 'currentAuthorityTitle' attribute type allows for the assignment
#    of an official or professional title to a current authority entry.
#
#    The attribute type 'title', as defined in Section 2.38 of [RFC4519],
#    is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.52
           NAME 'currentAuthorityTitle'
           DESC 'Title assigned to a current authority'
           SUP title
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "Chief Engineer"
#
# 2.3.53.  'currentAuthorityURI'
#
#    The 'currentAuthorityURI' attribute type allows for the assignment of
#    one (1) or more URI values to a current authority entry.
#
#    The attribute type 'labeledURI', as defined in [RFC2079], is a super
#    type of this attribute type.
#
#    A practical ABNF production for this attribute type is defined within
#    Appendix A of [RFC3986].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.53
           NAME 'currentAuthorityURI'
           DESC 'Uniform Resource Identifier for a current authority'
           SUP labeledURI
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "http://example.com Example", "http://example.com"
#
# 2.3.54.  'firstAuthority'
#
#    The 'firstAuthority' attribute type allows for the assignment of one
#    (1) or more DN values to a registration entry.
#
#    The value(s) of this attribute type are meant to refer to distinct
#    entries that contain previous authority information.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.54
           NAME 'firstAuthority'
           DESC 'DN of a previous authority'
           SUP distinguishedName
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "registrantID=XYZ,ou=Registrants,o=rA"
#
# 2.3.55.  'c-firstAuthority'
#
#    The 'c-firstAuthority' attribute type is the manifestation of the
#    'firstAuthority' attribute type, defined in Section 2.3.54 with
#    Collective Attribute [RFC3671] support.
#
#    This attribute type should only be used in directory environments
#    which actively support and require [RFC3671] capabilities.
#
#    This attribute type MAY be present for entries that also bear
#    a 'firstAuthority' attribute type instance.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.55
           NAME 'c-firstAuthority'
           DESC 'Collective DN of a previous authority'
           SUP distinguishedName
           COLLECTIVE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "registrantID=XYZ,ou=Registrants,o=rA"
#
# 2.3.56.  'firstAuthorityStartTimestamp'
#
#    The 'firstAuthorityStartTimestamp' attribute type allows for the
#    assignment of a generalized timestamp value to a previous authority,
#    indicative of the date and time at which the previous authority was
#    officiated.
#
#    A practical ABNF production for this attribute type is defined within
#    Section 3.3.13 of [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.56
           NAME 'firstAuthorityStartTimestamp'
           DESC 'Previous registration authority commencement timestamp'
           EQUALITY generalizedTimeMatch
           ORDERING generalizedTimeOrderingMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.24
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "20130105135904Z"
#
# 2.3.57.  'firstAuthorityEndTimestamp'
#
#    The 'firstAuthorityEndTimestamp' attribute type allows the assignment
#    of a generalized timestamp value to a previous authority, indicative
#    of the date and time at which an entity's authoritative role was, or
#    will be, terminated.
#
#    A practical ABNF production for this attribute type is defined within
#    Section 3.3.13 of [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.57
           NAME 'firstAuthorityEndTimestamp'
           DESC 'Previous registration authority termination timestamp'
           EQUALITY generalizedTimeMatch
           ORDERING generalizedTimeOrderingMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.24
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "20170528110555Z"
#
# 2.3.58.  'firstAuthorityCommonName'
#
#    The 'firstAuthorityCommonName' attribute type allows for a common
#    name to be assigned to a previous registration authority entry.
#
#    The attribute type 'cn', as defined in Section 2.3 of [RFC4519], is
#    a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.58
           NAME 'firstAuthorityCommonName'
           DESC 'Common Name for a previous authority'
           SUP cn
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "Jesse Coretta", "Jane Smith"
#
# 2.3.59.  'firstAuthorityCountryCode'
#
#    The 'firstAuthorityCountryCode' attribute type allows for a country
#    code to be assigned to a previous registration authority entry.
#
#    The attribute type 'c', as defined in Section 2.2 of [RFC4519], is a
#    super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.4 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.59
           NAME 'firstAuthorityCountryCode'
           DESC 'Country Code for a previous authority'
           SUP c
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "US", "CA"
#
# 2.3.60.  'firstAuthorityCountryName'
#
#    The 'firstAuthorityCountryName' attribute type allows for a country
#    name to be assigned to a previous registration authority entry.
#
#    The attribute type 'co', as defined in Section 2.4 of [RFC4519], is a
#    super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.60
           NAME 'firstAuthorityCountryName'
           DESC 'Country name for a previous authority'
           SUP co
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "United States", "Canada"
#
# 2.3.61.  'firstAuthorityEmail'
#
#    The 'firstAuthorityEmail' attribute type allows for the assignment
#    of an email address to a previous registration authority entry.
#
#    The attribute type 'mail', as defined in Section 2.16 of [RFC4524],
#    is a super type of this attribute type.
#
#    A practical ABNF production can be found in Section 3.2 of [RFC4517],
#    labeled 'IA5String'.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.61
           NAME 'firstAuthorityEmail'
           DESC 'Email address for a previous authority'
           SUP mail
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "jesse.coretta@icloud.com"
#
# 2.3.62.  'firstAuthorityFax'
#
#    The 'firstAuthorityFax' attribute type allows for the assignment of a
#    facsimile telephone number to a previous authority entry.
#
#    The attribute type 'facsimileTelephoneNumber', as defined in Section
#    2.10 of [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production can be found within Section 3.3.11 of
#    [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.62
           NAME 'firstAuthorityFax'
           DESC 'Facsimile telephone number for a previous authority'
           SUP facsimileTelephoneNumber
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "+11234567890"
#
# 2.3.63.  'firstAuthorityLocality'
#
#    The 'firstAuthorityLocality' attribute type allows the assignment of
#    a locality name to a previous authority entry.
#
#    The attribute type 'l', as defined in Section 2.16 of [RFC4519], is
#    a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.63
           NAME 'firstAuthorityLocality'
           DESC 'Locality name for a previous authority'
           SUP l
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "Palm Springs", "Anna Maria Island"
#
# 2.3.64.  'firstAuthorityMobile'
#
#    The 'firstAuthorityMobile' attribute type allows for the assignment
#    of a mobile telephone number to a previous authority entry.
#
#    The attribute type 'mobile', as defined in Section 2.18 of [RFC4524],
#    is a super type of this attribute type.
#
#    A practical ABNF production can be found in Section 3.2 of [RFC4517],
#    labeled 'PrintableString'.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.64
           NAME 'firstAuthorityMobile'
           DESC 'Mobile telephone number for a previous authority'
           SUP mobile
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "+11234567890"
#
# 2.3.65.  'firstAuthorityOrg'
#
#    The 'firstAuthorityOrg' attribute type allows for the assignment
#    of an organization name to a previous authority entry.
#
#    The attribute type 'o', as defined in Section 2.19 of [RFC4519], is
#    a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.65
           NAME 'firstAuthorityOrg'
           DESC 'Organization name for a previous authority'
           SUP o
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "Acme, Co."
#
# 2.3.66.  'firstAuthorityPOBox'
#
#    The 'firstAuthorityPOBox' attribute type allows for the assignment
#    of a post office box number to a previous registration authority
#    entry.
#
#    The attribute type 'postOfficeBox', as defined in Section 2.25 of
#    [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.66
           NAME 'firstAuthorityPOBox'
           DESC 'Post office box number for a previous authority'
           SUP postOfficeBox
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "555", "475"
#
# 2.3.67.  'firstAuthorityPostalAddress'
#
#    The 'firstAuthorityPostalAddress' attribute type allows for the
#    assignment of a complete postal address to a previous registration
#    authority entry.  This single attribute may be used instead of other
#    individual address component attribute types, but will require field
#    parsing on the client side.
#
#    The attribute type 'postalAddress', as defined in Section 2.23 of
#    [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.28 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.67
           NAME 'firstAuthorityPostalAddress'
           DESC 'Full postal address for a previous authority'
           SUP postalAddress
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "1 Fake St$Anytown$CA$12345$US"
#
# 2.3.68.  'firstAuthorityPostalCode'
#
#    The 'firstAuthorityPostalCode' attribute type allows for the
#    assignment of a postal code to a previous registration authority
#    entry.
#
#    The attribute type 'postalCode', as defined in Section 2.23 of
#    [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.68
           NAME 'firstAuthorityPostalCode'
           DESC 'Postal code for a previous authority'
           SUP postalCode
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "92262", "34216"
#
# 2.3.69.  'firstAuthorityState'
#
#    The 'firstAuthorityState' attribute type allows for the assignment
#    of a state or province name to a previous registration authority
#    entry.
#
#    The attribute type 'st', as defined in Section 2.33 of [RFC4519], is
#    a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.69
           NAME 'firstAuthorityState'
           DESC 'State or province name for a previous authority'
           SUP st
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "California", "North Dakota"
#
# 2.3.70.  'firstAuthorityStreet'
#
#    The 'firstAuthorityStreet' attribute type allows for the assignment
#    of a street name and number to a previous authority entry.
#
#    The attribute type 'street', as defined in Section 2.34 of [RFC4519],
#    is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.70
           NAME 'firstAuthorityStreet'
           DESC 'Street name and number for a previous authority'
           SUP street
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "1 Fake Street"
#
# 2.3.71.  'firstAuthorityTelephone'
#
#    The 'firstAuthorityTelephone' attribute type allows the assignment of
#    a telephone number to a previous authority entry.
#
#    The attribute type 'telephoneNumber', as defined in Section 2.35 of
#    [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production can be found in Section 3.2 of [RFC4517],
#    labeled 'PrintableString'.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.71
           NAME 'firstAuthorityTelephone'
           DESC 'Telephone number for a previous authority'
           SUP telephoneNumber
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "+11234567890"
#
# 2.3.72.  'firstAuthorityTitle'
#
#    The 'firstAuthorityTitle' attribute type allows for the assignment
#    of an official or professional title to a previous authority entry.
#
#    The attribute type 'title', as defined in Section 2.38 of [RFC4519],
#    is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.72
           NAME 'firstAuthorityTitle'
           DESC 'Title of a previous authority'
           SUP title
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "Chief Engineer"
#
# 2.3.73.  'firstAuthorityURI'
#
#    The 'firstAuthorityURI' attribute type allows the assignment of one
#    (1) or more URI values to a previous authority entry.
#
#    The attribute type 'labeledURI', as defined in [RFC2079], is a super
#    type of this attribute type.
#
#    A practical ABNF production for this attribute type is defined within
#    Appendix A of [RFC3986].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.73
           NAME 'firstAuthorityURI'
           DESC 'Uniform Resource Identifier for a previous authority'
           SUP labeledURI
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "http://example.com Example", "http://example.com"
#
# 2.3.74.  'sponsor'
#
#    The 'sponsor' attribute type allows for the assignment of one (1)
#    or more DN values to a registration.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
#    The attribute type 'distinguishedName', as defined in Section 2.7 of
#    [RFC4519], is a super type of this attribute type.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.74
           NAME 'sponsor'
           DESC 'DN of a sponsoring authority'
           SUP distinguishedName
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "registrantID=XYZ,ou=Registrants,o=rA"
#
# 2.3.75.  'c-sponsor'
#
#    The 'c-sponsor' attribute type is the manifestation of the 'sponsor'
#    attribute type, defined in Section 2.3.74 with Collective Attribute
#    [RFC3671] support.
#
#    This attribute type should only be used in directory environments
#    which actively support and require [RFC3671] capabilities.
#
#    This attribute type MAY be present for entries that also bear
#    a 'sponsor' attribute type instance.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.75
           NAME 'c-sponsor'
           DESC 'Collective DN of a sponsoring authority'
           SUP distinguishedName
           COLLECTIVE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "registrantID=XYZ,ou=Registrants,o=rA"
#
# 2.3.76.  'sponsorStartTimestamp'
#
#    The 'sponsorStartTimestamp' attribute type allows for the assignment
#    of a generalized timestamp value to a sponsor entry, indicative of
#    the date and time at which sponsorship was, or will be, officiated.
#
#    A practical ABNF production for this attribute type is defined within
#    Section 3.3.13 of [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.76
           NAME 'sponsorStartTimestamp'
           DESC 'Sponsoring authority commencement timestamp'
           EQUALITY generalizedTimeMatch
           ORDERING generalizedTimeOrderingMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.24
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "20130105135904Z"
#
# 2.3.77.  'sponsorEndTimestamp'
#
#    The 'sponsorEndTimestamp' attribute type allows for the assignment
#    of a generalized timestamp value to a sponsor entry, indicative the
#    date and time at which sponsorship was, or will be, terminated.
#
#    A practical ABNF production for this attribute type is defined within
#    Section 3.3.13 of [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.77
           NAME 'sponsorEndTimestamp'
           DESC 'Sponsoring authority termination timestamp'
           EQUALITY generalizedTimeMatch
           ORDERING generalizedTimeOrderingMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.24
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "20170528110555Z"
#
# 2.3.78.  'sponsorCommonName'
#
#    The 'sponsorCommonName' attribute type allows for the assignment
#    of a common name to a sponsor.
#
#    The attribute type 'cn', as defined in Section 2.3 of [RFC4519],
#    is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.78
           NAME 'sponsorCommonName'
           DESC 'Common Name of a sponsor'
           SUP cn
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "Jane Sponsor"
#
# 2.3.79.  'sponsorCountryCode'
#
#    The 'sponsorCountryCode' attribute type allows for the assignment of
#    a two-letter country code to a sponsor.
#
#    The attribute type 'c', as defined in Section 2.2 of [RFC4519], is a
#    super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.4 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.79
           NAME 'sponsorCountryCode'
           DESC 'Country code for a sponsor'
           SUP c
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "US", "CA"
#
# 2.3.80.  'sponsorCountryName'
#
#    The 'sponsorCountryName' attribute type allows the assignment of a
#    country name to a sponsor.
#
#    The attribute type 'co', as defined in Section 2.4 of [RFC4524], is
#    a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.80
           NAME 'sponsorCountryName'
           DESC 'Country name for a sponsor'
           SUP co
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "United States", "Canada"
#
# 2.3.81.  'sponsorEmail'
#
#    The 'sponsorEmail' attribute type allows for the assignment of an
#    email address to a sponsor.
#
#    The attribute type 'mail', as defined in Section 2.16 of [RFC4524],
#    is a super type of this attribute type.
#
#    A practical ABNF production can be found in Section 3.2 of [RFC4517],
#    labeled 'IA5String'.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.81
           NAME 'sponsorEmail'
           DESC 'Email address for a sponsor'
           SUP mail
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "sponsor@example.com"
#
# 2.3.82.  'sponsorFax'
#
#    The 'sponsorFax' attribute type allows for the assignment of a
#    facsimile telephone number to a sponsor.
#
#    The attribute type 'facsimileTelephoneNumber', as defined in Section
#    2.10 of [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production can be found within Section 3.3.11 of
#    [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.82
           NAME 'sponsorFax'
           DESC 'Facsimile telephone number for a sponsor'
           SUP facsimileTelephoneNumber
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "+11234567890"
#
# 2.3.83.  'sponsorLocality'
#
#    The 'sponsorLocality' attribute type allows for the assignment of
#    a locality name to a sponsor.
#
#    The attribute type 'l', as defined in Section 2.16 of [RFC4519], is
#    a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.83
           NAME 'sponsorLocality'
           DESC 'Locality name for a sponsor'
           SUP l
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "Palm Springs", "Anna Maria Island"
#
# 2.3.84.  'sponsorMobile'
#
#    The 'sponsorMobile' attribute type allows for the assignment of a
#    mobile telephone number to a sponsor.
#
#    The attribute type 'mobile', as defined in Section 2.18 of [RFC4524],
#    is a super type of this attribute type.
#
#    A practical ABNF production can be found in Section 3.2 of [RFC4517],
#    labeled 'PrintableString'.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.84
           NAME 'sponsorMobile'
           DESC 'Mobile telephone number for a sponsor'
           SUP mobile
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "+11234567890"
#
# 2.3.85.  'sponsorOrg'
#
#    The 'sponsorOrg' attribute type allows for the assignment of an
#    organization name to a sponsor.
#
#    The attribute type 'o', as defined in Section 2.19 of [RFC4519],
#    is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.85
           NAME 'sponsorOrg'
           DESC 'Organization name for a sponsor'
           SUP o
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "Sponsor, Co."
#
# 2.3.86.  'sponsorPOBox'
#
#    The 'sponsorPOBox' attribute type allows for the assignment of a
#    post office box number to a sponsor.
#
#    The attribute type 'postOfficeBox', as defined in Section 2.25 of
#    [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.86
           NAME 'sponsorPOBox'
           DESC 'Post office box number for a sponsor'
           SUP postOfficeBox
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "555", "475"
#
# 2.3.87.  'sponsorPostalAddress'
#
#    The 'sponsorPostalAddress' attribute type allows for the assignment
#    of a complete postal address sponsor.  This single attribute may be
#    used instead of other individual address component attribute types,
#    but will require field parsing on the client side.
#
#    The attribute type 'postalAddress', as defined in Section 2.23 of
#    [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.28 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.87
           NAME 'sponsorPostalAddress'
           DESC 'Full postal address for a sponsor'
           SUP postalAddress
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "1 Fake St$Anytown$CA$12345$US"
#
# 2.3.88.  'sponsorPostalCode'
#
#    The 'sponsorPostalCode' attribute type allows for a postal code
#    to be assigned to a sponsor.
#
#    The attribute type 'postalCode', as defined in Section 2.23 of
#    [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.88
           NAME 'sponsorPostalCode'
           DESC 'Postal code for a sponsor'
           SUP postalCode
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "92262", "34216"
#
# 2.3.89.  'sponsorState'
#
#    The 'sponsorState' attribute type allows for the assignment of a
#    state or province name to a sponsor.
#
#    The attribute type 'st', as defined in Section 2.33 of [RFC4519],
#    is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.89
           NAME 'sponsorState'
           DESC 'State or province name for a sponsor'
           SUP st
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "California", "North Dakota"
#
# 2.3.90.  'sponsorStreet'
#
#    The 'sponsorStreet' attribute type allows for the assignment of a
#    street name and number to a sponsor.
#
#    The attribute type 'street', as defined in Section 2.34 of [RFC4519],
#    is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.90
           NAME 'sponsorStreet'
           DESC 'Street name and number for a sponsor'
           SUP street
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "1 Fake Street"
#
# 2.3.91.  'sponsorTelephone'
#
#    The 'sponsorTelephone' attribute type allows for the assignment of
#    a telephone number to a sponsor.
#
#    The attribute type 'telephoneNumber', as defined in Section 2.35 of
#    [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production can be found in Section 3.2 of [RFC4517],
#    labeled 'PrintableString'.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.91
           NAME 'sponsorTelephone'
           DESC 'Telephone number for a sponsor'
           SUP telephoneNumber
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "+11234567890"
#
# 2.3.92.  'sponsorTitle'
#
#    The 'sponsorTitle' attribute type allows for the assignment of an
#    official or professional title to a sponsor.
#
#    The attribute type 'title', as defined in Section 2.38 of [RFC4519],
#    is a super type of this attribute type.
#
#    A practical ABNF production is found in Section 3.3.6 in [RFC4517].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.92
           NAME 'sponsorTitle'
           DESC 'Title of a sponsor'
           SUP title
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "Executive Sponsor"
#
# 2.3.93.  'sponsorURI'
#
#    The 'sponsorURI' attribute type allows for the assignment of one (1)
#    or more URI values, each with an optional label, to a sponsor.
#
#    The attribute type 'labeledURI', as defined in [RFC2079], is a super
#    type of this attribute type.
#
#    A practical ABNF production for this attribute type is defined within
#    Appendix A of [RFC3986].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.93
           NAME 'sponsorURI'
           DESC 'Uniform Resource Identifier for a sponsor'
           SUP labeledURI
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "http://example.com Example", "http://example.com"
#
# 2.3.94.  'rADITProfile'
#
#    The 'rADITProfile' attribute type references entries which contain
#    optimal 'rADUAConfig' configuration settings for client consumption.
#
#    The attribute type 'distinguishedName', as defined in Section 2.7
#    of [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.94
           NAME 'rADITProfile'
           DESC 'Advertised RA DIT profile DNs served by an RA DSA'
           SUP distinguishedName
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "n=1,dc=example,dc=com" "n=1,n=4,n=1,n=6,n=3,n=1"
#
# 2.3.95.  'rARegistrationBase'
#
#    The 'rARegistrationBase' attribute type allows for the storage of one
#    (1) or more DNs that refer to entries under which 'registration'
#    entries reside.
#
#    The attribute type 'distinguishedName', as defined in Section 2.7
#    of [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.95
           NAME 'rARegistrationBase'
           DESC 'DN of a registration base within an RA DIT'
           SUP distinguishedName
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "ou=Registrations,o=rA"
#
# 2.3.96.  'rARegistrantBase'
#
#    The 'rARegistrantBase' attribute type allows for the storage of one
#    (1) or more LDAP DNs that refer to entries under which 'registrant'
#    entries reside.
#
#    The attribute type 'distinguishedName', as defined in Section 2.7
#    of [RFC4519], is a super type of this attribute type.
#
#    A practical ABNF production for this attribute type can be found in
#    Section 3 of [RFC4514].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.96
           NAME 'rARegistrantBase'
           DESC 'DN of a registrant base within an RA DIT'
           SUP distinguishedName
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "ou=Registrants,o=rA"
#
# 2.3.97.  'rADirectoryModel'
#
#    The 'rADirectoryModel' attribute type allows for the storage of a
#    numerical OID meant to declare the structural design of the RA DIT.
#
#    A practical ABNF production, labeled 'numericoid', is defined within
#    Section 1.4 of [RFC4512].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.97
           NAME 'rADirectoryModel'
           DESC 'Governing directory model OID'
           EQUALITY objectIdentifierMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.38
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "1.3.6.1.4.1.56521.101.3.3"
#
# 2.3.98.  'rAServiceMail'
#
#    The 'rAServiceMail' attribute type allows for the assignment of an
#    email address to an 'rADUAConfig' entry for the purpose of support
#    or error reporting.
#
#    The attribute type 'mail', as defined in Section 2.16 of [RFC4524],
#    is a super type of this attribute type.
#
#    A practical ABNF production can be found in Section 3.2 of [RFC4517],
#    labeled 'IA5String'.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.98
           NAME 'rAServiceMail'
           DESC 'Registration Authority contact email address'
           SUP mail
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "ra@example.com"
#
# 2.3.99.  'rAServiceURI'
#
#    The 'rAServiceURI' attribute type allows for the assignment of one
#    (1) or more URI values to an 'rADUAConfig' entry for the purpose of
#    directing users to an appropriate RA endpoint for request handling.
#
#    The attribute type 'labeledURI', as defined in [RFC2079], is a super
#    type of this attribute type.
#
#    A practical ABNF production for this attribute type is defined within
#    Appendix A of [RFC3986].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.99
           NAME 'rAServiceURI'
           DESC 'Registration Authority URI'
           SUP labeledURI
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "http://example.com/ra.html Registrations"
#
# 2.3.100.  'rATTL'
#
#    The 'rATTL' attribute type allows for the specification of a TTL
#    value, expressed in seconds.
#
#    A practical ABNF production, labeled 'number', can be found within
#    Section 1.4 of [RFC4512].
#
#    See the RADUA ID for details relating to client-driven entry
#    caching and practical value implementation.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.100
           NAME 'rATTL'
           DESC 'RA entry Time to Live, expressed in seconds'
           EQUALITY integerMatch
           ORDERING integerOrderingMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.27
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Examples: "0", "3600", "-1", "86400"
#
# 2.3.101.  'c-rATTL'
#
#    The 'c-rATTL' attribute type is the manifestation of the 'rATTL'
#    type, defined in Section 2.3.99 with Collective Attribute support.
#
#    This attribute type should only be used in directory environments
#    which actively support and require [RFC3671] capabilities.
#
#    See the RADUA ID for details relating to client-driven entry
#    caching and practical value implementation.
#
#    A practical ABNF production, labeled 'number', can be found within
#    Section 1.4 of [RFC4512].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.101
           NAME 'c-rATTL'
           DESC 'Collective RA entry Time to Live, expressed in seconds'
           EQUALITY integerMatch
           ORDERING integerOrderingMatch
           SYNTAX 1.3.6.1.4.1.1466.115.121.1.27
           COLLECTIVE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.3.102.  'registeredUUID'
#
#    The 'registeredUUID' attribute type assigns a hexadecimal UUID value
#    to a registration.
#
#    A practical ABNF production for this attribute type is defined within
#    Section 3 of [RFC4122].
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.102
           NAME 'registeredUUID'
           DESC 'X.667 registered UUID'
           EQUALITY UUIDMatch
           ORDERING UUIDOrderingMatch
           SYNTAX 1.3.6.1.1.16.1
           SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "e6bcf22c-00bf-4b3d-b11f-36ec0522aa93"
#
# 2.3.103. 'dotEncoding'
#
#    The 'dotEncoding' attributeType allows the storage of the raw ASN.1
#    encoding of the 'dotNotation' for a registration. Clients SHOULD
#    expect values of the 'dotEncoding' attribute type to manifest as
#    Base64-encoded content.
#
#    A practical ABNF production for this attribute type is as follows:
#
#      dotEncoding   = tag length bytes
#      tag           = %x06             ; ASN.1 OBJECT IDENTIFIER tag (6)
#      length        = 1*OCTET          ; number of (encoded) bytes
#      bytes         = 1*OCTET          ; encoded bytes
#
#    The 'OCTET' ABNF production is defined within RFC 4512 Section 1.4.
#
attributetype ( 1.3.6.1.4.1.56521.101.2.3.103
	   NAME 'dotEncoding'
	   DESC 'ASN.1 encoded numeric OID'
	   EQUALITY octetStringMatch
	   SYNTAX 1.3.6.1.4.1.1466.115.121.1.40
	   SINGLE-VALUE
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
#    Example: "BgEr"
#
# 2.4.  Matching Rule Uses
#
#    No LDAP Matching Rule Uses definitions are defined anywhere in this
#    ID series.
#
# 2.5.  Object Classes
#
#    The following subsections describe seventeen (17) LDAP object classes
#    defined within this ID.
#
# 2.5.1.  'registration'
#
#    The 'registration' ABSTRACT class, which is a sub type of 'top', per
#    Section 2.4.1 of [RFC4512], serves as the foundation for ALL entries
#    intended to represent OID registrations within the spirit of this ID.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.1
           NAME 'registration'
           DESC 'Abstract OID arc class'
           SUP top ABSTRACT
           MUST n
           MAY ( description $ seeAlso $ rATTL )
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.2.  'rootArc'
#
#    The 'rootArc' STRUCTURAL class is meant to define a maximum of three
#    (3) root registrations for an RA DIT, per ITU-T Rec. [X.660].
#
#    Entries that bear this class SHALL ONLY represent the following root
#    registrations:
#
#      - ITU-T ({itu-t(0)})
#      - ISO ({iso(1)})
#      - Joint-ISO-ITU-T ({joint-iso-itu-t(2)})
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.2
           NAME 'rootArc'
           DESC 'X.660, cl. A.2: root arc class'
           SUP registration STRUCTURAL
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.3.  'arc'
#
#    The 'arc' STRUCTURAL object class extends a collection of attribute
#    types for use when allocating subordinate registrations in an RA DIT.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.3
           NAME 'arc'
           DESC 'X.660, cl. A.3-A.5: subordinate arc class'
           SUP registration STRUCTURAL
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.4. 'x660Context'
#
#    The 'x660Context' AUXILIARY class extends a collection of attribute
#    types derived from Rec. ITU-T [X.660].
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.4
           NAME 'x660Context'
           DESC 'X.660 contextual class'
           SUP registration AUXILIARY
           MAY ( additionalUnicodeValue $
                 currentAuthority $
                 firstAuthority $
                 secondaryIdentifier $
                 sponsor $
                 standardizedNameForm $
                 unicodeValue )
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.5. 'x667Context'
#
#    The 'x667Context' AUXILIARY class extends a single attribute type,
#    'registeredUUID' as defined in Section 2.3.101, for use within the
#    context of an ITU-T Rec. [X.667] (UUID) registration.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.5
           NAME 'x667Context'
           DESC 'X.667 contextual class'
           SUP registration AUXILIARY
           MUST registeredUUID
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.6. 'x680Context'
#
#    The 'x680Context' AUXILIARY class extends a collection of attribute
#    types derived from ITU-T Rec. [X.680].
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.6
           NAME 'x680Context'
           DESC 'X.680 contextual class'
           SUP registration AUXILIARY
           MAY ( aSN1Notation $
                 dotNotation $
                 identifier $
                 iRI $
                 nameAndNumberForm )
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.7.  'x690Context'
#
#    The 'x690Context' AUXILIARY class extends a collection of attribute
#    types derived from ITU-T Rec. [X.690].  Additional encoding types
#    that may be added in future revisions or extensions of this ID will
#    be specified in the MAY clause of this class.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.7
           NAME 'x690Context'
           DESC 'X.690 contextual class'
           SUP registration AUXILIARY
           MAY dotEncoding
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.8.  'iTUTRegistration'
#
#    The 'iTUTRegistration' AUXILIARY class labels the registration as
#    belonging to the International Telecommunications Union (ITU-T) in
#    subordinate form.
#
#    It is NOT RECOMMENDED to assign this class to any entry which bears
#    the 'rootArc' STRUCTURAL object class, per in Section 2.5.1.  This
#    class SHALL NOT appear on entries bearing the 'iSORegistration' (per
#    Section 2.5.9) or 'jointISOITUTRegistration' (per Section 2.5.10)
#    class.
#
#    The 'x660Context' (Section 2.5.4) and 'x680Context' (Section 2.5.6)
#    are super classes of this class.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.8
           NAME 'iTUTRegistration'
           DESC 'X.660, cl. A.2: ITU-T'
           SUP ( x660Context $
                 x680Context $
                 x690Context )
           AUXILIARY
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.9.  'iSORegistration'
#
#    The 'iSORegistration' AUXILIARY class labels the OID registration as
#    belonging to the International Organization for Standardization (ISO)
#    in subordinate form.
#
#    It is NOT RECOMMENDED to assign this class to any entry which bears
#    the 'rootArc' STRUCTURAL object class, per Section 2.5.1.  This class
#    SHALL NOT appear on entries bearing the 'iTUTRegistration' (defined
#    in Section 2.5.8) or 'jointISOITUTRegistration' (per Section 2.5.10)
#    class.
#
#    The 'x660Context' (Section 2.5.4) and 'x680Context' (Section 2.5.6)
#    are super classes of this class.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.9
           NAME 'iSORegistration'
           DESC 'X.660, cl. A.2: ISO'
           SUP ( x660Context $
                 x680Context $
                 x690Context )
           AUXILIARY
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.10.  'jointISOITUTRegistration'
#
#    The 'jointISOITUTRegistration' AUXILIARY class labels a registration
#    as being jointly administered by the International Organization for
#    Standardization (ISO) and the International Telecommunications Union
#    (ITU-T) in cooperative fashion.
#
#    It is NOT RECOMMENDED to assign this class to any entry which bears
#    the 'rootArc' STRUCTURAL object class, per Section 2.5.1.  This class
#    SHALL NOT appear on entries bearing the 'iTUTRegistration' (defined
#    in Section 2.5.8) or 'iSORegistration' (per Section 2.5.9) class.
#
#    This class extends use of the 'longArc' attribute type, as defined in
#    Section 2.3.20.
#
#    The 'x660Context' (Section 2.5.4) and 'x680Context' (Section 2.5.6)
#    are super classes of this class.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.10
           NAME 'jointISOITUTRegistration'
           DESC 'X.660, cl. A.2: Joint ISO/ITU-T Administration'
           SUP ( x660Context $
                 x680Context $
                 x690Context )
           AUXILIARY
           MAY longArc
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.11.  'spatialContext'
#
#    The 'spatialContext' AUXILIARY class extends vertical and horizontal
#    associative attribute types intended to describe the logical position
#    of a registration in relation to adjacent registrations.
#
#    Use of this class is entirely OPTIONAL, but can greatly aid in the
#    creation of navigational interfaces which allow both horizontal and
#    vertical movement across entire ancestries.
#
#    See the RADIT ID for further details on this topic.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.11
           NAME 'spatialContext'
           DESC 'Logical spatial orientation and association class'
           SUP registration AUXILIARY
           MAY ( topArc $ supArc $ subArc $
                 minArc $ maxArc $
                 leftArc $ rightArc )
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.12.  'registrationSupplement'
#
#    The 'registrationSupplement' AUXILIARY class extends a variety of
#    miscellaneous and non-standard attribute types for OPTIONAL.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.12
           NAME 'registrationSupplement'
           DESC 'Supplemental registration class'
           SUP registration AUXILIARY
           MAY ( discloseTo $ isFrozen $ isLeafNode $
                 registrationClassification $
                 registrationCreated $
                 registrationInformation $
                 registrationModified $
                 registrationRange $
                 registrationStatus $
                 registrationURI )
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.13.  'firstAuthorityContext'
#
#    The 'firstAuthorityContext' AUXILIARY class extends attribute types
#    intended to store previous authority information.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.13
           NAME 'firstAuthorityContext'
           DESC 'Initial registration authority class'
           SUP top AUXILIARY
           MAY ( firstAuthorityCommonName $
                 firstAuthorityCountryCode $
                 firstAuthorityCountryName $
                 firstAuthorityCountryName $
                 firstAuthorityEmail $
                 firstAuthorityEndTimestamp $
                 firstAuthorityFax $
                 firstAuthorityLocality $
                 firstAuthorityMobile $
                 firstAuthorityOrg $
                 firstAuthorityPostalCode $
                 firstAuthorityStartTimestamp $
                 firstAuthorityState $
                 firstAuthorityStreet $
                 firstAuthorityTelephone $
                 firstAuthorityTitle $
                 firstAuthorityURI )
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.14.  'currentAuthorityContext'
#
#    The 'currentAuthorityContext' AUXILIARY class extends attribute types
#    intended to store current authority information.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.14
           NAME 'currentAuthorityContext'
           DESC 'Current registration authority class'
           SUP top AUXILIARY
           MAY ( currentAuthorityCommonName $
                 currentAuthorityCountryCode $
                 currentAuthorityCountryName $
                 currentAuthorityCountryName $
                 currentAuthorityEmail $
                 currentAuthorityFax $
                 currentAuthorityLocality $
                 currentAuthorityMobile $
                 currentAuthorityOrg $
                 currentAuthorityPostalCode $
                 currentAuthorityStartTimestamp $
                 currentAuthorityState $
                 currentAuthorityStreet $
                 currentAuthorityTelephone $
                 currentAuthorityTitle $
                 currentAuthorityURI )
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.15.  'sponsorContext'
#
#    The 'currentAuthorityContext' AUXILIARY class extends attribute types
#    intended to store sponsoring authority information.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.15
           NAME 'sponsorContext'
           DESC 'Registration sponsoring authority class'
           SUP top AUXILIARY
           MAY ( sponsorCommonName $
                 sponsorCountryCode $
                 sponsorCountryName $
                 sponsorCountryName $
                 sponsorEmail $
                 sponsorEndTimestamp $
                 sponsorFax $
                 sponsorLocality $
                 sponsorMobile $
                 sponsorOrg $
                 sponsorPostalCode $
                 sponsorStartTimestamp $
                 sponsorState $
                 sponsorStreet $
                 sponsorTelephone $
                 sponsorTitle $
                 sponsorURI )
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.16.  'registrant'
#
#    The 'registrant' object class allows for current, previous (first)
#    and/or sponsorship data to be stored within an entry.
#
#    This object class is STRUCTURAL, and cannot be combined with the
#    'registration' object class within a single entry.
#
#    Use of the 'registrant' object class within an RA DIT implies use
#    of so-called dedicated authority entries, as described within the
#    RADIT ID.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.16
           NAME 'registrant'
           DESC 'Generalized auxiliary class for registrant data'
           SUP top STRUCTURAL
           MUST registrantID
           MAY ( description $ seeAlso $ rATTL )
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.5.17.  'rADUAConfig'
#
#    The 'rADUAConfig' object class allows for the storage of so-called
#    auto-discovered attribute types meant to guide the RA DUA in its
#    attempt to access registration and/or registrant information stored
#    within the RA DIT served by the RA DSA.
#
objectclass ( 1.3.6.1.4.1.56521.101.2.5.17
           NAME 'rADUAConfig'
           DESC 'RA DUA configuration advertisement class'
           SUP top AUXILIARY
           MAY ( description $
                 rADITProfile $
                 rADirectoryModel $
                 rARegistrationBase $
                 rARegistrantBase $
                 rAServiceMail $
                 rAServiceURI $
                 rATTL )
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.6.  DIT Content Rules
#
#    No DIT Content Rules are officially defined anywhere in this ID
#    series, however the following definitions have been added for the
#    sake of testing and examples within the go-schemax package.
#
# 2.6.1.  'rootArcContent'
#
ditcontentrule ( 1.3.6.1.4.1.56521.101.2.5.3
	NAME 'arcContent'
	DESC 'arc entry content rule'
	AUX ( x660Context
            $ x667Context
            $ x680Context
            $ x690Context )
	MUST ( n $
	       iRI $
	       identifier $
	       unicodeValue $
	       aSN1Notation )
	MAY ( secondaryIdentifier $
	      additionalUnicodeValue $
	      standardizedNameForm $
	      nameAndNumberForm )
	NOT dotNotation
        X-ORIGIN 'draft-coretta-oiddir-schema; unofficial supplement'
        X-WARNING 'UNOFFICIAL' )
#
#
# 2.7.  Name Forms
#
#    This section defines three (3) Name Forms for use in the composition
#    of any DIT Structure Rules required, if supported by the X.500/LDAP
#    product(s) in question.  See Section 4.1.7.2 of [RFC4512] for
#    details.
#
# 2.7.1.  'nRootArcForm'
#
#    The 'nRootArcForm' name form is devised to control the RDN of the
#    entries bearing the 'rootArc' STRUCTURAL object class.  Within the
#    RECOMMENDED three dimensional model of this ID series, use the 'n'
#    (Number Form) attribute type is the preferred RDN component.
#
nameform ( 1.3.6.1.4.1.56521.101.2.7.1
           NAME 'nRootArcForm'
           DESC 'root arc name form for a number form RDN'
           OC rootArc
           MUST n
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.7.2.  'nArcForm'
#
#    The 'nArcForm' name form is devised to control the RDN of the entries
#    bearing the 'arc' STRUCTURAL object class.  Within the RECOMMENDED
#    three dimensional model of this ID series, the 'n' (Number Form)
#    attribute type is the preferred RDN component.
#
nameform ( 1.3.6.1.4.1.56521.101.2.7.2
           NAME 'nArcForm'
           DESC 'arc name form for a number form RDN'
           OC arc
           MUST n
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.7.3.  'dotNotationArcForm'
#
#    The 'dotNotationArcForm' name form is devised to control the RDN
#    of the bearings bearing the 'arc' STRUCTURAL object class.  Within
#    the STRONGLY DISCOURAGED two dimensional model of this ID series,
#    the 'dotNotation' (numeric OID) attribute type is the preferred RDN
#    component.
#
nameform ( 1.3.6.1.4.1.56521.101.2.7.3
           NAME 'dotNotationArcForm'
           DESC 'arc name form for a numeric OID RDN'
           OC arc
           MUST dotNotation
	   X-ORIGIN 'draft-coretta-oiddir-schema' )
#
# 2.8.  DIT Structure Rules
#
#    No DIT structure rules are officially defined anywhere in this ID
#    series, however the following dITStructureRule definitions have
#    been added for the sake of testing and examples within the go-schemax
#    package.
#
# 2.8.1.  'rootArcStructure'
#
ditstructurerule ( 20
	NAME 'rootArcStructure'
	DESC 'structure rule for root arc entries in any dimension; FOR DEMONSTRATION USE ONLY'
	FORM nRootArcForm
	X-ORIGIN 'draft-coretta-oiddir-schema; unofficial supplement' )
#
# 2.8.2.  'arcStructure'
#
ditstructurerule ( 21
	NAME 'arcStructure'
	DESC 'structure rule for three dimensional arc entries; FOR DEMONSTRATION USE ONLY'
	FORM nArcForm
	SUP 0
	X-ORIGIN 'draft-coretta-oiddir-schema; unofficial supplement' )
#
# 2.8.3  'dotNotArcStructure'
#
ditstructurerule ( 22
	NAME 'dotNotArcStructure'
	DESC 'structure rule for two dimensional arc entries; FOR DEMONSTRATION USE ONLY'
	FORM dotNotationArcForm
	SUP 0
	X-ORIGIN 'draft-coretta-oiddir-schema; unofficial supplement' )
#
# END OF DEFINITIONS
`)

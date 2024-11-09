package schemax

import (
	"github.com/JesseCoretta/go-shifty"
	"github.com/JesseCoretta/go-stackage"
)

const (
	UserApplicationsUsage     uint = iota // RFC 4512 § 4.1.2
	DirectoryOperationUsage               // RFC 4512 § 4.1.2
	DistributedOperationUsage             // RFC 4512 § 4.1.2
	DSAOperationUsage                     // RFC 4512 § 4.1.2
)

const (
	StructuralKind uint = iota // RFC 4512 § 2.4.2, 4.1.1
	AbstractKind               // RFC 4512 § 2.4.1, 4.1.1
	AuxiliaryKind              // RFC 4512 § 2.4.3, 4.1.1
)

/*
Options wraps an instance of [shifty.BitValue] allowing clean and simple
bit shifting/unshifting to effect changes to a [Schema]'s behavior.

Instances of this type are accessed and managed via the [Schema.Options]
method.

Users should NOT attempt to instantiate instances of this type manually
for any reason.
*/
type Options shifty.BitValue

/*
Macros contains user-defined macros for parsing OID-referencing macros
used within definitions, such as those found within RFC 2307.

Instances of this type are accessed and managed via the [Schema.Macros]
method.

Users should NOT attempt to instantiate instances of this type manually
for any reason.
*/
type Macros struct {
	macros
}

type macros map[string]string

/*
Option represents a single configuration parameter within an instance of
[Options].
*/
type Option uint16

const (
	// HangingIndents will cause all eligible string
	// processing to include hanging intends for a
	// given definition, as opposed to all definitions
	// occupying a single line each.
	HangingIndents Option = 1 << iota

	// SortExtensions will cause all ANTLR-based parsing
	// operations to sort all extensions alphabetically
	// according to the respective XString field value.
	//
	// Note this does not influence definition Extensions
	// that were manually crafted by the user.
	SortExtensions

	// SortLists will cause all ANTLR-based parsing operations
	// to sort all extensions alphabetically according to the
	// principal identifying string value of a slice member.
	// This will influence stacks that represent multi-valued
	// clauses like MAY, MUST, APPLIES and others.
	//
	// Note this does not influence slice members that were
	// entered manually by the user.
	SortLists

	// AllowOverride declares that any instance of Definition
	// may be replaced by way of its Replace method. This is
	// reflected in all stacks in which the Definition resides.
	AllowOverride

	// AllowReindexedStructureRules allows the push process of
	// dITStructureRule instances into a Schema to reindex, or
	// renumber, the new rule if it conflicts with an existing
	// rule. This option exists because dITStructureRules are
	// especially prone to uniqueness violations due to their
	// numbering system lacking the dynamics normally enjoyed
	// through numeric OIDs.
	//
	// Note that it may be necessary to re-apply any custom
	// stringers of a static nature that are assigned to any
	// related dITStructureRules following a positive reindexing.
	// This may include subordinate rules.
	AllowReindexedStructureRules

	// As-of-yet unused bit settings
	//_                    //    16
	//_                    //    32
	//_                    //    64
	//_                    //   128
	//_                    //   256
	//_                    //   512
	//_                    //  1024
	//_                    //  2048
	//_                    //  4096
	//_                    //  8192
	//_                    // 16384
	//_                    // 32768
)

/*
SyntaxQualifier is an optional closure function or method signature
which may be honored by the end user for value verification controls
within [LDAPSyntax] instances.
*/
type SyntaxQualifier func(any) error

/*
AssertionMatcher is an optional closure function or method signature which
may be honored by the end user for assertion match controls with respect
to a [MatchingRule] instance. This signature is used in up to three (3)
distinct scenarios:

  - Equality matching  (e.g.: jesse=Jesse / caseIgnoreMatch)
  - Substring matching (e.g.: jess*=Jesse / caseIgnoreSubstringMatch)
  - Ordering matching  (e.g.: 20210814135117Z>=20210428081809Z / generalizedTimeOrderingMatch)

The main use case for instances of this type is to define the basis of
comparing or evaluating two (2) values with respect to the [MatchingRule]
indicated.
*/
type AssertionMatcher func(any, any) error

/*
ValueQualifier is an optional closure function of method signature which
may be honored by the end user for enhanced value interrogation of any
given value with respect to the [AttributeType] instance to which it is
assigned.

This allows the implementation (or simulation) of common directory server
features -- such as regular expression processing -- that allow for value
constraints beyond, or instead of, those made possible through use of an
[LDAPSyntax] alone. This is particularly useful within organizations or
other bodies in which only specific, specialized values are to be allowed.
*/
type ValueQualifier func(any) error

/*
Stringer is an optional function or method signature which allows user
controlled string representation per schema definition.

The package-default stringer is static in nature, having relied upon
[text/template] to generate a one-time string representation of the
receiver instance for performance reasons. The string representation
value is essentially embedded as a return value, used for every call
of String.

The drawback of this approach is that subsequent changes to one or
more schema definitions instances after a string was generated will
not be visible until the default stringer is reapplied by way of the
SetStringer method being executed in niladic fashion. This results
in a fresh call to the [text/template] system, producing the updated
value.

The SetStringer method can be executed not only on individual definitions,
but on entire stacks as well. While convenient, the operation may be rather
costly on particularly large stacks.

To mitigate this drawback entirely, users may author custom stringers of
their own design that are dynamic -- in that they compose a string value
for every call of String -- as opposed to static.
*/
type Stringer func() string

/*
oIDList implements oidlist per § 4.1 of RFC 4512.  Instances
of this type need not be handled by users directly.

	oidlist = oid *( WSP DOLLAR WSP oid )
*/
type oIDList stackage.Stack

/*
RuleIDList implements ruleidlist per [§ 4.1.7.1 of RFC 4512].
Instances of this type need not be handled by users directly.

	ruleidlist = ruleid *( SP ruleid )

[§ 4.1.7.1 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.7.1
*/
type RuleIDList stackage.Stack

/*
QuotedDescriptorList implements qdescrlist per [§ 4.1 of RFC 4512].
Instances of this type need not be handled by users directly.

	qdescrlist = [ qdescr *( SP qdescr ) ]

[§ 4.1 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1
*/
type QuotedDescriptorList stackage.Stack

/*
QuotedStringList implements qdstringlist per [§ 4.1 of RFC 4512].
Instances of this type need not be handled by users directly.

	qdstringlist = [ qdstring *( SP qdstring ) ]

[§ 4.1 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1
*/
type QuotedStringList stackage.Stack

/*
collection implements a common [stackage.Stack] type alias and is
not an RFC 4512 construct.  Instances of this type need not be
handled by users directly.
*/
type collection stackage.Stack

/*
Extensions implements extensions as defined in [§ 4.1 of RFC 4512]:

	extensions = *( SP xstring SP qdstrings )

[§ 4.1 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1
*/
type Extensions stackage.Stack

/*
Extension is the singular (slice) form of [Extensions], and contains the following:

  - One (1) instance of string (XString), declaring the effective "X-" name, AND ...
  - One (1) [QuotedStringList] stack instance, containing one (1) or more "qdstringlist" values

The ABNF production for "xstring", per [§ 4.1 of RFC 4512], is as follows:

	xstring	     = "X" HYPHEN 1*( ALPHA / HYPHEN / USCORE )

The ABNF production for "qdstringlist", per [§ 4.1 of RFC 4512], is as follows:

	qdstringlist = [ qdstring *( SP qdstring ) ]
	qdstring     = SQUOTE dstring SQUOTE
	dstring      = 1*( QS / QQ / QUTF8 )   ; escaped UTF-8 string

[§ 4.1 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1
*/
type Extension struct {
	*extension
}

type extension struct {
	XString string
	Values  QuotedStringList

	hindent  bool
	sortExts bool
	stringer Stringer
}

/*
Schema is a practical implementation of a 'subschemaSubentry'
in that individual definitions are accessible and collectively
define the elements available for use in populating a directory.

See the following methods to access any desired [Definition]
qualifier types:

  - [Schema.LDAPSyntaxes]
  - [Schema.MatchingRules]
  - [Schema.AttributeTypes]
  - [Schema.MatchingRuleUses]
  - [Schema.ObjectClasses]
  - [Schema.DITContentRules]
  - [Schema.NameForms]
  - [Schema.DITStructureRules]
*/
type Schema collection

type (
	ObjectClasses     collection // RFC 4512 § 4.2.1
	AttributeTypes    collection // RFC 4512 § 4.2.2
	MatchingRules     collection // RFC 4512 § 4.2.3
	MatchingRuleUses  collection // RFC 4512 § 4.2.4
	LDAPSyntaxes      collection // RFC 4512 § 4.2.5
	DITContentRules   collection // RFC 4512 § 4.2.6
	DITStructureRules collection // RFC 4512 § 4.2.7
	NameForms         collection // RFC 4512 § 4.2.8
)

/*
DefinitionMap implements a convenient map-based [Definition] type.  Use
of this type is normally indicated in external processing scenarios,
such as templating.

Note that, due to the underlying map instance from which this type
extends, ordering of clauses (e.g.: NAME vs. DESC) cannot be guaranteed.
*/
type DefinitionMap map[string][]string

/*
DefinitionMaps implements slices of [DefinitionMap] instances, collectively
representing an entire type-specific stack (e.g.: [AttributeTypes]).
*/
type DefinitionMaps []DefinitionMap

/*
AttributeType implements [§ 4.1.2 of RFC 4512].

	AttributeTypeDescription = LPAREN WSP
	    numericoid                    ; object identifier
	    [ SP "NAME" SP qdescrs ]      ; short names (descriptors)
	    [ SP "DESC" SP qdstring ]     ; description
	    [ SP "OBSOLETE" ]             ; not active
	    [ SP "SUP" SP oid ]           ; supertype
	    [ SP "EQUALITY" SP oid ]      ; equality matching rule
	    [ SP "ORDERING" SP oid ]      ; ordering matching rule
	    [ SP "SUBSTR" SP oid ]        ; substrings matching rule
	    [ SP "SYNTAX" SP noidlen ]    ; value syntax
	    [ SP "SINGLE-VALUE" ]         ; single-value
	    [ SP "COLLECTIVE" ]           ; collective
	    [ SP "NO-USER-MODIFICATION" ] ; not user modifiable
	    [ SP "USAGE" SP usage ]       ; usage
	    extensions WSP RPAREN         ; extensions

	usage = "userApplications"     /  ; user
	        "directoryOperation"   /  ; directory operational
	        "distributedOperation" /  ; DSA-shared operational
	        "dSAOperation"            ; DSA-specific operational

From clause 13.4.8 of [ITU-T Rec. X.501]:

	ATTRIBUTE ::= CLASS {
		&derivation ATTRIBUTE OPTIONAL,
		&Type OPTIONAL, -- either &Type or &derivation required
		&equality-match MATCHING-RULE OPTIONAL,
		&ordering-match MATCHING-RULE OPTIONAL,
		&substrings-match MATCHING-RULE OPTIONAL,
		&single-valued BOOLEAN DEFAULT FALSE,
		&collective BOOLEAN DEFAULT FALSE,
		&dummy BOOLEAN DEFAULT FALSE,
		-- operational extensions
		&no-user-modification BOOLEAN DEFAULT FALSE,
		&usage AttributeUsage DEFAULT userApplications,
		&ldapSyntax SYNTAX-NAME.&id OPTIONAL,
		&ldapName SEQUENCE SIZE(1..MAX) OF UTF8String OPTIONAL,
		&ldapDesc UTF8String OPTIONAL,
		&obsolete BOOLEAN DEFAULT FALSE,
		&id OBJECT IDENTIFIER UNIQUE }

	WITH SYNTAX {
		[SUBTYPE OF &derivation]
		[WITH SYNTAX &Type]
		[EQUALITY MATCHING RULE &equality-match]
		[ORDERING MATCHING RULE &ordering-match]
		[SUBSTRINGS MATCHING RULE &substrings-match]
		[SINGLE VALUE &single-valued]
		[COLLECTIVE &collective]
		[DUMMY &dummy]
		[NO USER MODIFICATION &no-user-modification]
		[USAGE &usage]
		[LDAP-SYNTAX &ldapSyntax]
		[LDAP-NAME &ldapName]
		[LDAP-DESC &ldapDesc]
		[OBSOLETE &obsolete]
		ID &id }

		AttributeUsage ::= ENUMERATED {
			userApplications (0),
			directoryOperation (1),
			distributedOperation (2),
			dSAOperation (3),
			... }

[§ 4.1.2 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.2
[ITU-T Rec. X.501]: https://www.itu.int/rec/T-REC-X.501
*/
type AttributeType struct {
	*attributeType
}

type attributeType struct {
	OID        string
	Macro      []string
	Desc       string
	Name       QuotedDescriptorList
	Obsolete   bool
	Single     bool
	Collective bool
	NoUserMod  bool
	SuperType  AttributeType
	Equality   MatchingRule
	Ordering   MatchingRule
	Substring  MatchingRule
	Syntax     LDAPSyntax
	MUB        uint
	Usage      uint
	Extensions Extensions

	schema   Schema
	stringer Stringer
	valQual  ValueQualifier
	data     any
}

/*
DITContentRule implements [§ 4.1.6 of RFC 4512].

	DITContentRuleDescription = LPAREN WSP
	    numericoid                 ; object identifier
	    [ SP "NAME" SP qdescrs ]   ; short names (descriptors)
	    [ SP "DESC" SP qdstring ]  ; description
	    [ SP "OBSOLETE" ]          ; not active
	    [ SP "AUX" SP oids ]       ; auxiliary object classes
	    [ SP "MUST" SP oids ]      ; attribute types
	    [ SP "MAY" SP oids ]       ; attribute types
	    [ SP "NOT" SP oids ]       ; attribute types
	    extensions WSP RPAREN      ; extensions

From clause 13.8.2 of [ITU-T Rec. X.501]:

	CONTENT-RULE ::= CLASS {
		&structuralClass OBJECT-CLASS.&id UNIQUE,
		&Auxiliaries OBJECT-CLASS OPTIONAL,
		&Mandatory ATTRIBUTE OPTIONAL,
		&Optional ATTRIBUTE OPTIONAL,
		&Precluded ATTRIBUTE OPTIONAL }

	WITH SYNTAX {
		STRUCTURAL OBJECT-CLASS &structuralClass
		[AUXILIARY OBJECT-CLASSES &Auxiliaries]
		[MUST CONTAIN &Mandatory]
		[MAY CONTAIN &Optional]
		[MUST-NOT CONTAIN &Precluded] }

[§ 4.1.6 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.6
[ITU-T Rec. X.501]: https://www.itu.int/rec/T-REC-X.501
*/
type DITContentRule struct {
	*dITContentRule
}

type dITContentRule struct {
	Desc       string
	Name       QuotedDescriptorList
	Obsolete   bool
	OID        ObjectClass
	Macro      []string
	Aux        ObjectClasses
	Must       AttributeTypes
	May        AttributeTypes
	Not        AttributeTypes
	Extensions Extensions

	schema   Schema
	stringer Stringer
	data     any
}

/*
DITStructureRule implements [§ 4.1.7.1 of RFC 4512].

	DITStructureRuleDescription = LPAREN WSP
	    ruleid                     ; rule identifier
	    [ SP "NAME" SP qdescrs ]   ; short names (descriptors)
	    [ SP "DESC" SP qdstring ]  ; description
	    [ SP "OBSOLETE" ]          ; not active
	    SP "FORM" SP oid           ; NameForm
	    [ SP "SUP" SP ruleids ]    ; superior rules <NOTE: SEE ERRATA # 7896>
	    extensions WSP RPAREN      ; extensions

	ruleids = ruleid / ( LPAREN WSP ruleidlist WSP RPAREN )
	ruleidlist = ruleid *( SP ruleid )
	ruleid = number

From clause 13.7.6 of [ITU-T Rec. X.501]:

	DITStructureRule ::= SEQUENCE {
		ruleIdentifier RuleIdentifier,
		-- shall be unique within the scope of the subschema
		nameForm NAME-FORM.&id,
		superiorStructureRules SET SIZE (1..MAX) OF RuleIdentifier OPTIONAL,
		... }

	RuleIdentifier ::= INTEGER

	STRUCTURE-RULE ::= CLASS {
		&nameForm NAME-FORM,
		&SuperiorStructureRules STRUCTURE-RULE.&id OPTIONAL,
		&id RuleIdentifier }

	WITH SYNTAX {
		NAME FORM &nameForm
		[SUPERIOR RULES &SuperiorStructureRules]
		ID &id }

[§ 4.1.7.1 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.7.1
[ITU-T Rec. X.501]: https://www.itu.int/rec/T-REC-X.501
*/
type DITStructureRule struct {
	*dITStructureRule
}

type dITStructureRule struct {
	ID         uint
	Desc       string
	Name       QuotedDescriptorList
	Obsolete   bool
	Form       NameForm
	SuperRules DITStructureRules
	Extensions Extensions

	schema   Schema
	stringer Stringer
	data     any
}

type extensions map[string]QuotedStringList

/*
LDAPSyntax implements [§ 4.1.5 of RFC 4512].

	SyntaxDescription = LPAREN WSP
	    numericoid                 ; object identifier
	    [ SP "DESC" SP qdstring ]  ; description
	    extensions WSP RPAREN      ; extensions

From clause 13.12 of [ITU-T Rec. X.501]:

	SYNTAX-NAME ::= CLASS {
		&desc UTF8String,
		&Type,
		&id OBJECT IDENTIFIER UNIQUE }

	WITH SYNTAX {
		DESC &desc
		DIRECTORY SYNTAX &Type
		ID &id }

[§ 4.1.5 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.5
[ITU-T Rec. X.501]: https://www.itu.int/rec/T-REC-X.501
*/
type LDAPSyntax struct {
	*lDAPSyntax
}

type lDAPSyntax struct {
	OID        string
	Macro      []string
	Desc       string
	Extensions Extensions

	schema   Schema
	stringer Stringer
	synQual  SyntaxQualifier
	data     any
}

/*
MatchingRule implements [§ 4.1.3 of RFC 4512].

	MatchingRuleDescription = LPAREN WSP
	    numericoid                 ; object identifier
	    [ SP "NAME" SP qdescrs ]   ; short names (descriptors)
	    [ SP "DESC" SP qdstring ]  ; description
	    [ SP "OBSOLETE" ]          ; not active
	    SP "SYNTAX" SP numericoid  ; assertion syntax
	    extensions WSP RPAREN      ; extensions

From clause 13.5.2 of [ITU-T Rec. X.501]:

	MATCHING-RULE ::= CLASS {
		&ParentMatchingRules MATCHING-RULE OPTIONAL,
		&AssertionType OPTIONAL,
		&uniqueMatchIndicator ATTRIBUTE OPTIONAL,
		&ldapSyntax SYNTAX-NAME.&id OPTIONAL,
		&ldapName SEQUENCE SIZE(1..MAX) OF UTF8String OPTIONAL,
		&ldapDesc UTF8String OPTIONAL,
		&id OBJECT IDENTIFIER UNIQUE }

	WITH SYNTAX {
		[PARENT &ParentMatchingRules]
		[SYNTAX &AssertionType]
		[UNIQUE-MATCH-INDICATOR &uniqueMatchIndicator]
		[LDAP-SYNTAX &ldapSyntax]
		[LDAP-NAME &ldapName]
		[LDAP-DESC &ldapDesc]
		ID &id }

[§ 4.1.3 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.3
[ITU-T Rec. X.501]: https://www.itu.int/rec/T-REC-X.501
*/
type MatchingRule struct {
	*matchingRule
}

type matchingRule struct {
	OID        string
	Macro      []string
	Name       QuotedDescriptorList
	Desc       string
	Obsolete   bool
	Syntax     LDAPSyntax
	Extensions Extensions

	schema   Schema
	stringer Stringer
	assMatch AssertionMatcher
	data     any
}

/*
MatchingRuleUse implements [§ 4.1.4 of RFC 4512].

	MatchingRuleUseDescription = LPAREN WSP
	    numericoid                 ; object identifier
	    [ SP "NAME" SP qdescrs ]   ; short names (descriptors)
	    [ SP "DESC" SP qdstring ]  ; description
	    [ SP "OBSOLETE" ]          ; not active
	    SP "APPLIES" SP oids       ; attribute types
	    extensions WSP RPAREN      ; extensions

From clause 13.6.2 of [ITU-T Rec. X.501]:

	MAPPING-BASED-MATCHING
	{SelectedBy, BOOLEAN:combinable, MappingResult, OBJECT IDENTIFIER:matchingRule} ::= CLASS {
		&selectBy SelectedBy OPTIONAL,
		&ApplicableTo ATTRIBUTE,
		&subtypesIncluded BOOLEAN DEFAULT TRUE,
		&combinable BOOLEAN(combinable),
		&mappingResults MappingResult OPTIONAL,
		&userControl BOOLEAN DEFAULT FALSE,
		&exclusive BOOLEAN DEFAULT TRUE,
		&matching-rule MATCHING-RULE.&id(matchingRule),
		&id OBJECT IDENTIFIER UNIQUE }

	WITH SYNTAX {
		[SELECT BY &selectBy]
		APPLICABLE TO &ApplicableTo
		[SUBTYPES INCLUDED &subtypesIncluded]
		COMBINABLE &combinable
		[MAPPING RESULTS &mappingResults]
		[USER CONTROL &userControl]
		[EXCLUSIVE &exclusive]
		MATCHING RULE &matching-rule
		ID &id }

[§ 4.1.4 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.4
[ITU-T Rec. X.501]: https://www.itu.int/rec/T-REC-X.501
*/
type MatchingRuleUse struct {
	*matchingRuleUse
}

type matchingRuleUse struct {
	OID        MatchingRule
	Name       QuotedDescriptorList
	Desc       string
	Obsolete   bool
	Applies    AttributeTypes
	Extensions Extensions

	schema   Schema
	stringer Stringer
	data     any
}

/*
NameForm implements [§ 4.1.7.2 of RFC 4512].

	NameFormDescription = LPAREN WSP
	    numericoid                 ; object identifier
	    [ SP "NAME" SP qdescrs ]   ; short names (descriptors)
	    [ SP "DESC" SP qdstring ]  ; description
	    [ SP "OBSOLETE" ]          ; not active
	    SP "OC" SP oid             ; structural object class
	    SP "MUST" SP oids          ; attribute types
	    [ SP "MAY" SP oids ]       ; attribute types
	    extensions WSP RPAREN      ; extensions

From clause 13.7.3 of [ITU-T Rec. X.501]:

	NAME-FORM ::= CLASS {
		&namedObjectClass OBJECT-CLASS,
		&MandatoryAttributes ATTRIBUTE,
		&OptionalAttributes ATTRIBUTE OPTIONAL,
		&ldapName SEQUENCE SIZE(1..MAX) OF UTF8String OPTIONAL,
		&ldapDesc UTF8String OPTIONAL,
		&id OBJECT IDENTIFIER UNIQUE }

	WITH SYNTAX {
		NAMES &namedObjectClass
		WITH ATTRIBUTES &MandatoryAttributes
		[AND OPTIONALLY &OptionalAttributes]
		[LDAP-NAME &ldapName]
		[LDAP-DESC &ldapDesc]
		ID &id }

[§ 4.1.7.2 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.7.2
[ITU-T Rec. X.501]: https://www.itu.int/rec/T-REC-X.501
*/
type NameForm struct {
	*nameForm
}

type nameForm struct {
	OID        string
	Macro      []string
	Desc       string
	Name       QuotedDescriptorList
	Obsolete   bool
	Structural ObjectClass
	Must       AttributeTypes
	May        AttributeTypes
	Extensions Extensions

	schema   Schema
	stringer Stringer
	data     any
}

/*
ObjectClass implements [§ 4.1.1 of RFC 4512].

	ObjectClassDescription = LPAREN WSP
	    numericoid                 ; object identifier
	    [ SP "NAME" SP qdescrs ]   ; short names (descriptors)
	    [ SP "DESC" SP qdstring ]  ; description
	    [ SP "OBSOLETE" ]          ; not active
	    [ SP "SUP" SP oids ]       ; superior object classes
	    [ SP kind ]                ; kind of class
	    [ SP "MUST" SP oids ]      ; attribute types
	    [ SP "MAY" SP oids ]       ; attribute types
	    extensions WSP RPAREN

	kind = "ABSTRACT" / "STRUCTURAL" / "AUXILIARY"

From clause 13.4.8 of [ITU-T Rec. X.501]:

	OBJECT-CLASS ::= CLASS {
		&Superclasses OBJECT-CLASS OPTIONAL,
		&kind ObjectClassKind DEFAULT structural,
		&MandatoryAttributes ATTRIBUTE OPTIONAL,
		&OptionalAttributes ATTRIBUTE OPTIONAL,
		&ldapName SEQUENCE SIZE(1..MAX) OF UTF8String OPTIONAL,
		&ldapDesc UTF8String OPTIONAL,
		&id OBJECT IDENTIFIER UNIQUE }

	WITH SYNTAX {
		[SUBCLASS OF &Superclasses]
		[KIND &kind]
		[MUST CONTAIN &MandatoryAttributes]
		[MAY CONTAIN &OptionalAttributes]
		[LDAP-NAME &ldapName]
		[LDAP-DESC &ldapDesc]
		ID &id }

[§ 4.1.1 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.1
[ITU-T Rec. X.501]: https://www.itu.int/rec/T-REC-X.501
*/
type ObjectClass struct {
	*objectClass
}

type objectClass struct {
	OID          string
	Macro        []string
	Desc         string
	Name         QuotedDescriptorList
	Obsolete     bool
	SuperClasses ObjectClasses
	Kind         uint
	Must         AttributeTypes
	May          AttributeTypes
	Extensions   Extensions

	schema   Schema
	stringer Stringer
	data     any
}

/*
Counters is a simple struct type defined to store the current number
of definition instances within an instance of [Schema].

Instances of this type are NOT thread-safe.
*/
type Counters struct {
	LS int
	MR int
	AT int
	MU int
	OC int
	DC int
	NF int
	DS int
}

/*
Inventory is a type alias of map[string][]string, and is used to
provide a simple manifest of all members of a collection type,
such as [LDAPSyntaxes].  This can be useful during activities such
as templating.

Unlike the [DefinitionMap] type, this type is only used to manifest
the most basic details of a collection of definitions, namely the
numerical and textual identifiers present.

Keys represent the numerical identifier for a definition, whether a
numeric OID or integer rule ID.  In the case of [DITStructureRule]
definitions, the rule ID -- an unsigned integer -- is used. In all
other cases, a numeric OID is used.

Values represent the NAME or DESC by which the definition is known.

Note that DESC is only used in the case of [LDAPSyntax] instances, and
NAME is used for all definition types *except* [LDAPSyntax].
*/
type Inventory map[string][]string

/*
Definition is an interface type used to allow basic interaction with
any of the following instance types:

  - [LDAPSyntax]
  - [MatchingRule]
  - [AttributeType]
  - [MatchingRuleUse]
  - [ObjectClass]
  - [DITContentRule]
  - [NameForm]
  - [DITStructureRule]

Not all methods extended through instances of these types are made
available through this interface.
*/
type Definition interface {
	// NumericOID returns the string representation of the ASN.1
	// numeric OID value of the instance.  If the underlying type
	// is a DITStructureRule, no value is returned as this type
	// of definition does not bear an OID.
	NumericOID() string

	// Data returns the underlying user-assigned value present
	// within the receiver instance.
	Data() any

	// Parse returns an error following an attempt read the string
	// input value into the receiver instance.
	//
	// Please note:
	//
	//   - The receiver MUST be initialized and associated with an
	//     appropriate Schema instance
	//   - MatchingRuleUse instances are NOT eligible for this method
	Parse(string) error

	// Marshal returns an error following an attempt to read an
	// instance of DefinitionMap or map[string]any into a qualifying
	// instance.
	Marshal(any) error

	// Name returns the first string NAME value present within
	// the underlying Name stack instance.  A zero
	// string is returned if no names were set, or if the given
	// type instance is LDAPSyntax, which does not bear a name.
	Name() string

	// Names returns the underlying instance of QuotedDescriptorList.
	// If executed upon an instance of LDAPSyntax, an empty instance
	// is returned, as LDAPSyntaxes do not bear names.
	Names() QuotedDescriptorList

	// IsZero returns a Boolean value indicative of nilness
	// with respect to the embedded type instance.
	IsZero() bool

	// Compliant returns a Boolean value indicative of compliance
	// with respect to relevant RFCs, such as RFC 4512.
	Compliant() bool

	// String returns the complete string representation of the
	// underlying definition type per § 4.1.x of RFC 4512.
	String() string

	// Description returns the DESC clause of the underlying
	// definition type, else a zero string if undefined.
	Description() string

	// Schema returns the associated Schema instance with which
	// the receiver instance is associated.
	Schema() Schema

	// IsIdentifiedAs returns a Boolean value indicative of
	// whether the input string represents an identifying
	// value for the underlying definition.
	//
	// Acceptable input values vary based on the underlying
	// type.  See the individual method notes for details.
	//
	// Case-folding of input values is not significant in
	// the matching process, regardless of underlying type.
	IsIdentifiedAs(string) bool

	// Map returns a DefinitionMap instance based upon the
	// the contents and state of the receiver instance.
	Map() DefinitionMap

	// Obsolete returns a Boolean value indicative of the
	// condition of definition obsolescence. Executing this
	// method upon an LDAPSyntax receiver will always return
	// false, as the condition of obsolescence does not
	// apply to this definition type.
	Obsolete() bool

	// Type returns the string literal name for the receiver
	// instance. For example, if the receiver is an LDAPSyntax
	// instance, the return value is "ldapSyntax".
	Type() string

	// Extensions returns the underlying Extension instance value.
	// An empty Extensions instance is returned if no extensions
	// were set.
	Extensions() Extensions

	// setOID empties the definition Macro slice following completion
	// of the resolution process of a numeric OID during the parsing
	// process.
	setOID(string)

	// macro allows private access to the underlying Macro field
	// present within (nearly) all Definition qualifier types. This
	// is used for a low-cyclo means of resolving macros to actual
	// numeric OIDs during the parsing phase.
	macro() []string
}

/*
Definitions is an interface type used to allow basic interaction with
any of the following stack types:

  - [LDAPSyntaxes]
  - [MatchingRules]
  - [AttributeTypes]
  - [MatchingRuleUses]
  - [ObjectClasses]
  - [DITContentRules]
  - [NameForms]
  - [DITStructureRules]

Not all methods extended through instances of these types are made
available through this interface.
*/
type Definitions interface {
	// Len returns the integer length of the receiver instance,
	// indicating the number of definitions residing within.
	Len() int

	// IsZero returns a Boolean value indicative of nilness.
	IsZero() bool

	// Compliant returns a Boolean value indicative of all slices
	// being in compliance with respect to relevant RFCs, such as
	// RFC 4512.
	Compliant() bool

	// Contains returns a Boolean value indicative of whether the
	// specified string value represents the RFC 4512 OID of a
	// Definition qualifier found within the receiver instance.
	Contains(string) bool

	// String returns the string representation of the underlying
	// stack instance.  Note that the manner in which the instance
	// was initialized will influence the resulting output.
	String() string

	// Maps returns slices of DefinitionMap instances, each of which
	// are expressions of actual Definition qualifiers found within
	// the receiver instance.
	Maps() DefinitionMaps

	// Inventory returns an instance of Inventory containing
	// the fundamental numerical and textual identifiers for
	// the contents of any given definition collection. Note
	// that the semantics for representation within instances
	// of this type vary among definition types.
	Inventory() Inventory

	// Type returns the string literal name for the receiver instance.
	// For example, if the receiver is LDAPSyntaxes, "ldapSyntaxes"
	// is the return value.  This is useful in case switching scenarios.
	Type() string

	// Push returns an error instance following an attempt to push
	// an instance of any into the receiver stack instance.  The
	// appropriate instance type depends on the nature of the
	// underlying stack instance.
	Push(any) error

	// cast is a private method used to reduce the cyclomatic penalties
	// normally incurred during the handling the eight (8) distinct RFC
	// 4512 definition stack types. cast is meant to allow this package
	// to access certain stackage.Stack methods that we have not wrapped
	// explicitly.
	cast() stackage.Stack
}

/*
TODO: Deprecated: get rid of me using Options.HangingIndents
*/
func hindent(use bool) (x string) {
	x = string(rune(10)) + `    `
	if !use {
		x = ` `
	}

	return
}

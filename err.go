package schemax

/*
err.go contains a quick error wrapper and some predefined errors meant
for use within this package as well as by end-users writing closures.
*/

import "errors"

var (
	ErrNamingViolationMissingMust  error = errors.New("Naming violation; required attribute type not used")
	ErrNamingViolationUnsanctioned error = errors.New("Naming violation; unsanctioned attribute type used")
	ErrNamingViolationChildlessSSR error = errors.New("Naming violation; childless superior structure rule")
	ErrNamingViolationBadClassAttr error = errors.New("Naming violation; named object class does not facilitate one or more attribute types present")
	ErrDuplicateDef                error = errors.New("Cannot parse or load; duplicate definition")
	ErrNilSyntaxQualifier          error = errors.New("No SyntaxQualifier instance assigned to LDAPSyntax")
	ErrNilValueQualifier           error = errors.New("No ValueQualifier instance assigned to AttributeType")
	ErrNilAssertionMatcher         error = errors.New("No AssertionMatcher instance assigned to MatchingRule")
	ErrNilReceiver                 error = errors.New("Receiver instance is nil")
	ErrNilInput                    error = errors.New("Input instance is nil")
	ErrNilDef                      error = errors.New("Referenced definition is nil or not specified")
	ErrNilSchemaRef                error = errors.New("Receiver instance lacks a Schema reference")
	ErrDefNonCompliant             error = errors.New("Definition failed compliance checks")
	ErrInvalidInput                error = errors.New("Input instance not compatible")
	ErrInvalidNames                error = errors.New("One or more names (descr) are invalid")
	ErrInvalidOID                  error = errors.New("Numeric OID is invalid")
	ErrInvalidRuleID               error = errors.New("Invalid integer identifier for structure rule")
	ErrInvalidSyntax               error = errors.New("Value does not meet the prescribed syntax qualifications")
	ErrInvalidValue                error = errors.New("Value does not meet the prescribed value qualifications")
	ErrNoMatch                     error = errors.New("Values do not match according to prescribed assertion match")
	ErrInvalidType                 error = errors.New("Incompatible type for operation")
	ErrTypeAssert                  error = errors.New("Type assertion failed")
	ErrNotUnique                   error = errors.New("Definition is already defined")
	ErrNotEqual                    error = errors.New("Values are not equal")
	ErrMissingNumericOID           error = errors.New("Missing or invalid numeric OID for definition")
	ErrInvalidDNOrFlatInt          error = errors.New("Invalid DN or flattened integer")
	ErrIncompatStructuralClass     error = errors.New("Incompatible structural class for target")

	ErrSuperTypeNotFound     error = errors.New("SUP AttributeType not found")
	ErrOrderingRuleNotFound  error = errors.New("ORDERING MatchingRule not found")
	ErrSubstringRuleNotFound error = errors.New("SUBSTR MatchingRule not found")
	ErrEqualityRuleNotFound  error = errors.New("EQUALITY MatchingRule not found")

	ErrAttributeTypeNotFound    error = errors.New("AttributeType not found")
	ErrObjectClassNotFound      error = errors.New("ObjectClass not found")
	ErrNameFormNotFound         error = errors.New("NameForm not found")
	ErrMatchingRuleNotFound     error = errors.New("MatchingRule not found")
	ErrMatchingRuleUseNotFound  error = errors.New("MatchingRuleUse not found")
	ErrLDAPSyntaxNotFound       error = errors.New("LDAPSyntax not found")
	ErrDITContentRuleNotFound   error = errors.New("DITContentRule not found")
	ErrDITStructureRuleNotFound error = errors.New("DITStructureRule not found")
)

var mkerr func(string) error = errors.New

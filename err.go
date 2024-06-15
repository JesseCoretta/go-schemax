package schemax

/*
err.go contains a quick error wrapper and some predefined errors meant
for use within this package as well as by end-users writing closures.
*/

import "errors"

var (
	ErrNilSyntaxQualifier  error = errors.New("No SyntaxQualifier instance assigned to LDAPSyntax")
	ErrNilValueQualifier   error = errors.New("No ValueQualifier instance assigned to AttributeType")
	ErrNilAssertionMatcher error = errors.New("No AssertionMatcher instance assigned to MatchingRule")
	ErrNilReceiver         error = errors.New("Receiver instance is nil")
	ErrNilInput            error = errors.New("Input instance is nil")
	ErrNilDef              error = errors.New("Referenced definition is nil or not specified")
	ErrNilSchemaRef        error = errors.New("Receiver instance lacks a Schema reference")
	ErrDefNonCompliant     error = errors.New("Definition failed compliancy checks")
	ErrInvalidInput        error = errors.New("Input instance not compatible")
	ErrInvalidSyntax       error = errors.New("Value does not meet the prescribed syntax qualifications")
	ErrInvalidValue        error = errors.New("Value does not meet the prescribed value qualifications")
	ErrNoMatch             error = errors.New("Values do not match according to prescribed assertion match")
	ErrInvalidType         error = errors.New("Incompatible type for operation")
	ErrTypeAssert          error = errors.New("Type assertion failed")
	ErrNotUnique           error = errors.New("Definition is already defined")
	ErrMissingNumericOID   error = errors.New("Missing or invalid numeric OID for definition")

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

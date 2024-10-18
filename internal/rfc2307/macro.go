package rfc2307

/*
Macros defines a single name-to-numericOID macro allowing support for
RFC 2307 definitions which do not specify explicit numeric OIDs.
*/
var Macros map[string]string = map[string]string{
	`nisSchema`: `1.3.6.1.1.1`,
}

# go-schemax

Package schemax incorporates a powerful [RFC 4512](https://www.rfc-editor.org/rfc/rfc4512.txt) parser, wrapped with convenient, reflective features for creating and interrogating directory schemas.

Requires Go version 1.22 or higher.

[![Go Report Card](https://goreportcard.com/badge/JesseCoretta/go-schemax)](https://goreportcard.com/report/github.com/JesseCoretta/go-schemax) [![codecov](https://codecov.io/gh/JesseCoretta/go-schemax/graph/badge.svg?token=6P4ZUQ3IGP)](https://codecov.io/gh/JesseCoretta/go-schemax) [![CodeQL](https://github.com/JesseCoretta/go-schemax/workflows/CodeQL/badge.svg)](https://github.com/JesseCoretta/go-schemax/actions/workflows/github-code-scanning/codeql) [![Reference](https://pkg.go.dev/badge/github.com/JesseCoretta/go-schemax.svg)](https://pkg.go.dev/github.com/JesseCoretta/go-schemax) [![License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://github.com/JesseCoretta/go-schemax/blob/main/LICENSE) [![Help Animals](https://img.shields.io/badge/help_animals-gray?label=%F0%9F%90%BE%20%F0%9F%98%BC%20%F0%9F%90%B6&labelColor=yellow)](https://github.com/JesseCoretta/JesseCoretta/blob/main/DONATIONS.md)

## License

The schemax package is available under the terms of the MIT license.  For further details, see the LICENSE file within the root of the repository.

## Releases

Two (2) releases are available for end-users:

| Version | Notes |
| :----- | :--- |
| 1.1.6 | Legacy, custom parser |
| >= 1.5.0 | Current, ANTLR parser |

## History of schemax

The goal of schemax has always been to provide a reliable parsing subsystem for directory schema definitions that allows transformation into usable Go objects.

The original design of schemax (version < 1.5.0) involved a custom-made parser. While this design performed remarkably well for years, it was not without its shortcomings. 

The newly released build of schemax involves the import of an ANTLR4-based [RFC 4512](https://www.rfc-editor.org/rfc/rfc4512.txt) lexer/parser solution. This is made possible using a newly released "sister" package -- [`go-antlr4512`](https://github.com/JesseCoretta/go-antlr4512) -- which handles all of the low-level ANTLR actions such as tokenization.

Therefore, the new build of schemax is of a simpler fundamental design thanks to offloading the bulk of the parser to another package. This also keeps all code-grading penalties (due to ANTLR's characteristically high cyclomatic factors) confined elsewhere, and allows schemax to focus on extending the slick features users have come to expect.

Users who are only interested in _tokenization_ and do not require the advanced features of this package should consider use of [`go-antlr4512`](https://github.com/JesseCoretta/go-antlr4512) exclusively.

## The Parser

The (ANTLR) parsing subsystem imported by the aforementioned sister package is flexible in terms of the following:

  - Presence of header, footer and line-terminating Bash comments surrounding a given definition is acceptable
    - Note that comments are entirely _discarded_ by ANTLR
  - Support for (escaped!) `'` and `\` characters within quoted strings ('this isn\\'t a bad example')
  - Support for linebreaks within definitions
  - Definition prefixing allows variations of the standard [RFC 4512](https://www.rfc-editor.org/rfc/rfc4512.txt) "labels" during file and directory parsing
    - "`attributeTypes`", "`attributeType`" and other variations are permitted for `AttributeType` definitions
  - Definition delimitation -- using colon (`:`), equals (`=`) or whitespace (` `, `\t`) of any sensible combination -- are permitted for the purpose of separating a definition prefix (label) from its definition statement
    - "attributeTypes: ...", "attributeType=...", "attributeType ..." are valid expressions
  - Multiple files are joined using an [ASCII](## "American Standard Code for Information Interchange") [#10](## '0x0a') during **directory** parsing
    - Users need not worry about adding a trailing newline to each file to be read; schemax will do this for you if needed

## File and Directory Readers

The legacy release branches of schemax did not offer a robust file and directory parsing solution, rather it focused on the byte representations of a given definition and the tokens derived therein, leaving it to the end-user to devise a delivery method.

The new (>=1.5.0) release branches introduce proper  `ParseRaw`, `ParseFile` and `ParseDirectory` methods that greatly simplify use of this package in the midst of an established schema "library".  For example:

```
func main() {
	// Create a brand new schema *and* load it with
	// standard RFC-sourced definitions.
	mySchema := NewSchema()

	// If your organization has any number of its own
	// custom schema definitions, there are three (3)
	// ways you could load them, each of which are
	// covered below.

	// By directory: a directory structure -- which may
	// or may not contain subdirectories of its own --
	// containing one or more ".schema" files, named in
	// such a way that they remain naturally ordered in
	// terms of super type, super class, and super rule
	// dependencies.
	schemaDir := "/home/you/ds/schemas"
	if err := mySchema.ParseDirectory(schemaDir); err != nil {
		fmt.Println(err)
		return
	}

	// By file: a single file, which MUST end in ".schema",
	// read using the Schema.ParseFile method.  Note the same
	// dependency considerations described in the previous
	// "directory" example shall apply here.
	schemaFile := "/home/you/other.schema"
	if err := mySchema.ParseFile(schemaFile); err != nil {
		fmt.Println(err)
		return
	}

	// By bytes: a series of bytes previously read from a file
	// or other source can be submitted to the Schema.ParseRaw
	// method. Again, the same dependency considerations noted
	// above shall apply here.
	schemaBytes := []byte{...contents of some .schema file...}
	if err := mySchema.ParseRaw(schemaBytes); err != nil {
		fmt.Println(err)
		return
	}

	// Take a snapshot of the current definition counts by
	// category:
	fmt.Printf("%#v\n", mySchema.Counters()
	// Output: schemax.Counters{LS:67, MR:44, AT:317, MU:32, OC:80, DC:1, NF:13, DS:13}
}
```

Though the `ParseFile` function operates identically to the above-demonstrated `ParseDirectory` function, it is important to order the respective files and directories according to any applicable dependencies.  In other words, if "fileB.schema" requires definitions from "fileA.schema", "fileA.schema" must be parsed first.

Sub-directories encountered shall be traversed indefinitely and in their natural order according to name. Files encountered through directory traversal shall only be read and parsed IF the extension is ".schema".  This prevents other files -- such as text or `README.md` files -- from interfering with the parsing process needlessly. The same considerations related to ordering of directories by name applies to individual ".schema" files.

The general rule-of-thumb is suggests that if the `ls -l` Bash command _consistently_ lists the indicated schema files in correct order, _and_ assuming those files contain properly ordered and well-formed definitions, the parsing process should work nicely.

The `ParseRaw` method is subject to the same conditions related to the order of dependent definitions.

## The Schema Itself

The `Schema` type defined within this package is a [`stackage.Stack`](https://pkg.go.dev/github.com/JesseCoretta/go-stackage#Stack) derivative type. An instance of a `Schema` can manifest in any of the following manners:

  - As an empty (unpopulated) `Schema`, initialized by way of the `NewEmptySchema` function
  - As a basic (minimally populated) `Schema`, initialized by way of the `NewBasicSchema` function
  - As a complete (fully populated) `Schema`, initialized by way of the `NewSchema` function

There are certain scenarios which call for one of the above initialization procedures:

  - An empty `Schema` is ideal for LDAP professionals, and allows for the creation of a `Schema` of particularly narrow focus for R&D, testing or product development
  - A basic `Schema` resembles the foundational (starting) `Schema` context observed in most directory server products, in that it comes "pre-loaded" with official `LDAPSyntax` and `MatchingRule` definitions -- but few to no `AttributeTypes` -- making it a most suitable empty canvas upon which a new `Schema` may be devised from scratch
  - A full `Schema` is the most obvious choice for "Quick Start" scenarios, in that a `Schema` is produced containing a very large portion of the standard `AttributeType` and `ObjectClass` definitions used in the wild by most (if not all) directory products

Regardless of the content present, a given `Schema` is capable of storing definitions from all eight (8) [RFC 4512](https://www.rfc-editor.org/rfc/rfc4512.txt) "categories".  These are known as "collections", and are stored in nested [`stackage.Stack`](https://pkg.go.dev/github.com/JesseCoretta/go-stackage#Stack) derivative types, accessed using any of the following methods:

  - `Schema.LDAPSyntaxes`
  - `Schema.MatchingRules`
  - `Schema.AttributeTypes`
  - `Schema.MatchingRuleUses`
  - `Schema.ObjectClasses`
  - `Schema.DITContentRules`
  - `Schema.NameForms`
  - `Schema.DITStructureRules`

Definition instances produced by way of parsing -- namely using one of the `Schema.Parse<Type>` methods-- will automatically gain internal access to the `Schema` instance in which it is stored.

However, definitions produced manually by way of the various `Set<Item>` methods or by way of localized `Parse` method  extended through types defined within this package will require manual execution of the `SetSchema` method, using the intended `Schema` instance as the input argument.  Ideally this should occur early in the definition composition.

In either case, this internal reference is used for seamless verification of any reference, such as an `LDAPSyntax`, when introduced to a given type instance. This ensures definition pointer references remain valid.

## Closure Methods

This package is closure-friendly with regards to user-authored closure functions or methods meant to perform specific tasks:

  - Assertion matching, by way of an instance of `MatchingRule` applicable to two assertion values within a `AssertionMatcher` closure (i.e.: is "value1" equal to "value2"?)
  - Syntax qualification, by way of an instance of `LDAPSyntax` to be honored by a value within a `SyntaxQualifier` closure (i.e.: does value qualify for specified syntax?)
  - General-use value qualification, by way of an instance of `AttributeType` to be analyzed in specialized scenarios within a `ValueQualifier` closure (i.e: company/user-specific value processing)
  - Definition string representation, through assignment of a custom `Stringer` closure to eligible definition instances

Understand that assertion, syntax and general-use qualifying closures are entirely user-defined; this package does not provide such predefined instances itself, leaving that to the user or another package which may be imported and used in a "pluggable" manner in this context.

See [RFC 4517](https://www.rfc-editor.org/rfc/rfc4517.txt), et al, for some practical guidelines relating to certain syntax and assertion matching procedures that may guide users in creating such closures.

This package does, however, include a default `Stringer`, which can be invoked for an instance simply by running the instance's `SetStringer` method in niladic form.

## Fluent Methods

This package extends fluent methods that are write-based in nature. Typically these methods are prefaced with `Set` or `Push`. This means such methods may be "chained" together using the standard Go command "." delimiter.

Fluency does not extend to methods that are interrogative in nature, in that they return `bool`, `string` or `error` values. Fluency also precludes use of the `Definition` interface due to unique return signatures.

## Built-In Definitions

The following table describes the contents and coverage of the so-called "built-in" schema definitions, all of which are sourced from recognized RFCs only. These can be imported en masse by users, or in piece-meal fashion. At present, the library contains more than four hundred such definitions.

Note that no `dITContentRule` definitions exist in any RFC at this time, thus none are available for import.

| DOCUMENT | [LS](## "LDAP Syntaxes")  | [MR](## "Matching Rules")  | [AT](## "Attribute Types")  | [OC](## "Object Classes")  | [DC](## "DIT Content Rules")  | [NF](## "Name Forms")  | [DS](## "DIT Structure Rules")  |
| -------- | :----: | :----: | :----: | :----: | :----: | :----: | :----:  |
| [![RFC 2307](https://img.shields.io/badge/RFC-2307-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc2307)  |  ✅  |  ✅  |  ✅  |  ✅  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |
| [![RFC 2589](https://img.shields.io/badge/RFC-2589-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc2589)  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ✅  |  ✅  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |
| [![RFC 2798](https://img.shields.io/badge/RFC-2798-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc2798)  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ✅  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |
| [![RFC 3045](https://img.shields.io/badge/RFC-3045-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc3045)  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ✅  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |
| [![RFC 3671](https://img.shields.io/badge/RFC-3671-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc3671)  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ✅  |  ✅  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |
| [![RFC 3672](https://img.shields.io/badge/RFC-3672-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc3672)  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ✅  |  ✅  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |
| [![RFC 4403](https://img.shields.io/badge/RFC-4403-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc4403)  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ✅  |  ✅  |  ⁿ/ₐ  |  ✅  |  ✅  |
| [![RFC 4512](https://img.shields.io/badge/RFC-4512-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc4512)  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ✅  |  ✅  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |
| [![RFC 4517](https://img.shields.io/badge/RFC-4517-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc4517)  |  ✅  |  ✅  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |
| [![RFC 4519](https://img.shields.io/badge/RFC-4519-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc4519)  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ✅  |  ✅  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |
| [![RFC 4523](https://img.shields.io/badge/RFC-4523-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc4523)  |  ✅  |  ✅  |  ✅  |  ✅  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |
| [![RFC 4524](https://img.shields.io/badge/RFC-4524-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc4524)  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ✅  |  ✅  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |
| [![RFC 4530](https://img.shields.io/badge/RFC-4530-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc4530)  |  ✅  |  ✅  |  ✅  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |
| [![RFC 5020](https://img.shields.io/badge/RFC-5020-blue?cacheSeconds=500000)](https://datatracker.ietf.org/doc/html/rfc5020)  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ✅  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |  ⁿ/ₐ  |

## A note about test latency

The go-schemax package contains well over two hundred unit tests/examples. When performing a full run of `go test`, it takes approxiately one (1) second to complete.  The reason it is so slow is due to the elaborate nature in which some of the tests are conducted.

go-schemax is written to be extremely resilient in terms of operational stability. Many safeguards and balance-checks are performed at many junctures. This creates test coverage difficulties, meaning that certain error conditions simply cannot be triggered through ordinary means because the "higher layers" do a (really) good job at preventing such conditions. The result is an unflattering code coverage grade.

Therefore, rather extraordinary measures are required to mitigate coverage issues, such as re-enveloping (or re-casting) instances to try and "trick" the underlying `stackage.Stack` type into accepting something undesirable. This, and many other bizarre procedures, can be observed in the `_codecov` test functions found throughout the various `*_test.go` files.

Even worse, some of the `Example` functions perform **repeated imports of the ENTIRE SCHEMA LIBRARY** -- just for the benefit of certain examples. As such, a lot more is happening within these tests than a user would normally trigger under ordinary conditions. Understand this is intentional, and is meant to ensure that go-schemax remains quick and snappy regardless.

As such, latency in tests run by this package probably won't be observed under ordinary (real-life) circumstances.

package schemax

/*
Marshal takes the provided schema definition (def) and attempts to marshal it into x.  x MUST be one of the following types:

• AttributeType

• ObjectClass

• LDAPSyntax

• MatchingRule

• MatchingRuleUse

• DITContentRule

• DITStructureRule

• NameForm

Should any validation errors occur, a non-nil instance of error is returned.

Note that it is far more convenient to use the Subschema.Marshal wrapper, as it only requires a single argument (the raw definition).
*/
func Marshal(raw string, x interface{},
	atc AttributeTypeCollection,
	occ ObjectClassCollection,
	lsc LDAPSyntaxCollection,
	mrc MatchingRuleCollection,
	mruc MatchingRuleUseCollection,
	dcrc DITContentRuleCollection,
	dsrc DITStructureRuleCollection,
	nfc NameFormCollection) (err error) {
	// I am so sorry.

	if len(raw) == 0 {
		return raise(emptyDefinition, "no length")
	}

	// Remove all outer WHSP, collapse all successive inner
	// WHSP to single space, and purge all linebreaks.
	raw = sanitize(raw)

	var mac *Macros
	if atc != nil {
		mac = atc.(*AttributeTypes).macros
	}

	def, ok := newDefinition(x, mac)
	if !ok {
		return raise(invalidMarshal, "newDefinition: assembly failure")
	}

	id, rest, ok := parse(raw)
	if !ok {
		return raise(invalidMarshal, invalidOID.Error())
	}

	if _, assert := x.(*DITStructureRule); assert {
		def.values[0].Set(valueOf(NewRuleID(id[0])))
	} else {
		// here we parse the OID, which is a constant
		// for almost every definition type (except
		// dITStructureRule instances)
		isnumoid := isNumericalOID(id[0])
		if mac.IsZero() {
			if !isnumoid {
				return raise(invalidMarshal, "unresolvable alias '%s' (nil manifest)", id[0])
			}
			def.values[0].Set(valueOf(NewOID(id[0])))
		} else {
			oid, ok := mac.Resolve(id[0])
			if !ok {
				return raise(invalidMarshal, "unresolvable alias '%s'", id[0])
			}
			def.values[0].Set(valueOf(oid))
		}
	}

	// Now we'll parse all KEY WHSP VALUE [VALUE...] instances
	for {
		if len(rest) <= 1 || err != nil {
			break
		}

		// parseDefLabel receives one chunk of information
		// (which should be a single, raw schema definition).
		// We then attempt to extract a "label" from this,
		// and then parse the remainder ("rest") based on
		// the appropriate known actions for said type of
		// value (e.g.: `NAME` vs. `SYNTAX`).
		var label []string
		if label, rest, ok = parse_definition_label(rest); !ok {
			return raise(invalidLabel,
				"failed parse for label was: '%s', raw def: '%s'",
				label, rest)
		} else {
			idx := def.lfindex(label[0])
			if idx == -1 {
				return raise(invalidLabel,
					"failed index localization (lfindex) for label: '%s', raw def: '%s'",
					label[0], rest)
			}

			var value []string
			if value, rest, ok = def.meths[idx](rest); !ok {
				return raise(invalidValue,
					"failed value localization for label: '%s' (deflabel:%s), raw def: '%s'",
					label[0], def.labels[idx], rest)
			}

			switch def.labels[idx] {
			case `KIND`:
				err = def.setKind(value[0], idx)
			case `EXT`:
				err = def.setExtensions(label[0], value, idx)
			case `NAME`:
				err = def.setName(idx, value...)
			case `DESC`:
				err = def.setDesc(idx, value[0])
			case `BOOLS`:
				err = def.setdefinitionFlags(label[0], x)
			case `USAGE`:
				err = def.setUsage(value[0], idx)
			case `FORM`:
				err = def.setNameForm(nfc, x, value[0], idx)
			case `MAY`:
				err = def.setPermittedAttributeTypes(atc, x, value, idx)
			case `NOT`:
				err = def.setProhibitedAttributeTypes(atc, x, value, idx)
			case `MUST`:
				err = def.setRequiredAttributeTypes(atc, x, value, idx)
			case `APPLIES`:
				err = def.setApplicableAttributeTypes(atc, x, value, idx)
			case `AUX`:
				err = def.setAuxiliaryObjectClasses(occ, x, value, idx)
			case `OC`:
				err = def.setStructuralObjectClass(occ, x, value[0], idx)
			case `SUP`:
				switch def.definitionType() {
				case `oc`:
					err = def.setSuperiorObjectClasses(occ, x, value, idx)
				case `sat`, `at`:
					err = def.setSuperiorAttributeType(atc, x, value[0], idx)
				case `dsr`:
					err = def.setSuperiorDITStructureRules(dsrc, x, value, idx)
				}
			case `SYNTAX`:
				err = def.setSyntax(lsc, x, value[0], idx)
			case `EQUALITY`, `SUBSTR`, `SUBSTRINGS`, `ORDERING`:
				err = def.setEqSubOrd(mrc, x, value[0], idx)
			default:
				return raise(invalidLabel,
					"Field '%s'(def.label:'%s') unhandled; would have set '%s', raw def: '%s')",
					label[0], def.labels[idx], value, rest)
			}
		}
		label = []string{} // reset our label
	}

	if err != nil {
		return
	}

	// Our instance has been populated with the
	// marshaled bytes. Now we conduct validation
	// checks to ensure said bytes were sane.
	switch tv := x.(type) {
	case *LDAPSyntax:
		// we'll take an extra step to identify
		// any syntax that is considered to be
		// human readable either through a value
		// of 'FALSE' for the X-NOT-HUMAN-READABLE
		// well-known extension, or absence of said
		// extension altogether.
		if tv.Extensions.Exists(`X-NOT-HUMAN-READABLE`) {
			if strInSlice(`FALSE`, tv.Extensions[`X-NOT-HUMAN-READABLE`]) {
				tv.flags.set(HumanReadable)
			}
		} else {
			tv.flags.set(HumanReadable)
		}
		err = tv.validate()
	case *AttributeType:
		err = tv.validate()
	case *ObjectClass:
		err = tv.validate()
	case *NameForm:
		err = tv.validate()
	case *MatchingRule:
		err = tv.validate()
	case *MatchingRuleUse:
		err = tv.validate()
	case *DITContentRule:
		err = tv.validate()
	case *DITStructureRule:
		err = tv.validate()
	default:
		err = raise(unexpectedType,
			"Validator for %T not yet implemented", tv)
	}

	return
}

/*
DefinitionUnmarshalFunc is a first-class "closure" function intended for use in situations where it is desirable to format a given definition during the unmarshal process, i.e.: to add indents and linebreaks.

The string input argument defines a specifier to be printed just before the definition value. This is used to declare the type of definition being defined in a schema, e.g.: "attributetype". This is particularly useful when interacting with different LDAP DSA products, as such specifiers do vary between implementations.  One real-world example of this is the difference between OpenLDAP and Netscape directory subschema subentries -- namely "attributetype" vs. "attributetypes:". Providing a zero length string argument results in no specifier being printed. The user-provided specifier case is preserved when used (neither "folding" nor normalization will occur).

During the unmarshaling or unsafe stringifaction processes, users may choose to:

• Perform NO formatting whatsoever, producing definitions that span only a single line, or ...

• Use the package-provided formatting closure function appropriate for the definition type, or ...

• Define a custom unmarshal function that honors the defined signature of this type

By default, NO special formatting is performed during unmarshaling or unsafe stringification of definitions.
*/
type DefinitionUnmarshalFunc func() (string, error)

/*
Unmarshal takes an instance of one (1) of the following types and (if valid) and returns the textual form of the definition:

• ObjectClass

• AttributeType

• LDAPSyntax

• MatchingRule

• MatchingRuleUse

• DITContentRule

• DITStructureRule

• NameForm

Should any validation errors occur, a non-nil instance of error is returned.
*/
func Unmarshal(x interface{}) (def string, err error) {
	switch tv := x.(type) {
	case *ObjectClass:
		def, err = tv.unmarshal()
	case *AttributeType:
		def, err = tv.unmarshal()
	case *LDAPSyntax:
		def, err = tv.unmarshal()
	case *MatchingRule:
		def, err = tv.unmarshal()
	case *MatchingRuleUse:
		def, err = tv.unmarshal()
	case *DITContentRule:
		def, err = tv.unmarshal()
	case *DITStructureRule:
		def, err = tv.unmarshal()
	case *NameForm:
		def, err = tv.unmarshal()
	default:
		err = raise(invalidUnmarshal,
			"unknown or unsupported type %T", tv)
	}

	if err != nil {
		err = raise(invalidUnmarshal, err.Error())
	} else if len(def) == 0 {
		err = raise(invalidUnmarshal,
			"zero-length definition returned from Unmarshal (of %T)", x)
	}

	return
}

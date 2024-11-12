package schemax

import (
	"sort"

	"github.com/JesseCoretta/go-antlr4512"
)

var (
	new4512Schema func() antlr4512.Schema                          = antlr4512.NewSchema
	parseLS       func(string) (antlr4512.LDAPSyntax, error)       = antlr4512.ParseLDAPSyntax
	parseMR       func(string) (antlr4512.MatchingRule, error)     = antlr4512.ParseMatchingRule
	parseAT       func(string) (antlr4512.AttributeType, error)    = antlr4512.ParseAttributeType
	parseMU       func(string) (antlr4512.MatchingRuleUse, error)  = antlr4512.ParseMatchingRuleUse
	parseOC       func(string) (antlr4512.ObjectClass, error)      = antlr4512.ParseObjectClass
	parseDC       func(string) (antlr4512.DITContentRule, error)   = antlr4512.ParseDITContentRule
	parseNF       func(string) (antlr4512.NameForm, error)         = antlr4512.ParseNameForm
	parseDS       func(string) (antlr4512.DITStructureRule, error) = antlr4512.ParseDITStructureRule
)

/*
ParseLDAPSyntax returns an error following an attempt to parse raw into an
instance of [LDAPSyntax] and append it to the [Schema.LDAPSyntaxes] stack.
*/
func (r Schema) ParseLDAPSyntax(raw string) error {
	def, err := parseLS(raw)
	if err == nil {
		var _def LDAPSyntax
		if _def, err = r.marshalLS(def); err == nil {
			r.LDAPSyntaxes().push(_def)
		}
	}

	return err
}

/*
ParseMatchingRule returns an error following an attempt to parse raw into an
instance of [MatchingRule] and append it to the [Schema.MatchingRules] stack.
*/
func (r Schema) ParseMatchingRule(raw string) error {
	def, err := parseMR(raw)
	if err == nil {
		var _def MatchingRule
		if _def, err = r.marshalMR(def); err == nil {
			r.MatchingRules().push(_def)
		}
	}

	return err
}

/*
ParseMatchingRuleUse returns an error following an attempt to parse raw into an
instance of [MatchingRuleUse] and append it to the [Schema.MatchingRuleUses] stack.
*/
func (r Schema) ParseMatchingRuleUse(raw string) error {
	def, err := parseMU(raw)
	if err == nil {
		var _def MatchingRuleUse
		if _def, err = r.marshalMU(def); err == nil {
			r.MatchingRuleUses().push(_def)
		}
	}

	return err
}

/*
ParseAttributeType returns an error following an attempt to parse raw into an
instance of [AttributeType] and append it to the [Schema.AttributeTypes] stack.
*/
func (r Schema) ParseAttributeType(raw string) error {
	def, err := parseAT(raw)
	if err == nil {
		var _def AttributeType
		if _def, err = r.marshalAT(def); err == nil {
			r.AttributeTypes().push(_def)
		}
	}

	return err
}

/*
ParseObjectClass returns an error following an attempt to parse raw into an
instance of [ObjectClass] and append it to the [Schema.ObjectClasses] stack.
*/
func (r Schema) ParseObjectClass(raw string) error {
	def, err := parseOC(raw)
	if err == nil {
		var _def ObjectClass
		if _def, err = r.marshalOC(def); err == nil {
			r.ObjectClasses().push(_def)
		}
	}

	return err
}

/*
ParseDITContentRule returns an error following an attempt to parse raw into an
instance of [DITContentRule] and append it to the [Schema.DITContentRules] stack.
*/
func (r Schema) ParseDITContentRule(raw string) error {
	def, err := parseDC(raw)
	if err == nil {
		var _def DITContentRule
		if _def, err = r.marshalDC(def); err == nil {
			r.DITContentRules().push(_def)
		}
	}

	return err
}

/*
ParseNameForm returns an error following an attempt to parse raw into an
instance of [NameForm] and append it to the [Schema.NameForms] stack.
*/
func (r Schema) ParseNameForm(raw string) error {
	def, err := parseNF(raw)
	if err == nil {
		var _def NameForm
		if _def, err = r.marshalNF(def); err == nil {
			r.NameForms().push(_def)
		}
	}

	return err
}

/*
ParseDITStructureRule returns an error following an attempt to parse raw into an
instance of [DITStructureRule] and append it to the [Schema.DITStructureRules] stack.
*/
func (r Schema) ParseDITStructureRule(raw string) error {
	def, err := parseDS(raw)
	if err == nil {
		var _def DITStructureRule
		if _def, err = r.marshalDS(def); err == nil {
			r.DITStructureRules().push(_def)
		}
	}

	return err
}

/*
incorporate returns an error following an attempt to marshal the contents
of s into r. The concept of "incorporation" is another term for post-parsing
"hand-off" from antlr4512 to schemax, and is the last of two (2) phases
involving the antlr4512 package.

Empty slice types within s shall not result in an error.
*/
func (r Schema) incorporate(s antlr4512.Schema) (err error) {
	for _, err := range []error{
		r.incorporateLS(s.LS),
		r.incorporateMR(s.MR),
		r.incorporateAT(s.AT),
		r.incorporateMU(s.MU),
		r.incorporateOC(s.OC),
		r.incorporateDC(s.DC),
		r.incorporateNF(s.NF),
		r.incorporateDS(s.DS),
	} {
		if err != nil {
			break
		}
	}

	return
}

func (r Schema) incorporateLS(s antlr4512.LDAPSyntaxes) (err error) {
	for i := 0; i < len(s) && err == nil; i++ {
		var def LDAPSyntax
		if def, err = r.marshalLS(s[i]); err == nil {
			r.LDAPSyntaxes().push(def)
		}
	}

	return
}

func (r Schema) marshalLS(s antlr4512.LDAPSyntax) (def LDAPSyntax, err error) {
	// try to resolve a macro only if no
	// numeric OID is present.
	if len(s.Macro) == 2 && len(s.OID) == 0 {
		mc, found := r.Macros().Resolve(s.Macro[0])
		if found {
			s.OID = mc + `.` + s.Macro[1]
		}
	}

	if !isNumericOID(s.OID) {
		err = ErrMissingNumericOID
		return
	}

	// fail attempts to marshal a duplicate definition.
	if lup := r.LDAPSyntaxes().get(s.OID); !lup.IsZero() {
		err = ErrDuplicateDef
		return
	}

	_def := newLDAPSyntax()
	_def.OID = s.OID
	_def.Desc = s.Desc
	_def.schema = r
	_def.Extensions.setDefinition(LDAPSyntax{_def})

	// Marshal our extensions, taking our sorting preference into account
	marshalExt(s.Extensions, _def.Extensions,
		r.Options().Positive(SortExtensions))

	str, err := _def.prepareString() // perform one-time text/template op
	if err == nil {
		// Save the stringer
		_def.stringer = func() string {
			// Return a preserved value.
			return str
		}
		def = LDAPSyntax{_def}
	}

	return
}

func (r Schema) incorporateMR(s antlr4512.MatchingRules) (err error) {
	for i := 0; i < len(s) && err == nil; i++ {
		var def MatchingRule
		if def, err = r.marshalMR(s[i]); err == nil {
			r.MatchingRules().push(def)
		}
	}

	return
}

func (r Schema) marshalMR(s antlr4512.MatchingRule) (def MatchingRule, err error) {
	// try to resolve a macro only if no
	// numeric OID is present.
	if len(s.Macro) == 2 && len(s.OID) == 0 {
		mc, found := r.Macros().Resolve(s.Macro[0])
		if found {
			s.OID = mc + `.` + s.Macro[1]
		}
	}

	// fail attempts to marshal a duplicate definition.
	if lup := r.MatchingRules().get(s.OID); !lup.IsZero() {
		err = ErrDuplicateDef
		return
	}

	if !isNumericOID(s.OID) {
		err = ErrMissingNumericOID
		return
	} else if !isNumericOID(s.Syntax) {
		// try to resolve a quotedDescriptor to an OID
		// by scanning the DESC values of each syntax
		// within the receiver. Use of Whitespace and
		// case folding are not significant in the
		// matching process.
		if syn := r.LDAPSyntaxes().get(s.Syntax); !syn.IsZero() {
			// overwrite the DESC reference of syntax
			// with the intended syntax's numeric OID.
			s.Syntax = syn.NumericOID()
		}
	}

	syn := r.LDAPSyntaxes().get(s.Syntax)
	if syn.IsZero() {
		// throw an error due to bad syntax ref
		err = mkerr(ErrLDAPSyntaxNotFound.Error() + `(` + s.Syntax + `)`)
		return
	}

	_def := newMatchingRule()
	_def.OID = s.OID
	_def.Desc = s.Desc
	_def.Obsolete = s.Obsolete
	_def.schema = r
	_def.Syntax = syn
	_def.Extensions.setDefinition(MatchingRule{_def})

	// Marshal our extensions, taking our sorting preference into account
	marshalExt(s.Extensions, _def.Extensions,
		r.Options().Positive(SortExtensions))

	for _, name := range s.Name {
		_def.Name.push(name)
	}

	str, err := _def.prepareString() // perform one-time text/template op
	if err == nil {
		// Save the stringer
		_def.stringer = func() string {
			// Return a preserved value.
			return str
		}
		def = MatchingRule{_def}
	}

	return
}

func (r Schema) incorporateMU(s antlr4512.MatchingRuleUses) (err error) {
	for i := 0; i < len(s) && err == nil; i++ {
		var def MatchingRuleUse
		if def, err = r.marshalMU(s[i]); err == nil {
			r.MatchingRuleUses().push(def)
		}
	}

	return
}

func (r Schema) marshalMU(s antlr4512.MatchingRuleUse) (def MatchingRuleUse, err error) {
	if !isNumericOID(s.OID) {
		err = ErrMissingNumericOID
		return
	}

	lup := r.MatchingRules().get(s.OID)
	if lup.IsZero() {
		// silently ignore attempts to reference a bogus matchingRule OID
		return
	}

	// fail attempts to marshal a duplicate definition.
	if lup := r.MatchingRuleUses().get(s.OID); !lup.IsZero() {
		err = ErrDuplicateDef
		return
	}

	_def := newMatchingRuleUse()
	_def.OID = lup
	_def.Desc = s.Desc
	_def.Obsolete = s.Obsolete
	_def.schema = r
	_def.Extensions.setDefinition(MatchingRuleUse{_def})

	sortL := r.Options().Positive(SortLists)

	for i := 0; i < len(sortList(s.Applies, sortL)); i++ {
		_at := s.Applies[i]
		at := r.AttributeTypes().get(_at)
		if at.IsZero() {
			err = ErrAttributeTypeNotFound
			return
		}
		_def.Applies.push(at)
	}

	// Marshal our extensions, taking our sorting preference into account
	marshalExt(s.Extensions, _def.Extensions,
		r.Options().Positive(SortExtensions))

	for _, name := range s.Name {
		_def.Name.push(name)
	}

	str, err := _def.prepareString() // perform one-time text/template op
	if err == nil {
		// Save the stringer
		_def.stringer = func() string {
			// Return a preserved value.
			return str
		}
		def = MatchingRuleUse{_def}
	}

	return
}

func (r Schema) incorporateAT(s antlr4512.AttributeTypes) (err error) {
	for i := 0; i < len(s) && err == nil; i++ {
		var def AttributeType
		if def, err = r.marshalAT(s[i]); err == nil {
			r.AttributeTypes().push(def)
		}
	}

	return
}

func (r Schema) marshalAT(s antlr4512.AttributeType) (def AttributeType, err error) {
	// try to resolve a macro only if no numeric OID is present.
	if s.OID = handleMacro(r, s.Macro, s.OID); !isNumericOID(s.OID) {
		err = ErrMissingNumericOID
		return
	} else if lup := r.AttributeTypes().get(s.OID); !lup.IsZero() {
		// fail on attempts to marshal a duplicate definition.
		err = ErrDuplicateDef
		return
	}

	// alloc
	_def := newAttributeType()

	// Take any primitive values we can
	_def.OID = s.OID
	_def.Desc = s.Desc
	_def.Obsolete = s.Obsolete
	_def.NoUserMod = s.Immutable
	_def.Collective = s.Collective
	_def.Single = s.Single
	_def.MUB = s.MUB
	_def.schema = r
	_def.Extensions.setDefinition(AttributeType{_def})

	// avoid bogus Boolean states
	if _def.Single && _def.Collective {
		err = mkerr("AttributeType is both single-valued and collective; aborting (" + s.OID + `)`)
		return
	}

	// verify and set an LDAPSyntax if specified
	if syn := s.Syntax; len(syn) > 0 {
		llup := r.LDAPSyntaxes().get(syn)
		if llup.IsZero() {
			// throw an error due to bad syntax ref
			err = mkerr(ErrLDAPSyntaxNotFound.Error() + `(` + syn + `)`)
			return
		}
		_def.Syntax = llup
	}

	// Attempt to marshal any (EQ/ORD/SS) matching rules
	// from the antlr instance into the destination instance
	if err = _def.marshalMRs(s); err != nil {
		return
	}

	// verify and set supertype, if specified.
	if _sup := s.SuperType; len(_sup) > 0 {
		sup := r.AttributeTypes().get(_sup)
		if sup.IsZero() {
			err = mkerr(ErrAttributeTypeNotFound.Error() + `( supertype: ` + _sup + `)`)
			return
		}
		_def.SuperType = sup
	}

	// marshal our intended attribute usage
	_def.marshalUsage(s)

	// Marshal our extensions, taking our sorting preference into account
	marshalExt(s.Extensions, _def.Extensions,
		r.Options().Positive(SortExtensions))

	// Process and set any names
	for _, name := range s.Name {
		_def.Name.push(name)
	}

	str, err := _def.prepareString() // perform one-time text/template op
	if err == nil {
		// Save the stringer
		_def.stringer = func() string {
			// Return a preserved value.
			return str
		}
		def = AttributeType{_def}
	}

	return
}

func (r *attributeType) marshalMRs(s antlr4512.AttributeType) (err error) {
	// verify any non-zero matchingrule references
	for idx, mrl := range []string{
		s.Equality,  // EQUALITY
		s.Substring, // SUBSTR
		s.Ordering,  // ORDERING
	} {
		if len(mrl) > 0 {
			// If mrl was non-zero, lookup using its
			// text value (name or OID). Fail if not
			// found, assign to appropriate field
			// otherwise.
			llup := r.schema.MatchingRules().get(mrl)
			if llup.IsZero() {
				err = mkerr(ErrMatchingRuleNotFound.Error() + `(` + mrl + `)`)
				break
			}

			switch idx {
			case 1:
				r.Substring = llup
			case 2:
				r.Ordering = llup
			default:
				r.Equality = llup
			}
		}
	}

	return
}

func (r *attributeType) marshalUsage(s antlr4512.AttributeType) {
	// verify and set usage IF NOT default ("userApplications")
	if len(s.Usage) > 0 {
		switch lc(s.Usage) {
		case `directoryoperation`:
			r.Usage = DirectoryOperationUsage
		case `distributedoperation`:
			r.Usage = DistributedOperationUsage
		case `dsaoperation`:
			r.Usage = DSAOperationUsage
		}
	}
}

func (r Schema) incorporateOC(s antlr4512.ObjectClasses) (err error) {
	for i := 0; i < len(s) && err == nil; i++ {
		var def ObjectClass
		if def, err = r.marshalOC(s[i]); err == nil {
			r.ObjectClasses().push(def)
		}
	}

	return
}

func (r Schema) marshalOC(s antlr4512.ObjectClass) (def ObjectClass, err error) {
	// try to resolve a macro only if no numeric OID is present.
	if s.OID = handleMacro(r, s.Macro, s.OID); !isNumericOID(s.OID) {
		err = ErrMissingNumericOID
		return
	}

	// fail attempts to marshal a duplicate definition.
	if lup := r.ObjectClasses().get(s.OID); !lup.IsZero() {
		err = ErrDuplicateDef
		return
	}

	_def := newObjectClass()
	_def.OID = s.OID
	_def.Desc = s.Desc
	_def.Obsolete = s.Obsolete
	_def.schema = r
	_def.Extensions.setDefinition(ObjectClass{_def})

	sortL := r.Options().Positive(SortLists)

	// verify and set kind if not default of STRUCTURAL
	switch k := lc(s.Kind); k {
	case `structural`, ``:
		_def.Kind = StructuralKind
	case `auxiliary`:
		_def.Kind = AuxiliaryKind
	case `abstract`:
		_def.Kind = AbstractKind
	default:
		err = mkerr("Invalid kind for objectClass: " + k)
		return
	}

	for _, must := range sortList(s.Must, sortL) {
		m := r.AttributeTypes().get(must)
		if m.IsZero() {
			err = mkerr("Unknown AttributeType for MUST clause: " + must)
			return
		}
		_def.Must.push(m)
	}

	for _, may := range sortList(s.May, sortL) {
		m := r.AttributeTypes().get(may)
		if m.IsZero() {
			err = mkerr("Unknown AttributeType for MAY clause: " + may)
			return
		}
		_def.May.push(m)
	}

	for _, sup := range sortList(s.SuperClasses, sortL) {
		m := r.ObjectClasses().get(sup)
		if m.IsZero() {
			err = mkerr("Unknown SuperClass: " + sup)
			return
		}
		_def.SuperClasses.push(m)
	}

	// Marshal our extensions, taking our sorting preference into account
	marshalExt(s.Extensions, _def.Extensions,
		r.Options().Positive(SortExtensions))

	// Process and set any names
	for _, name := range s.Name {
		_def.Name.push(name)
	}

	str, err := _def.prepareString() // perform one-time text/template op
	if err == nil {
		// Save the stringer
		_def.stringer = func() string {
			// Return a preserved value.
			return str
		}
		def = ObjectClass{_def}
	}

	return
}

func (r Schema) incorporateDC(s antlr4512.DITContentRules) (err error) {
	for i := 0; i < len(s) && err == nil; i++ {
		var def DITContentRule
		if def, err = r.marshalDC(s[i]); err == nil {
			r.DITContentRules().push(def)
		}
	}

	return
}

func (r Schema) checkDCOID(s antlr4512.DITContentRule) (err error) {
	if soc := r.ObjectClasses().Get(s.OID); soc.IsZero() {
		err = mkerr(ErrObjectClassNotFound.Error() + `( superclass: ` + s.OID + `)`)
	} else if lup := r.DITContentRules().get(s.OID); !lup.IsZero() {
		// fail attempts to marshal a duplicate definition.
		err = ErrDuplicateDef
	}

	return
}

func (r Schema) marshalDC(s antlr4512.DITContentRule) (def DITContentRule, err error) {
	if !isNumericOID(s.OID) {
		err = ErrMissingNumericOID
		return
	}

	_def := newDITContentRule()

	if err = r.checkDCOID(s); err != nil {
		return
	}
	soc := r.ObjectClasses().Get(s.OID)

	_def.OID = soc
	_def.Desc = s.Desc
	_def.Obsolete = s.Obsolete
	_def.schema = r
	_def.Extensions.setDefinition(DITContentRule{_def})

	sortL := r.Options().Positive(SortLists)

	for _, must := range sortList(s.Must, sortL) {
		m := r.AttributeTypes().get(must)
		if m.IsZero() {
			err = mkerr("Unknown AttributeType for MUST clause: " + must)
			return
		}
		_def.Must.push(m)
	}

	for _, may := range sortList(s.May, sortL) {
		m := r.AttributeTypes().get(may)
		if m.IsZero() {
			err = mkerr("Unknown AttributeType for MAY clause: " + may)
			return
		}
		_def.May.push(m)
	}

	for _, not := range sortList(s.Not, sortL) {
		m := r.AttributeTypes().get(not)
		if m.IsZero() {
			err = mkerr("Unknown AttributeType for NOT clause: " + not)
			return
		}
		_def.Not.push(m)
	}

	for _, aux := range sortList(s.Aux, sortL) {
		m := r.ObjectClasses().get(aux)
		if m.IsZero() {
			err = mkerr("Unknown ObjectClass for AUX clause: " + aux)
			return
		}
		_def.Aux.push(m)
	}

	// Marshal our extensions, taking our sorting preference into account
	marshalExt(s.Extensions, _def.Extensions,
		r.Options().Positive(SortExtensions))

	for _, name := range s.Name {
		_def.Name.push(name)
	}

	str, err := _def.prepareString() // perform one-time text/template op
	if err == nil {
		// Save the stringer
		_def.stringer = func() string {
			// Return a preserved value.
			return str
		}

		def = DITContentRule{_def}
		if !def.Compliant() {
			err = ErrDefNonCompliant
		}
	}

	return
}

func (r Schema) incorporateNF(s antlr4512.NameForms) (err error) {
	for i := 0; i < len(s) && err == nil; i++ {
		var def NameForm
		if def, err = r.marshalNF(s[i]); err == nil {
			r.NameForms().push(def)
		}
	}

	return
}

func (r Schema) marshalNF(s antlr4512.NameForm) (def NameForm, err error) {
	if !isNumericOID(s.OID) {
		err = ErrMissingNumericOID
		return
	}

	// fail attempts to marshal a duplicate definition.
	if lup := r.NameForms().get(s.OID); !lup.IsZero() {
		err = ErrDuplicateDef
		return
	}

	_def := newNameForm()
	_def.OID = s.OID
	_def.Desc = s.Desc
	_def.Obsolete = s.Obsolete
	_def.schema = r
	_def.Extensions.setDefinition(NameForm{_def})

	sortL := r.Options().Positive(SortLists)

	oc := r.ObjectClasses().get(s.OC)
	if oc.IsZero() {
		err = mkerr(ErrObjectClassNotFound.Error() + `( structural: ` + s.OC + `)`)
		return
	}
	_def.Structural = oc

	for _, must := range sortList(s.Must, sortL) {
		m := r.AttributeTypes().get(must)
		if m.IsZero() {
			err = mkerr("Unknown AttributeType for MUST clause: " + must)
			return
		}
		_def.Must.push(m)
	}

	for _, may := range sortList(s.May, sortL) {
		m := r.AttributeTypes().get(may)
		if m.IsZero() {
			err = mkerr("Unknown AttributeType for MAY clause: " + may)
			return
		}
		_def.May.push(m)
	}

	// Marshal our extensions, taking our sorting preference into account
	marshalExt(s.Extensions, _def.Extensions,
		r.Options().Positive(SortExtensions))

	for _, name := range s.Name {
		_def.Name.push(name)
	}

	str, err := _def.prepareString() // perform one-time text/template op
	if err == nil {
		// Save the stringer
		_def.stringer = func() string {
			// Return a preserved value.
			return str
		}
		def = NameForm{_def}
	}

	return
}

func (r Schema) incorporateDS(s antlr4512.DITStructureRules) (err error) {
	for i := 0; i < len(s) && err == nil; i++ {
		var def DITStructureRule
		if def, err = r.marshalDS(s[i]); err == nil {
			r.DITStructureRules().push(def)
		}
	}

	return
}

func (r Schema) marshalDS(s antlr4512.DITStructureRule) (def DITStructureRule, err error) {
	var ruleid uint
	var ok bool
	if ruleid, ok = atoui(s.ID); !ok {
		err = mkerr("Invalid structure rule ID " + s.ID)
		return
	}

	// fail attempts to marshal a duplicate definition.
	if lup := r.DITStructureRules().get(ruleid); !lup.IsZero() {
		err = ErrDuplicateDef
		return
	}

	_def := newDITStructureRule()
	_def.ID = ruleid
	_def.Desc = s.Desc
	_def.Obsolete = s.Obsolete
	_def.schema = r
	_def.Extensions.setDefinition(DITStructureRule{_def})

	sortL := r.Options().Positive(SortLists)

	nf := r.NameForms().get(s.Form)
	if nf.IsZero() {
		err = mkerr(ErrNameFormNotFound.Error() + `(` + s.Form + `)`)
		return
	}
	_def.Form = nf

	for _, sup := range sortList(s.SuperRules, sortL) {
		m := r.DITStructureRules().get(sup)
		if m.IsZero() {
			err = mkerr("Unknown rule for SUP clause: " + sup)
			return
		}
		_def.SuperRules.push(m)
	}

	// Marshal our extensions, taking our sorting preference into account
	marshalExt(s.Extensions, _def.Extensions,
		r.Options().Positive(SortExtensions))

	for _, name := range s.Name {
		_def.Name.push(name)
	}

	str, err := _def.prepareString() // perform one-time text/template op
	if err == nil {
		// Save the stringer
		_def.stringer = func() string {
			// Return a preserved value.
			return str
		}
		def = DITStructureRule{_def}
	}

	return
}

func sortList(list []string, sortBySlice bool) []string {
	if !sortBySlice {
		return list
	}

	sort.Strings(list)
	return list
}

/*
marshalExt funnels mext into ext with or without sorting.
*/
func marshalExt(mext map[string][]string, ext Extensions, sortByXStr bool) {
	if !sortByXStr {
		for k, v := range mext {
			ext.Set(k, v...)
		}
		return
	}

	// Get a list of keys from mext
	keys := make([]string, 0, len(mext))
	for k := range mext {
		keys = append(keys, k)
	}

	// Sort keys alphabetically
	sort.Strings(keys)

	// Use new sorting scheme to influence
	// Extension initialization within ext.
	for _, k := range keys {
		ext.Set(k, mext[k]...)
	}
}

func handleMacro(r Schema, m []string, o string) (resv string) {
	if len(o) > 0 {
		resv = o
		return
	}

	if len(m) == 2 && len(o) == 0 {
		mc, found := r.Macros().Resolve(m[0])
		if found {
			resv = mc + `.` + m[1]
		}
	}

	return
}

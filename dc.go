package schemax

/*
NewDITContentRules initializes a new [DITContentRules] instance.
*/
func NewDITContentRules() DITContentRules {
	r := DITContentRules(newCollection(``))
	r.cast().SetPushPolicy(r.canPush)

	return r
}

/*
Description returns the underlying (optional) descriptive text
assigned to the receiver instance.
*/
func (r DITContentRule) Description() (desc string) {
	if !r.IsZero() {
		desc = r.dITContentRule.Desc
	}

	return
}

/*
E returns the underlying error instance.
*/
func (r DITContentRule) E() (err error) {
	if !r.IsZero() {
		err = r.dITContentRule.err
	} else {
		err = ErrNilReceiver
	}

	return
}

/*
DITContentRules returns the [DITContentRules] instance from within
the receiver instance.
*/
func (r Schema) DITContentRules() (dcs DITContentRules) {
	slice, _ := r.cast().Index(dITContentRulesIndex)
	dcs, _ = slice.(DITContentRules)
	return
}

/*
NewDITContentRule initializes and returns a new instance of [DITContentRule],
ready for manual assembly.  This method need not be used when creating
new [DITContentRule] instances by way of parsing, as that is handled on an
internal basis.

Use of this method does NOT automatically push the return instance into
the [Schema.DITContentRules] stack; this is left to the user.

Unlike the package-level [NewDITContentRule] function, this method will
automatically reference its originating [Schema] instance (the receiver).
This negates the need for manual use of the [DITContentRule.SetSchema]
method.

This is the recommended means of creating a new [DITContentRule] instance
wherever a single [Schema] is being used, which represents most use cases.
*/
func (r Schema) NewDITContentRule() DITContentRule {
	return NewDITContentRule().SetSchema(r)
}

/*
NewDITContentRule initializes and returns a new instance of [DITContentRule],
ready for manual assembly.  This method need not be used when creating
new [DITContentRule] instances by way of parsing, as that is handled on an
internal basis.

Use of this function does not automatically reference the "parent" [Schema]
instance, leaving it up to the user to invoke the [DITContentRule.SetSchema]
method manually.

When interacting with a single [Schema] instance, which represents most use
cases, use of the [Schema.NewDITContentRule] method is PREFERRED over use of
this package-level function.

However certain migration efforts, schema audits and other such activities
may require distinct associations of [DITContentRule] instances with specific
[Schema] instances. Use of this function allows the user to specify the
appropriate [Schema] instance at a later point for a specific instance of
a [DITContentRule] instance.
*/
func NewDITContentRule() DITContentRule {
	dc := DITContentRule{newDITContentRule()}
	dc.dITContentRule.Extensions.setDefinition(dc)
	return dc
}

func newDITContentRule() *dITContentRule {
	return &dITContentRule{
		Name:       NewName(),
		Aux:        NewObjectClassOIDList(`AUX`),
		Must:       NewAttributeTypeOIDList(`MUST`),
		May:        NewAttributeTypeOIDList(`MAY`),
		Not:        NewAttributeTypeOIDList(`NOT`),
		Extensions: NewExtensions(),
	}
}

/*
SetData assigns x to the receiver instance. This is a general-use method and has no
specific intent beyond convenience. The contents may be subsequently accessed via the
[DITContentRule.Data] method.

This is a fluent method.
*/
func (r DITContentRule) SetData(x any) DITContentRule {
	if !r.IsZero() {
		r.dITContentRule.setData(x)
	}

	return r
}

func (r *dITContentRule) setData(x any) {
	r.data = x
}

/*
Data returns the underlying value (x) assigned to the receiver's data storage field. Data
can be set within the receiver instance by way of the [DITContentRule.SetData] method.
*/
func (r DITContentRule) Data() (x any) {
	if !r.IsZero() {
		x = r.dITContentRule.data
	}

	return
}

/*
SetSchema assigns an instance of [Schema] to the receiver instance.  This allows
internal verification of certain actions without the need for user input of
an instance of [Schema] manually at each juncture.

Note that the underlying [Schema] instance is automatically set when creating
instances of this type by way of parsing, as well as if the receiver instance
was initialized using the [Schema.NewDITContentRule] method.

This is a fluent method.
*/
func (r DITContentRule) SetSchema(schema Schema) DITContentRule {
	if !r.IsZero() {
		r.dITContentRule.setSchema(schema)
	}

	return r
}

func (r *dITContentRule) setSchema(schema Schema) {
	r.schema = schema
}

/*
Schema returns the [Schema] instance associated with the receiver instance.
*/
func (r DITContentRule) Schema() (s Schema) {
	if !r.IsZero() {
		s = r.dITContentRule.getSchema()
	}

	return
}

func (r *dITContentRule) getSchema() (s Schema) {
	if r != nil {
		s = r.schema
	}

	return
}

/*
Inventory returns an instance of [Inventory] which represents the current
inventory of [DITContentRule] instances within the receiver.

The keys are numeric OIDs, while the values are zero (0) or more string
slices, each representing a name by which the definition is known.
*/
func (r DITContentRules) Inventory() (inv Inventory) {
	inv = make(Inventory, 0)
	for i := 0; i < r.len(); i++ {
		def := r.index(i)
		inv[def.NumericOID()] = def.Names().List()
	}

	return
}

/*
Maps returns slices of [DefinitionMap] instances.
*/
func (r DITContentRules) Maps() (defs DefinitionMaps) {
	defs = make(DefinitionMaps, r.Len())
	for i := 0; i < r.Len(); i++ {
		defs[i] = r.Index(i).Map()
	}

	return
}

/*
Map marshals the receiver instance into an instance of
[DefinitionMap].
*/
func (r DITContentRule) Map() (def DefinitionMap) {
	if r.IsZero() {
		return
	}

	var auxs []string
	for i := 0; i < r.Aux().Len(); i++ {
		m := r.Aux().Index(i)
		auxs = append(auxs, m.OID())
	}

	var musts []string
	for i := 0; i < r.Must().Len(); i++ {
		m := r.Must().Index(i)
		musts = append(musts, m.OID())
	}

	var mays []string
	for i := 0; i < r.May().Len(); i++ {
		m := r.May().Index(i)
		mays = append(mays, m.OID())
	}

	var nots []string
	for i := 0; i < r.Not().Len(); i++ {
		m := r.Not().Index(i)
		nots = append(nots, m.OID())
	}

	def = make(DefinitionMap, 0)
	def[`NUMERICOID`] = []string{r.NumericOID()}
	def[`NAME`] = r.Names().List()
	def[`DESC`] = []string{r.Description()}
	def[`OBSOLETE`] = []string{bool2str(r.Obsolete())}
	def[`AUX`] = auxs
	def[`MUST`] = musts
	def[`MAY`] = mays
	def[`NOT`] = nots
	def[`TYPE`] = []string{r.Type()}
	def[`RAW`] = []string{r.String()}

	// copy our extensions from receiver r
	// into destination def.
	def = mapTransferExtensions(r, def)

	// Clean up any empty fields
	def.clean()

	return
}

/*
Type returns the string literal "dITContentRule".
*/
func (r DITContentRule) Type() string {
	return `dITContentRule`
}

func (r dITContentRule) Type() string {
	return `dITContentRule`
}

/*
Type returns the string literal "dITContentRules".
*/
func (r DITContentRules) Type() string {
	return `dITContentRules`
}

/*
Len returns the current integer length of the receiver instance.
*/
func (r DITContentRules) Len() int {
	return r.len()
}

func (r DITContentRules) len() int {
	return r.cast().Len()
}

/*
String is a stringer method that returns the string representation
of the receiver instance.
*/
func (r DITContentRules) String() string {
	return r.cast().String()
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r DITContentRules) IsZero() bool {
	return r.cast().IsZero()
}

/*
Index returns the instance of [DITContentRule] found within the
receiver stack instance at index N.  If no instance is found at
the index specified, a zero [DITContentRule] instance is returned.
*/
func (r DITContentRules) Index(idx int) DITContentRule {
	return r.index(idx)
}

func (r DITContentRules) index(idx int) (dc DITContentRule) {
	slice, found := r.cast().Index(idx)
	if found {
		if _dc, ok := slice.(DITContentRule); ok {
			dc = _dc
		}
	}

	return
}

func (r DITContentRule) macro() (m []string) {
	if !r.IsZero() {
		m = r.dITContentRule.Macro
	}

	return
}

func (r DITContentRules) oIDsStringer(_ ...any) (present string) {
	var _present []string
	for i := 0; i < r.len(); i++ {
		_present = append(_present, r.index(i).OID())
	}

	switch len(_present) {
	case 0:
		break
	case 1:
		present = _present[0]
	default:
		padchar := string(rune(32))
		if !r.cast().IsPadded() {
			padchar = ``
		}

		joined := join(_present, padchar+`$`+padchar)
		present = "(" + padchar + joined + padchar + ")"
	}

	return
}

/*
Identifier returns the output of [DITContentRule.OID] and is merely used
for an interface-friendly means of obtaining the principal identifier
of the receiver instance.
*/
func (r DITContentRule) Identifier() (ident string) { return r.OID() }

/*
IsIdentifiedAs returns a Boolean value indicative of whether id matches
either the numericOID or descriptor of the receiver instance.  Case is
not significant in the matching process.
*/
func (r DITContentRule) IsIdentifiedAs(id string) (ident bool) {
	if !r.IsZero() {
		ident = id == r.NumericOID() || r.dITContentRule.Name.contains(id)
	}

	return
}

/*
Obsolete returns a Boolean value indicative of definition obsolescence.
*/
func (r DITContentRule) Obsolete() (o bool) {
	if !r.IsZero() {
		o = r.dITContentRule.Obsolete
	}

	return
}

/*
OID returns the string representation of an OID -- which is either a
numeric OID or descriptor -- that is held by the receiver instance.
*/
func (r DITContentRule) OID() (oid string) {
	if !r.IsZero() {
		oid = r.NumericOID() // default
		if r.dITContentRule.Name.len() > 0 {
			oid = r.dITContentRule.Name.index(0)
		}
	}

	return
}

/*
Aux returns an [ObjectClasses] containing zero (0) or more auxiliary
[ObjectClass] definitions for use within entries governed by this rule.
*/
func (r DITContentRule) Aux() (aux ObjectClasses) {
	if !r.IsZero() {
		aux = r.dITContentRule.Aux
	}

	return
}

/*
SetAux assigns the provided input [ObjectClass] instance(s) -- which must
be AUXILIARY via the [AuxiliaryKind] constant -- to the receiver's AUX clause.

This is a fluent method.
*/
func (r DITContentRule) SetAux(m ...any) DITContentRule {
	if !r.IsZero() {
		r.dITContentRule.setAux(m...)
	}

	return r
}

func (r *dITContentRule) setAux(m ...any) {
	var err error
	for i := 0; i < len(m) && err == nil; i++ {
		var oc ObjectClass
		switch tv := m[i].(type) {
		case string:
			oc = r.schema.ObjectClasses().get(tv)
		case ObjectClass:
			oc = tv
		default:
			err = ErrInvalidType
		}

		if err == nil && oc.Kind() == AuxiliaryKind {
			r.Aux.Push(oc)
		}
	}

	if err != nil {
		r.err = err
	}
}

/*
Compliant returns a Boolean value indicative of every [DITContentRule]
returning a compliant response from the [DITContentRule.Compliant] method.
*/
func (r DITContentRules) Compliant() bool {
	var act int
	for i := 0; i < r.Len(); i++ {
		if r.Index(i).Compliant() {
			act++
		}
	}

	return act == r.Len()
}

/*
StructuralClass returns the STRUCTURAL [ObjectClass] set within the
receiver instance, or a zero instance if unset.

This method is essentially the inverse of [ObjectClass.EnforcedBy].
*/
func (r DITContentRule) StructuralClass() (soc ObjectClass) {
	if !r.IsZero() {
		soc = r.dITContentRule.OID
	}

	return
}

/*
Compliant returns a Boolean value indicative of the receiver being fully
compliant per the required clauses of [§ 4.1.6 of RFC 4512]:

  - Numeric OID must relate to a predefined [ObjectClass] in the associated [Schema] instance
  - [ObjectClass] referenced by OID must be STRUCTURAL
  - [ObjectClass] referenced by OID must be COMPLIANT itself
  - Collective [AttributeType] instances are permitted, but not verified as they are never present in any [ObjectClass]
  - MUST, MAY and NOT clause [AttributeType] instances are limited to those present in the [ObjectClass] super chain
  - No conflicting clause values (e.g.: cannot forbid (NOT) a required type (MUST)), with emphasis on related [DITStructureRule] FORM ([NameForm]) instances.

[§ 4.1.6 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.6
*/
func (r DITContentRule) Compliant() bool {
	if r.IsZero() {
		return false
	}

	structural := r.StructuralClass()
	if !structural.Compliant() || structural.Kind() != StructuralKind {
		return false
	}

	// verify all AUX clause members are valid
	musts := NewAttributeTypeOIDList() // from STRUCTURAL class MUST clause
	mays := NewAttributeTypeOIDList()  // from STRUCTURAL class MAY clause

	smust := structural.Must()
	for i := 0; i < smust.Len(); i++ {
		musts.Push(smust.Index(i))
	}

	smay := structural.May()
	for i := 0; i < smay.Len(); i++ {
		mays.Push(smay.Index(i))
	}

	// add ABSTRACT requirements/allowances
	abstracts := structural.SuperClasses()
	for i := 0; i < abstracts.Len(); i++ {
		if abs := abstracts.Index(i); abs.Kind() == 1 {
			ms := abs.Must()
			my := abs.May()
			for j := 0; j < ms.Len(); j++ {
				musts.Push(ms.Index(j))
			}
			for j := 0; j < my.Len(); j++ {
				mays.Push(my.Index(j))
			}
		}
	}

	for _, boo := range []bool{
		r.auxComply(musts, mays),
		r.notComply(musts, mays),
		r.mustComply(musts, mays),
		r.mayComply(musts, mays),
	} {
		if !boo {
			return boo
		}
	}

	// if any dITStructureRule definitions exist,
	// make sure they don't produce a MUST/NOT
	// conflict.
	ok := r.dsrComply(structural)
	if ok {
		r.dITContentRule.err = nil
	}

	return ok
}

func (r DITContentRule) dsrComply(structural ObjectClass) bool {
	// In the event a matching dITStructureRule exists
	// whose FORM bears the same structural class OID
	// as the receiver, make sure the rules do not
	// produce a MUST/NOT conflict.
	dsr := r.schema.DITStructureRules()
	for i := 0; i < dsr.Len(); i++ {
		form := dsr.Index(i).Form()
		if form.OC().NumericOID() == structural.NumericOID() {
			// We found a matching dITStructureRule. We want to
			// be sure that none of its nameForm's MUST clause
			// members are present in the receiver's NOT clause.
			clause := form.Must()
			for i := 0; i < clause.Len(); i++ {
				if r.Not().Contains(clause.Index(i).OID()) {
					return false
				}
			}
			break // per X.501, only one rule applies per schema
		}
	}

	return true
}

func (r DITContentRule) auxComply(must, may AttributeTypes) bool {
	var aux ObjectClasses = r.Aux()
	for i := 0; i < aux.Len(); i++ {
		aoc := aux.Index(i)
		if !aoc.Compliant() || aoc.Kind() != AuxiliaryKind {
			return false
		}

		_must := aoc.Must()
		for j := 0; j < _must.Len(); j++ {
			must.Push(_must.Index(j))
		}

		_may := aoc.May()
		for j := 0; j < _may.Len(); j++ {
			may.Push(_may.Index(j))
		}
	}

	return true
}

func (r DITContentRule) notComply(must, may AttributeTypes) bool {
	rnots := r.Not()
	for i := 0; i < rnots.Len(); i++ {
		no := rnots.Index(i)
		if no.Collective() {
			continue
		}
		if must.Contains(no.OID()) && !may.Contains(no.OID()) {
			// Cannot preclude a MUST or MAY
			return false
		}
	}

	return true
}

func (r DITContentRule) mustComply(must, may AttributeTypes) bool {
	rmusts := r.Must()
	for i := 0; i < rmusts.Len(); i++ {
		cmust := rmusts.Index(i)
		if cmust.Collective() {
			continue
		}
		if !must.Contains(cmust.OID()) && !may.Contains(cmust.OID()) {
			// can't require an unauthorized attribute
			return false
		}
	}

	return true
}

func (r DITContentRule) mayComply(must, may AttributeTypes) bool {
	rmays := r.May()
	for i := 0; i < rmays.Len(); i++ {
		cmay := rmays.Index(i)
		if cmay.Collective() {
			continue
		}
		if must.Contains(cmay.OID()) {
			// can't require a MAY
			return false
		}
		if !may.Contains(cmay.OID()) {
			// can't allow an unauthorized attribute
			return false
		}
	}

	return true
}

/*
Marshal returns an error following an attempt to marshal the contents of
def, which may be either a [DefinitionMap] or map[string]any instance.

The receiver instance must be initialized prior to use of this method
using the [Schema.NewDITContentRule] method.
*/
func (r DITContentRule) Marshal(def any) error {
	m, err := getMarshalMap(r, def)
	if err != nil {
		return err
	}

	for k, v := range m {
		switch key := uc(k); key {
		case `NAME`:
			switch tv := v.(type) {
			case string:
				r.SetName(tv)
			case []string:
				r.SetName(tv...)
			}
		case `NUMERICOID`, `DESC`:
			z := map[string]func(string) DITContentRule{
				`DESC`:       r.SetDescription,
				`NUMERICOID`: r.SetNumericOID,
			}
			switch tv := v.(type) {
			case string:
				z[k](tv)
			case []string:
				z[k](tv[0])
			}
		case `OBSOLETE`:
			r.marshalBoolean(v)
		case `MUST`, `MAY`, `NOT`, `AUX`:
			r.marshalMulti(key, v)
		default:
			r.marshalExt(key, v)
		}
	}

	if !r.Compliant() {
		return ErrDefNonCompliant
	}
	r.SetStringer()

	return nil
}

func (r DITContentRule) marshalBoolean(v any) {
	switch tv := v.(type) {
	case string:
		if eq(tv, `TRUE`) {
			r.SetObsolete()
		}
	case []string:
		if eq(tv[0], `TRUE`) {
			r.SetObsolete()
		}
	case bool:
		if tv {
			r.SetObsolete()
		}
	}
}

func (r DITContentRule) marshalExt(key string, v any) {
	if hasPfx(key, `X-`) {
		switch tv := v.(type) {
		case string:
			r.SetExtension(key, tv)
		case []string:
			r.SetExtension(key, tv...)
		}
	}
}

func (r DITContentRule) marshalMulti(k string, v any) {
	z := map[string]func(...any) DITContentRule{
		`MUST`: r.SetMust,
		`MAY`:  r.SetMay,
		`NOT`:  r.SetNot,
		`AUX`:  r.SetAux,
	}
	switch tv := v.(type) {
	case []string:
		for i := 0; i < len(tv); i++ {
			z[k](tv[i])
		}
	case string:
		z[k](tv)
	}
}

/*
Parse returns an error following an attempt to parse raw into the receiver
instance.

Note that the receiver MUST possess a [Schema] reference prior to the execution
of this method.

Also note that successful execution of this method does NOT automatically push
the receiver into any [DITContentRules] stack, nor does it automatically execute
the [DITContentRule.SetStringer] method, leaving these tasks to the user.  If the
automatic handling of these tasks is desired, see the [Schema.ParseDITContentRule]
method as an alternative.
*/
func (r DITContentRule) Parse(raw string) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
		return
	}

	if r.getSchema().IsZero() {
		err = ErrNilSchemaRef
		return
	}

	err = r.dITContentRule.parse(raw)

	return
}

func (r *dITContentRule) parse(raw string) error {
	// parseLS wraps the antlr4512 DITContentRule parser/lexer
	mp, err := parseDC(raw)
	if err == nil {
		// We received the parsed data from ANTLR (mp).
		// Now we need to marshal it into the receiver.
		var def DITContentRule
		if def, err = r.schema.marshalDC(mp); err == nil {
			r.OID = def.dITContentRule.OID
			_r := DITContentRule{r}
			_r.replace(def)
		}
	}

	return err
}

/*
SetName assigns the provided names to the receiver instance.

Name instances must conform to RFC 4512 descriptor format but
need not be quoted.

This is a fluent method.
*/
func (r DITContentRule) SetName(x ...string) DITContentRule {
	if !r.IsZero() {
		r.dITContentRule.setName(x...)
	}

	return r
}

func (r *dITContentRule) setName(x ...string) {
	b4 := r.Name.Len()
	for i := 0; i < len(x); i++ {
		r.Name.Push(x[i])
	}

	if r.Name.Len()-len(x) != b4 {
		r.err = ErrInvalidNames
	}
}

/*
SetObsolete sets the receiver instance to OBSOLETE if not already set. Note that obsolescence cannot be unset.

This is a fluent method.
*/
func (r DITContentRule) SetObsolete() DITContentRule {
	if !r.IsZero() {
		r.dITContentRule.setObsolete()
	}

	return r
}

func (r *dITContentRule) setObsolete() {
	if !r.Obsolete {
		r.Obsolete = true
	}
}

/*
SetDescription parses desc into the underlying DESC clause within the
receiver instance.  Although a RFC 4512-compliant QuotedString is
required, the outer single-quotes need not be specified literally.

This is a fluent method.
*/
func (r DITContentRule) SetDescription(desc string) DITContentRule {
	if !r.IsZero() {
		r.dITContentRule.setDescription(desc)
	}

	return r
}

func (r *dITContentRule) setDescription(desc string) {
	if len(desc) < 3 {
		return
	}

	if rune(desc[0]) == rune(39) {
		desc = desc[1:]
	}

	if rune(desc[len(desc)-1]) == rune(39) {
		desc = desc[:len(desc)-1]
	}

	r.Desc = desc

	return
}

/*
SetExtension assigns key x to value xstrs within the receiver's underlying
[Extensions] instance.

This is a fluent method.
*/
func (r DITContentRule) SetExtension(x string, xstrs ...string) DITContentRule {
	if !r.IsZero() {
		r.dITContentRule.setExtension(x, xstrs...)
	}

	return r
}

func (r *dITContentRule) setExtension(x string, xstrs ...string) {
	r.Extensions.Set(x, xstrs...)
}

/*
SetNumericOID allows the manual assignment of a numeric OID to the
receiver instance if the following are all true:

  - The input id value is a syntactically valid numeric OID
  - The receiver does not already possess a numeric OID

This is a fluent method.
*/
func (r DITContentRule) SetNumericOID(id string) DITContentRule {
	if !r.IsZero() {
		r.dITContentRule.setNumericOID(id)
	}

	return r
}

func (r *dITContentRule) setNumericOID(id string) {
	if isNumericOID(id) {
		// only set an OID when the receiver
		// lacks one (iow: no modifications)
		if r.OID.IsZero() {
			oc := r.schema.ObjectClasses().Get(id)
			if oc.Kind() == StructuralKind {
				r.OID = oc
			} else {
				r.err = ErrIncompatStructuralClass
			}
		}
	} else {
		r.err = ErrInvalidOID
	}

	return
}

/*
Must returns an [AttributeTypes] containing zero (0) or more required
[AttributeType] definitions for use within entries governed by this rule.
*/
func (r DITContentRule) Must() (must AttributeTypes) {
	if !r.IsZero() {
		must = r.dITContentRule.Must
	}

	return
}

/*
SetMust assigns the provided input [AttributeType] instance(s) to the
receiver's MUST clause.

This is a fluent method.
*/
func (r DITContentRule) SetMust(m ...any) DITContentRule {
	if !r.IsZero() {
		r.dITContentRule.setMust(m...)
	}

	return r
}

func (r *dITContentRule) setMust(m ...any) {
	var err error
	for i := 0; i < len(m) && err == nil; i++ {
		var at AttributeType
		switch tv := m[i].(type) {
		case string:
			at = r.schema.AttributeTypes().get(tv)
		case AttributeType:
			at = tv
		default:
			err = ErrInvalidType
		}

		if err == nil && !at.IsZero() {
			r.Must.Push(at)
		}
	}

	if err != nil {
		r.err = err
	}
}

/*
May returns an AttributeTypes containing zero (0) or more allowed
[AttributeType] definitions for use within entries governed by this rule.
*/
func (r DITContentRule) May() (may AttributeTypes) {
	if !r.IsZero() {
		may = r.dITContentRule.May
	}

	return
}

/*
SetMay assigns the provided input [AttributeType] instance(s) to the
receiver's MAY clause.

This is a fluent method.
*/
func (r DITContentRule) SetMay(m ...any) DITContentRule {
	if !r.IsZero() {
		r.dITContentRule.setMay(m...)
	}

	return r
}

func (r *dITContentRule) setMay(m ...any) {
	var err error
	for i := 0; i < len(m) && err == nil; i++ {
		var at AttributeType
		switch tv := m[i].(type) {
		case string:
			at = r.schema.AttributeTypes().get(tv)
		case AttributeType:
			at = tv
		default:
			err = ErrInvalidType
		}

		if err == nil && !at.IsZero() {
			r.May.Push(at)
		}
	}

	if err != nil {
		r.err = err
	}
}

/*
Not returns an AttributeTypes containing zero (0) or more [AttributeType]
definitions disallowed for use within entries governed by this rule.
*/
func (r DITContentRule) Not() (not AttributeTypes) {
	if !r.IsZero() {
		not = r.dITContentRule.Not
	}

	return
}

/*
SetNot assigns the provided input [AttributeType] instance(s) to the
receiver's NOT clause.

This is a fluent method.
*/
func (r DITContentRule) SetNot(m ...any) DITContentRule {
	if !r.IsZero() {
		r.dITContentRule.setNot(m...)
	}

	return r
}

func (r *dITContentRule) setNot(m ...any) {
	var err error
	for i := 0; i < len(m) && err == nil; i++ {
		var at AttributeType
		switch tv := m[i].(type) {
		case string:
			at = r.schema.AttributeTypes().get(tv)
		case AttributeType:
			at = tv
		default:
			err = ErrInvalidType
		}

		if err == nil && !at.IsZero() {
			r.Not.Push(at)
		}
	}

	if err != nil {
		r.err = err
	}
}

/*
SetStringer allows the assignment of an individual [Stringer] function or
method to all [DITContentRule] slices within the receiver stack instance.

Input of zero (0) variadic values, or an explicit nil, will overwrite all
preexisting stringer functions with the internal closure default, which is
based upon a one-time use of the [text/template] package by all receiver
slice instances.

Input of a non-nil closure function value will overwrite all preexisting
stringers.

This is a fluent method and may be used multiple times.
*/
func (r DITContentRules) SetStringer(function ...Stringer) DITContentRules {
	for i := 0; i < r.Len(); i++ {
		def := r.Index(i)
		def.SetStringer(function...)
	}

	return r
}

/*
SetStringer allows the assignment of an individual [Stringer] function
or method to the receiver instance.

Input of zero (0) variadic values, or an explicit nil, will overwrite any
preexisting stringer function with the internal closure default, which is
based upon a one-time use of the [text/template] package by the receiver
instance.

Input of a non-nil closure function value will overwrite any preexisting
stringer.

This is a fluent method and may be used multiple times.
*/
func (r DITContentRule) SetStringer(function ...Stringer) DITContentRule {
	if r.Compliant() {
		r.dITContentRule.setStringer(function...)
	}

	return r
}

func (r *dITContentRule) setStringer(function ...Stringer) {
	var stringer Stringer
	if len(function) > 0 {
		stringer = function[0]
	}

	if stringer == nil {
		str, err := r.prepareString() // perform one-time text/template op
		if err == nil {
			// Save the stringer
			r.stringer = func() string {
				// Return a preserved value.
				return str
			}
		} else {
			r.err = err
		}
		return
	}

	r.stringer = stringer
}

/*
XOrigin returns an instance of [DITContentRules] containing only definitions
which bear the X-ORIGIN value of x. Case is not significant in the matching
process, nor is whitespace (e.g.: RFC 4517 vs. RFC4517).
*/
func (r DITContentRules) XOrigin(x string) (defs DITContentRules) {
	defs = NewDITContentRules()
	for i := 0; i < r.Len(); i++ {
		def := r.Index(i)
		if xo, found := def.Extensions().Get(`X-ORIGIN`); found {
			if xo.Contains(x) {
				defs.push(def)
			}
		}
	}

	return
}

/*
String is a stringer method that returns the string representation
of the receiver instance.
*/
func (r DITContentRule) String() (dcr string) {
	if !r.IsZero() {
		if r.dITContentRule.stringer != nil {
			dcr = r.dITContentRule.stringer()
		}
	}

	return
}

func (r *dITContentRule) prepareString() (str string, err error) {
	buf := newBuf()
	// Create a new template instance bearing the "Type" value
	// string literal as its name.  Declare custom templating
	// functions enveloped within a template.FuncMap instance.
	t := newTemplate(r.Type()).
		Funcs(funcMap(map[string]any{
			// ExtensionSet used by all definitions for
			// "X-" extensions (e.g.: X-ORIGIN, X-SUBSTR)
			`ExtensionSet`: r.Extensions.tmplFunc,
			// StructuralOID refers to the OID shared with the
			// structural OC upon which this rule is based
			`StructuralOID`: r.OID.NumericOID,
			// AuxLen refers to the number of AUXILIARY ObjectClasses
			// present within the rule's AUX clause.
			`AuxLen`: r.Aux.len,
			// MayLen refers to the number of AttributeTypes
			// present within the rule's MAY clause.
			`MayLen`: r.May.len,
			// MustLen refers to the number of AttributeTypes
			// present within the rule's MUST clause.
			`MustLen`: r.Must.len,
			// NotLen refers to the number of AttributeTypes
			// present within the rule's NOT clause.
			`NotLen`: r.Not.len,
			// Obsolete indicates definition obsolescence
			// in the form of a Boolean value.
			`Obsolete`: func() bool { return r.Obsolete },
		}))

	// Parse raw template into *template.Template instance
	// and ensure the raw directives represent legal, well
	// formed and error-free templating instructions.
	if t, err = t.Parse(dITContentRuleTmpl); err == nil {
		// Execute the now-verified template, funnel our
		// needed objects and settings values into an
		// anonymous struct instance.
		if err = t.Execute(buf, struct {
			Definition *dITContentRule
			HIndent    string
		}{
			Definition: r,
			HIndent:    hindent(r.schema.Options().Positive(HangingIndents)),
		}); err == nil {
			// Dump our templated output from
			// the *bytes.Buffer instance into
			// the return (string) value.
			str = buf.String()
		}
	}

	if err != nil {
		r.err = err
	}

	return
}

/*
Name returns the string form of the principal name of the receiver instance, if set.
*/
func (r DITContentRule) Name() (id string) {
	if !r.IsZero() {
		id = r.dITContentRule.Name.index(0)
	}

	return
}

/*
Names returns the underlying instance of [QuotedDescriptorList] from within
the receiver.
*/
func (r DITContentRule) Names() (names QuotedDescriptorList) {
	if !r.IsZero() {
		names = r.dITContentRule.Name
	}
	return
}

/*
Extensions returns the [Extensions] instance -- if set -- within
the receiver.
*/
func (r DITContentRule) Extensions() (e Extensions) {
	if !r.IsZero() {
		e = r.dITContentRule.Extensions
	}

	return
}

/*
NumericOID returns the string representation of the numeric OID
held by the receiver instance.
*/
func (r DITContentRule) NumericOID() (noid string) {
	if !r.IsZero() {
		if soid := r.dITContentRule.OID; !soid.IsZero() {
			noid = soid.NumericOID()
		}
	}

	return
}

func (r DITContentRule) setOID(_ string) {}

/*
Push returns an error following an attempt to push a [DITContentRule]
into the receiver stack instance.
*/
func (r DITContentRules) Push(dc any) error {
	return r.push(dc)
}

func (r DITContentRules) push(x any) (err error) {
	switch tv := x.(type) {
	case DITContentRule:
		if !tv.Compliant() {
			err = ErrDefNonCompliant
			break
		}
		r.cast().Push(tv)
	default:
		err = ErrInvalidType
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r DITContentRule) IsZero() bool {
	return r.dITContentRule == nil
}

/*
Contains calls [DITStructureRules.Get] to return a Boolean value indicative of
a successful, non-zero retrieval of a [DITStructureRule] instance -- matching
the provided id -- from within the receiver stack instance.
*/
func (r DITContentRules) Contains(id string) bool {
	return r.contains(id)
}

func (r DITContentRules) contains(id string) bool {
	return !r.get(id).IsZero()
}

/*
Get returns an instance of [DITContentRule] based upon a search for id within
the receiver stack instance.

The return instance, if not nil, was retrieved based upon a textual match of
the principal identifier of a [DITContentRule] and the provided id.

The return instance is nil if no match was made.

Case is not significant in the matching process.
*/
func (r DITContentRules) Get(id string) DITContentRule {
	return r.get(id)
}

func (r DITContentRules) get(id string) (dc DITContentRule) {
	for i := 0; i < r.len() && dc.IsZero(); i++ {
		if _dc := r.index(i); !_dc.IsZero() {
			if _dc.NumericOID() == id {
				dc = _dc
			} else if _dc.dITContentRule.Name.contains(id) {
				dc = _dc
			}
		}
	}

	return
}

// stackage closure func - do not exec directly.
func (r DITContentRules) canPush(x ...any) (err error) {
	if len(x) == 0 {
		return
	}

	for i := 0; i < len(x) && err == nil; i++ {
		dc, ok := x[i].(DITContentRule)
		if !ok || dc.IsZero() {
			err = ErrTypeAssert
		} else if tst := r.get(dc.NumericOID()); !tst.IsZero() {
			err = mkerr(ErrNotUnique.Error() + ": " + dc.Type() + `, ` + dc.NumericOID())
		}
	}

	return
}

/*
Replace overrides the receiver with x. Both must bear an identical
numeric OID and x MUST be compliant.

Note that the relevant [Schema] instance must be configured to allow
definition override by way of the [AllowOverride] bit setting.  See
the [Schema.Options] method for a means of accessing the settings
value.

Note that this method does not reallocate a new pointer instance
within the [DITContentRule] envelope type, thus all references to the
receiver instance within various stacks will be preserved.

This is a fluent method.
*/
func (r DITContentRule) Replace(x DITContentRule) DITContentRule {
	if !r.Schema().Options().Positive(AllowOverride) {
		return r
	}

	if !r.IsZero() && x.Compliant() {
		r.dITContentRule.replace(x)
	}

	return r
}

func (r *dITContentRule) replace(x DITContentRule) {
	if r.OID.NumericOID() == `` {
		r.err = ErrMissingNumericOID
		return
	} else if r.OID.NumericOID() != x.NumericOID() {
		r.err = ErrInvalidOID
		return
	}

	r.OID = x.dITContentRule.OID
	r.Name = x.dITContentRule.Name
	r.Desc = x.dITContentRule.Desc
	r.Obsolete = x.dITContentRule.Obsolete
	r.Must = x.dITContentRule.Must
	r.May = x.dITContentRule.May
	r.Not = x.dITContentRule.Not
	r.Aux = x.dITContentRule.Aux
	r.Extensions = x.dITContentRule.Extensions
	r.data = x.dITContentRule.data
	r.schema = x.dITContentRule.schema
	r.stringer = x.dITContentRule.stringer
	r.data = x.dITContentRule.data
}

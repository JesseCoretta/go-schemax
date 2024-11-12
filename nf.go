package schemax

/*
NewNameForms initializes a new [NameForms] instance.
*/
func NewNameForms() NameForms {
	r := NameForms(newCollection(``))
	r.cast().SetPushPolicy(r.canPush)

	return r
}

/*
NameForms returns the [NameForms] instance from within the receiver
instance.
*/
func (r Schema) NameForms() (nfs NameForms) {
	slice, _ := r.cast().Index(nameFormsIndex)
	nfs, _ = slice.(NameForms)
	return
}

/*
Maps returns slices of [DefinitionMap] instances.
*/
func (r NameForms) Maps() (defs DefinitionMaps) {
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
func (r NameForm) Map() (def DefinitionMap) {
	if r.IsZero() {
		return
	}

	var musts []string
	for i := 0; i < r.Must().Len(); i++ {
		m := r.Must().Index(i)
		musts = append(musts, m.OID())
	}

	// if either of these evaluate as true,
	// we shall not pass.
	if len(musts) == 0 {
		return
	} else if r.OC().IsZero() {
		return
	}

	var mays []string
	for i := 0; i < r.May().Len(); i++ {
		m := r.May().Index(i)
		mays = append(mays, m.OID())
	}

	def = make(DefinitionMap, 0)
	def[`NUMERICOID`] = []string{r.NumericOID()}
	def[`NAME`] = r.Names().List()
	def[`DESC`] = []string{r.Description()}
	def[`OBSOLETE`] = []string{bool2str(r.Obsolete())}
	def[`OC`] = []string{r.OC().OID()}
	def[`MUST`] = musts
	def[`MAY`] = mays
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
E returns the underlying error instance.
*/
func (r NameForm) E() (err error) {
	if !r.IsZero() {
		err = r.nameForm.err
	} else {
		err = ErrNilReceiver
	}

	return
}

/*
NewNameForm initializes and returns a new instance of [NameForm],
ready for manual assembly.  This method need not be used when creating
new [NameForm] instances by way of parsing, as that is handled on an
internal basis.

Use of this method does NOT automatically push the return instance into
the [Schema.NameForms] stack; this is left to the user.

Unlike the package-level [NewNameForm] function, this method will
automatically reference its originating [Schema] instance (the receiver).
This negates the need for manual use of the [NameForm.SetSchema]
method.

This is the recommended means of creating a new [NameForm] instance
wherever a single [Schema] is being used, which represents most use cases.
*/
func (r Schema) NewNameForm() NameForm {
	a := NewNameForm()
	a.SetSchema(r)
	return a
}

/*
NewNameForm initializes and returns a new instance of [NameForm],
ready for manual assembly.  This method need not be used when creating
new [NameForm] instances by way of parsing, as that is handled on an
internal basis.

Use of this function does not automatically reference the "parent" [Schema]
instance, leaving it up to the user to invoke the [NameForm.SetSchema]
method manually.

When interacting with a single [Schema] instance, which represents most use
cases, use of the [Schema.NewNameForm] method is PREFERRED over use of
this package-level function.

However certain migration efforts, schema audits and other such activities
may require distinct associations of [NameForm] instances with specific
[Schema] instances. Use of this function allows the user to specify the
appropriate [Schema] instance at a later point for a specific instance of
a [NameForm] instance.
*/
func NewNameForm() NameForm {
	nf := NameForm{newNameForm()}
	nf.nameForm.Extensions.setDefinition(nf)
	return nf
}

func newNameForm() *nameForm {
	return &nameForm{
		Name:       NewName(),
		Must:       NewAttributeTypeOIDList(`MUST`),
		May:        NewAttributeTypeOIDList(`MAY`),
		Extensions: NewExtensions(),
	}
}

/*
Replace overrides the receiver with x. Both must bear an identical
numeric OID and x MUST be compliant.

Note that the relevant [Schema] instance must be configured to allow
definition override by way of the [AllowOverride] bit setting.  See
the [Schema.Options] method for a means of accessing the settings
value.

Note that this method does not reallocate a new pointer instance
within the [NameForm] envelope type, thus all references to the
receiver instance within various stacks will be preserved.

This is a fluent method.
*/
func (r NameForm) Replace(x NameForm) NameForm {
	if !r.Schema().Options().Positive(AllowOverride) {
		return r
	}

	if !r.IsZero() && x.Compliant() {
		r.nameForm.replace(x)
	}

	return r
}

func (r *nameForm) replace(x NameForm) {
	if r.OID == `` {
		r.err = ErrMissingNumericOID
		return
	} else if r.OID != x.NumericOID() {
		r.err = ErrInvalidOID
		return
	}

	r.OID = x.nameForm.OID
	r.Macro = x.nameForm.Macro
	r.Name = x.nameForm.Name
	r.Desc = x.nameForm.Desc
	r.Obsolete = x.nameForm.Obsolete
	r.Structural = x.nameForm.Structural
	r.Must = x.nameForm.Must
	r.May = x.nameForm.May
	r.Extensions = x.nameForm.Extensions
	r.data = x.nameForm.data
	r.schema = x.nameForm.schema
	r.stringer = x.nameForm.stringer
	r.data = x.nameForm.data
}

/*
EnforcedBy returns an instance of [DITStructureRules] containing all
[DITStructureRule] instances which enforce the receiver instance.  A
return value with an integer length of zero (0) indicates that there
are no [DITStructureRule] instances which bear the receiver instance
through the 'FORM' clause at present.
*/
func (r NameForm) EnforcedBy() (dsr DITStructureRules) {
	if !r.Schema().IsZero() {
		dsr = NewDITStructureRules()
		rules := r.Schema().DITStructureRules()
		for i := 0; i < rules.Len(); i++ {
			if ds := rules.Index(i); ds.Form().IsIdentifiedAs(r.OID()) {
				dsr.Push(ds)
			}
		}
	}

	return
}

/*
Identifier returns the output of [NameForm.OID] and is merely used for
an interface-friendly means of obtaining the principal identifier of
the receiver instance.
*/
func (r NameForm) Identifier() (ident string) { return r.OID() }

/*
IsIdentifiedAs returns a Boolean value indicative of whether id matches
either the numericOID or descriptor of the receiver instance.  Case is
not significant in the matching process.
*/
func (r NameForm) IsIdentifiedAs(id string) (ident bool) {
	if !r.IsZero() {
		ident = id == r.NumericOID() || r.nameForm.Name.contains(id)
	}

	return
}

/*
Compliant returns a Boolean value indicative of every [NameForm]
returning a compliant response from the [NameForm.Compliant] method.
*/
func (r NameForms) Compliant() bool {
	var act int
	for i := 0; i < r.Len(); i++ {
		if r.Index(i).Compliant() {
			act++
		}
	}

	return act == r.Len()
}

/*
Compliant returns a Boolean value indicative of the receiver being fully
compliant per the required clauses of [§ 4.1.7.2 of RFC 4512]:

  - Numeric OID is present and valid
  - Structural [ObjectClass] is compliant
  - MUST clause has at least one (1) [AttributeType]
  - MAY and MUST clause [AttributeType] instances are compliant

[§ 4.1.7.2 of RFC 4512]: https://datatracker.ietf.org/doc/html/rfc4512#section-4.1.7.2
*/
func (r NameForm) Compliant() bool {
	if r.IsZero() {
		return false
	}

	if !isNumericOID(r.nameForm.OID) {
		return false
	}

	if !r.OC().Compliant() {
		return false
	}

	var (
		mct  int
		may  AttributeTypes = r.May()
		must AttributeTypes = r.Must()
	)

	for i := 0; i < must.Len(); i++ {
		if !must.Index(i).Compliant() {
			return false
		}
		mct++
	}

	for i := 0; i < may.Len(); i++ {
		if !may.Index(i).Compliant() {
			return false
		}
	}

	ok := mct > 0
	if ok {
		r.nameForm.err = nil
	}

	return ok
}

/*
SetData assigns x to the receiver instance. This is a general-use method and has no
specific intent beyond convenience. The contents may be subsequently accessed via the
[NameForm.Data] method.

This is a fluent method.
*/
func (r NameForm) SetData(x any) NameForm {
	if !r.IsZero() {
		r.nameForm.setData(x)
	}

	return r
}

func (r *nameForm) setData(x any) {
	r.data = x
}

/*
Data returns the underlying value (x) assigned to the receiver's data storage field. Data
can be set within the receiver instance by way of the [NameForm.SetData] method.
*/
func (r NameForm) Data() (x any) {
	if !r.IsZero() {
		x = r.nameForm.data
	}

	return
}

/*
SetSchema assigns an instance of [Schema] to the receiver instance.  This allows
internal verification of certain actions without the need for user input of
an instance of [Schema] manually at each juncture.

Note that the underlying [Schema] instance is automatically set when creating
instances of this type by way of parsing, as well as if the receiver instance
was initialized using the [Schema.NewNameForm] method.

This is a fluent method.
*/
func (r NameForm) SetSchema(schema Schema) NameForm {
	if !r.IsZero() {
		r.nameForm.setSchema(schema)
	}

	return r
}

func (r *nameForm) setSchema(schema Schema) {
	r.schema = schema
}

/*
Schema returns the [Schema] instance associated with the receiver instance.
*/
func (r NameForm) Schema() (s Schema) {
	if !r.IsZero() {
		s = r.nameForm.getSchema()
	}

	return
}

func (r *nameForm) getSchema() (s Schema) {
	if r != nil {
		s = r.schema
	}

	return
}

/*
SetName allows the manual assignment of one (1) or more RFC 4512-compliant
descriptor values by which the receiver instance is to be known.  This will
append to -- not replace -- any preexisting names.

This is a fluent method.
*/
func (r NameForm) SetName(name ...string) NameForm {
	if !r.IsZero() {
		r.nameForm.setName(name...)
	}

	return r
}

func (r *nameForm) setName(x ...string) {
	b4 := r.Name.Len()
	for i := 0; i < len(x); i++ {
		r.Name.Push(x[i])
	}

	if r.Name.Len()-len(x) != b4 {
		r.err = ErrInvalidNames
	}
}

/*
SetNumericOID allows the manual assignment of a numeric OID to the
receiver instance if the following are all true:

  - The input id value is a syntactically valid numeric OID
  - The receiver does not already possess a numeric OID

This is a fluent method.
*/
func (r NameForm) SetNumericOID(id string) NameForm {
	if !r.IsZero() {
		r.nameForm.setNumericOID(id)
	}

	return r
}

func (r *nameForm) setNumericOID(id string) {
	if isNumericOID(id) {
		// only set an OID when the receiver
		// lacks one (iow: no modifications)
		if len(r.OID) == 0 {
			r.OID = id
		}
	} else {
		r.err = ErrInvalidOID
	}

	return
}

/*
OID returns the string representation of the OID -- whether numeric or
descriptor -- held by the receiver instance.  Note that a principal name
is requested first and, if not found, the numeric OID is returned as a
fallback.  The return value should NEVER be zero for a valid receiver.
*/
func (r NameForm) OID() (oid string) {
	if !r.IsZero() {
		oid = r.NumericOID() // default
		if r.nameForm.Name.len() > 0 {
			oid = r.nameForm.Name.index(0)
		}
	}

	return
}

/*
SetObsolete will declare the receiver instance as OBSOLETE.  Calls to
this method will not operate in a toggling manner in that there is no
way to "unset" a state of obsolescence.

This is a fluent method.
*/
func (r NameForm) SetObsolete() NameForm {
	if !r.IsZero() {
		r.nameForm.setObsolete()
	}

	return r
}

func (r *nameForm) setObsolete() {
	if !r.Obsolete {
		r.Obsolete = true
	}
}

/*
Obsolete returns a Boolean value indicative of definition obsolescence.
*/
func (r NameForm) Obsolete() (o bool) {
	if !r.IsZero() {
		o = r.nameForm.Obsolete
	}

	return
}

/*
Extensions returns the [Extensions] instance -- if set -- within the
receiver instance.

Additions to the [Extensions] instance returned may be effected using
the [Extensions.Set] method. Existing [Extensions] cannot be altered.
*/
func (r NameForm) Extensions() (e Extensions) {
	if !r.IsZero() {
		e = r.nameForm.Extensions
	}

	return
}

/*
May returns an [AttributeTypes] containing zero (0) or more allowed
[AttributeType] definitions for use with this class.

Additions to the [AttributeTypes] instance returned may be effected
using the [AttributeTypes.Push] method.  Existing [AttributeType]
instances cannot be removed.
*/
func (r NameForm) May() (may AttributeTypes) {
	if !r.IsZero() {
		may = r.nameForm.May
	}

	return
}

/*
SetMay assigns the provided input [AttributeType] instance(s) to the
receiver's MAY clause.

This is a fluent method.
*/
func (r NameForm) SetMay(m ...any) NameForm {
	if !r.IsZero() {
		r.nameForm.setMay(m...)
	}

	return r
}

func (r *nameForm) setMay(m ...any) {
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
SetExtension assigns key x to value xstrs within the receiver's underlying
[Extensions] instance.

This is a fluent method.
*/
func (r NameForm) SetExtension(x string, xstrs ...string) NameForm {
	if !r.IsZero() {
		r.nameForm.setExtension(x, xstrs...)
	}

	return r
}

func (r *nameForm) setExtension(x string, xstrs ...string) {
	r.Extensions.Set(x, xstrs...)
}

/*
Must returns an [AttributeTypes] containing zero (0) or more required
[AttributeType] definitions for use with this class.

Additions to the [AttributeTypes] instance returned may be effected
using the [AttributeTypes.Push] method.  Existing [AttributeType]
instances cannot be removed.
*/
func (r NameForm) Must() (must AttributeTypes) {
	if !r.IsZero() {
		must = r.nameForm.Must
	}

	return
}

/*
SetMust assigns the provided input [AttributeType] instance(s) to the
receiver's MUST clause.

This is a fluent method.
*/
func (r NameForm) SetMust(m ...any) NameForm {
	if !r.IsZero() {
		r.nameForm.setMust(m...)
	}

	return r
}

func (r *nameForm) setMust(m ...any) {
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
NumericOID returns the string representation of the numeric OID
held by the receiver instance.
*/
func (r NameForm) NumericOID() (noid string) {
	if !r.IsZero() {
		noid = r.nameForm.OID
	}

	return
}

/*
Name returns the string form of the principal name of the receiver instance, if set.
*/
func (r NameForm) Name() (id string) {
	if !r.IsZero() {
		id = r.nameForm.Name.index(0)
	}

	return
}

/*
Names returns the underlying instance of [QuotedDescriptorList] from within the
receiver instance.
*/
func (r NameForm) Names() (names QuotedDescriptorList) {
	if !r.IsZero() {
		names = r.nameForm.Name
	}

	return
}

/*
SetStringer allows the assignment of an individual [Stringer] function or
method to all [NameForm] slices within the receiver stack instance.

Input of zero (0) variadic values, or an explicit nil, will overwrite all
preexisting stringer functions with the internal closure default, which is
based upon a one-time use of the [text/template] package by all receiver
slice instances.

Input of a non-nil closure function value will overwrite all preexisting
stringers.

This is a fluent method and may be used multiple times.
*/
func (r NameForms) SetStringer(function ...Stringer) NameForms {
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
func (r NameForm) SetStringer(function ...Stringer) NameForm {
	if r.Compliant() {
		r.nameForm.setStringer(function...)
	}

	return r
}

func (r *nameForm) setStringer(function ...Stringer) {
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
XOrigin returns an instance of [NameForms] containing only definitions
which bear the X-ORIGIN value of x. Case is not significant in the matching
process, nor is whitespace (e.g.: RFC 4517 vs. RFC4517).
*/
func (r NameForms) XOrigin(x string) (defs NameForms) {
	defs = NewNameForms()
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
func (r NameForm) String() (nf string) {
	if !r.IsZero() {
		if r.nameForm.stringer != nil {
			nf = r.nameForm.stringer()
		}
	}

	return
}

func (r NameForm) macro() (m []string) {
	if !r.IsZero() {
		m = r.nameForm.Macro
	}

	return
}

func (r NameForm) setOID(x string) {
	if !r.IsZero() {
		r.nameForm.OID = x
	}
}

/*
prepareString returns a string an an error indicative of an attempt
to represent the receiver instance as a string using [text/template].
*/
func (r *nameForm) prepareString() (str string, err error) {
	buf := newBuf()
	t := newTemplate(`nameForm`).
		Funcs(funcMap(map[string]any{
			`ExtensionSet`: r.Extensions.tmplFunc,
			`MayLen`:       r.May.len,
			`Obsolete`:     func() bool { return r.Obsolete },
		}))

	if t, err = t.Parse(nameFormTmpl); err == nil {
		if err = t.Execute(buf, struct {
			Definition *nameForm
			HIndent    string
		}{
			Definition: r,
			HIndent:    hindent(r.schema.Options().Positive(HangingIndents)),
		}); err == nil {
			str = buf.String()
		}
	}

	if err != nil {
		r.err = err
	}

	return
}

/*
Description returns the underlying (optional) descriptive text
assigned to the receiver instance.
*/
func (r NameForm) Description() (desc string) {
	if !r.IsZero() {
		desc = r.nameForm.Desc
	}

	return
}

/*
SetDescription parses desc into the underlying DESC clause within the
receiver instance.  Although a RFC 4512-compliant QuotedString is
required, the outer single-quotes need not be specified literally.
*/
func (r NameForm) SetDescription(desc string) NameForm {
	if !r.IsZero() {
		r.nameForm.setDescription(desc)
	}

	return r
}

func (r *nameForm) setDescription(desc string) {
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
Marshal returns an error following an attempt to marshal the contents of
def, which may be either a [DefinitionMap] or map[string]any instance.

The receiver instance must be initialized prior to use of this method
using the [Schema.NewNameForm] method.
*/
func (r NameForm) Marshal(def any) error {
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
		case `DESC`, `NUMERICOID`:
			z := map[string]func(string) NameForm{
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
		case `OC`:
			r.marshalOC(v)
		case `MUST`, `MAY`:
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

func (r NameForm) marshalOC(v any) {
	switch tv := v.(type) {
	case []string:
		r.SetOC(tv[0])
	case string:
		r.SetOC(tv)
	}
}

func (r NameForm) marshalBoolean(v any) {
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

func (r NameForm) marshalExt(key string, v any) {
	if hasPfx(key, `X-`) {
		switch tv := v.(type) {
		case string:
			r.SetExtension(key, tv)
		case []string:
			r.SetExtension(key, tv...)
		}
	}
}

func (r NameForm) marshalMulti(k string, v any) {
	z := map[string]func(...any) NameForm{
		`MUST`: r.SetMust,
		`MAY`:  r.SetMay,
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
the receiver into any [NameForms] stack, nor does it automatically execute
the [NameForm.SetStringer] method, leaving these tasks to the user.  If the
automatic handling of these tasks is desired, see the [Schema.ParseNameForm]
method as an alternative.
*/
func (r NameForm) Parse(raw string) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
		return
	}

	if r.getSchema().IsZero() {
		err = ErrNilSchemaRef
		return
	}

	err = r.nameForm.parse(raw)

	return
}

func (r *nameForm) parse(raw string) error {
	// parseLS wraps the antlr4512 NameForm parser/lexer
	mp, err := parseNF(raw)
	if err == nil {
		// We received the parsed data from ANTLR (mp).
		// Now we need to marshal it into the receiver.
		var def NameForm
		if def, err = r.schema.marshalNF(mp); err == nil {
			r.OID = def.nameForm.OID
			_r := NameForm{r}
			_r.replace(def)
		}
	}

	return err
}

/*
Inventory returns an instance of [Inventory] which represents the current
inventory of [NameForm] instances within the receiver.

The keys are numeric OIDs, while the values are zero (0) or more string
slices, each representing a name by which the definition is known.
*/
func (r NameForms) Inventory() (inv Inventory) {
	inv = make(Inventory, 0)
	for i := 0; i < r.len(); i++ {
		def := r.index(i)
		inv[def.NumericOID()] = def.Names().List()
	}

	return
}

/*
Type returns the string literal "nameForm".
*/
func (r NameForm) Type() string {
	return `nameForm`
}

/*
Type returns the string literal "nameForms".
*/
func (r NameForms) Type() string {
	return `nameForms`
}

// stackage closure func - do not exec directly.
func (r NameForms) canPush(x ...any) (err error) {
	if len(x) == 0 {
		return
	}

	for i := 0; i < len(x) && err == nil; i++ {
		nf, ok := x[i].(NameForm)
		if !ok || nf.IsZero() {
			err = ErrTypeAssert
		} else if tst := r.get(nf.NumericOID()); !tst.IsZero() {
			err = mkerr(ErrNotUnique.Error() + ": " + nf.Type() + `, ` + nf.NumericOID())
		}
	}

	return
}

/*
Len returns the current integer length of the receiver instance.
*/
func (r NameForms) Len() int {
	return r.len()
}

func (r NameForms) len() int {
	return r.cast().Len()
}

/*
String is a stringer method that returns the string representation
of the receiver instance.
*/
func (r NameForms) String() string {
	return r.cast().String()
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r NameForms) IsZero() bool {
	return r.cast().IsZero()
}

/*
Index returns the instance of [NameForm] found within the receiver
stack instance at index N.  If no instance is found at the index
specified, a zero [NameForm] instance is returned.
*/
func (r NameForms) Index(idx int) NameForm {
	return r.index(idx)
}

func (r NameForms) index(idx int) (nf NameForm) {
	slice, found := r.cast().Index(idx)
	if found {
		if _nf, ok := slice.(NameForm); ok {
			nf = _nf
		}
	}

	return
}

/*
Push returns an error following an attempt to push a [NameForm]
instance into the receiver stack instance.
*/
func (r NameForms) Push(nf any) error {
	return r.push(nf)
}

func (r NameForms) push(x any) (err error) {
	switch tv := x.(type) {
	case NameForm:
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
Contains calls [NameForms.Get] to return a Boolean value indicative
of a successful, non-zero retrieval of a [NameForm] instance -- matching
the provided id -- from within the receiver stack instance.
*/
func (r NameForms) Contains(id string) bool {
	return r.contains(id)
}

func (r NameForms) contains(id string) bool {
	return !r.get(id).IsZero()
}

/*
Get returns an instance of [NameForm] based upon a search for id within
the receiver stack instance.

The return instance, if not nil, was retrieved based upon a textual
match of the principal identifier of a [NameForm] and the provided id.

The return instance is nil if no match was made.

Case is not significant in the matching process.
*/
func (r NameForms) Get(id string) NameForm {
	return r.get(id)
}

func (r NameForms) get(id string) (nf NameForm) {
	for i := 0; i < r.len() && nf.IsZero(); i++ {
		if _nf := r.index(i); !_nf.IsZero() {
			if _nf.nameForm.OID == id {
				nf = _nf
			} else if _nf.nameForm.Name.contains(id) {
				nf = _nf
			}
		}
	}

	return
}

/*
OC returns the STRUCTURAL [ObjectClass] specified within the receiver instance.
*/
func (r NameForm) OC() ObjectClass {
	return r.nameForm.Structural
}

/*
SetOC sets the structural class of the receiver to the value provided. Valid input
types are string, to represent an RFC 4512 OID residing in the underlying [Schema]
instance, or an actual structural [ObjectClass] instance already obtained or crafted.

This is a fluent method.
*/
func (r NameForm) SetOC(x any) NameForm {
	if !r.IsZero() {
		r.nameForm.setOC(x)
	}

	return r
}

func (r *nameForm) setOC(x any) {
	var oc ObjectClass
	switch tv := x.(type) {
	case string:
		if !r.schema.IsZero() {
			oc = r.schema.ObjectClasses().get(tv)
		}
	case ObjectClass:
		oc = tv
	}

	if !oc.IsZero() && oc.Kind() == StructuralKind {
		r.Structural = oc
	} else {
		r.err = ErrIncompatStructuralClass
	}
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r NameForm) IsZero() bool {
	return r.nameForm == nil
}

/*
LoadNameForms returns an error following an attempt to load all
built-in [NameForm] slices into the receiver instance.
*/
func (r Schema) LoadNameForms() error {
	return r.loadNameForms()
}

func (r Schema) loadNameForms() (err error) {
	if !r.IsZero() {
		funks := []func() error{
			r.loadRFC2377NameForms,
			r.loadRFC4403NameForms,
			r.loadX501NameForms,
		}

		for i := 0; i < len(funks) && err == nil; i++ {
			err = funks[i]()
		}
	}

	return
}

/*
LoadRFC2377NameForms returns an error following an attempt to
load all RFC 2377 [NameForm] slices into the receiver instance.
*/
func (r Schema) LoadRFC2377NameForms() error {
	return r.loadRFC2377NameForms()
}

func (r Schema) loadRFC2377NameForms() (err error) {

	var i int
	for i = 0; i < len(rfc2377NameForms) && err == nil; i++ {
		at := rfc2377NameForms[i]
		err = r.ParseNameForm(string(at))
	}

	if want := rfc2377NameForms.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC2377 NameForms parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4403NameForms returns an error following an attempt to
load all RFC 4403 [NameForm] slices into the receiver instance.
*/
func (r Schema) LoadRFC4403NameForms() error {
	return r.loadRFC4403NameForms()
}

func (r Schema) loadRFC4403NameForms() (err error) {

	var i int
	for i = 0; i < len(rfc4403NameForms) && err == nil; i++ {
		at := rfc4403NameForms[i]
		err = r.ParseNameForm(string(at))
	}

	if want := rfc4403NameForms.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4403 NameForms parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadX501NameForm returns an error following an attempt to load all
[ITU-T Rec. X.501] [NameForm] slices into the receiver instance.

[ITU-T Rec. X.501]: https://www.itu.int/rec/T-REC-X.501
*/
func (r Schema) LoadX501NameForms() error {
	return r.loadX501NameForms()
}

func (r Schema) loadX501NameForms() (err error) {

	var i int
	for i = 0; i < len(x501NameForms) && err == nil; i++ {
		nf := x501NameForms[i]
		err = r.ParseNameForm(string(nf))
	}

	if want := x501NameForms.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of X.501 NameForms parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

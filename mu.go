package schemax

/*
NewMatchingRuleUses initializes a new [MatchingRuleUses] instance.
*/
func NewMatchingRuleUses() MatchingRuleUses {
	r := MatchingRuleUses(newCollection(`matchingRuleUses`))
	r.cast().SetPushPolicy(r.canPush)

	return r
}

/*
NewMatchingRuleUse initializes and returns a new instance of [MatchingRuleUse],
ready for manual assembly.  This method need not be used when creating
new [MatchingRuleUse] instances by way of parsing, as that is handled on an
internal basis.

Use of this method does NOT automatically push the return instance into
the [Schema.MatchingRuleUses] stack; this is left to the user.

Unlike the package-level [NewMatchingRuleUse] function, this method will
automatically reference its originating [Schema] instance (the receiver).
This negates the need for manual use of the [MatchingRuleUse.SetSchema]
method.

This is the recommended means of creating a new [MatchingRuleUse] instance
wherever a single [Schema] is being used, which represents most use cases.
*/
func (r Schema) NewMatchingRuleUse() MatchingRuleUse {
	return NewMatchingRuleUse().SetSchema(r)
}

/*
NewMatchingRuleUse initializes and returns a new instance of [MatchingRuleUse],
ready for manual assembly.  This method need not be used when creating
new [MatchingRuleUse] instances by way of parsing, as that is handled on an
internal basis.

Use of this function does not automatically reference the "parent" [Schema]
instance, leaving it up to the user to invoke the [MatchingRuleUse.SetSchema]
method manually.

When interacting with a single [Schema] instance, which represents most use
cases, use of the [Schema.NewMatchingRuleUse] method is PREFERRED over use of
this package-level function.

However certain migration efforts, schema audits and other such activities
may require distinct associations of [MatchingRuleUse] instances with specific
[Schema] instances. Use of this function allows the user to specify the
appropriate [Schema] instance at a later point for a specific instance of
a [MatchingRuleUse] instance.
*/
func NewMatchingRuleUse() MatchingRuleUse {
	mu := MatchingRuleUse{newMatchingRuleUse()}
	mu.matchingRuleUse.Extensions.setDefinition(mu)
	return mu
}

func newMatchingRuleUse() *matchingRuleUse {
	return &matchingRuleUse{
		Name:       NewName(),
		Applies:    NewAttributeTypeOIDList(`APPLIES`),
		Extensions: NewExtensions(),
	}
}

/*
Parse always returns an error, as parsing does not apply to [MatchingRuleUse]
instances, as they are meant to be dynamically generated -- not parsed.

See the [Schema.UpdateMatchingRuleUses] method for a means of refreshing
the current [Schema.MatchingRuleUses] stack based upon the [AttributeType]
instances present within the [Schema] instance at runtime.
*/
func (r MatchingRuleUse) Parse(raw string) error {
	return mkerr("Parsing is not applicable to a MatchingRuleUse")
}

/*
Marshal returns an error following an attempt to marshal the contents of
def, which may be either a [DefinitionMap] or map[string]any instance.

The receiver instance must be initialized prior to use of this method
using the [Schema.NewMatchingRule] method.
*/
func (r MatchingRuleUse) Marshal(def any) error {
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
			z := map[string]func(string) MatchingRuleUse{
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
		case `APPLIES`:
			r.marshalApplies(v)
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

func (r MatchingRuleUse) marshalBoolean(v any) {
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

func (r MatchingRuleUse) marshalExt(key string, v any) {
	if hasPfx(key, `X-`) {
		switch tv := v.(type) {
		case string:
			r.SetExtension(key, tv)
		case []string:
			r.SetExtension(key, tv...)
		}
	}
}

func (r MatchingRuleUse) marshalApplies(v any) {
	switch tv := v.(type) {
	case []string:
		for i := 0; i < len(tv); i++ {
			r.SetApplies(tv[i])
		}
	case string:
		r.SetApplies(tv)
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
within the [MatchingRuleUse] envelope type, thus all references to the
receiver instance within various stacks will be preserved.

This is a fluent method.
*/
func (r MatchingRuleUse) Replace(x MatchingRuleUse) MatchingRuleUse {
	if !r.Schema().Options().Positive(AllowOverride) {
		return r
	}

	if !r.IsZero() && r.Compliant() {
		r.matchingRuleUse.replace(x)
	}

	return r
}

func (r *matchingRuleUse) replace(x MatchingRuleUse) {
	if r.OID.IsZero() {
		return
	} else if r.OID.NumericOID() != x.NumericOID() {
		return
	}

	r.OID = x.matchingRuleUse.OID
	r.Name = x.matchingRuleUse.Name
	r.Desc = x.matchingRuleUse.Desc
	r.Obsolete = x.matchingRuleUse.Obsolete
	r.Applies = x.matchingRuleUse.Applies
	r.Extensions = x.matchingRuleUse.Extensions
	r.data = x.matchingRuleUse.data
	r.schema = x.matchingRuleUse.schema
	r.stringer = x.matchingRuleUse.stringer
	r.data = x.matchingRuleUse.data
}

/*
Compliant returns a Boolean value indicative of every [MatchingRuleUse]
returning a compliant response from the [MatchingRuleUse.Compliant] method.
*/
func (r MatchingRuleUses) Compliant() bool {
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
compliant per the required clauses of [§ 4.1.4 of RFC 4512]:

  - Numeric OID must be present and valid
  - Numeric OID must correlate to a previously registered [MatchingRule]

[§ 4.1.4 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.4
*/
func (r MatchingRuleUse) Compliant() bool {
	if r.IsZero() {
		return false
	}

	var (
		appl AttributeTypes = r.Applies()
		act  int            // Applied AttributeType count
	)

	for i := 0; i < appl.Len(); i++ {
		if appl.Index(i).Compliant() {
			act++
		}
	}

	if r.Schema().MatchingRules().get(r.NumericOID()).IsZero() {
		return false
	}

	return act == r.Applies().Len()
}

/*
SetData assigns x to the receiver instance. This is a general-use method and has no
specific intent beyond convenience. The contents may be subsequently accessed via the
[MatchingRuleUse.Data] method.

This is a fluent method.
*/
func (r MatchingRuleUse) SetData(x any) MatchingRuleUse {
	if !r.IsZero() {
		r.matchingRuleUse.setData(x)
	}

	return r
}

func (r *matchingRuleUse) setData(x any) {
	r.data = x
}

/*
Data returns the underlying value (x) assigned to the receiver's data storage field. Data
can be set within the receiver instance by way of the [MatchingRuleUse.SetData] method.
*/
func (r MatchingRuleUse) Data() (x any) {
	if !r.IsZero() {
		x = r.matchingRuleUse.data
	}

	return
}

/*
SetStringer allows the assignment of an individual [Stringer] function or
method to all [MatchingRuleUse] slices within the receiver stack instance.

Input of zero (0) variadic values, or an explicit nil, will overwrite all
preexisting stringer functions with the internal closure default, which is
based upon a one-time use of the [text/template] package by all receiver
slice instances.

Input of a non-nil closure function value will overwrite all preexisting
stringers.

This is a fluent method and may be used multiple times.
*/
func (r MatchingRuleUses) SetStringer(function ...Stringer) MatchingRuleUses {
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
func (r MatchingRuleUse) SetStringer(function ...Stringer) MatchingRuleUse {
	if r.Compliant() {
		r.matchingRuleUse.setStringer(function...)
	}

	return r
}

func (r *matchingRuleUse) setStringer(function ...Stringer) {
	var stringer Stringer
	if len(function) > 0 {
		stringer = function[0]
	}

	if stringer == nil {
		// no user provided closure means we
		// defer to a general use stringer.
		str, err := r.prepareString() // perform one-time text/template op
		if err == nil {
			// Save the stringer
			r.stringer = func() string {
				// Return a preserved value.
				return str
			}
		}
		return
	}

	// assign user-provided closure
	r.stringer = stringer
}

/*
XOrigin returns an instance of [MatchingRuleUses] containing only definitions
which bear the X-ORIGIN value of x. Case is not significant in the matching
process, nor is whitespace (e.g.: RFC 4517 vs. RFC4517).
*/
func (r MatchingRuleUses) XOrigin(x string) (defs MatchingRuleUses) {
	defs = NewMatchingRuleUses()
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
String is a stringer method that returns the string representation of
the receiver instance.  A zero-value indicates an invalid receiver, or
that the [ObjectClass.SetStringer] method was not used during MANUAL
composition of the receiver.
*/
func (r MatchingRuleUse) String() (def string) {
	if !r.IsZero() {
		if r.matchingRuleUse.stringer != nil {
			def = r.matchingRuleUse.stringer()
		}
	}

	return
}

/*
Obsolete returns a Boolean value indicative of definition obsolescence.
*/
func (r MatchingRuleUse) Obsolete() (o bool) {
	if !r.IsZero() {
		o = r.matchingRuleUse.Obsolete
	}

	return
}

/*
SetObsolete sets the receiver instance to OBSOLETE if not already set.

Obsolescence cannot be unset.

This is a fluent method.
*/
func (r MatchingRuleUse) SetObsolete() MatchingRuleUse {
	if !r.IsZero() {
		r.matchingRuleUse.setObsolete()
	}

	return r
}

func (r *matchingRuleUse) setObsolete() {
	if !r.Obsolete {
		r.Obsolete = true
	}
}

/*
Applies returns an [AttributeTypes] instance containing pointer references
to all [AttributeType] instances to which the receiver applies.
*/
func (r MatchingRuleUse) Applies() (aa AttributeTypes) {
	if !r.IsZero() {
		aa = r.matchingRuleUse.Applies
	}

	return
}

/*
SetApplies assigns the provided input values as applied [AttributeType]
instances advertised through the receiver instance.

This is a fluent method.
*/
func (r MatchingRuleUse) SetApplies(m ...any) MatchingRuleUse {
	if !r.IsZero() {
		r.matchingRuleUse.setApplies(m...)
	}

	return r
}

func (r *matchingRuleUse) setApplies(m ...any) {
	var err error
	for i := 0; i < len(m) && err == nil; i++ {
		var at AttributeType
		switch tv := m[i].(type) {
		case string:
			at = r.schema.AttributeTypes().get(tv)
		case AttributeType:
			at = tv
		default:
			err = ErrInvalidInput
		}

		if err == nil && at.Compliant() {
			r.Applies.Push(at)
		}
	}
}

/*
SetSchema assigns an instance of [Schema] to the receiver instance.  This allows
internal verification of certain actions without the need for user input of
an instance of [Schema] manually at each juncture.

Note that the underlying [Schema] instance is automatically set when creating
instances of this type by way of parsing, as well as if the receiver instance
was initialized using the [Schema.NewMatchingRuleUse] method.

This is a fluent method.
*/
func (r MatchingRuleUse) SetSchema(schema Schema) MatchingRuleUse {
	if !r.IsZero() {
		r.matchingRuleUse.setSchema(schema)
	}

	return r
}

func (r *matchingRuleUse) setSchema(schema Schema) {
	r.schema = schema
}

/*
Schema returns the [Schema] instance associated with the receiver instance.
*/
func (r MatchingRuleUse) Schema() (s Schema) {
	if !r.IsZero() {
		s = r.matchingRuleUse.getSchema()
	}

	return
}

func (r *matchingRuleUse) getSchema() (s Schema) {
	if r != nil {
		s = r.schema
	}

	return
}

/*
IsIdentifiedAs returns a Boolean value indicative of whether id matches
either the numericOID or descriptor of the receiver instance.  Case is
not significant in the matching process.
*/
func (r MatchingRuleUse) IsIdentifiedAs(id string) (ident bool) {
	if !r.IsZero() {
		ident = id == r.NumericOID() || r.matchingRuleUse.Name.contains(id)
	}

	return
}

/*
MatchingRuleUses returns the [MatchingRuleUses] instance from
within the receiver instance.
*/
func (r Schema) MatchingRuleUses() (mus MatchingRuleUses) {
	slice, _ := r.cast().Index(matchingRuleUsesIndex)
	mus, _ = slice.(MatchingRuleUses)
	return
}

/*
prepareString returns a string an an error indicative of an attempt
to represent the receiver instance as a string using [text/template].
*/
func (r *matchingRuleUse) prepareString() (str string, err error) {
	buf := newBuf()
	t := newTemplate(r.Type()).
		Funcs(funcMap(map[string]any{
			`MatchingRuleOID`: r.OID.NumericOID,
			`ExtensionSet`:    r.Extensions.tmplFunc,
			`Applied`:         r.Applies.String,
			`Obsolete`:        func() bool { return r.Obsolete },
		}))
	if t, err = t.Parse(matchingRuleUseTmpl); err == nil {
		if err = t.Execute(buf, struct {
			Definition *matchingRuleUse
			HIndent    string
		}{
			Definition: r,
			HIndent:    hindent(r.schema.Options().Positive(HangingIndents)),
		}); err == nil {
			str = buf.String()
		}
	}

	return
}

/*
Maps returns slices of [DefinitionMap] instances.
*/
func (r MatchingRuleUses) Maps() (defs DefinitionMaps) {
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
func (r MatchingRuleUse) Map() (def DefinitionMap) {
	if r.IsZero() {
		return
	}

	var applies []string
	for i := 0; i < r.Applies().Len(); i++ {
		m := r.Applies().Index(i)
		applies = append(applies, m.OID())
	}

	def = make(DefinitionMap, 0)
	def[`NUMERICOID`] = []string{r.NumericOID()}
	def[`NAME`] = r.Names().List()
	def[`DESC`] = []string{r.Description()}
	def[`OBSOLETE`] = []string{bool2str(r.Obsolete())}
	def[`APPLIES`] = applies
	def[`RAW`] = []string{r.String()}
	def[`TYPE`] = []string{r.Type()}

	// copy our extensions from receiver r
	// into destination def.
	def = mapTransferExtensions(r, def)

	// Clean up any empty fields
	def.clean()

	return
}

/*
Inventory returns an instance of [Inventory] which represents the current
inventory of [MatchingRuleUse] instances within the receiver.

The keys are numeric OIDs, while the values are zero (0) or more string
slices, each representing a name by which the definition is known.
*/
func (r MatchingRuleUses) Inventory() (inv Inventory) {
	inv = make(Inventory, 0)
	for i := 0; i < r.len(); i++ {
		def := r.index(i)
		inv[def.NumericOID()] = def.Names().List()

	}

	return
}

/*
Type returns the string literal "matchingRuleUse".
*/
func (r MatchingRuleUse) Type() string {
	return r.matchingRuleUse.Type()
}

func (r matchingRuleUse) Type() string {
	return `matchingRuleUse`
}

/*
Type returns the string literal "matchingRuleUses".
*/
func (r MatchingRuleUses) Type() string {
	return `matchingRuleUses`
}

func (r MatchingRuleUses) prepareStrings() (err error) {
	for i := 0; i < r.Len() && err == nil; i++ {
		mu := r.index(i)
		var str string
		str, err = mu.matchingRuleUse.prepareString()
		if err == nil {
			mu.matchingRuleUse.stringer = func() string {
				return str
			}
		}
	}

	return
}

/*
OID returns the string representation of an OID -- which is either
a numeric OID or descriptor -- that refers to the [MatchingRule]
upon which the receiver instance is based.
*/
func (r MatchingRuleUse) OID() (oid string) {
	if !r.IsZero() {
		oid = r.NumericOID() // default
		if r.matchingRuleUse.Name.len() > 0 {
			oid = r.matchingRuleUse.Name.index(0)
		}
	}

	return
}

/*
Extensions returns the [Extensions] instance -- if set -- within
the receiver.
*/
func (r MatchingRuleUse) Extensions() (e Extensions) {
	if !r.IsZero() {
		e = r.matchingRuleUse.Extensions
	}

	return
}

/*
SetExtension assigns key x to value xstrs within the receiver's underlying
[Extensions] instance.

This is a fluent method.
*/
func (r MatchingRuleUse) SetExtension(x string, xstrs ...string) MatchingRuleUse {
	if !r.IsZero() {
		r.matchingRuleUse.setExtension(x, xstrs...)
	}

	return r
}

func (r *matchingRuleUse) setExtension(x string, xstrs ...string) {
	r.Extensions.Set(x, xstrs...)
}

/*
NumericOID returns the string representation of the numeric OID
held by the receiver instance.
*/
func (r MatchingRuleUse) NumericOID() (noid string) {
	if !r.IsZero() {
		noid = r.matchingRuleUse.OID.NumericOID()
	}

	return
}

/*
SetNumericOID allows the manual assignment of a numeric OID to the
receiver instance if the following are all true:

  - The input id value is a syntactically valid numeric OID
  - The input id relates to a known [MatchingRule] instance within the associated [Schema]
  - The receiver does not already possess a numeric OID

This is a fluent method.
*/
func (r MatchingRuleUse) SetNumericOID(id string) MatchingRuleUse {
	if !r.IsZero() {
		r.matchingRuleUse.setNumericOID(id)
	}

	return r
}

func (r *matchingRuleUse) setNumericOID(id string) {
	mr := r.schema.MatchingRules().Get(id)
	// only set an OID when the receiver
	// lacks one (iow: no modifications)
	// and when the MR has been found.
	if !mr.IsZero() && r.OID.IsZero() {
		r.OID = mr
	}

	return
}

/*
Name returns the string form of the principal name of the receiver instance, if set.
*/
func (r MatchingRuleUse) Name() (id string) {
	if !r.IsZero() {
		id = r.matchingRuleUse.Name.index(0)
	}

	return
}

/*
Names returns the underlying instance of [QuotedDescriptorList] from
within the receiver.
*/
func (r MatchingRuleUse) Names() (names QuotedDescriptorList) {
	if !r.IsZero() {
		names = r.matchingRuleUse.Name
	}

	return
}

/*
SetName assigns the provided names to the receiver instance.

Name instances must conform to RFC 4512 descriptor format but
need not be quoted.

This is a fluent method.
*/
func (r MatchingRuleUse) SetName(x ...string) MatchingRuleUse {
	if !r.IsZero() {
		r.matchingRuleUse.setName(x...)
	}

	return r
}

func (r *matchingRuleUse) setName(x ...string) {
	for i := 0; i < len(x); i++ {
		r.Name.Push(x[i])
	}
}

/*
Description returns the underlying (optional) descriptive text
assigned to the receiver instance.
*/
func (r MatchingRuleUse) Description() (desc string) {
	if !r.IsZero() {
		desc = r.matchingRuleUse.Desc
	}
	return
}

/*
SetDescription parses desc into the underlying DESC clause within the
receiver instance.  Although a RFC 4512-compliant QuotedString is
required, the outer single-quotes need not be specified literally.
*/
func (r MatchingRuleUse) SetDescription(desc string) MatchingRuleUse {
	if !r.IsZero() {
		r.matchingRuleUse.setDescription(desc)
	}

	return r
}

func (r *matchingRuleUse) setDescription(desc string) {
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
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r MatchingRuleUse) IsZero() bool {
	return r.matchingRuleUse == nil
}

func (r MatchingRuleUses) canPush(x ...any) (err error) {
	if len(x) == 0 {
		return
	}

	for i := 0; i < len(x) && err == nil; i++ {
		mu, ok := x[i].(MatchingRuleUse)
		if !ok || mu.IsZero() {
			err = ErrTypeAssert
		} else if tst := r.get(mu.NumericOID()); !tst.IsZero() {
			err = mkerr(ErrNotUnique.Error() + ": " + mu.Type() + `, ` + mu.NumericOID())
		}
	}

	return
}

/*
Len returns the current integer length of the receiver instance.
*/
func (r MatchingRuleUses) Len() int {
	return r.len()
}

func (r MatchingRuleUses) len() int {
	return r.cast().Len()
}

/*
String is a stringer method that returns the string representation
of the receiver instance.
*/
func (r MatchingRuleUses) String() string {
	return r.cast().String()
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r MatchingRuleUses) IsZero() bool {
	return r.cast().IsZero()
}

/*
Index returns the instance of [MatchingRuleUse] found within the
receiver stack instance at index N.  If no instance is found at
the index specified, a zero [MatchingRuleUse] instance is returned.
*/
func (r MatchingRuleUses) Index(idx int) MatchingRuleUse {
	return r.index(idx)
}

func (r MatchingRuleUses) index(idx int) (mu MatchingRuleUse) {
	slice, found := r.cast().Index(idx)
	if found {
		if _mu, ok := slice.(MatchingRuleUse); ok {
			mu = _mu
		}
	}

	return
}

/*
Push returns an error following an attempt to push a [MatchingRuleUse]
instance into the receiver instance.
*/
func (r MatchingRuleUses) Push(mu any) error {
	return r.push(mu)
}

func (r MatchingRuleUses) push(x any) (err error) {
	switch tv := x.(type) {
	case MatchingRuleUse:
		if tv.IsZero() {
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
Contains calls [MatchingRuleUses.Get] to return a Boolean value indicative of
a successful, non-zero retrieval of a [MatchingRuleUse] instance -- matching
the provided id -- from within the receiver stack instance.
*/
func (r MatchingRuleUses) Contains(id string) bool {
	return r.contains(id)
}

func (r MatchingRuleUses) contains(id string) bool {
	return !r.get(id).IsZero()
}

/*
Get returns an instance of [MatchingRuleUse] based upon a search for id within
the receiver stack instance.

The return instance, if not nil, was retrieved based upon a textual match of
the principal identifier of a [MatchingRuleUse] and the provided id.

The return instance is nil if no match was made.

Case is not significant in the matching process.
*/
func (r MatchingRuleUses) Get(id string) MatchingRuleUse {
	return r.get(id)
}

func (r MatchingRuleUses) get(id string) (mu MatchingRuleUse) {
	for i := 0; i < r.len() && mu.IsZero(); i++ {
		if _mu := r.index(i); !_mu.IsZero() {
			if _mu.IsIdentifiedAs(id) {
				mu = _mu
			}
		}
	}

	return
}

/*
makeMatchingRuleUse fashions and returns a new MatchingRuleUse instance
based on the contents of the receiver instance.  The returned instance,
assuming a nil error condition, may have its APPLIES clause populated
with "users" (AttributeType instances) of the indicated matchingRule.
*/
func (r MatchingRule) makeMatchingRuleUse() (mu MatchingRuleUse, err error) {
	if !r.Compliant() {
		err = ErrDefNonCompliant
		return
	}

	_mu := newMatchingRuleUse()
	_mu.Name = r.matchingRule.Name
	_mu.OID = r
	_mu.schema = r.Schema()
	_mu.Extensions = r.matchingRule.Extensions
	mu = MatchingRuleUse{_mu}
	mu.SetStringer()

	return
}

/*
updateMatchingRuleUses returns an instance of error following an attempt
to refresh the collection of MatchingRuleUse instances within the
receiver to include input variable "at" wherever appropriate.
*/
func (r Schema) updateMatchingRuleUses(ats AttributeTypes) (err error) {
	if ats.len() == 0 {
		return
	}

	for i := 0; i < ats.len() && err == nil; i++ {
		if at := ats.index(i); !at.IsZero() {
			for _, funk := range []func(AttributeType) error{
				r.updateEqualityUses,
				r.updateSubstringUses,
				r.updateOrderingUses,
			} {
				if err = funk(at); err != nil {
					break
				}
			}
		}
	}

	if err == nil {
		err = r.MatchingRuleUses().prepareStrings()
	}

	return
}

// updateEqualityUses is called by updateMatchingRuleUses and will
// extract any equality matchingRule details from the input instance
// of AttributeType and store information about this association in
// the receiver's MU stack field.  An error is returned should any
// part of this process fail.
func (r Schema) updateEqualityUses(at AttributeType) (err error) {
	if eqty := at.Equality(); !eqty.IsZero() {
		mu := r.MatchingRuleUses().get(eqty.NumericOID())

		// If the MatchingRuleUse instance does not exist,
		// create it now.
		if mu.IsZero() {
			_eqy := r.MatchingRules().get(eqty.NumericOID())
			if _eqy.IsZero() {
				err = ErrEqualityRuleNotFound
			} else if mu, err = _eqy.makeMatchingRuleUse(); err == nil {
				r.MatchingRuleUses().push(mu)
			}
		}

		// Append the AttributeType instance to the
		// MatchingRuleUse.Applies stack, assuming
		// no errors occur.
		if err == nil {
			mu.setApplies(at)
		}
	}

	return
}

// updateSubstringUses is called by updateMatchingRuleUses and will
// extract any substring matchingRule details from the input instance
// of AttributeType and store information about this association in
// the receiver's MU stack field.  An error is returned should any
// part of this process fail.
func (r Schema) updateSubstringUses(at AttributeType) (err error) {
	if substr := at.Substring(); !substr.IsZero() {
		mu := r.MatchingRuleUses().get(substr.NumericOID())

		// If the MatchingRuleUse instance does not exist,
		// create it now.
		if mu.IsZero() {
			_substr := r.MatchingRules().get(substr.NumericOID())
			if _substr.IsZero() {
				err = ErrSubstringRuleNotFound
			} else if mu, err = _substr.makeMatchingRuleUse(); err == nil {
				r.MatchingRuleUses().push(mu)
			}
		}

		// Append the AttributeType instance to the
		// MatchingRuleUse.Applies stack, assuming
		// no errors occur.
		if err == nil {
			mu.setApplies(at)
		}
	}

	return
}

// updateOrderingUses is called by updateMatchingRuleUses and will
// extract any ordering matchingRule details from the input instance
// of AttributeType and store information about this association in
// the receiver's MU stack field.  An error is returned should any
// part of this process fail.
func (r Schema) updateOrderingUses(at AttributeType) (err error) {
	if order := at.Ordering(); !order.IsZero() {
		mu := r.MatchingRuleUses().get(order.NumericOID())

		// If the MatchingRuleUse instance does not exist,
		// create it now.
		if mu.IsZero() {
			_order := r.MatchingRules().get(order.NumericOID())
			if _order.IsZero() {
				err = ErrOrderingRuleNotFound
			} else if mu, err = _order.makeMatchingRuleUse(); err == nil {
				r.MatchingRuleUses().push(mu)
			}
		}

		// Append the AttributeType instance to the
		// MatchingRuleUse.Applies stack, assuming
		// no errors occur.
		if err == nil {
			mu.setApplies(at)
		}
	}

	return
}

func (r MatchingRuleUse) setOID(_ string) {}
func (r MatchingRuleUse) macro() []string { return []string{} }

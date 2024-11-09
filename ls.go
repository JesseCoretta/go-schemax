package schemax

/*
NewLDAPSyntaxes initializes a new [LDAPSyntaxes] instance.
*/
func NewLDAPSyntaxes() LDAPSyntaxes {
	r := LDAPSyntaxes(newCollection(`ldapSyntaxes`))
	r.cast().SetPushPolicy(r.canPush)

	return r
}

/*
IsIdentifiedAs returns a Boolean value indicative of whether id matches
either the numericOID or description of the receiver instance.  Case is
not significant in the matching process.
*/
func (r LDAPSyntax) IsIdentifiedAs(id string) (ident bool) {
	if !r.IsZero() {
		ident = id == r.NumericOID() || eq(r.lDAPSyntax.Desc, id)
	}

	return
}

/*
Inventory returns an instance of [Inventory] which represents the current
inventory of [LDAPSyntaxes] instances within the receiver.

The keys are numeric OIDs, while the values are zero (0) or more string
slices, each representing a name by which the definition is known.
*/
func (r LDAPSyntaxes) Inventory() (inv Inventory) {
	inv = make(Inventory, 0)
	for i := 0; i < r.len(); i++ {
		def := r.index(i)
		inv[def.NumericOID()] = []string{def.Description()}

	}

	return
}

/*
Len returns the current integer length of the receiver instance.
*/
func (r LDAPSyntaxes) Len() int {
	return r.len()
}

func (r LDAPSyntaxes) len() int {
	return r.cast().Len()
}

/*
Contains calls [MatchingRules.Get] to return a Boolean value indicative of
a successful, non-zero retrieval of a [MatchingRule] instance -- matching
the provided id -- from within the receiver stack instance.
*/
func (r LDAPSyntaxes) Contains(id string) bool {
	return r.contains(id)
}

func (r LDAPSyntaxes) contains(id string) bool {
	return !r.get(id).IsZero()
}

/*
Marshal returns an error following an attempt to marshal the contents of
def, which may be either a [DefinitionMap] or map[string]any instance.

The receiver instance must be initialized prior to use of this method
using the [Schema.NewLDAPSyntax] method.
*/
func (r LDAPSyntax) Marshal(def any) error {
	m, err := getMarshalMap(r, def)
	if err != nil {
		return err
	}

	for k, v := range m {
		switch key := uc(k); key {
		case `DESC`:
			switch tv := v.(type) {
			case string:
				r.SetDescription(tv)
			case []string:
				r.SetDescription(tv[0])
			}
		case `NUMERICOID`:
			switch tv := v.(type) {
			case string:
				r.SetNumericOID(tv)
			case []string:
				r.SetNumericOID(tv[0])
			}
		default:
			if hasPfx(key, `X-`) {
				switch tv := v.(type) {
				case string:
					r.SetExtension(key, tv)
				case []string:
					r.SetExtension(key, tv...)
				}
			}
		}
	}

	if !r.Compliant() {
		return ErrDefNonCompliant
	}
	r.SetStringer()

	return nil
}

/*
SetData assigns x to the receiver instance. This is a general-use method and has no
specific intent beyond convenience. The contents may be subsequently accessed via the
[LDAPSyntax.Data] method.

This is a fluent method.
*/
func (r LDAPSyntax) SetData(x any) LDAPSyntax {
	if !r.IsZero() {
		r.lDAPSyntax.setData(x)
	}

	return r
}

func (r *lDAPSyntax) setData(x any) {
	r.data = x
}

/*
Data returns the underlying value (x) assigned to the receiver's data storage field. Data
can be set within the receiver instance by way of the [LDAPSyntax.SetData] method.
*/
func (r LDAPSyntax) Data() (x any) {
	if !r.IsZero() {
		x = r.lDAPSyntax.data
	}

	return
}

/*
SetSchema assigns an instance of [Schema] to the receiver instance.  This allows
internal verification of certain actions without the need for user input of
an instance of [Schema] manually at each juncture.

Note that the underlying [Schema] instance is automatically set when creating
instances of this type by way of parsing, as well as if the receiver instance
was initialized using the [Schema.NewLDAPSyntax] method.

This is a fluent method.
*/
func (r LDAPSyntax) SetSchema(schema Schema) LDAPSyntax {
	if !r.IsZero() {
		r.lDAPSyntax.setSchema(schema)
	}

	return r
}

func (r *lDAPSyntax) setSchema(schema Schema) *lDAPSyntax {
	r.schema = schema
	return r
}

/*
Schema returns the [Schema] instance associated with the receiver instance.
*/
func (r LDAPSyntax) Schema() (s Schema) {
	if !r.IsZero() {
		s = r.lDAPSyntax.getSchema()
	}

	return
}

func (r *lDAPSyntax) getSchema() (s Schema) {
	if r != nil {
		s = r.schema
	}

	return
}

/*
Compliant returns a Boolean value indicative of every [LDAPSyntax]
returning a compliant response from the [LDAPSyntax.Compliant] method.
*/
func (r LDAPSyntaxes) Compliant() bool {
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
compliant per the required clauses of [§ 4.1.5 of RFC 4512]:

  - Numeric OID must be present and valid

[§ 4.1.5 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.5
*/
func (r LDAPSyntax) Compliant() bool {
	if r.IsZero() {
		return false
	}

	return isNumericOID(r.lDAPSyntax.OID)
}

/*
SetStringer allows the assignment of an individual [Stringer] function or
method to all [LDAPSyntax] slices within the receiver stack instance.

Input of zero (0) variadic values, or an explicit nil, will overwrite all
preexisting stringer functions with the internal closure default, which is
based upon a one-time use of the [text/template] package by all receiver
slice instances.

Input of a non-nil closure function value will overwrite all preexisting
stringers.

This is a fluent method and may be used multiple times.
*/
func (r LDAPSyntaxes) SetStringer(function ...Stringer) LDAPSyntaxes {
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
func (r LDAPSyntax) SetStringer(function ...Stringer) LDAPSyntax {
	if r.Compliant() {
		r.lDAPSyntax.setStringer(function...)
	}

	return r
}

func (r *lDAPSyntax) setStringer(function ...Stringer) {
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
XOrigin returns an instance of [LDAPSyntaxes] containing only definitions
which bear the X-ORIGIN value of x. Case is not significant in the matching
process, nor is whitespace (e.g.: RFC 4517 vs. RFC4517).
*/
func (r LDAPSyntaxes) XOrigin(x string) (defs LDAPSyntaxes) {
	defs = NewLDAPSyntaxes()
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
HumanReadable returns a Boolean value indicative of whether the receiver
is considered a HUMAN-READABLE (HR) syntax.

To that end, the determination is made based on the presence of the
'X-NOT-HUMAN-READABLE' [Extension] bearing a value of 'TRUE'. In this
state, a value of false is returned, indicating the [LDAPSyntax] is NOT
considered HUMAN-READABLE.  A value of false is also returned if the
receiver is nil.

All other states return true.
*/
func (r LDAPSyntax) HumanReadable() bool {
	if r.IsZero() {
		return false
	}

	ext, _ := r.Extensions().Get(`X-NOT-HUMAN-READABLE`)
	if ext.IsZero() {
		return true
	}

	return !ext.Contains(`TRUE`) // invert evaluation
}

/*
String is a stringer method that returns the string representation
of the receiver instance.
*/
func (r LDAPSyntax) String() (def string) {
	if !r.IsZero() {
		if r.lDAPSyntax.stringer != nil {
			def = r.lDAPSyntax.stringer()
		}
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r LDAPSyntax) IsZero() bool {
	return r.lDAPSyntax == nil
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r LDAPSyntaxes) IsZero() bool {
	return r.cast().IsZero()
}

/*
OID is an alias for the [LDAPSyntax.NumericOID] method, as
[LDAPSyntax] instances do not allow for the assignment of
names.  This method exists solely to satisfy Go's interface
signature requirements.
*/
func (r LDAPSyntax) OID() string { return r.NumericOID() }

/*
Name returns an empty string, as [LDAPSyntax] definitions do not
bear names.  This method exists only to satisfy Go's interface
signature requirements.
*/
func (r LDAPSyntax) Name() string { return `` }

/*
Names returns an empty instance of [QuotedDescriptorList], as names do not
apply to [LDAPSyntax] definitions.  This method exists only to satisfy
Go's interface signature requirements.
*/
func (r LDAPSyntax) Names() QuotedDescriptorList { return QuotedDescriptorList{} }

/*
String returns the string representation of the receiver instance.
*/
func (r LDAPSyntaxes) String() string {
	return r.cast().String()
}

func (r LDAPSyntax) macro() (m []string) {
	if !r.IsZero() {
		m = r.lDAPSyntax.Macro
	}

	return
}

func (r LDAPSyntax) setOID(x string) {
	if !r.IsZero() {
		r.lDAPSyntax.OID = x
	}
}

func (r LDAPSyntaxes) canPush(x ...any) (err error) {
	if len(x) == 0 {
		return
	}

	for i := 0; i < len(x) && err == nil; i++ {
		instance := x[i]
		if ls, ok := instance.(LDAPSyntax); !ok || ls.IsZero() {
			err = ErrTypeAssert
		} else if tst := r.get(ls.NumericOID()); !tst.IsZero() {
			err = mkerr(ErrNotUnique.Error() + ": " + ls.Type() + `, ` + ls.NumericOID())
		}
	}

	return
}

/*
SetNumericOID allows the manual assignment of a numeric OID to the
receiver instance if the following are all true:

  - The input id value is a syntactically valid numeric OID
  - The receiver does not already possess a numeric OID

This is a fluent method.
*/
func (r LDAPSyntax) SetNumericOID(id string) LDAPSyntax {
	if !r.IsZero() {
		r.lDAPSyntax.setNumericOID(id)
	}

	return r
}

func (r *lDAPSyntax) setNumericOID(id string) {
	if isNumericOID(id) {
		// only set an OID when the receiver
		// lacks one (iow: no modifications)
		if len(r.OID) == 0 {
			r.OID = id
		}
	}

	return
}

/*
Extensions returns the [Extensions] instance -- if set -- within
the receiver.
*/
func (r LDAPSyntax) Extensions() (e Extensions) {
	if !r.IsZero() {
		e = r.lDAPSyntax.Extensions
	}

	return
}

/*
SetExtension assigns key x to value xstrs within the receiver's underlying
[Extensions] instance.

This is a fluent method.
*/
func (r LDAPSyntax) SetExtension(x string, xstrs ...string) LDAPSyntax {
	if !r.IsZero() {
		r.lDAPSyntax.setExtension(x, xstrs...)
	}

	return r
}

func (r *lDAPSyntax) setExtension(x string, xstrs ...string) {
	r.Extensions.Set(x, xstrs...)
}

/*
Maps returns slices of [DefinitionMap] instances.
*/
func (r LDAPSyntaxes) Maps() (defs DefinitionMaps) {
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
func (r LDAPSyntax) Map() (def DefinitionMap) {
	if r.IsZero() {
		return
	}

	def = make(DefinitionMap, 0)
	def[`NUMERICOID`] = []string{r.NumericOID()}
	def[`DESC`] = []string{r.Description()}
	def[`TYPE`] = []string{r.Type()}
	def[`HR`] = []string{bool2str(r.HumanReadable())}
	def[`RAW`] = []string{r.String()}

	// copy our extensions from receiver r
	// into destination def.
	def = mapTransferExtensions(r, def)

	// Clean up any empty fields
	def.clean()

	return
}

/*
Type returns the string literal "ldapSyntax".
*/
func (r LDAPSyntax) Type() string {
	return `ldapSyntax`
}

/*
Type returns the string literal "ldapSyntaxes".
*/
func (r LDAPSyntaxes) Type() string {
	return `ldapSyntaxes`
}

/*
Description returns the underlying (optional) descriptive text
assigned to the receiver instance.
*/
func (r LDAPSyntax) Description() (desc string) {
	if !r.IsZero() {
		desc = r.lDAPSyntax.Desc
	}

	return
}

/*
SetDescription parses desc into the underlying DESC clause within the
receiver instance.  Although a RFC 4512-compliant QuotedString is
required, the outer single-quotes need not be specified literally.

This is a fluent method.
*/
func (r LDAPSyntax) SetDescription(desc string) LDAPSyntax {
	if !r.IsZero() {
		r.lDAPSyntax.setDescription(desc)
	}

	return r
}

func (r *lDAPSyntax) setDescription(desc string) {
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
Obsolete only returns a false Boolean value, as definition obsolescence
does not apply to [LDAPSyntax] definitions.  This method exists only to
satisfy Go's interface signature requirements.
*/
func (r LDAPSyntax) Obsolete() bool { return false }

/*
Get returns an instance of [LDAPSyntax] based upon a search for id within
the receiver stack instance.

The return instance, if not nil, was retrieved based upon a textual match of
either of the following:

  - the numeric OID of an [LDAPSyntax] and the provided id
  - the description text -- minus whitespace -- and the provided id

The return instance is nil if no match was made.

Case is not significant in the matching process.
*/
func (r LDAPSyntaxes) Get(id string) LDAPSyntax {
	return r.get(id)
}

func (r LDAPSyntaxes) get(id string) (ls LDAPSyntax) {
	for i := 0; i < r.len() && ls.IsZero(); i++ {
		if _ls := r.index(i); !_ls.IsZero() {
			if eq(_ls.lDAPSyntax.OID, id) {
				ls = _ls
			} else {
				if desc := repAll(_ls.lDAPSyntax.Desc, ` `, ``); eq(desc, id) {
					ls = _ls
				}
			}
		}
	}

	return
}

/*
NumericOID returns the string representation of the numeric OID
held by the receiver instance.
*/
func (r LDAPSyntax) NumericOID() (noid string) {
	if !r.IsZero() {
		noid = r.lDAPSyntax.OID
	}

	return
}

/*
Index returns the instance of [LDAPSyntax] found within the
receiver stack instance at index N.  If no instance is found at
the index specified, a zero [LDAPSyntax] instance is returned.
*/
func (r LDAPSyntaxes) Index(idx int) LDAPSyntax {
	return r.index(idx)
}

func (r LDAPSyntaxes) index(idx int) (ls LDAPSyntax) {
	slice, found := r.cast().Index(idx)
	if found {
		if _ls, ok := slice.(LDAPSyntax); ok {
			ls = _ls
		}
	}

	return
}

/*
Push returns an error following an attempt to push an [LDAPSyntax]
into the receiver stack instance.
*/
func (r LDAPSyntaxes) Push(ls any) error {
	return r.push(ls)
}

func (r LDAPSyntaxes) push(x any) (err error) {
	switch tv := x.(type) {
	case LDAPSyntax:
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
Parse returns an error following an attempt to parse raw into the receiver
instance.

Note that the receiver MUST possess a [Schema] reference prior to the execution
of this method.

Also note that successful execution of this method does NOT automatically push
the receiver into any [LDAPSyntaxes] stack, nor does it automatically execute
the [LDAPSyntax.SetStringer] method, leaving these tasks to the user.  If the
automatic handling of these tasks is desired, see the [Schema.ParseLDAPSyntax]
method as an alternative.
*/
func (r LDAPSyntax) Parse(raw string) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
		return
	}

	if r.getSchema().IsZero() {
		err = ErrNilSchemaRef
		return
	}

	err = r.lDAPSyntax.parse(raw)

	return
}

func (r *lDAPSyntax) parse(raw string) error {
	// parseLS wraps the antlr4512 LDAPSyntax parser/lexer
	mp, err := parseLS(raw)
	if err == nil {
		// We received the parsed data from ANTLR (mp).
		// Now we need to marshal it into the receiver.
		var def LDAPSyntax
		if def, err = r.schema.marshalLS(mp); err == nil {
			r.OID = def.NumericOID()
			_r := LDAPSyntax{r}
			_r.replace(def)
		}
	}

	return err
}

/*
NewLDAPSyntax initializes and returns a new instance of [LDAPSyntax],
ready for manual assembly.  This method need not be used when creating
new [LDAPSyntax] instances by way of parsing, as that is handled on an
internal basis.

Use of this method does NOT automatically push the return instance into
the [Schema.LDAPSyntaxes] stack; this is left to the user.

Unlike the package-level [NewLDAPSyntax] function, this method will
automatically reference its originating [Schema] instance (the receiver).
This negates the need for manual use of the [LDAPSyntax.SetSchema]
method.

This is the recommended means of creating a new [LDAPSyntax] instance
wherever a single [Schema] is being used, which represents most use cases.
*/
func (r Schema) NewLDAPSyntax() LDAPSyntax {
	return NewLDAPSyntax().SetSchema(r)
}

/*
NewLDAPSyntax initializes and returns a new instance of [LDAPSyntax],
ready for manual assembly.  This method need not be used when creating
new [LDAPSyntax] instances by way of parsing, as that is handled on an
internal basis.

Use of this function does not automatically reference the "parent" [Schema]
instance, leaving it up to the user to invoke the [LDAPSyntax.SetSchema]
method manually.

When interacting with a single [Schema] instance, which represents most use
cases, use of the [Schema.NewLDAPSyntax] method is PREFERRED over use of
this package-level function.

However certain migration efforts, schema audits and other such activities
may require distinct associations of [LDAPSyntax] instances with specific
[Schema] instances. Use of this function allows the user to specify the
appropriate [Schema] instance at a later point for a specific instance of
an [LDAPSyntax] instance.
*/
func NewLDAPSyntax() LDAPSyntax {
	ls := LDAPSyntax{newLDAPSyntax()}
	ls.lDAPSyntax.Extensions.setDefinition(ls)
	return ls
}

func newLDAPSyntax() *lDAPSyntax {
	return &lDAPSyntax{
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
within the [LDAPSyntax] envelope type, thus all references to the
receiver instance within various stacks will be preserved.

This is a fluent method.
*/
func (r LDAPSyntax) Replace(x LDAPSyntax) LDAPSyntax {
	if !r.Schema().Options().Positive(AllowOverride) {
		return r
	}

	if !r.IsZero() && x.Compliant() {
		r.lDAPSyntax.replace(x)
	}

	return r
}

func (r *lDAPSyntax) replace(x LDAPSyntax) {
	if r.OID == `` {
		return
	} else if r.OID != x.NumericOID() {
		return
	}

	r.OID = x.lDAPSyntax.OID
	r.Desc = x.lDAPSyntax.Desc
	r.Extensions = x.lDAPSyntax.Extensions
	r.data = x.lDAPSyntax.data
	r.schema = x.lDAPSyntax.schema
	r.stringer = x.lDAPSyntax.stringer
	r.synQual = x.lDAPSyntax.synQual
	r.data = x.lDAPSyntax.data
}

/*
prepareString returns a string an an error indicative of an attempt
to represent the receiver instance as a string using [text/template].
*/
func (r *lDAPSyntax) prepareString() (str string, err error) {
	buf := newBuf()
	t := newTemplate(`ldapSyntax`).
		Funcs(funcMap(map[string]any{
			`ExtensionSet`: r.Extensions.tmplFunc,
		}))

	if t, err = t.Parse(lDAPSyntaxTmpl); err == nil {
		if err = t.Execute(buf, struct {
			Definition *lDAPSyntax
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
LDAPSyntaxes returns the [LDAPSyntaxes] instance from within the
receiver instance.
*/
func (r Schema) LDAPSyntaxes() (lss LDAPSyntaxes) {
	slice, _ := r.cast().Index(ldapSyntaxesIndex)
	lss, _ = slice.(LDAPSyntaxes)
	return
}

/*
SetSyntaxQualifier assigns an instance of [SyntaxQualifier] to the receiver
instance. A nil value may be passed to disable syntax checking capabilities.

See the [LDAPSyntax.QualifySyntax] method for details on making active use
of the [SyntaxQualifier] capabilities.

This is a fluent method.
*/
func (r LDAPSyntax) SetSyntaxQualifier(function SyntaxQualifier) LDAPSyntax {
	if !r.IsZero() {
		r.lDAPSyntax.setSyntaxQualifier(function)
	}

	return r
}

func (r *lDAPSyntax) setSyntaxQualifier(function SyntaxQualifier) {
	r.synQual = function
}

/*
QualifySyntax returns an error instance following an analysis of the
input value using the [SyntaxQualifier] instance previously assigned to
the receiver instance.

If a [SyntaxQualifier] is not assigned to the receiver instance, the
[ErrNilSyntaxQualifier] error is returned if and when this method is
executed. Otherwise, an error is returned based on the custom
[SyntaxQualifier] error handler devised within the user provided
closure.

A nil error return always indicates valid input value syntax.

See the [LDAPSyntax.SetSyntaxQualifier] method for information regarding
the assignment of an instance of [SyntaxQualifier] to the receiver.
*/
func (r LDAPSyntax) QualifySyntax(value any) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
	} else if r.lDAPSyntax.synQual == nil {
		err = ErrNilSyntaxQualifier
	} else {
		err = r.lDAPSyntax.synQual(value)
	}

	return
}

/*
LoadLDAPSyntaxes returns an error following an attempt to load all
built-in [LDAPSyntax] definitions into the receiver instance.
*/
func (r Schema) LoadLDAPSyntaxes() error {
	return r.loadSyntaxes()
}

/*
loadSyntaxes returns an error following an attempt to load all
built-in ldapSyntax definitions found within this package into
the receiver instance.
*/
func (r Schema) loadSyntaxes() (err error) {
	if !r.IsZero() {
		funks := []func() error{
			r.loadRFC4517Syntaxes,
			r.loadRFC4523Syntaxes,
			r.loadRFC4530Syntaxes,
			r.loadRFC2307Syntaxes,
		}

		for i := 0; i < len(funks) && err == nil; i++ {
			err = funks[i]()
		}
	}

	return
}

/*
LoadRFC2307Syntaxes returns an error following an attempt to load
all RFC 2307 [LDAPSyntax] slices into the receiver instance.
*/
func (r Schema) LoadRFC2307Syntaxes() error {
	return r.loadRFC2307Syntaxes()
}

func (r Schema) loadRFC2307Syntaxes() (err error) {
	for k, v := range rfc2307Macros {
		r.Macros().Set(k, v)
	}

	var i int
	for i = 0; i < len(rfc2307LDAPSyntaxes) && err == nil; i++ {
		ls := rfc2307LDAPSyntaxes[i]
		err = r.ParseLDAPSyntax(string(ls))
	}

	if want := rfc2307LDAPSyntaxes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC2307 LDAPSyntaxes parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4517Syntaxes returns an error following an attempt to load
all RFC 4517 [LDAPSyntax] slices into the receiver instance.
*/
func (r Schema) LoadRFC4517Syntaxes() error {
	return r.loadRFC4517Syntaxes()
}

func (r Schema) loadRFC4517Syntaxes() (err error) {

	var i int
	for i = 0; i < len(rfc4517LDAPSyntaxes) && err == nil; i++ {
		ls := rfc4517LDAPSyntaxes[i]
		err = r.ParseLDAPSyntax(string(ls))
	}

	if want := rfc4517LDAPSyntaxes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4517 LDAPSyntaxes parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4523Syntaxes returns an error following an attempt to load
all RFC 4523 [LDAPSyntax] slices into the receiver instance.
*/
func (r Schema) LoadRFC4523Syntaxes() error {
	return r.loadRFC4523Syntaxes()
}

func (r Schema) loadRFC4523Syntaxes() (err error) {

	var i int
	for i = 0; i < len(rfc4523LDAPSyntaxes) && err == nil; i++ {
		ls := rfc4523LDAPSyntaxes[i]
		err = r.ParseLDAPSyntax(string(ls))
	}

	if want := rfc4523LDAPSyntaxes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4523 LDAPSyntaxes parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4530Syntaxes returns an error following an attempt to load
all RFC 4530 [LDAPSyntax] slices into the receiver instance.
*/
func (r Schema) LoadRFC4530Syntaxes() error {
	return r.loadRFC4530Syntaxes()
}

func (r Schema) loadRFC4530Syntaxes() (err error) {

	var i int
	for i = 0; i < len(rfc4530LDAPSyntaxes) && err == nil; i++ {
		ls := rfc4530LDAPSyntaxes[i]
		err = r.ParseLDAPSyntax(string(ls))
	}

	if want := rfc4530LDAPSyntaxes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4530 LDAPSyntaxes parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

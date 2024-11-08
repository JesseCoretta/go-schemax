package schemax

/*
NewAttributeType initializes and returns a new instance of [AttributeType],
ready for manual assembly.  This method need not be used when creating
new [AttributeType] instances by way of parsing, as that is handled on an
internal basis.

Use of this method does NOT automatically push the return instance into
the [Schema.AttributeTypes] stack; this is left to the user.

Unlike the package-level [NewAttributeType] function, this method will
automatically reference its originating [Schema] instance (the receiver).
This negates the need for manual use of the [AttributeType.SetSchema]
method.

This is the recommended means of creating a new [AttributeType] instance
wherever a single [Schema] is being used, which represents most use cases.
*/
func (r Schema) NewAttributeType() AttributeType {
	return NewAttributeType().SetSchema(r)
}

/*
NewAttributeType initializes and returns a new instance of [AttributeType],
ready for manual assembly.  This method need not be used when creating
new [AttributeType] instances by way of parsing, as that is handled on an
internal basis.

Use of this function does not automatically reference the "parent" [Schema]
instance, leaving it up to the user to invoke the [AttributeType.SetSchema]
method manually.

When interacting with a single [Schema] instance, which represents most use
cases, use of the [Schema.NewAttributeType] method is PREFERRED over use of
this package-level function.

However certain migration efforts, schema audits and other such activities
may require distinct associations of [AttributeType] instances with specific
[Schema] instances. Use of this function allows the user to specify the
appropriate [Schema] instance at a later point for a specific instance of
an [AttributeType] instance.
*/
func NewAttributeType() AttributeType {
	at := AttributeType{newAttributeType()}
	at.attributeType.Extensions.setDefinition(at)
	return at
}

/*
newAttributeType initializes a new instance of the private (embedded)
type *attributeType.
*/
func newAttributeType() *attributeType {
	return &attributeType{
		Name:       NewName(),
		Extensions: NewExtensions(),
	}
}

/*
Parse returns an error following an attempt to parse raw into the receiver
instance.

Note that the receiver MUST possess a [Schema] reference prior to the execution
of this method.

Also note that successful execution of this method does NOT automatically push
the receiver into any [AttributeTypes] stack, nor does it automatically execute
the [AttributeType.SetStringer] method, leaving these tasks to the user.  If the
automatic handling of these tasks is desired, see the [Schema.ParseAttributeType]
method as an alternative.
*/
func (r AttributeType) Parse(raw string) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
		return
	}

	if r.getSchema().IsZero() {
		err = ErrNilSchemaRef
		return
	}

	err = r.attributeType.parse(raw)

	return
}

func (r *attributeType) parse(raw string) error {
	// parseAT wraps the antlr4512 AttributeType parser/lexer
	mp, err := parseAT(raw)
	if err == nil {
		// We received the parsed data from ANTLR (mp).
		// Now we need to marshal it into the receiver.
		var def AttributeType
		if def, err = r.schema.marshalAT(mp); err == nil {
			r.OID = def.NumericOID()
			_r := AttributeType{r}
			_r.replace(def)
		}
	}

	return err
}

/*
Replace overrides the receiver with x. Both must bear an identical
numeric OID and x MUST be compliant.

Note that the relevant [Schema] instance must be configured to allow
definition override by way of the [AllowOverride] bit setting.  See
the [Schema.Options] method for a means of accessing the settings
value.

Note that this method does not reallocate a new pointer instance
within the [AttributeType] envelope type, thus all references to the
receiver instance within various stacks will be preserved.

This is a fluent method.
*/
func (r AttributeType) Replace(x AttributeType) AttributeType {
	if !r.Schema().Options().Positive(AllowOverride) {
		return r
	}

	if !r.IsZero() && x.Compliant() {
		r.attributeType.replace(x)
	}

	return r
}

func (r *attributeType) replace(x AttributeType) {
	if r.OID == `` {
		return
	} else if r.OID != x.NumericOID() {
		return
	}

	r.OID = x.attributeType.OID
	r.Macro = x.attributeType.Macro
	r.Name = x.attributeType.Name
	r.Desc = x.attributeType.Desc
	r.Obsolete = x.attributeType.Obsolete
	r.MUB = x.attributeType.MUB
	r.Single = x.attributeType.Single
	r.Collective = x.attributeType.Collective
	r.NoUserMod = x.attributeType.NoUserMod
	r.SuperType = x.attributeType.SuperType
	r.Equality = x.attributeType.Equality
	r.Substring = x.attributeType.Substring
	r.Ordering = x.attributeType.Ordering
	r.Syntax = x.attributeType.Syntax
	r.Usage = x.attributeType.Usage
	r.Extensions = x.attributeType.Extensions
	r.data = x.attributeType.data
	r.schema = x.attributeType.schema
	r.stringer = x.attributeType.stringer
	r.valQual = x.attributeType.valQual
	r.data = x.attributeType.data
}

/*
SetSchema assigns an instance of [Schema] to the receiver instance.  This allows
internal verification of certain actions without the need for user input of
an instance of [Schema] manually at each juncture.

Note that the underlying [Schema] instance is automatically set when creating
instances of this type by way of parsing, as well as if the receiver instance
was initialized using the [Schema.NewAttributeType] method.

This is a fluent method.
*/
func (r AttributeType) SetSchema(schema Schema) AttributeType {
	if !r.IsZero() {
		r.attributeType.setSchema(schema)
	}

	return r
}

func (r *attributeType) setSchema(schema Schema) {
	r.schema = schema
}

/*
Schema returns the [Schema] instance associated with the receiver instance.
*/
func (r AttributeType) Schema() (s Schema) {
	if !r.IsZero() {
		s = r.attributeType.getSchema()
	}

	return
}

func (r *attributeType) getSchema() (s Schema) {
	if r != nil {
		s = r.schema
	}

	return
}

/*
SetData assigns x to the receiver instance. This is a general-use method and has no
specific intent beyond convenience. The contents may be subsequently accessed via the
[AttributeType.Data] method.

This is a fluent method.
*/
func (r AttributeType) SetData(x any) AttributeType {
	if !r.IsZero() {
		r.attributeType.setData(x)
	}

	return r
}

func (r *attributeType) setData(x any) {
	r.data = x
}

/*
Data returns the underlying value (x) assigned to the receiver's data storage field. Data
can be set within the receiver instance by way of the [AttributeType.SetData] method.
*/
func (r AttributeType) Data() (x any) {
	if !r.IsZero() {
		x = r.attributeType.data
	}

	return
}

/*
QualifySyntax wraps [LDAPSyntax.QualifySyntax].
*/
func (r AttributeType) QualifySyntax(value any) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
	} else if esyn := r.EffectiveSyntax(); esyn.IsZero() {
		err = ErrNilDef
	} else {
		err = esyn.QualifySyntax(value)
	}

	return
}

/*
EqualityAssertion returns an error following an attempt to perform an
EQUALITY assertion match upon value1 and value2 using the effective
[MatchingRule] honored by the receiver instance.
*/
func (r AttributeType) EqualityAssertion(value1, value2 any) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
	} else if emr := r.EffectiveEquality(); emr.IsZero() {
		err = ErrNilDef
	} else {
		err = emr.Assertion(value1, value2)
	}

	return
}

/*
SubstringAssertion returns an error following an attempt to perform an
SUBSTR assertion match upon value1 and value2 using the effective
[MatchingRule] honored by the receiver instance.
*/
func (r AttributeType) SubstringAssertion(value1, value2 any) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
	} else if emr := r.EffectiveSubstring(); emr.IsZero() {
		err = ErrNilDef
	} else {
		err = emr.Assertion(value1, value2)
	}

	return
}

/*
OrderingAssertion returns an error following an attempt to perform an
ORDERING assertion match upon value1 and value2 using the effective
[MatchingRule] honored by the receiver instance.
*/
func (r AttributeType) OrderingAssertion(value1, value2 any) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
	} else if emr := r.EffectiveOrdering(); emr.IsZero() {
		err = ErrNilDef
	} else {
		err = emr.Assertion(value1, value2)
	}

	return
}

/*
QualifyValue returns an error instance following an analysis of the input
value using the [ValueQualifier] instance previously assigned to the receiver
instance.

If a [ValueQualifier] is not assigned to the receiver instance, the
[ErrNilValueQualifier] error is returned if and when this method is
executed. Otherwise, an error is returned based on the custom
[ValueQualifier] error handler devised by the author.

A nil error return always indicates valid input value syntax.

See the [AttributeType.SetValueQualifier] for information regarding
the assignment of an instance of [ValueQualifier] to the receiver.
*/
func (r AttributeType) QualifyValue(value any) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
	} else if r.attributeType.valQual == nil {
		err = ErrNilValueQualifier
	} else {
		err = r.attributeType.valQual(value)
	}

	return
}

/*
SetValueQualifier assigns an instance of [ValueQualifier] to the receiver
instance. A nil value may be passed to disable general value qualification
capabilities.

See the [AttributeType.QualifyValue] method for details on making active
use of the [ValueQualifier] capabilities.

This is a fluent method.
*/
func (r AttributeType) SetValueQualifier(function ValueQualifier) AttributeType {
	if !r.IsZero() {
		r.attributeType.setValueQualifier(function)
	}

	return r
}

func (r *attributeType) setValueQualifier(function ValueQualifier) {
	r.valQual = function
}

func (r AttributeType) schema() (s Schema) {
	if !r.IsZero() {
		s = r.attributeType.schema
	}

	return
}

/*
Inventory returns an instance of [Inventory] which represents the current
inventory of [AttributeType] instances within the receiver.

The keys are numeric OIDs, while the values are zero (0) or more string
slices, each representing a name by which the definition is known.
*/
func (r AttributeTypes) Inventory() (inv Inventory) {
	inv = make(Inventory, 0)
	for i := 0; i < r.len(); i++ {
		def := r.index(i)
		inv[def.NumericOID()] = def.Names().List()

	}

	return
}

// stackage closure func - do not exec directly (use String method)
func (r AttributeTypes) oIDsStringer(_ ...any) string {
	slice := r.index(0)
	hd := slice.Schema().Options().Positive(HangingIndents)
	id := r.cast().ID()
	if hd && id != `at_oidlist` {
		return r.oIDsStringerPretty(len(id))
	}

	return r.oIDsStringerStd()
}

/*
factory default stackage closure func for oid lists - do not exec directly.
*/
func (r AttributeTypes) oIDsStringerStd(_ ...any) string {
	var _present []string
	for i := 0; i < r.len(); i++ {
		_present = append(_present, r.index(i).OID())
	}

	present := newStringBuilder()

	switch len(_present) {
	case 0:
	case 1:
		present.WriteString(_present[0])
	default:
		padchar := string(rune(32))
		if !r.cast().IsPadded() {
			padchar = ``
		}

		joined := join(_present, padchar+`$`+padchar)
		present.WriteString(`(` + padchar + joined + padchar + `)`)
	}

	return present.String()
}

/*
prettified stackage closure func for oid lists - do not exec directly.

prepare a custom [stackage.PresentationPolicy] instance for our input
[QuotedDescriptorList] stack to convert the following:

	( cn $ sn $ l $ c $ st )

... into ...

	( cn
	$ sn
	$ l
	$ c
	$ st )

This has no effect if the stack has only one member, producing something
like:

	cn
*/
func (r AttributeTypes) oIDsStringerPretty(lead int) string {
	L := r.Len()
	switch L {
	case 0:
		return ``
	case 1:
		return r.Index(0).OID()
	}

	present := newStringBuilder()
	num := lead + 5
	for idx := 0; idx < L; idx++ {
		sl := r.Index(idx).OID()
		if idx == 0 {
			present.WriteString(`( ` + sl)
			present.WriteRune(rune(10))
			continue
		}

		for i := 0; i < num; i++ {
			present.WriteRune(rune(32))
		}

		present.WriteString(`$ ` + sl)
		if idx == L-1 {
			present.WriteString(` )`)
			break // no newline on last line
		}

		present.WriteRune(rune(10))
	}

	return present.String()
}

// stackage closure func - do not exec directly.
func (r AttributeTypes) canPush(x ...any) (err error) {
	if len(x) == 0 {
		return
	}

	for i := 0; i < len(x) && err == nil; i++ {
		instance := x[i]
		if at, ok := instance.(AttributeType); !ok || at.IsZero() {
			err = ErrTypeAssert
		} else {
			if tst := r.get(at.NumericOID()); !tst.IsZero() {
				err = mkerr(ErrNotUnique.Error() + ": " + at.Type() + `, ` + at.NumericOID())
			}
		}
	}

	return
}

/*
AttributeTypes returns the [AttributeTypes] instance from within
the receiver instance.
*/
func (r Schema) AttributeTypes() (ats AttributeTypes) {
	slice, _ := r.cast().Index(attributeTypesIndex)
	ats, _ = slice.(AttributeTypes)
	return
}

/*
Len returns the current integer length of the receiver instance.
*/
func (r AttributeTypes) Len() int {
	return r.len()
}

func (r AttributeTypes) len() int {
	return r.cast().Len()
}

/*
String is a stringer method that returns the string representation
of the receiver instance.
*/
func (r AttributeTypes) String() string {
	return r.cast().String()
}

/*
Index returns the instance of [AttributeType] found within the
receiver stack instance at index N.  If no instance is found at
the index specified, a zero [AttributeType] instance is returned.
*/
func (r AttributeTypes) Index(idx int) AttributeType {
	return r.index(idx)
}

func (r AttributeTypes) index(idx int) (at AttributeType) {
	slice, found := r.cast().Index(idx)
	if found {
		if _at, ok := slice.(AttributeType); ok {
			at = _at
		}
	}

	return
}

/*
Push returns an error following an attempt to push an [AttributeType]
into the receiver stack instance.
*/
func (r AttributeTypes) Push(at any) error {
	return r.push(at)
}

func (r AttributeTypes) push(x any) (err error) {
	switch tv := x.(type) {
	case AttributeType:
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
Contains calls [AttributeTypes.Get] to return a Boolean value indicative
of a successful, non-zero retrieval of an [AttributeType] instance -- matching
the provided id -- from within the receiver stack instance.
*/
func (r AttributeTypes) Contains(id string) bool {
	return r.contains(id)
}

func (r AttributeTypes) contains(id string) bool {
	return !r.get(id).IsZero()
}

/*
Get returns an instance of [AttributeType] based upon a search for id within
the receiver stack instance.

The return instance, if not nil, was retrieved based upon a textual match of
the principal identifier of an [AttributeType] and the provided id.

The return instance is nil if no match was made.

Case is not significant in the matching process.
*/
func (r AttributeTypes) Get(id string) AttributeType {
	return r.get(id)
}

func (r AttributeTypes) get(id string) (at AttributeType) {
	for i := 0; i < r.len() && at.IsZero(); i++ {
		if _at := r.index(i); !_at.IsZero() {
			if _at.attributeType.OID == id {
				at = _at
			} else if _at.attributeType.Name.contains(id) {
				at = _at
			}
		}
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r AttributeType) IsZero() bool {
	return r.attributeType == nil
}

/*
IsIdentifiedAs returns a Boolean value indicative of whether id matches
either the numericOID or descriptor of the receiver instance.  Case is
not significant in the matching process.
*/
func (r AttributeType) IsIdentifiedAs(id string) (ident bool) {
	if !r.IsZero() {
		ident = id == r.NumericOID() || r.attributeType.Name.contains(id)
	}

	return
}

/*
NoUserModification returns a Boolean value indicative of whether the receiver
instance's NO-USER-MODIFICATION clause evaluates as TRUE.  In such a scenario,
only the DSA may manage such values.
*/
func (r AttributeType) NoUserModification() (o bool) {
	if !r.IsZero() {
		o = r.attributeType.NoUserMod
	}

	return
}

/*
Collective returns a Boolean value indicative of whether the receiver
is COLLECTIVE.  A value of true is mutually exclusive of SINGLE-VALUE'd
[AttributeType] instances.
*/
func (r AttributeType) Collective() (o bool) {
	if !r.IsZero() {
		o = r.attributeType.Collective
	}

	return
}

/*
SingleValue returns a Boolean value indicative of whether the receiver
is set to only allow one (1) value to be assigned to an entry using this
type.  A value of true is mutually exclusive of COLLECTIVE [AttributeType]
instances.
*/
func (r AttributeType) SingleValue() (o bool) {
	if !r.IsZero() {
		o = r.attributeType.Single
	}

	return
}

/*
Obsolete returns a Boolean value indicative of definition obsolescence.
*/
func (r AttributeType) Obsolete() (o bool) {
	if !r.IsZero() {
		o = r.attributeType.Obsolete
	}

	return
}

/*
Name returns the string form of the principal name of the receiver instance, if set.
*/
func (r AttributeType) Name() (id string) {
	if !r.IsZero() {
		id = r.attributeType.Name.index(0)
	}

	return
}

/*
Names returns the underlying instance of [QuotedDescriptorList] from within
the receiver.
*/
func (r AttributeType) Names() (names QuotedDescriptorList) {
	return r.attributeType.Name
}

/*
NewAttributeTypes initializes and returns a new [AttributeTypes] instance,
configured to allow the storage of all [AttributeType] instances.
*/
func NewAttributeTypes() AttributeTypes {
	r := AttributeTypes(newCollection(`attributeTypes`))
	r.cast().SetPushPolicy(r.canPush)

	return r
}

/*
NewAttributeTypeOIDList initializes and returns a new [AttributeTypes] that has
been configured to allow the storage of arbitrary [AttributeType] instances.
*/
func NewAttributeTypeOIDList(label ...string) AttributeTypes {
	name := `at_oidlist`
	if len(label) > 0 {
		name = label[0]
	}

	r := AttributeTypes(newOIDList(name))
	r.cast().
		SetPushPolicy(r.canPush).
		SetPresentationPolicy(r.oIDsStringer)

	return r
}

/*
OID returns the string representation of an OID -- which is either a
numeric OID or descriptor -- that is held by the receiver instance.
*/
func (r AttributeType) OID() (oid string) {
	if !r.IsZero() {
		oid = r.NumericOID() // default
		if r.attributeType.Name.len() > 0 {
			oid = r.attributeType.Name.index(0)
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
func (r AttributeType) SetNumericOID(id string) AttributeType {
	if !r.IsZero() {
		r.attributeType.setNumericOID(id)
	}

	return r
}

func (r *attributeType) setNumericOID(id string) {
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
SetExtension assigns key x to value xstrs within the receiver's underlying
[Extensions] instance.

This is a fluent method.
*/
func (r AttributeType) SetExtension(x string, xstrs ...string) AttributeType {
	if !r.IsZero() {
		r.attributeType.setExtension(x, xstrs...)
	}

	return r
}

func (r *attributeType) setExtension(x string, xstrs ...string) {
	r.Extensions.Set(x, xstrs...)
}

/*
Extensions returns the [Extensions] instance -- if set -- within
the receiver.
*/
func (r AttributeType) Extensions() (e Extensions) {
	if !r.IsZero() {
		e = r.attributeType.Extensions
	}

	return
}

/*
SetName assigns the provided names to the receiver instance.

Name instances must conform to RFC 4512 descriptor format but
need not be quoted.

This is a fluent method.
*/
func (r AttributeType) SetName(x ...string) AttributeType {
	if !r.IsZero() {
		r.attributeType.setName(x...)
	}

	return r
}

func (r *attributeType) setName(x ...string) {
	for i := 0; i < len(x); i++ {
		r.Name.Push(x[i])
	}
}

/*
NumericOID returns the string representation of the numeric OID
held by the receiver instance.
*/
func (r AttributeType) NumericOID() (noid string) {
	if !r.IsZero() {
		noid = r.attributeType.OID
	}

	return
}

/*
Description returns the underlying (optional) descriptive text
assigned to the receiver instance.
*/
func (r AttributeType) Description() (desc string) {
	if !r.IsZero() {
		desc = r.attributeType.Desc
	}

	return
}

/*
SetDescription parses desc into the underlying DESC clause within the
receiver instance.  Although a RFC 4512-compliant QuotedString is
required, the outer single-quotes need not be specified literally.

This is a fluent method.
*/
func (r AttributeType) SetDescription(desc string) AttributeType {
	if !r.IsZero() {
		r.attributeType.setDescription(desc)
	}

	return r
}

func (r *attributeType) setDescription(desc string) {
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
SetStringer allows the assignment of an individual [Stringer] function or
method to all [AttributeType] slices within the receiver stack instance.

Input of zero (0) variadic values, or an explicit nil, will overwrite all
preexisting stringer functions with the internal closure default, which is
based upon a one-time use of the [text/template] package by all receiver
slice instances.

Input of a non-nil closure function value will overwrite all preexisting
stringers.

This is a fluent method and may be used multiple times.
*/
func (r AttributeTypes) SetStringer(function ...Stringer) AttributeTypes {
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
func (r AttributeType) SetStringer(function ...Stringer) AttributeType {
	if r.Compliant() {
		r.attributeType.setStringer(function...)
	}

	return r
}

func (r *attributeType) setStringer(function ...Stringer) {
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
				// Return a static value.
				return str
			}
		}
		return
	}

	// assign user-provided closure
	r.stringer = stringer
}

/*
XOrigin returns an instance of [AttributeTypes] containing only definitions
which bear the X-ORIGIN value of x. Case is not significant in the matching
process, nor is whitespace (e.g.: RFC 4517 vs. RFC4517).
*/
func (r AttributeTypes) XOrigin(x string) (defs AttributeTypes) {
	defs = NewAttributeTypes()
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
prepareString returns a string an an error indicative of an attempt
to represent the receiver instance as a string using [text/template].
*/
func (r *attributeType) prepareString() (str string, err error) {
	buf := newBuf()
	t := newTemplate(`attributeType`).
		Funcs(funcMap(map[string]any{
			`Substring`:    func() string { return r.Substring.OID() },
			`Ordering`:     func() string { return r.Ordering.OID() },
			`Equality`:     func() string { return r.Equality.OID() },
			`Syntax`:       func() string { return r.Syntax.NumericOID() },
			`MUB`:          func() string { return `{` + uitoa(r.MUB) + `}` },
			`SuperType`:    func() string { return r.SuperType.OID() },
			`ExtensionSet`: r.Extensions.tmplFunc,
			`Obsolete`:     func() bool { return r.Obsolete },
			`IsSingleVal`:  func() bool { return r.Single },
			`Collective`:   func() bool { return r.Collective },
			`IsNoUserMod`:  func() bool { return r.NoUserMod },
			`Usage`:        func() string { return AttributeType{r}.Usage() },
		}))

	if t, err = t.Parse(attributeTypeTmpl); err == nil {
		if err = t.Execute(buf, struct {
			Definition *attributeType
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
String is a stringer method that returns the string representation of
the receiver instance.  A zero-value indicates an invalid receiver, or
that the [AttributeType.SetStringer] method was not used during MANUAL
composition of the receiver.
*/
func (r AttributeType) String() (def string) {
	if !r.IsZero() {
		if r.attributeType.stringer != nil {
			def = r.attributeType.stringer()
		}
	}

	return
}

/*
MinimumUpperBounds returns the unsigned integer form of the receiver's
"size limit", if set. A non-zero value indicates a specific maximum
has been declared.

Note that this method can only be used upon a receiver instance that
has an effective [LDAPSyntax] known to be human-readable.
*/
func (r AttributeType) MinimumUpperBounds() (mub uint) {
	if !r.IsZero() {
		mub = r.attributeType.MUB
	}

	return
}

/*
SetMinimumUpperBounds assigns the provided value -- which may be an int
or uint -- to the receiver instance.

This is a fluent method.
*/
func (r AttributeType) SetMinimumUpperBounds(mub any) AttributeType {
	if !r.IsZero() {
		if syn := r.EffectiveSyntax(); syn.HumanReadable() {
			r.attributeType.setMinimumUpperBounds(mub)
		}
	}

	return r
}

func (r *attributeType) setMinimumUpperBounds(mub any) {
	switch tv := mub.(type) {
	case int:
		if tv >= 0 {
			r.MUB = uint(tv)
		}
	case uint:
		r.MUB = tv
	}
}

/*
Syntax returns the [LDAPSyntax] instance held (directly) by the receiver.
*/
func (r AttributeType) Syntax() (syntax LDAPSyntax) {
	if !r.IsZero() {
		syntax = r.attributeType.Syntax
	}

	return
}

/*
EffectiveSyntax returns the [LDAPSyntax] instance held (directly or indirectly)
by the receiver.

If the receiver does not directly reference an [LDAPSyntax] instance, this method
walks the super type chain until it encounters a superior [AttributeType] which
directly references an [LDAPSyntax].

If the return instance returns a Boolean value of true following an execution
of [LDAPSyntax.IsZero], this means there is NO effective [LDAPSyntax] specified
anywhere in the receiver's super type chain.
*/
func (r AttributeType) EffectiveSyntax() LDAPSyntax {
	if syntax := r.Syntax(); !syntax.IsZero() {
		return syntax
	} else if r.SuperType().IsZero() {
		return syntax
	}

	return r.SuperType().EffectiveSyntax()
}

/*
SetSyntax assigns x to the receiver instance as an instance of [LDAPSyntax].

This is a fluent method.
*/
func (r AttributeType) SetSyntax(x any) AttributeType {
	if !r.IsZero() {
		r.attributeType.setSyntax(x)
	}

	return r
}

func (r *attributeType) setSyntax(x any) {
	var def LDAPSyntax
	switch tv := x.(type) {
	case string:
		if !r.schema.IsZero() {
			def = r.schema.LDAPSyntaxes().get(tv)
		}
	case LDAPSyntax:
		def = tv
	}

	if def.Compliant() {
		r.Syntax = def
	}
}

/*
Equality returns the underlying instance of [MatchingRule] if set within
the receiver instance. If unset, a zero instance is returned.
*/
func (r AttributeType) Equality() (eql MatchingRule) {
	if !r.IsZero() {
		eql = r.attributeType.Equality
	}

	return
}

/*
EffectiveEquality returns the EQUALITY [MatchingRule] instance held (directly or
indirectly) by the receiver.

If the receiver does not directly reference a [MatchingRule] instance, this method
walks the super type chain until it encounters a superior [AttributeType] which
directly references an EQUALITY [MatchingRule].

If the return instance returns a Boolean value of true following an execution of
[MatchingRule.IsZero], this means there is NO effective EQUALITY [MatchingRule]
specified anywhere in the receiver's super type chain.
*/
func (r AttributeType) EffectiveEquality() MatchingRule {
	if mr := r.Equality(); !mr.IsZero() {
		return mr
	} else if r.SuperType().IsZero() {
		// end of the chain
		return mr
	}

	return r.SuperType().EffectiveEquality()
}

/*
SetEquality sets the equality [MatchingRule] of the receiver instance.

Input x may be the string representation of the desired equality matchingRule,
by name or numeric OID.  This requires an underlying instance of [Schema] set
within the receiver instance. The [Schema] instance is expected to possess
the referenced [MatchingRule], else a zero instance will ultimately be returned
and discarded.

Input x may also be a bonafide, non-zero instance of an equality [MatchingRule].
This requires no underlying [Schema] instance, but may allow the introduction of
bogus [MatchingRule] instances.

This is a fluent method.
*/
func (r AttributeType) SetEquality(x any) AttributeType {
	if !r.IsZero() {
		r.attributeType.setEquality(x)
	}

	return r
}

func (r *attributeType) setEquality(x any) {
	var def MatchingRule
	switch tv := x.(type) {
	case string:
		if !r.schema.IsZero() {
			def = r.schema.MatchingRules().get(tv)
		}
	case MatchingRule:
		def = tv
	}

	if def.Compliant() {
		r.Equality = def
	}
}

/*
Substring returns the underlying instance of [MatchingRule] if set within
the receiver instance. If unset, a zero instance is returned.
*/
func (r AttributeType) Substring() (sub MatchingRule) {
	if !r.IsZero() {
		sub = r.attributeType.Substring
	}

	return
}

/*
EffectiveSubstring returns the SUBSTR [MatchingRule] instance held (directly or
indirectly) by the receiver.

If the receiver does not directly reference a [MatchingRule] instance, this method
walks the super type chain until it encounters a superior [AttributeType] which
directly references a SUBSTR [MatchingRule].

If the return instance returns a Boolean value of true following an execution of
[MatchingRule.IsZero], this means there is NO effective SUBSTR [MatchingRule]
specified anywhere in the receiver's super type chain.
*/
func (r AttributeType) EffectiveSubstring() MatchingRule {
	if mr := r.Substring(); !mr.IsZero() {
		return mr
	} else if r.SuperType().IsZero() {
		// end of the chain
		return mr
	}

	return r.SuperType().EffectiveSubstring()
}

/*
SetSubstring sets the substring [MatchingRule] of the receiver instance.

Input x may be the string representation of the desired substring matchingRule,
by name or numeric OID.  This requires an underlying instance of [Schema] set
within the receiver instance. The [Schema] instance is expected to possess
the referenced [MatchingRule], else a zero instance will ultimately be returned
and discarded.

Input x may also be a bonafide, non-zero instance of a substring [MatchingRule].
This requires no underlying [Schema] instance, but may allow the introduction of
bogus [MatchingRule] instances.

This is a fluent method.
*/
func (r AttributeType) SetSubstring(x any) AttributeType {
	if !r.IsZero() {
		r.attributeType.setSubstring(x)
	}

	return r
}

func (r *attributeType) setSubstring(x any) {
	var def MatchingRule
	switch tv := x.(type) {
	case string:
		if !r.schema.IsZero() {
			def = r.schema.MatchingRules().get(tv)
		}
	case MatchingRule:
		def = tv
	}

	if def.Compliant() {
		r.Substring = def
	}
}

/*
Ordering returns the underlying instance of [MatchingRule] if set within
the receiver instance. If unset, a zero instance is returned.
*/
func (r AttributeType) Ordering() (ord MatchingRule) {
	if !r.IsZero() {
		ord = r.attributeType.Ordering
	}

	return
}

/*
EffectiveOrdering returns the ORDERING [MatchingRule] instance held (directly or
indirectly) by the receiver.

If the receiver does not directly reference a [MatchingRule] instance, this method
walks the super type chain until it encounters a superior [AttributeType] which
directly references an ORDERING [MatchingRule].

If the return instance returns a Boolean value of true following an execution of
[MatchingRule.IsZero], this means there is NO effective ORDERING [MatchingRule]
specified anywhere in the receiver's super type chain.
*/
func (r AttributeType) EffectiveOrdering() MatchingRule {
	if mr := r.Ordering(); !mr.IsZero() {
		return mr
	} else if r.SuperType().IsZero() {
		// end of the chain
		return mr
	}

	return r.SuperType().EffectiveOrdering()
}

/*
SetOrdering sets the ordering [MatchingRule] of the receiver instance.

Input x may be the string representation of the desired ordering matchingRule,
by name or numeric OID.  This requires an underlying instance of [Schema] set
within the receiver instance. The [Schema] instance is expected to possess
the referenced [MatchingRule], else a zero instance will ultimately be returned
and discarded.

Input x may also be a bonafide, non-zero instance of an ordering [MatchingRule].
This requires no underlying [Schema] instance, but may allow the introduction of
bogus [MatchingRule] instances.

This is a fluent method.
*/
func (r AttributeType) SetOrdering(x any) AttributeType {
	if !r.IsZero() {
		r.attributeType.setOrdering(x)
	}

	return r
}

func (r *attributeType) setOrdering(x any) {
	var def MatchingRule
	switch tv := x.(type) {
	case string:
		if !r.schema.IsZero() {
			def = r.schema.MatchingRules().get(tv)
		}
	case MatchingRule:
		def = tv
	}

	if def.Compliant() {
		r.Ordering = def
	}
}

/*
SuperType returns the underlying instance of [AttributeType] if set within
the receiver instance as its super type. If unset, a zero instance is returned.
*/
func (r AttributeType) SuperType() (sup AttributeType) {
	if !r.IsZero() {
		sup = r.attributeType.SuperType
	}

	return
}

/*
SubTypes returns an instance of [AttributeTypes] containing slices of
[AttributeType] instances that are direct subordinates to the receiver
instance. As such, this method is essentially the inverse of the
[AttributeType.SuperType] method.

The super chain is NOT traversed beyond immediate subordinate instances.

Note that the relevant [Schema] instance must have been set using the
[AttributeType.SetSchema] method prior to invocation of this method.
Should this requirement remain unfulfilled, the return instance will
be a zero instance.
*/
func (r AttributeType) SubTypes() (subs AttributeTypes) {
	if !r.IsZero() {
		subs = NewAttributeTypeOIDList()
		ats := r.schema().AttributeTypes()
		for i := 0; i < ats.Len(); i++ {
			typ := ats.Index(i)
			if typ.SuperType().NumericOID() == r.NumericOID() {
				subs.Push(typ)
			}
		}
	}

	return
}

/*
SuperChain returns an [AttributeTypes] stack of [AttributeType] instances
which make up the super type chain of the receiver instance.
*/
func (r AttributeType) SuperChain() (sups AttributeTypes) {
	if !r.IsZero() {
		sups = r.attributeType.SuperType.superChain()
	}

	return
}

func (r AttributeType) superChain() (sups AttributeTypes) {
	sups = NewAttributeTypes()
	sups.Push(r)
	if !r.SuperType().IsZero() {
		x := r.SuperType().SuperChain()
		for i := 0; i < x.Len(); i++ {
			sups.Push(x.Index(i))
		}
	}
	return
}

/*
SetSuperType sets the super type of the receiver to the value provided. Valid
input types are string, to represent an RFC 4512 OID residing in the underlying
[Schema] instance, or an actual [AttributeType] instance already obtained or crafted.

This is a fluent method.
*/
func (r AttributeType) SetSuperType(x any) AttributeType {
	if !r.IsZero() {
		r.attributeType.setSuperType(x)
	}

	return r
}

func (r *attributeType) setSuperType(x any) {
	var def AttributeType
	switch tv := x.(type) {
	case string:
		if !r.schema.IsZero() {
			def = r.schema.AttributeTypes().get(tv)
		}
	case AttributeType:
		def = tv
	}

	if def.Compliant() {
		r.SuperType = def
	}
}

/*
SetSingleValue assigns the input value to the underlying SINGLE-VALUE clause
within the receiver.

Input types may be bool, or string representations of bool. When strings
are used, case is not significant.

Note that a value of true will be ignored if the receiver is a collective
[AttributeType].

This is a fluent method.
*/
func (r AttributeType) SetSingleValue(x any) AttributeType {
	if !r.IsZero() {
		r.attributeType.setBoolean(`sv`, x)
	}

	return r
}

/*
SetCollective assigns the input value to the underlying COLLECTIVE clause
within the receiver.

Input types may be bool, or string representations of bool. When strings
are used, case is not significant.

Note that a value of true will be ignored if the receiver is a single-valued
[AttributeType].

This is a fluent method.
*/
func (r AttributeType) SetCollective(x any) AttributeType {
	if !r.IsZero() {
		r.attributeType.setBoolean(`c`, x)
	}

	return r
}

/*
SetNoUserModification assigns the input value to the underlying NO-USER-MODIFICATION clause
within the receiver.

Input types may be bool, or string representations of bool. When strings
are used, case is not significant.

This is a fluent method.
*/
func (r AttributeType) SetNoUserModification(x any) AttributeType {
	if !r.IsZero() {
		r.attributeType.setBoolean(`num`, x)
	}

	return r
}

/*
SetObsolete sets the receiver's obsolescence state to true.

Obsolescence cannot be unset.

This is a fluent method.
*/
func (r AttributeType) SetObsolete(x any) AttributeType {
	if !r.IsZero() {
		r.attributeType.setObsolete()
	}

	return r
}

func (r *attributeType) setObsolete() {
	if !r.Obsolete {
		r.Obsolete = true
	}
}

func (r *attributeType) setBoolean(t string, x any) {

	var Bool bool
	switch tv := x.(type) {
	case string:
		if eq(tv, `true`) {
			Bool = true
		}
	case bool:
		Bool = tv
	default:
		return
	}

	switch t {
	case `sv`:
		if !r.Collective {
			r.Single = Bool
		}
	case `c`:
		if !r.Single {
			r.Collective = Bool
		}
	case `num`:
		if !r.NoUserMod {
			r.NoUserMod = Bool
		}
	case `obs`:
		if !r.Obsolete {
			r.Obsolete = Bool
		}
	}
}

/*
SetUsage assigns the specified USAGE to the receiver. Input types
may be string, int or uint.

	1, or directoryOperation
	2, or distributedOperation
	3, or dSAOperation

Any other value results in assignment of the userApplications USAGE.

This is a fluent method.
*/
func (r AttributeType) SetUsage(u any) AttributeType {
	if !r.IsZero() {
		r.attributeType.setUsage(u)
	}

	return r
}

func (r *attributeType) setUsage(u any) {
	switch tv := u.(type) {
	case string:
		switch lc(tv) {
		case `directoryoperation`:
			r.Usage = DirectoryOperationUsage
		case `distributedoperation`:
			r.Usage = DistributedOperationUsage
		case `dsaoperation`:
			r.Usage = DSAOperationUsage
		default:
			r.Usage = UserApplicationsUsage
		}
	case uint:
		r.setUsage(int(tv))
	case int:
		switch tv {
		case 1:
			r.Usage = DirectoryOperationUsage
		case 2:
			r.Usage = DistributedOperationUsage
		case 3:
			r.Usage = DSAOperationUsage
		default:
			r.Usage = UserApplicationsUsage
		}
	}
}

/*
Type returns the string literal "attributeType".
*/
func (r AttributeType) Type() string {
	return `attributeType`
}

/*
Type returns the string literal "attributeTypes".
*/
func (r AttributeTypes) Type() string {
	return `attributeTypes`
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r AttributeTypes) IsZero() bool {
	return r.cast().IsZero()
}

/*
Usage returns the string representation of the underlying USAGE if set
within the receiver instance. If unset, a zero string -- which implies
use of the [UserApplicationsUsage] USAGE value by default -- is returned.
*/
func (r AttributeType) Usage() (usage string) {
	if !r.IsZero() {
		switch v := r.attributeType.Usage; int(v) {
		case 0:
			break // zero is default (userApplications)
		case 1:
			usage = `directoryOperation`
		case 2:
			usage = `distributedOperation`
		case 3:
			usage = `dSAOperation`
		}
	}

	return
}

/*
Compliant returns a Boolean value indicative of every [AttributeType]
returning a compliant response from the [AttributeType.Compliant] method.
*/
func (r AttributeTypes) Compliant() bool {
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
compliant per the required clauses of [§ 4.1.2 of RFC 4512]:

  - Numeric OID must be present and valid
  - Specified EQUALITY, SUBSTR and ORDERING [MatchingRule] instances must be COMPLIANT
  - Specified [LDAPSyntax] MUST be COMPLIANT
  - [AttributeType], if not COLLECTIVE, cannot extend from a COLLECTIVE super type

Additional consideration is given to RFC 3671 in that an [AttributeType]
shall not be both COLLECTIVE and SINGLE-VALUE'd.

[§ 4.1.2 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.2
*/
func (r AttributeType) Compliant() bool {
	if r.IsZero() {
		return false
	}

	if !isNumericOID(r.attributeType.OID) {
		return false
	}

	syn := r.schema().LDAPSyntaxes().get(r.Syntax().NumericOID())
	if !syn.IsZero() && !syn.Compliant() {
		return false
	}

	for _, mr := range []MatchingRule{
		r.schema().MatchingRules().get(r.Equality().NumericOID()),
		r.schema().MatchingRules().get(r.Ordering().NumericOID()),
		r.schema().MatchingRules().get(r.Substring().NumericOID()),
	} {
		if !mr.IsZero() && !mr.Compliant() {
			return false
		}
	}

	collective := r.Collective()
	sup := r.schema().AttributeTypes().get(r.SuperType().NumericOID())
	if !sup.IsZero() {
		if sup.Collective() && !collective {
			// Non collective types cannot derive from
			// collective super types.
			return false
		}

		if !sup.Compliant() {
			return false
		}
	}

	// Any combination of SV/C is permitted
	// EXCEPT for BOTH.  See RFC 3671.
	return !(r.SingleValue() && collective)
}

/*
Maps returns slices of [DefinitionMap] instances.
*/
func (r AttributeTypes) Maps() (defs DefinitionMaps) {
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
func (r AttributeType) Map() (def DefinitionMap) {
	if r.IsZero() {
		return
	}

	def = make(DefinitionMap, 0)
	def[`NUMERICOID`] = []string{r.NumericOID()}
	def[`NAME`] = r.Names().List()
	def[`DESC`] = []string{r.Description()}
	def[`OBSOLETE`] = []string{bool2str(r.Obsolete())}
	def[`SUP`] = []string{r.SuperType().OID()}
	def[`EQUALITY`] = []string{r.Equality().OID()}
	def[`SUBSTR`] = []string{r.Substring().OID()}
	def[`ORDERING`] = []string{r.Ordering().OID()}
	def[`SYNTAX`] = []string{r.Syntax().NumericOID()}
	def[`MUB`] = []string{itoa(int(r.MinimumUpperBounds()))}
	def[`SINGLE-VALUE`] = []string{bool2str(r.SingleValue())}
	def[`COLLECTIVE`] = []string{bool2str(r.Collective())}
	def[`NO-USER-MODIFICATION`] = []string{bool2str(r.NoUserModification())}
	def[`USAGE`] = []string{r.Usage()}
	def[`TYPE`] = []string{r.Type()}
	def[`RAW`] = []string{r.String()}

	// copy our extensions from receiver r
	// into destination def.
	def = mapTransferExtensions(r, def)

	// Clean up any empty fields
	def.clean()

	return
}

func (r AttributeType) macro() (m []string) {
	if !r.IsZero() {
		m = r.attributeType.Macro
	}

	return
}

func (r AttributeType) setOID(x string) {
	if !r.IsZero() {
		r.attributeType.OID = x
	}
}

/*
LoadAttributeTypes returns an error following an attempt to load all
built-in [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadAttributeTypes() error {
	return r.loadAttributeTypes()
}

/*
loadAttributeTypes returns an error following an attempt to load
all built-in [AttributeType] slices into the receiver instance.
*/
func (r Schema) loadAttributeTypes() (err error) {
	if !r.IsZero() {
		funks := []func() error{
			r.loadRFC4512AttributeTypes,
			r.loadX501AttributeTypes,
			r.loadRFC2079AttributeTypes,
			r.loadRFC2798AttributeTypes,
			r.loadRFC3045AttributeTypes,
			r.loadRFC3672AttributeTypes,
			r.loadRFC4519AttributeTypes,
			r.loadRFC2307AttributeTypes,
			r.loadRFC3671AttributeTypes,
			r.loadRFC4523AttributeTypes,
			r.loadRFC4524AttributeTypes,
			r.loadRFC4530AttributeTypes,
			r.loadRFC2589AttributeTypes,
			r.loadRFC4403AttributeTypes,
			r.loadRFC5020AttributeTypes,
		}

		for i := 0; i < len(funks) && err == nil; i++ {
			err = funks[i]()
		}
	}

	return
}

/*
LoadX501AttributeTypes returns an error following an attempt to load all
[ITU-T Rec. X.501] [AttributeType] slices into the receiver instance.

[ITU-T Rec. X.501]: https://www.itu.int/rec/T-REC-X.501
*/
func (r Schema) LoadX501AttributeTypes() error {
	return r.loadX501AttributeTypes()
}

func (r Schema) loadX501AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(x501AttributeTypes) && err == nil; i++ {
		at := x501AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := x501AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of X.501 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC2079AttributeTypes returns an error following an attempt to
load all RFC 2079 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC2079AttributeTypes() error {
	return r.loadRFC2079AttributeTypes()
}

func (r Schema) loadRFC2079AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(rfc2079AttributeTypes) && err == nil; i++ {
		at := rfc2079AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc2079AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC2079 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC2307AttributeTypes returns an error following an attempt to
load all RFC 2307 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC2307AttributeTypes() error {
	return r.loadRFC2307AttributeTypes()
}

func (r Schema) loadRFC2307AttributeTypes() (err error) {
	for k, v := range rfc2307Macros {
		r.Macros().Set(k, v)
	}

	var i int
	for i = 0; i < len(rfc2307AttributeTypes) && err == nil; i++ {
		at := rfc2307AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc2307AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC2307 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC2798AttributeTypes returns an error following an attempt to
load all RFC 2798 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC2798AttributeTypes() error {
	return r.loadRFC2798AttributeTypes()
}

func (r Schema) loadRFC2798AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(rfc2798AttributeTypes) && err == nil; i++ {
		at := rfc2798AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc2798AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC2798 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC3045AttributeTypes returns an error following an attempt to
load all RFC 3045 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC3045AttributeTypes() error {
	return r.loadRFC3045AttributeTypes()
}

func (r Schema) loadRFC3045AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(rfc3045AttributeTypes) && err == nil; i++ {
		at := rfc3045AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc3045AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC3045 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC3671AttributeTypes returns an error following an attempt to
load all RFC 3671 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC3671AttributeTypes() error {
	return r.loadRFC3671AttributeTypes()
}

func (r Schema) loadRFC3671AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(rfc3671AttributeTypes) && err == nil; i++ {
		at := rfc3671AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc3671AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC3671 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC3672AttributeTypes returns an error following an attempt to
load all RFC 3672 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC3672AttributeTypes() error {
	return r.loadRFC3672AttributeTypes()
}

func (r Schema) loadRFC3672AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(rfc3672AttributeTypes) && err == nil; i++ {
		at := rfc3672AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc3672AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC3672 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4512AttributeTypes returns an error following an attempt to
load all RFC 4512 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC4512AttributeTypes() error {
	return r.loadRFC4512AttributeTypes()
}

func (r Schema) loadRFC4512AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(rfc4512AttributeTypes) && err == nil; i++ {
		at := rfc4512AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc4512AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4512 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4519AttributeTypes returns an error following an attempt to
load all RFC 4519 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC4519AttributeTypes() error {
	return r.loadRFC4519AttributeTypes()
}

func (r Schema) loadRFC4519AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(rfc4519AttributeTypes) && err == nil; i++ {
		at := rfc4519AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc4519AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4519 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4523AttributeTypes returns an error following an attempt to
load all RFC 4523 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC4523AttributeTypes() error {
	return r.loadRFC4523AttributeTypes()
}

func (r Schema) loadRFC4523AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(rfc4523AttributeTypes) && err == nil; i++ {
		at := rfc4523AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc4523AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4523 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4524AttributeTypes returns an error following an attempt to
load all RFC 4524 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC4524AttributeTypes() error {
	return r.loadRFC4524AttributeTypes()
}

func (r Schema) loadRFC4524AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(rfc4524AttributeTypes) && err == nil; i++ {
		at := rfc4524AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc4524AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4524 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4530AttributeTypes returns an error following an attempt to
load all RFC 4530 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC4530AttributeTypes() error {
	return r.loadRFC4530AttributeTypes()
}

func (r Schema) loadRFC4530AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(rfc4530AttributeTypes) && err == nil; i++ {
		at := rfc4530AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc4530AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4530 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC2589AttributeTypes returns an error following an attempt to
load all RFC 2589 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC2589AttributeTypes() error {
	return r.loadRFC2589AttributeTypes()
}

func (r Schema) loadRFC2589AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(rfc2589AttributeTypes) && err == nil; i++ {
		at := rfc2589AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc2589AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC2589 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC5020AttributeTypes returns an error following an attempt to
load all RFC 5020 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC5020AttributeTypes() error {
	return r.loadRFC5020AttributeTypes()
}

func (r Schema) loadRFC5020AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(rfc5020AttributeTypes) && err == nil; i++ {
		at := rfc5020AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc5020AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC5020 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4403AttributeTypes returns an error following an attempt to
load all RFC 4403 [AttributeType] slices into the receiver instance.
*/
func (r Schema) LoadRFC4403AttributeTypes() error {
	return r.loadRFC4403AttributeTypes()
}

func (r Schema) loadRFC4403AttributeTypes() (err error) {

	var i int
	for i = 0; i < len(rfc4403AttributeTypes) && err == nil; i++ {
		at := rfc4403AttributeTypes[i]
		err = r.ParseAttributeType(string(at))
	}

	if want := rfc4403AttributeTypes.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4403 AttributeTypes parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

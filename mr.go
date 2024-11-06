package schemax

/*
NewMatchingRules initializes a new [MatchingRules] instance.
*/
func NewMatchingRules() MatchingRules {
	r := MatchingRules(newCollection(`matchingRules`))
	r.cast().SetPushPolicy(r.canPush)

	return r
}

/*
NewMatchingRule initializes and returns a new instance of [MatchingRule],
ready for manual assembly.  This method need not be used when creating
new [MatchingRule] instances by way of parsing, as that is handled on an
internal basis.

Use of this method does NOT automatically push the return instance into
the [Schema.MatchingRules] stack; this is left to the user.

Unlike the package-level [NewMatchingRule] function, this method will
automatically reference its originating [Schema] instance (the receiver).
This negates the need for manual use of the [MatchingRule.SetSchema]
method.

This is the recommended means of creating a new [MatchingRule] instance
wherever a single [Schema] is being used, which represents most use cases.
*/
func (r Schema) NewMatchingRule() MatchingRule {
	return NewMatchingRule().SetSchema(r)
}

/*
NewMatchingRule initializes and returns a new instance of [MatchingRule],
ready for manual assembly.  This method need not be used when creating
new [MatchingRule] instances by way of parsing, as that is handled on an
internal basis.

Use of this function does not automatically reference the "parent" [Schema]
instance, leaving it up to the user to invoke the [MatchingRule.SetSchema]
method manually.

When interacting with a single [Schema] instance, which represents most use
cases, use of the [Schema.NewMatchingRule] method is PREFERRED over use of
this package-level function.

However certain migration efforts, schema audits and other such activities
may require distinct associations of [MatchingRule] instances with specific
[Schema] instances. Use of this function allows the user to specify the
appropriate [Schema] instance at a later point for a specific instance of
a [MatchingRule] instance.
*/
func NewMatchingRule() MatchingRule {
	mr := MatchingRule{newMatchingRule()}
	mr.matchingRule.Extensions.setDefinition(mr)
	return mr
}

func newMatchingRule() *matchingRule {
	return &matchingRule{
		Name:       NewName(),
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
within the [MatchingRule] envelope type, thus all references to the
receiver instance within various stacks will be preserved.

This is a fluent method.
*/
func (r MatchingRule) Replace(x MatchingRule) MatchingRule {
	if !r.Schema().Options().Positive(AllowOverride) {
		return r
	}

	if !r.IsZero() && x.Compliant() {
		r.matchingRule.replace(x)
	}

	return r
}

func (r *matchingRule) replace(x MatchingRule) {
	if r.OID == `` {
		return
	} else if r.OID != x.NumericOID() {
		return
	}

	r.OID = x.matchingRule.OID
	r.Macro = x.matchingRule.Macro
	r.Name = x.matchingRule.Name
	r.Desc = x.matchingRule.Desc
	r.Obsolete = x.matchingRule.Obsolete
	r.Syntax = x.matchingRule.Syntax
	r.Extensions = x.matchingRule.Extensions
	r.data = x.matchingRule.data
	r.schema = x.matchingRule.schema
	r.stringer = x.matchingRule.stringer
	r.data = x.matchingRule.data
	r.assMatch = x.matchingRule.assMatch
}

/*
Parse returns an error following an attempt to parse raw into the receiver
instance.

Note that the receiver MUST possess a [Schema] reference prior to the execution
of this method.

Also note that successful execution of this method does NOT automatically push
the receiver into any [MatchingRules] stack, nor does it automatically execute
the [MatchingRule.SetStringer] method, leaving these tasks to the user.  If the
automatic handling of these tasks is desired, see the [Schema.ParseMatchingRule]
method as an alternative.
*/
func (r MatchingRule) Parse(raw string) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
		return
	}

	if r.getSchema().IsZero() {
		err = ErrNilSchemaRef
		return
	}

	err = r.matchingRule.parse(raw)

	return
}

func (r *matchingRule) parse(raw string) error {
	// parseMR wraps the antlr4512 MatchingRule parser/lexer
	mp, err := parseMR(raw)
	if err == nil {
		// We received the parsed data from ANTLR (mp).
		// Now we need to marshal it into the receiver.
		var def MatchingRule
		if def, err = r.schema.marshalMR(mp); err == nil {
			r.OID = def.NumericOID()
			_r := MatchingRule{r}
			_r.replace(def)
		}
	}

	return err
}

/*
NumericOID returns the string representation of the numeric OID
held by the receiver instance.
*/
func (r MatchingRule) NumericOID() (noid string) {
	if !r.IsZero() {
		noid = r.matchingRule.OID
	}

	return
}

/*
SetData assigns x to the receiver instance. This is a general-use method and has no
specific intent beyond convenience. The contents may be subsequently accessed via the
[MatchingRule.Data] method.

This is a fluent method.
*/
func (r MatchingRule) SetData(x any) MatchingRule {
	if !r.IsZero() {
		r.matchingRule.setData(x)
	}

	return r
}

func (r *matchingRule) setData(x any) {
	r.data = x
}

/*
Data returns the underlying value (x) assigned to the receiver's data storage field. Data
can be set within the receiver instance by way of the [MatchingRule.SetData] method.
*/
func (r MatchingRule) Data() (x any) {
	if !r.IsZero() {
		x = r.matchingRule.data
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
func (r MatchingRule) SetNumericOID(id string) MatchingRule {
	if !r.IsZero() {
		r.matchingRule.setNumericOID(id)
	}

	return r
}

func (r *matchingRule) setNumericOID(id string) {
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
OID returns the string representation of an OID -- which is either a
numeric OID or descriptor -- that is held by the receiver instance.
*/
func (r MatchingRule) OID() (oid string) {
	if !r.IsZero() {
		oid = r.NumericOID() // default
		if r.matchingRule.Name.len() > 0 {
			oid = r.matchingRule.Name.index(0)
		}
	}

	return
}

/*
Syntax returns the [LDAPSyntax] reference held by the receiver instance.
*/
func (r MatchingRule) Syntax() (syntax LDAPSyntax) {
	if !r.IsZero() {
		syntax = r.matchingRule.Syntax
	}

	return
}

/*
SetSyntax assigns x to the receiver instance as an instance of [LDAPSyntax].

This is a fluent method.
*/
func (r MatchingRule) SetSyntax(x any) MatchingRule {
	if !r.IsZero() {
		r.matchingRule.setSyntax(x)
	}

	return r
}

/*
Assertion returns an error instance following an analysis of the two
input values provided in the context of an assertion match based on
the receiver instance.

If an [AssertionMatcher] is not assigned to the receiver instance, the
[ErrNilAssertionMatcher] error is returned if and when this method is
executed. Otherwise, an error is returned based on the custom
[AssertionMatcher] error handler devised within the user provided
closure.

A nil error return always indicates valid input value syntax.

See the [MatchingRule.SetAssertionMatcher] for information regarding the
assignment of an instance of [AssertionMatcher] to the receiver.
*/
func (r MatchingRule) Assertion(value1, value2 any) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
	} else if r.matchingRule.assMatch == nil {
		err = ErrNilAssertionMatcher
	} else {
		err = r.matchingRule.assMatch(value1, value2)
	}

	return
}

/*
SetAssertionMatcher assigns an instance of [AssertionMatcher] to the receiver
instance. A nil value may be passed to assertion matching capabilities.

See the [MatchingRule.Assertion] method for details on making active use of
the [AssertionMatcher] capabilities.

This is a fluent method.
*/
func (r MatchingRule) SetAssertionMatcher(function AssertionMatcher) MatchingRule {
	if !r.IsZero() {
		r.matchingRule.setAssertionMatcher(function)
	}

	return r
}

func (r *matchingRule) setAssertionMatcher(function AssertionMatcher) {
	r.assMatch = function
}

func (r *matchingRule) setSyntax(x any) {
	var def LDAPSyntax
	switch tv := x.(type) {
	case string:
		if !r.schema.IsZero() {
			def = r.schema.LDAPSyntaxes().get(tv)
		}
	case LDAPSyntax:
		def = tv
	}

	if !def.IsZero() {
		r.Syntax = def
	}
}

/*
SetSchema assigns an instance of [Schema] to the receiver instance.  This allows
internal verification of certain actions without the need for user input of
an instance of [Schema] manually at each juncture.

Note that the underlying [Schema] instance is automatically set when creating
instances of this type by way of parsing, as well as if the receiver instance
was initialized using the [Schema.NewMatchingRule] method.

This is a fluent method.
*/
func (r MatchingRule) SetSchema(schema Schema) MatchingRule {
	if !r.IsZero() {
		r.matchingRule.setSchema(schema)
	}

	return r
}

func (r *matchingRule) setSchema(schema Schema) {
	r.schema = schema
}

/*
Schema returns the [Schema] instance associated with the receiver instance.
*/
func (r MatchingRule) Schema() (s Schema) {
	if !r.IsZero() {
		s = r.matchingRule.getSchema()
	}

	return
}

func (r *matchingRule) getSchema() (s Schema) {
	if r != nil {
		s = r.schema
	}

	return
}

/*
SetDescription parses desc into the underlying DESC clause within the
receiver instance.  Although a RFC 4512-compliant QuotedString is
required, the outer single-quotes need not be specified literally.

This is a fluent method.
*/
func (r MatchingRule) SetDescription(desc string) MatchingRule {
	if !r.IsZero() {
		r.matchingRule.setDescription(desc)
	}

	return r
}

func (r *matchingRule) setDescription(desc string) {
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
Description returns the underlying (optional) descriptive text
assigned to the receiver instance.
*/
func (r MatchingRule) Description() (desc string) {
	if !r.IsZero() {
		desc = r.matchingRule.Desc
	}
	return
}

/*
IsIdentifiedAs returns a Boolean value indicative of whether id matches
either the numericOID or descriptor of the receiver instance.  Case is
not significant in the matching process.
*/
func (r MatchingRule) IsIdentifiedAs(id string) (ident bool) {
	if !r.IsZero() {
		ident = id == r.NumericOID() || r.matchingRule.Name.contains(id)
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r MatchingRule) IsZero() bool {
	return r.matchingRule == nil
}

/*
MatchingRules returns the [MatchingRules] instance from within
the receiver instance.
*/
func (r Schema) MatchingRules() (mrs MatchingRules) {
	slice, _ := r.cast().Index(matchingRulesIndex)
	mrs, _ = slice.(MatchingRules)
	return
}

/*
Maps returns slices of [DefinitionMap] instances.
*/
func (r MatchingRules) Maps() (defs DefinitionMaps) {
	defs = make(DefinitionMaps, r.Len())
	for i := 0; i < r.Len(); i++ {
		defs[i] = r.Index(i).Map()
	}

	return
}

/*
Map marshals the receiver instance into an instance of
DefinitionMap.
*/
func (r MatchingRule) Map() (def DefinitionMap) {
	if r.IsZero() {
		return
	}

	def = make(DefinitionMap, 0)
	def[`NUMERICOID`] = []string{r.NumericOID()}
	def[`NAME`] = r.Names().List()
	def[`DESC`] = []string{r.Description()}
	def[`OBSOLETE`] = []string{bool2str(r.Obsolete())}
	def[`SYNTAX`] = []string{r.Syntax().NumericOID()}
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
Obsolete returns a Boolean value indicative of definition obsolescence.
*/
func (r MatchingRule) Obsolete() (o bool) {
	if !r.IsZero() {
		o = r.matchingRule.Obsolete
	}

	return
}

/*
SetObsolete sets the receiver instance to OBSOLETE if not already set. Note that obsolescence cannot be unset.

This is a fluent method.
*/
func (r MatchingRule) SetObsolete() MatchingRule {
	if !r.IsZero() {
		r.matchingRule.setObsolete()
	}

	return r
}

func (r *matchingRule) setObsolete() {
	if !r.Obsolete {
		r.Obsolete = true
	}
}

/*
SetName assigns the provided names to the receiver instance.

Name instances must conform to RFC 4512 descriptor format but
need not be quoted.

This is a fluent method.
*/
func (r MatchingRule) SetName(x ...string) MatchingRule {
	if !r.IsZero() {
		r.matchingRule.setName(x...)
	}

	return r
}

func (r *matchingRule) setName(x ...string) {
	for i := 0; i < len(x); i++ {
		r.Name.Push(x[i])
	}
}

/*
Names returns the underlying instance of [QuotedDescriptorList] from within
the receiver.
*/
func (r MatchingRule) Names() (names QuotedDescriptorList) {
	if !r.IsZero() {
		names = r.matchingRule.Name
	}

	return
}

/*
Name returns the string form of the principal name of the receiver instance, if set.
*/
func (r MatchingRule) Name() (id string) {
	if !r.IsZero() {
		id = r.matchingRule.Name.index(0)
	}

	return
}

/*
SetExtension assigns key x to value xstrs within the receiver's underlying
[Extensions] instance.

This is a fluent method.
*/
func (r MatchingRule) SetExtension(x string, xstrs ...string) MatchingRule {
	if !r.IsZero() {
		r.matchingRule.setExtension(x, xstrs...)
	}

	return r
}

func (r *matchingRule) setExtension(x string, xstrs ...string) {
	r.Extensions.Set(x, xstrs...)
}

/*
Extensions returns the [Extensions] instance -- if set -- within
the receiver.
*/
func (r MatchingRule) Extensions() (e Extensions) {
	if !r.IsZero() {
		e = r.matchingRule.Extensions
	}

	return
}

func (r MatchingRule) macro() (m []string) {
	if !r.IsZero() {
		m = r.matchingRule.Macro
	}

	return
}

func (r MatchingRule) setOID(x string) {
	if !r.IsZero() {
		r.matchingRule.OID = x
	}
}

/*
LoadMatchingRules returns an error following to attempt to load all
built-in [MatchingRule] definitions into the receiver instance.
*/
func (r Schema) LoadMatchingRules() error {
	return r.loadMatchingRules()
}

/*
loadMatchingRules returns an error following an attempt to load
all matchingRule definitions found within this package into the
receiver instance.
*/
func (r Schema) loadMatchingRules() (err error) {
	if !r.IsZero() {
		funks := []func() error{
			r.loadRFC2307MatchingRules,
			r.loadRFC4517MatchingRules,
			r.loadRFC4523MatchingRules,
			r.loadRFC4530MatchingRules,
		}

		for i := 0; i < len(funks) && err == nil; i++ {
			err = funks[i]()
		}
	}

	return
}

/*
LoadRFC2307MatchingRules returns an error following an attempt to load
all RFC 2307 [MatchingRule] slices into the receiver instance.
*/
func (r Schema) LoadRFC2307MatchingRules() error {
	return r.loadRFC2307MatchingRules()
}

func (r Schema) loadRFC2307MatchingRules() (err error) {
	for k, v := range rfc2307Macros {
		r.Macros().Set(k, v)
	}

	var i int
	for i = 0; i < len(rfc2307MatchingRules) && err == nil; i++ {
		mr := rfc2307MatchingRules[i]
		err = r.ParseMatchingRule(string(mr))
	}

	if want := rfc2307MatchingRules.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC2307 MatchingRules parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4517MatchingRules returns an error following an attempt to load
all RFC 4517 [MatchingRule] slices into the receiver instance.
*/
func (r Schema) LoadRFC4517MatchingRules() error {
	return r.loadRFC4517MatchingRules()
}

func (r Schema) loadRFC4517MatchingRules() (err error) {

	var i int
	for i = 0; i < len(rfc4517MatchingRules) && err == nil; i++ {
		mr := rfc4517MatchingRules[i]
		err = r.ParseMatchingRule(string(mr))
	}

	if want := rfc4517MatchingRules.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4517 MatchingRules parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4523MatchingRules returns an error following an attempt to load
all RFC 4523 [MatchingRule] slices into the receiver instance.
*/
func (r Schema) LoadRFC4523MatchingRules() error {
	return r.loadRFC4523MatchingRules()
}

func (r Schema) loadRFC4523MatchingRules() (err error) {

	var i int
	for i = 0; i < len(rfc4523MatchingRules) && err == nil; i++ {
		mr := rfc4523MatchingRules[i]
		err = r.ParseMatchingRule(string(mr))
	}

	if want := rfc4523MatchingRules.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4523 MatchingRules parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4530MatchingRules returns an error following an attempt to load
all RFC 4530 [MatchingRule] slices into the receiver instance.
*/
func (r Schema) LoadRFC4530MatchingRules() error {
	return r.loadRFC4530MatchingRules()
}

func (r Schema) loadRFC4530MatchingRules() (err error) {

	var i int
	for i = 0; i < len(rfc4530MatchingRules) && err == nil; i++ {
		mr := rfc4530MatchingRules[i]
		err = r.ParseMatchingRule(string(mr))
	}

	if want := rfc4530MatchingRules.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4530 MatchingRules parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
prepareString returns a string an an error indicative of an attempt
to represent the receiver instance as a string using [text/template].
*/
func (r *matchingRule) prepareString() (str string, err error) {
	buf := newBuf()
	t := newTemplate(r.Type()).
		Funcs(funcMap(map[string]any{
			`Syntax`:       func() string { return r.Syntax.NumericOID() },
			`ExtensionSet`: r.Extensions.tmplFunc,
			`Obsolete`:     func() bool { return r.Obsolete },
		}))
	if t, err = t.Parse(matchingRuleTmpl); err == nil {
		if err = t.Execute(buf, struct {
			Definition *matchingRule
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

// stackage closure func - do not exec directly.
func (r MatchingRules) canPush(x ...any) (err error) {
	if len(x) == 0 {
		return
	}

	for i := 0; i < len(x) && err == nil; i++ {
		instance := x[i]
		if mr, ok := instance.(MatchingRule); !ok || mr.IsZero() {
			err = ErrTypeAssert
		} else if tst := r.get(mr.NumericOID()); !tst.IsZero() {
			err = mkerr(ErrNotUnique.Error() + ": " + mr.Type() + `, ` + mr.NumericOID())
		}
	}

	return
}

/*
Len returns the current integer length of the receiver instance.
*/
func (r MatchingRules) Len() int {
	return r.len()
}

func (r MatchingRules) len() int {
	return r.cast().Len()
}

/*
String returns the string representation of the receiver instance.
*/
func (r MatchingRules) String() string {
	return r.cast().String()
}

/*
Compliant returns a Boolean value indicative of every [MatchingRule]
returning a compliant response from the [MatchingRule.Compliant] method.
*/
func (r MatchingRules) Compliant() bool {
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
compliant per the required clauses of [§ 4.1.3 of RFC 4512]:

  - Numeric OID must be present and valid

[§ 4.1.3 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.3
*/
func (r MatchingRule) Compliant() bool {
	if r.IsZero() {
		return false
	}

	if !isNumericOID(r.NumericOID()) {
		return false
	}

	return r.Syntax().Compliant()
}

/*
SetStringer allows the assignment of an individual [Stringer] function or
method to all [MatchingRule] slices within the receiver stack instance.

Input of zero (0) variadic values, or an explicit nil, will overwrite all
preexisting stringer functions with the internal closure default, which is
based upon a one-time use of the [text/template] package by all receiver
slice instances.

Input of a non-nil closure function value will overwrite all preexisting
stringers.

This is a fluent method and may be used multiple times.
*/
func (r MatchingRules) SetStringer(function ...Stringer) MatchingRules {
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
func (r MatchingRule) SetStringer(function ...Stringer) MatchingRule {
	if r.Compliant() {
		r.matchingRule.setStringer(function...)
	}

	return r
}

func (r *matchingRule) setStringer(function ...Stringer) {
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
XOrigin returns an instance of [MatchingRules] containing only definitions
which bear the X-ORIGIN value of x. Case is not significant in the matching
process, nor is whitespace (e.g.: RFC 4517 vs. RFC4517).
*/
func (r MatchingRules) XOrigin(x string) (defs MatchingRules) {
	defs = NewMatchingRules()
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
func (r MatchingRule) String() (def string) {
	if !r.IsZero() {
		if r.matchingRule.stringer != nil {
			def = r.matchingRule.stringer()
		}
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r MatchingRules) IsZero() bool {
	return r.cast().IsZero()
}

/*
Index returns the instance of [MatchingRule] found within the
receiver stack instance at index N.  If no instance is found at
the index specified, a zero [MatchingRule] instance is returned.
*/
func (r MatchingRules) Index(idx int) MatchingRule {
	return r.index(idx)
}

func (r MatchingRules) index(idx int) (mr MatchingRule) {
	slice, found := r.cast().Index(idx)
	if found {
		if _mr, ok := slice.(MatchingRule); ok {
			mr = _mr
		}
	}

	return
}

/*
Push returns an error following an attempt to push a MatchingRule
into the receiver stack instance.
*/
func (r MatchingRules) Push(mr any) error {
	return r.push(mr)
}

func (r MatchingRules) push(x any) (err error) {
	switch tv := x.(type) {
	case MatchingRule:
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
Contains calls [MatchingRules.Get] to return a Boolean value indicative of
a successful, non-zero retrieval of a [MatchingRule] instance -- matching
the provided id -- from within the receiver stack instance.
*/
func (r MatchingRules) Contains(id string) bool {
	return r.contains(id)
}

func (r MatchingRules) contains(id string) bool {
	return !r.get(id).IsZero()
}

/*
Type returns the string literal "matchingRule".
*/
func (r MatchingRule) Type() string {
	return r.matchingRule.Type()
}

func (r matchingRule) Type() string {
	return `matchingRule`
}

/*
Type returns the string literal "matchingRules".
*/
func (r MatchingRules) Type() string {
	return `matchingRules`
}

/*
Inventory returns an instance of [Inventory] which represents the current
inventory of [MatchingRule] instances within the receiver.

The keys are numeric OIDs, while the values are zero (0) or more string
slices, each representing a name by which the definition is known.
*/
func (r MatchingRules) Inventory() (inv Inventory) {
	inv = make(Inventory, 0)
	for i := 0; i < r.len(); i++ {
		def := r.index(i)
		inv[def.NumericOID()] = def.Names().List()

	}

	return
}

/*
Get returns an instance of [MatchingRule] based upon a search for id within
the receiver stack instance.

The return instance, if not nil, was retrieved based upon a textual match of
the principal identifier of a [MatchingRule] and the provided id.

The return instance is nil if no match was made.

Case is not significant in the matching process.
*/
func (r MatchingRules) Get(id string) MatchingRule {
	return r.get(id)
}

func (r MatchingRules) get(id string) (mr MatchingRule) {
	for i := 0; i < r.len() && mr.IsZero(); i++ {
		if _mr := r.index(i); !_mr.IsZero() {
			if _mr.IsIdentifiedAs(id) {
				mr = _mr
			}
		}
	}

	return
}

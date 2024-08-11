package schemax

/*
NewDITStructureRules initializes a new [DITStructureRules] instance.
*/
func NewDITStructureRules() DITStructureRules {
	r := DITStructureRules(newCollection(``))
	r.cast().SetPushPolicy(r.canPush)

	return r
}

/*
NewDITStructureRuleIDList initializes a new [RuleIDList] instance and casts it
as a [DITStructureRules] instance.

This is mainly used to define a series of superior [DITStructureRule] instances
specified by a subordinate instance of [DITStructureRule].
*/
func NewDITStructureRuleIDList() DITStructureRules {
	r := DITStructureRules(newRuleIDList(``))
	r.cast().
		SetPushPolicy(r.canPush).
		SetPresentationPolicy(r.iDsStringer)

	return r
}

/*
DITStructureRules returns the [DITStructureRules] instance from
within the receiver instance.
*/
func (r Schema) DITStructureRules() (dss DITStructureRules) {
	slice, _ := r.cast().Index(dITStructureRulesIndex)
	dss, _ = slice.(DITStructureRules)
	return
}

func (r DITStructureRule) schema() (s Schema) {
	if !r.IsZero() {
		s = r.dITStructureRule.schema
	}

	return
}

/*
Replace overrides the receiver with x. Both must bear an identical
numeric rule ID and x MUST be compliant.

Note that the relevant [Schema] instance must be configured to allow
definition override by way of the [AllowOverride] bit setting.  See
the [Schema.Options] method for a means of accessing the settings
value.

Note that this method does not reallocate a new pointer instance
within the [DITStructureRule] envelope type, thus all references to the
receiver instance within various stacks will be preserved.

This is a fluent method.
*/
func (r DITStructureRule) Replace(x DITStructureRule) DITStructureRule {
	if !r.Schema().Options().Positive(AllowOverride) {
		return r
	}

	if !r.IsZero() {
		r.dITStructureRule.replace(x)
	}

	return r
}

func (r *dITStructureRule) replace(x DITStructureRule) {
	if r == nil {
		r = newDITStructureRule()
	} else if r.ID != x.RuleID() {
		return
	}

	r.ID = x.dITStructureRule.ID
	r.Name = x.dITStructureRule.Name
	r.Desc = x.dITStructureRule.Desc
	r.Form = x.dITStructureRule.Form
	r.Obsolete = x.dITStructureRule.Obsolete
	r.SuperRules = x.dITStructureRule.SuperRules
	r.Extensions = x.dITStructureRule.Extensions
	r.stringer = x.dITStructureRule.stringer
	r.schema = x.dITStructureRule.schema
	r.data = x.dITStructureRule.data
}

/*
SetData assigns x to the receiver instance. This is a general-use method and has no
specific intent beyond convenience. The contents may be subsequently accessed via the
[DITStructureRule.Data] method.

This is a fluent method.
*/
func (r DITStructureRule) SetData(x any) DITStructureRule {
	if !r.IsZero() {
		r.dITStructureRule.setData(x)
	}

	return r
}

func (r *dITStructureRule) setData(x any) {
	r.data = x
}

/*
Data returns the underlying value (x) assigned to the receiver's data storage field. Data
can be set within the receiver instance by way of the [DITStructureRule.SetData] method.
*/
func (r DITStructureRule) Data() (x any) {
	if !r.IsZero() {
		x = r.dITStructureRule.data
	}

	return
}

/*
SetSchema assigns an instance of [Schema] to the receiver instance.  This allows
internal verification of certain actions without the need for user input of
an instance of [Schema] manually at each juncture.

Note that the underlying [Schema] instance is automatically set when creating
instances of this type by way of parsing, as well as if the receiver instance
was initialized using the [Schema.NewDITStructureRule] method.

This is a fluent method.
*/
func (r DITStructureRule) SetSchema(schema Schema) DITStructureRule {
	if !r.IsZero() {
		r.dITStructureRule.setSchema(schema)
	}

	return r
}

func (r *dITStructureRule) setSchema(schema Schema) {
	r.schema = schema
}

/*
Schema returns the [Schema] instance associated with the receiver instance.
*/
func (r DITStructureRule) Schema() (s Schema) {
	if !r.IsZero() {
		s = r.dITStructureRule.getSchema()
	}

	return
}

func (r *dITStructureRule) getSchema() (s Schema) {
	if r != nil {
		s = r.schema
	}

	return
}

/*
NumericOID returns an empty string, as [DITStructureRule] definitions
do not bear numeric OIDs.  This method exists only to satisfy Go
interface requirements.
*/
func (r DITStructureRule) NumericOID() string { return `` }

/*
Obsolete returns a Boolean value indicative of definition obsolescence.
*/
func (r DITStructureRule) Obsolete() (o bool) {
	if !r.IsZero() {
		o = r.dITStructureRule.Obsolete
	}

	return
}

/*
Compliant returns a Boolean value indicative of every [DITStructureRule]
returning a compliant response from the [DITStructureRule.Compliant] method.
*/
func (r DITStructureRules) Compliant() bool {
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
compliant per the required clauses of [§ 4.1.7.1 of RFC 4512]:

  - "rule ID" must be specified in the form of an unsigned integer of any magnitude
  - FORM clause MUST refer to a known [NameForm] instance within the associated [Schema] instance
  - FORM clause MUST refer to a COMPLIANT [NameForm]
  - FORM must not violate, or be violated by, a relevant [DITContentRule] within the associated [Schema] instance

[§ 4.1.7.1 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.7.1
*/
func (r DITStructureRule) Compliant() bool {
	// presence of ruleid is guaranteed via
	// uint default, no need to check.

	if r.IsZero() {
		return false
	}

	// obtain nameForm and verify as compliant.
	form := r.Form()
	if !form.Compliant() {
		return false
	}

	// attempt to call the dITContentRule which
	// shares the same OID as the nameForm's
	// structural class.  If zero, we can bail
	// right now, as the upcoming section does
	// not apply.
	dc := r.schema().DITContentRules().Get(form.OC().OID())
	if dc.IsZero() {
		return true
	}

	// We found a matching dITContentRule. We want to
	// be sure that none of the nameForm's MUST clause
	// members are present in the dITContentRule's NOT
	// clause
	clause := form.Must()
	for i := 0; i < clause.Len(); i++ {
		if dc.Not().Contains(clause.Index(i).OID()) {
			return false
		}
	}

	return true
}

/*
Type returns the string literal "dITStructureRule".
*/
func (r DITStructureRule) Type() string {
	return `dITStructureRule`
}

func (r dITStructureRule) Type() string {
	return `dITStructureRule`
}

/*
Type returns the string literal "dITStructureRules".
*/
func (r DITStructureRules) Type() string {
	return `dITStructureRules`
}

/*
SuperRules returns a [DITStructureRules] containing zero (0) or more
superior [DITStructureRule] instances from which the receiver extends.
*/
func (r DITStructureRule) SuperRules() (sup DITStructureRules) {
	if !r.IsZero() {
		sup = r.dITStructureRule.SuperRules
	}

	return
}

/*
SubRules returns an instance of [DITStructureRules] containing slices of
[DITStructureRule] instances that are direct subordinates to the receiver
instance. As such, this method is essentially the inverse of the
[DITStructureRule.SuperRules] method.

The super chain is NOT traversed beyond immediate subordinate instances.

Note that the relevant [Schema] instance must have been set using the
[DITStructureRule.SetSchema] method prior to invocation of this method.
Should this requirement remain unfulfilled, the return instance will
be a zero instance.
*/
func (r DITStructureRule) SubRules() (subs DITStructureRules) {
	if !r.IsZero() {
		subs = NewDITStructureRuleIDList()
		dsrs := r.schema().DITStructureRules()
		for i := 0; i < dsrs.Len(); i++ {
			typ := dsrs.Index(i)
			supers := typ.SuperRules()
			if got := supers.Get(r.RuleID()); !got.IsZero() {
				subs.Push(typ)
			}
		}
	}

	return
}

/*
ID returns the string representation of the principal name OR rule ID
held by the receiver instance.
*/
func (r DITStructureRule) ID() (id string) {
	if !r.IsZero() {
		_id := uitoa(r.RuleID())
		if id = r.Name(); len(id) == 0 {
			id = _id
		}
	}

	return
}

/*
IsIdentifiedAs returns a Boolean value indicative of whether id matches
either the numericOID or descriptor of the receiver instance.  Case is
not significant in the matching process.
*/
func (r DITStructureRule) IsIdentifiedAs(id string) (ident bool) {
	if !r.IsZero() {
		ident = id == uitoa(r.RuleID()) ||
			r.dITStructureRule.Name.contains(id)
	}

	return
}

/*
Form returns the underlying instance of [NameForm] set within the
receiver. If unset, a zero instance is returned.
*/
func (r DITStructureRule) Form() (nf NameForm) {
	if !r.IsZero() {
		nf = r.dITStructureRule.Form
	}

	return
}

/*
Extensions returns the [Extensions] instance -- if set -- within
the receiver.
*/
func (r DITStructureRule) Extensions() (e Extensions) {
	if !r.IsZero() {
		e = r.dITStructureRule.Extensions
	}

	return
}

/*
RuleID returns the unsigned integer identifier held by the
receiver instance.
*/
func (r DITStructureRule) RuleID() (id uint) {
	if !r.IsZero() {
		id = r.dITStructureRule.ID
	}

	return
}

/*
Name returns the string form of the principal name of the receiver instance, if set.
*/
func (r DITStructureRule) Name() (id string) {
	if !r.IsZero() {
		id = r.dITStructureRule.Name.index(0)
	}

	return
}

/*
Names returns the underlying instance of [QuotedDescriptorList] from
within the receiver.
*/
func (r DITStructureRule) Names() (names QuotedDescriptorList) {
	if !r.IsZero() {
		names = r.dITStructureRule.Name
	}

	return
}

/*
Description returns the underlying (optional) descriptive text
assigned to the receiver instance.
*/
func (r DITStructureRule) Description() (desc string) {
	if !r.IsZero() {
		desc = r.dITStructureRule.Desc
	}

	return
}

/*
SetStringer allows the assignment of an individual [Stringer] function or
method to all [DITStructureRule] slices within the receiver stack instance.

Input of zero (0) variadic values, or an explicit nil, will overwrite all
preexisting stringer functions with the internal closure default, which is
based upon a one-time use of the [text/template] package by all receiver
slice instances.

Input of a non-nil closure function value will overwrite all preexisting
stringers.

This is a fluent method and may be used multiple times.
*/
func (r DITStructureRules) SetStringer(function ...Stringer) DITStructureRules {
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
func (r DITStructureRule) SetStringer(function ...Stringer) DITStructureRule {
	if r.Compliant() {
		r.dITStructureRule.setStringer(function...)
	}

	return r
}

func (r *dITStructureRule) setStringer(function ...Stringer) {
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
		}
	} else {
		r.stringer = stringer
	}
}

/*
SetName assigns the provided names to the receiver instance.

Name instances must conform to RFC 4512 descriptor format but
need not be quoted.

This is a fluent method.
*/
func (r DITStructureRule) SetName(x ...string) DITStructureRule {
	if !r.IsZero() {
		r.dITStructureRule.setName(x...)
	}

	return r
}

func (r *dITStructureRule) setName(x ...string) {
	for i := 0; i < len(x); i++ {
		r.Name.Push(x[i])
	}
}

/*
SetObsolete sets the receiver instance to OBSOLETE if not already set. Note that obsolescence cannot be unset.

This is a fluent method.
*/
func (r DITStructureRule) SetObsolete() DITStructureRule {
	if !r.IsZero() {
		r.dITStructureRule.setObsolete()
	}

	return r
}

func (r *dITStructureRule) setObsolete() {
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
func (r DITStructureRule) SetDescription(desc string) DITStructureRule {
	if !r.IsZero() {
		r.dITStructureRule.setDescription(desc)
	}

	return r
}

func (r *dITStructureRule) setDescription(desc string) {
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
func (r DITStructureRule) SetExtension(x string, xstrs ...string) DITStructureRule {
	if !r.IsZero() {
		r.dITStructureRule.setExtension(x, xstrs...)
	}

	return r
}

func (r *dITStructureRule) setExtension(x string, xstrs ...string) {
	r.Extensions.Set(x, xstrs...)
}

/*
SetNumericOID allows the manual assignment of a numeric OID to the
receiver instance if the following are all true:

  - The input id value is a syntactically valid numeric OID
  - The receiver does not already possess a numeric OID

This is a fluent method.
*/
func (r DITStructureRule) SetRuleID(id any) DITStructureRule {
	if !r.IsZero() {
		r.dITStructureRule.setRuleID(id)
	}

	return r
}

func (r *dITStructureRule) setRuleID(x any) {
	switch tv := x.(type) {
	case uint64:
		r.ID = uint(tv)
	case uint:
		r.ID = tv
	case int:
		if tv >= 0 {
			r.ID = uint(tv)
		}
	case string:
		if z, ok := atoui(tv); ok {
			r.ID = z
		}
	}

	return
}

/*
SetSuperRule assigns the provided input [DITStructureRule] instance(s) to the
receiver's SUP clause.

This is a fluent method.
*/
func (r DITStructureRule) SetSuperRule(m ...any) DITStructureRule {
	if !r.IsZero() {
		r.dITStructureRule.setSuperRule(m...)
	}

	return r
}

func (r *dITStructureRule) setSuperRule(m ...any) {
	var err error
	for i := 0; i < len(m) && err == nil; i++ {
		var def DITStructureRule
		switch tv := m[i].(type) {
		case string:
			def = r.schema.DITStructureRules().get(tv)
		case DITStructureRule:
			def = tv
		default:
			err = ErrInvalidType
		}

		r.SuperRules.Push(def)
	}
}

/*
Parse returns an error following an attempt to parse raw into the receiver
instance.

Note that the receiver MUST possess a [Schema] reference prior to the execution
of this method.

Also note that successful execution of this method does NOT automatically push
the receiver into any [DITStructureRules] stack, nor does it automatically execute
the [DITStructureRule.SetStringer] method, leaving these tasks to the user.  If the
automatic handling of these tasks is desired, see the [Schema.ParseDITStructureRule]
method as an alternative.
*/
func (r DITStructureRule) Parse(raw string) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
		return
	}

	if r.getSchema().IsZero() {
		err = ErrNilSchemaRef
		return
	}

	err = r.dITStructureRule.parse(raw)

	return
}

func (r *dITStructureRule) parse(raw string) error {
	// parseLS wraps the antlr4512 DITStructureRule parser/lexer
	mp, err := parseDS(raw)
	if err == nil {
		// We received the parsed data from ANTLR (mp).
		// Now we need to marshal it into the receiver.
		var def DITStructureRule
		if def, err = r.schema.marshalDS(mp); err == nil {
			r.ID = def.RuleID()
			_r := DITStructureRule{r}
			_r.replace(def)
		}
	}

	return err
}

/*
SetForm assigns x to the receiver instance as an instance of [NameForm].

This is a fluent method.
*/
func (r DITStructureRule) SetForm(x any) DITStructureRule {
	if !r.IsZero() {
		r.dITStructureRule.setForm(x)
	}

	return r
}

func (r *dITStructureRule) setForm(x any) {
	var def NameForm
	switch tv := x.(type) {
	case string:
		if !r.schema.IsZero() {
			def = r.schema.NameForms().get(tv)
		}
	case NameForm:
		def = tv
	}

	if !def.IsZero() {
		r.Form = def
	}
}

/*
NewDITStructureRule initializes and returns a new instance of [DITStructureRule],
ready for manual assembly.  This method need not be used when creating
new [DITStructureRule] instances by way of parsing, as that is handled on an
internal basis.

Use of this method does NOT automatically push the return instance into
the [Schema.DITStructureRules] stack; this is left to the user.

Unlike the package-level [NewDITStructureRule] function, this method will
automatically reference its originating [Schema] instance (the receiver).
This negates the need for manual use of the [DITStructureRule.SetSchema]
method.

This is the recommended means of creating a new [DITStructureRule] instance
wherever a single [Schema] is being used, which represents most use cases.
*/
func (r Schema) NewDITStructureRule() DITStructureRule {
	return NewDITStructureRule().SetSchema(r)
}

/*
NewDITStructureRule initializes and returns a new instance of [DITStructureRule],
ready for manual assembly.  This method need not be used when creating
new [DITStructureRule] instances by way of parsing, as that is handled on an
internal basis.

Use of this function does not automatically reference the "parent" [Schema]
instance, leaving it up to the user to invoke the [DITStructureRule.SetSchema]
method manually.

When interacting with a single [Schema] instance, which represents most use
cases, use of the [Schema.NewDITStructureRule] method is PREFERRED over use of
this package-level function.

However certain migration efforts, schema audits and other such activities
may require distinct associations of [DITStructureRule] instances with specific
[Schema] instances. Use of this function allows the user to specify the
appropriate [Schema] instance at a later point for a specific instance of
a [DITStructureRule] instance.
*/
func NewDITStructureRule() DITStructureRule {
	ds := DITStructureRule{newDITStructureRule()}
	ds.dITStructureRule.Extensions.setDefinition(ds)
	return ds

}

func newDITStructureRule() *dITStructureRule {
	return &dITStructureRule{
		Name:       NewName(),
		SuperRules: NewDITStructureRuleIDList(),
		Extensions: NewExtensions(),
	}
}

/*
String is a stringer method that returns the string representation
of the receiver instance.
*/
func (r DITStructureRule) String() (dsr string) {
	if !r.IsZero() {
		if r.dITStructureRule.stringer != nil {
			dsr = r.dITStructureRule.stringer()
		}
	}

	return
}

/*
Inventory returns an instance of [Inventory] which represents the current
inventory of [DITStructureRule] instances within the receiver.

The keys are numeric OIDs, while the values are zero (0) or more string
slices, each representing a name by which the definition is known.
*/
func (r DITStructureRules) Inventory() (inv Inventory) {
	inv = make(Inventory, 0)
	for i := 0; i < r.len(); i++ {
		def := r.index(i)
		inv[itoa(int(def.RuleID()))] = def.Names().List()
	}

	return
}

func (r DITStructureRule) setOID(_ string) {}
func (r DITStructureRule) macro() []string { return []string{} }

// stackage closure func - do not exec directly (use String method)
func (r DITStructureRules) iDsStringer(_ ...any) (present string) {
	var _present []string
	for i := 0; i < r.len(); i++ {
		_present = append(_present, itoa(int(r.index(i).RuleID())))
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

		joined := join(_present, padchar+` `+padchar)
		present = `(` + padchar + joined + padchar + `)`
	}

	return
}

// stackage closure func - do not exec directly.
func (r DITStructureRules) canPush(x ...any) (err error) {
	if len(x) == 0 {
		return
	}

	for i := 0; i < len(x); i++ {
		instance := x[i]
		ds, ok := instance.(DITStructureRule)
		if !ok || ds.IsZero() {
			err = ErrTypeAssert
			break
		}

		// Check whether a dSR exists bearing the same
		// ruleid as the newly pushed candidate.
		if tst := r.get(ds.RuleID()); !tst.IsZero() {
			// If explicitly permitted by the schema config,
			// we'll re-index any dSRs bearing a numerical
			// ID that conflicts with a preexisting rule
			// within the stack.
			opts := ds.Schema().Options()
			if !opts.Positive(AllowReindexedStructureRules) {
				// Not allowed!
				err = mkerr(ErrNotUnique.Error() + ": " + ds.Type() +
					`, ` + uitoa(ds.RuleID()))
				break
			}

			// Don't reindex if there is also a naming conflict.
			if len(ds.Name()) > 0 && eq(tst.Name(), ds.Name()) {
				err = mkerr("dITStructureRule name/id conflict; cannot reindex (" +
					uitoa(ds.RuleID()) + ")")
				break
			}

			var next uint
			// Determine the next available ruleid
			if next, ok = r.nextIndex(); !ok {
				err = mkerr("reindex id overflow or uninitialized dITStructureRules")
				break
			}

			// Overwrite the ruleid previously in conflict
			ds.SetRuleID(next)

			// Let the administrator know that a reindex has
			// occurred.
			ds.SetExtension(`X-WARNING`, `REINDEXED`)

			// Update stringer (will clobber any CUSTOM)
			// stringer
			ds.setStringer()
		}
	}

	return
}

/*
nextIndex returns a uint alongside a Boolean value.

If ok evaluates as true, the uint instance represents the next-available
[DITStructureRule] ID that may be used in the event that re-indexing is
to be conducted.

If ok does not evaluate as true, this indicates an overflow would occur
and 0 is returned as a result. If the receiver is uninitialized, a similar
return ensues.

The analysis does not look for breaks in otherwise contiguous sequences
of numbers, rather it finds the highest number and checks whether an
overflow would occur if 1 were added to that number. If no overflow is
perceived, the new sum of the number is returned.
*/
func (r DITStructureRules) nextIndex() (idx uint, ok bool) {
	if !r.IsZero() {
		L := r.Len()
		if L == 0 {
			ok = true
			return
		}

		var indices []uint
		for i := 0; i < L; i++ {
			indices = append(indices, r.Index(i).RuleID())
		}

		var slice uint
		for i := 0; i < len(indices); i++ {
			slice = indices[i]
			for j := 0; j < len(indices); j++ {
				if slice < indices[j] {
					slice = indices[j]
					continue
				}
			}
		}

		if !(slice+1 <= ^uint(0)) {
			return
		}

		idx = slice + 1
		ok = true
	}

	return
}

/*
Len returns the current integer length of the receiver instance.
*/
func (r DITStructureRules) Len() int {
	return r.len()
}

func (r DITStructureRules) len() int {
	return r.cast().Len()
}

/*
String is a stringer method that returns the string representation
of the receiver instance.
*/
func (r DITStructureRules) String() string {
	return r.cast().String()
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r DITStructureRules) IsZero() bool {
	return r.cast().IsZero()
}

/*
Index returns the instance of [DITStructureRule] found within the
receiver stack instance at index N. If no instance is found at the
index specified, a zero [DITStructureRule] instance is returned.
*/
func (r DITStructureRules) Index(idx int) DITStructureRule {
	return r.index(idx)
}

func (r DITStructureRules) index(idx int) (ds DITStructureRule) {
	if slice, found := r.cast().Index(idx); found {
		if _ds, ok := slice.(DITStructureRule); ok {
			ds = _ds
		}
	}

	return
}

/*
Push returns an error following an attempt to push a [DITStructureRule]
into the receiver stack instance.
*/
func (r DITStructureRules) Push(ds any) error {
	return r.push(ds)
}

func (r DITStructureRules) push(x any) (err error) {
	switch tv := x.(type) {
	case DITStructureRule:
		if !tv.Compliant() {
			err = ErrDefNonCompliant
			break
		}
		r.cast().Push(tv)
		err = r.cast().Err()
	default:
		err = ErrInvalidType
	}

	return
}

/*
Contains calls [DITStructureRules.Get] to return a Boolean value indicative
of a successful, non-zero retrieval of a [DITStructureRules] instance --
matching the provided id -- from within the receiver stack instance.
*/
func (r DITStructureRules) Contains(id string) bool {
	return r.contains(id)
}

func (r DITStructureRules) contains(id string) bool {
	return !r.get(id).IsZero()
}

/*
Get returns an instance of [DITStructureRule] based upon a search for id
within the receiver stack instance.

The return instance, if not nil, was retrieved based upon a textual match
of the principal identifier of a [DITStructureRule] and the provided id.

The return instance is nil if no match was made.

Case is not significant in the matching process.
*/
func (r DITStructureRules) Get(id any) DITStructureRule {
	return r.get(id)
}

func (r DITStructureRules) get(id any) (ds DITStructureRule) {
	L := r.len()
	if L == 0 {
		return
	}

	var n uint
	var name string
	var named bool

	switch tv := id.(type) {
	case int:
		if tv < 0 {
			return
		}
		n = uint(tv)
	case uint:
		n = tv
	case string:
		// string may be a string
		// uint, or a name.
		if _n, err := atoi(tv); err != nil {
			named = true
			name = tv
		} else {
			return r.get(_n)
		}
	}

	for i := 0; i < L && ds.IsZero(); i++ {
		_ds := r.index(i)
		if named {
			if _ds.Names().Contains(name) && len(name) > 0 {
				ds = _ds
			}
		} else if _ds.RuleID() == n {
			ds = _ds
		}
	}

	return
}

/*
Maps returns slices of [DefinitionMap] instances.
*/
func (r DITStructureRules) Maps() (defs DefinitionMaps) {
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
func (r DITStructureRule) Map() (def DefinitionMap) {
	if r.IsZero() {
		return
	}

	var sups []string
	for i := 0; i < r.SuperRules().Len(); i++ {
		m := r.SuperRules().Index(i)
		sups = append(sups, itoa(int(m.RuleID())))
	}

	if r.Form().IsZero() {
		return
	}

	def = make(DefinitionMap, 0)
	def[`RULEID`] = []string{itoa(int(r.RuleID()))}
	def[`NAME`] = r.Names().List()
	def[`DESC`] = []string{r.Description()}
	def[`OBSOLETE`] = []string{bool2str(r.Obsolete())}
	def[`FORM`] = []string{r.Form().OID()}
	def[`SUP`] = sups
	def[`TYPE`] = []string{r.Type()}
	def[`RAW`] = []string{r.String()}

	// copy our extensions from receiver r
	// into destination def.
	def = mapTransferExtensions(r, def)

	// Clean up any empty fields
	def.clean()

	return
}

func (r *dITStructureRule) prepareString() (str string, err error) {
	buf := newBuf()
	t := newTemplate(r.Type()).
		Funcs(funcMap(map[string]any{
			`ExtensionSet`: r.Extensions.tmplFunc,
			// SuperLen refers to the integer number of
			// superior dITStructureRule instances held
			// by a dITStructureRule.
			`SuperLen`: r.SuperRules.len,
			`Obsolete`: func() bool { return r.Obsolete },
		}))

	if t, err = t.Parse(dITStructureRuleTmpl); err == nil {
		if err = t.Execute(buf, struct {
			Definition *dITStructureRule
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
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r DITStructureRule) IsZero() bool {
	return r.dITStructureRule == nil
}

/*
LoadDITStructureRules returns an error following an attempt to load all
built-in [DITStructureRule] slices into the receiver instance.
*/
func (r Schema) LoadDITStructureRules() error {
	return r.loadDITStructureRules()
}

func (r Schema) loadDITStructureRules() (err error) {
	if !r.IsZero() {
		funks := []func() error{
			r.loadRFC4403DITStructureRules,
		}

		for i := 0; i < len(funks) && err == nil; i++ {
			err = funks[i]()
		}
	}

	return
}

/*
LoadRFC4403DITStructureRules returns an error following an attempt to
load all RFC 4403 [DITStructureRule] slices into the receiver instance.
*/
func (r Schema) LoadRFC4403DITStructureRules() error {
	return r.loadRFC4403DITStructureRules()
}

func (r Schema) loadRFC4403DITStructureRules() (err error) {

	var i int
	for i = 0; i < len(rfc4403DITStructureRules) && err == nil; i++ {
		at := rfc4403DITStructureRules[i]
		err = r.ParseDITStructureRule(string(at))
	}

	if want := rfc4403DITStructureRules.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4403 DITStructureRules parsed: want " +
				itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

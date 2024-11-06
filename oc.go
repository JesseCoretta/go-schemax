package schemax

/*
NewObjectClass initializes and returns a new instance of [ObjectClass],
ready for manual assembly.  This method need not be used when creating
new [ObjectClass] instances by way of parsing, as that is handled on an
internal basis.

Use of this method does NOT automatically push the return instance into
the [Schema.ObjectClasses] stack; this is left to the user.

Unlike the package-level [NewObjectClass] function, this method will
automatically reference its originating [Schema] instance (the receiver).
This negates the need for manual use of the [ObjectClass.SetSchema]
method.

This is the recommended means of creating a new [ObjectClass] instance
wherever a single [Schema] is being used, which represents most use cases.
*/
func (r Schema) NewObjectClass() ObjectClass {
	return NewObjectClass().SetSchema(r)
}

/*
NewObjectClasses initializes and returns a new [ObjectClasses] instance,
configured to allow the storage of all [ObjectClass] instances.
*/
func NewObjectClasses() ObjectClasses {
	r := ObjectClasses(newCollection(``))
	r.cast().SetPushPolicy(r.canPush)

	return r
}

/*
NewObjectClassOIDList initializes and returns a new [ObjectClasses] that
has been configured to allow the storage of arbitrary [ObjectClass] instances.
*/
func NewObjectClassOIDList(label ...string) ObjectClasses {
	name := `oc_oidlist`
	if len(label) > 0 {
		name = label[0]
	}

	r := ObjectClasses(newOIDList(name))
	r.cast().
		SetPushPolicy(r.canPush).
		SetPresentationPolicy(r.oIDsStringer)

	return r
}

/*
NewObjectClass initializes and returns a new instance of [ObjectClass],
ready for manual assembly.  This method need not be used when creating
new [ObjectClass] instances by way of parsing, as that is handled on an
internal basis.

Use of this function does not automatically reference the "parent" [Schema]
instance, leaving it up to the user to invoke the [ObjectClass.SetSchema]
method manually.

When interacting with a single [Schema] instance, which represents most use
cases, use of the [Schema.NewObjectClass] method is PREFERRED over use of
this package-level function.

However certain migration efforts, schema audits and other such activities
may require distinct associations of [ObjectClass] instances with specific
[Schema] instances. Use of this function allows the user to specify the
appropriate [Schema] instance at a later point for a specific instance of
an [ObjectClass] instance.
*/
func NewObjectClass() ObjectClass {
	oc := ObjectClass{newObjectClass()}
	oc.objectClass.Extensions.setDefinition(oc)
	return oc
}

func newObjectClass() *objectClass {
	return &objectClass{
		Name:         NewName(),
		Must:         NewAttributeTypeOIDList(`MUST`),
		May:          NewAttributeTypeOIDList(`MAY`),
		SuperClasses: NewObjectClassOIDList(`SUP`),
		Extensions:   NewExtensions(),
	}
}

/*
Parse returns an error following an attempt to parse raw into the receiver
instance.

Note that the receiver MUST possess a [Schema] reference prior to the execution
of this method.

Also note that successful execution of this method does NOT automatically push
the receiver into any [ObjectClass] stack, nor does it automatically execute
the [ObjectClass.SetStringer] method, leaving these tasks to the user.  If the
automatic handling of these tasks is desired, see the [Schema.ParseObjectClass]
method as an alternative.
*/
func (r ObjectClass) Parse(raw string) (err error) {
	if r.IsZero() {
		err = ErrNilReceiver
		return
	}

	if r.getSchema().IsZero() {
		err = ErrNilSchemaRef
		return
	}

	err = r.objectClass.parse(raw)

	return
}

func (r *objectClass) parse(raw string) error {
	// parseMR wraps the antlr4512 ObjectClass parser/lexer
	mp, err := parseOC(raw)
	if err == nil {
		// We received the parsed data from ANTLR (mp).
		// Now we need to marshal it into the receiver.
		var def ObjectClass
		if def, err = r.schema.marshalOC(mp); err == nil {
			r.OID = def.NumericOID()
			_r := ObjectClass{r}
			_r.replace(def)
		}
	}

	return err
}

func (r ObjectClass) setOID(x string) {
	if !r.IsZero() {
		r.objectClass.OID = x
	}
}

func (r ObjectClass) macro() (m []string) {
	if !r.IsZero() {
		m = r.objectClass.Macro
	}

	return
}

/*
EnforcedBy returns the [DITContentRule] instance which bears the same
numeric OID held by the receiver instance, which must be a STRUCTURAL
[ObjectClass].

This method is essentially the inverse of [DITContentRule.StructuralClass].

A schema can only contain one (1) such [DITContentRule] per STRUCTURAL
[ObjectClass].

If the return instance is zero, this means that either the receiver is
not a STRUCTURAL [ObjectClass], or that no [DITContentRule] bearing the
receiver's numeric OID is currently in force.
*/
func (r ObjectClass) EnforcedBy() (dcr DITContentRule) {
	if !r.Schema().IsZero() && r.Kind() == StructuralKind {
		dcr = r.Schema().DITContentRules().Get(r.NumericOID())
	}

	return
}

/*
SetName assigns the provided names to the receiver instance.

Name instances must conform to RFC 4512 descriptor format but
need not be quoted.

This is a fluent method.
*/
func (r ObjectClass) SetName(x ...string) ObjectClass {
	if !r.IsZero() {
		r.objectClass.setName(x...)
	}

	return r
}

func (r *objectClass) setName(x ...string) {
	for i := 0; i < len(x); i++ {
		r.Name.Push(x[i])
	}
}

/*
SetData assigns x to the receiver instance. This is a general-use method and has no
specific intent beyond convenience. The contents may be subsequently accessed via the
[ObjectClass.Data] method.

This is a fluent method.
*/
func (r ObjectClass) SetData(x any) ObjectClass {
	if !r.IsZero() {
		r.objectClass.setData(x)
	}

	return r
}

func (r *objectClass) setData(x any) {
	r.data = x
}

/*
Data returns the underlying value (x) assigned to the receiver's data storage field. Data
can be set within the receiver instance by way of the [ObjectClass.SetData] method.
*/
func (r ObjectClass) Data() (x any) {
	if !r.IsZero() {
		x = r.objectClass.data
	}

	return
}

/*
SetExtension assigns key x to value xstrs within the receiver's underlying
[Extensions] instance.

This is a fluent method.
*/
func (r ObjectClass) SetExtension(x string, xstrs ...string) ObjectClass {
	if !r.IsZero() {
		r.objectClass.setExtension(x, xstrs...)
	}

	return r
}

func (r *objectClass) setExtension(x string, xstrs ...string) {
	r.Extensions.Set(x, xstrs...)
}

/*
SetKind assigns the abstraction of an objectClass "kind" to the receiver
instance. Valid input types are as follows:

  - Literal uint consts ([AbstractKind], [AuxiliaryKind], [StructuralKind])
  - int equivalents of uint consts
  - string names of known objectClass kinds

In the case of string names, case is not significant.

This is a fluent method.
*/
func (r ObjectClass) SetKind(k any) ObjectClass {
	if !r.IsZero() {
		r.objectClass.setKind(k)
	}

	return r
}

func (r *objectClass) setKind(k any) {
	switch tv := k.(type) {
	case string:
		switch lc(tv) {
		case `structural`, ``:
			r.Kind = StructuralKind
		case `abstract`:
			r.Kind = AbstractKind
		case `auxiliary`:
			r.Kind = AuxiliaryKind
		}
	case int:
		if 0 <= tv && tv <= 2 {
			r.Kind = uint(tv)
		}
	case uint:
		switch tv {
		case AbstractKind,
			AuxiliaryKind,
			StructuralKind:
			r.Kind = tv
		}
	}
}

/*
SetMust assigns the provided input [AttributeType] instance(s) to the
receiver's MUST clause.

This is a fluent method.
*/
func (r ObjectClass) SetMust(m ...any) ObjectClass {
	if !r.IsZero() {
		r.objectClass.setMust(m...)
	}

	return r
}

func (r *objectClass) setMust(m ...any) {
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

}

/*
SetMay assigns the provided input [AttributeType] instance(s) to the
receiver's MAY clause.

This is a fluent method.
*/
func (r ObjectClass) SetMay(m ...any) ObjectClass {
	if !r.IsZero() {
		r.objectClass.setMay(m...)
	}

	return r
}

func (r *objectClass) setMay(m ...any) {
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
}

/*
SetSuperClass appends the input value(s) to the the super classes stack within the
the receiver.  Valid input types are string, to represent an RFC 4512 OID residing
in the underlying [Schema] instance, or an bonafide [ObjectClass] instance already
obtained or crafted.

This is a fluent method.
*/
func (r ObjectClass) SetSuperClass(x ...any) ObjectClass {
	if !r.IsZero() {
		r.objectClass.setSuperClass(x...)
	}

	return r
}

func (r *objectClass) setSuperClass(x ...any) {
	var err error
	for i := 0; i < len(x) && err == nil; i++ {
		var sup ObjectClass
		switch tv := x[i].(type) {
		case string:
			sup = r.schema.ObjectClasses().get(tv)
		case ObjectClass:
			sup = tv
		default:
			err = ErrInvalidType
			continue
		}

		err = r.verifySuperClass(sup.objectClass)
		if err == nil && !sup.IsZero() {
			r.SuperClasses.push(sup)
		}
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
within the [ObjectClass] envelope type, thus all references to the
receiver instance within various stacks will be preserved.

This is a fluent method.
*/
func (r ObjectClass) Replace(x ObjectClass) ObjectClass {
	if !r.Schema().Options().Positive(AllowOverride) {
		return r
	}

	if !r.IsZero() && x.Compliant() {
		r.objectClass.replace(x)
	}

	return r
}

func (r *objectClass) replace(x ObjectClass) {
	if r.OID != x.NumericOID() {
		return
	}

	r.OID = x.objectClass.OID
	r.Macro = x.objectClass.Macro
	r.Name = x.objectClass.Name
	r.Desc = x.objectClass.Desc
	r.Obsolete = x.objectClass.Obsolete
	r.Kind = x.objectClass.Kind
	r.SuperClasses = x.objectClass.SuperClasses
	r.Must = x.objectClass.Must
	r.May = x.objectClass.May
	r.Extensions = x.objectClass.Extensions
	r.data = x.objectClass.data
	r.schema = x.objectClass.schema
	r.stringer = x.objectClass.stringer
	r.data = x.objectClass.data
}

/*
Maps returns slices of [DefinitionMap] instances.
*/
func (r ObjectClasses) Maps() (defs DefinitionMaps) {
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
func (r ObjectClass) Map() (def DefinitionMap) {
	if !r.Compliant() {
		return
	}

	var sups []string
	for i := 0; i < r.SuperClasses().Len(); i++ {
		m := r.SuperClasses().Index(i)
		sups = append(sups, m.OID())
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

	def = make(DefinitionMap, 0)
	def[`NUMERICOID`] = []string{r.NumericOID()}
	def[`NAME`] = r.Names().List()
	def[`DESC`] = []string{r.Description()}
	def[`OBSOLETE`] = []string{bool2str(r.Obsolete())}
	def[`SUP`] = sups
	def[`MUST`] = musts
	def[`MAY`] = mays
	def[`TYPE`] = []string{r.Type()}
	def[`RAW`] = []string{r.String()}

	switch r.Kind() {
	case AbstractKind:
		def[`KIND`] = []string{`ABSTRACT`}
	case AuxiliaryKind:
		def[`KIND`] = []string{`AUXILIARY`}
	default:
		def[`KIND`] = []string{`STRUCTURAL`}
	}

	// copy our extensions from receiver r
	// into destination def.
	def = mapTransferExtensions(r, def)

	// Clean up any empty fields
	def.clean()

	return
}

/*
Compliant returns a Boolean value indicative of every [ObjectClass]
returning a compliant response from the [ObjectClass.Compliant] method.
*/
func (r ObjectClasses) Compliant() bool {
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
compliant per the required clauses of [§ 4.1.1 of RFC 4512]:

  - Numeric OID must be present and valid
  - MUST and MAY clause [AttributeTypes] are compliant
  - MUST and MAY clause contains no Collective [AttributeType] instances

[§ 4.1.1 of RFC 4512]: https://rfc-editor.org/rfc/rfc4512.html#section-4.1.1
*/
func (r ObjectClass) Compliant() bool {
	if r.IsZero() {
		return false
	}

	var (
		may  AttributeTypes = r.May()
		must AttributeTypes = r.Must()
	)

	for i := 0; i < must.Len(); i++ {
		if !must.Index(i).Compliant() || must.Index(i).Collective() {
			return false
		}
	}

	for i := 0; i < may.Len(); i++ {
		if !may.Index(i).Compliant() || may.Index(i).Collective() {
			return false
		}
	}

	return isNumericOID(r.NumericOID())
}

/*
verifySuperClass returns an error following the execution of basic sanity
checks meant to assess the intended super type chain.
*/
func (r *objectClass) verifySuperClass(sup *objectClass) (err error) {
	renv := ObjectClass{r}
	senv := ObjectClass{sup}
	if r == nil || sup == nil {
		err = ErrNilInput
	}

	if renv.NumericOID() == senv.NumericOID() {
		// r and sup are the same class
		err = mkerr("objectClass is subordinate to itself: " +
			renv.NumericOID() + `==` + senv.NumericOID())
	} else if renv.SuperClassOf(senv) {
		// r is a super class of sup; can't have both!
		err = mkerr("cyclical superiority loop: " +
			renv.NumericOID() + ` <--> ` + senv.NumericOID())
	}

	return
}

/*
SuperClassOf returns a Boolean value indicative of r being a superior ("SUP")
[ObjectClass] of sub.

Note: this will trace all super class chains indefinitely and, thus, will
recognize any superior association without regard for "depth".
*/
func (r ObjectClass) SuperClassOf(sub ObjectClass) (sup bool) {
	dsups := sub.SuperClasses() // iterated and possibly traversed
	for i := 0; i < dsups.Len(); i++ {
		dsup := dsups.Index(i)
		if sup = dsup.NumericOID() == r.NumericOID(); sup {
			// direct (immediate) match
			break
		} else if sup = r.SuperClassOf(dsup); sup {
			// match by traversal
			break
		}
	}

	return
}

/*
Obsolete returns a Boolean value indicative of definition
obsolescence.
*/
func (r ObjectClass) Obsolete() (o bool) {
	if !r.IsZero() {
		o = r.objectClass.Obsolete
	}

	return
}

/*
SetObsolete sets the receiver instance to OBSOLETE if not already set.

Obsolescence cannot be unset.

This is a fluent method.
*/
func (r ObjectClass) SetObsolete() ObjectClass {
	if !r.IsZero() {
		r.objectClass.setObsolete()
	}

	return r
}

func (r *objectClass) setObsolete() {
	if !r.Obsolete {
		r.Obsolete = true
	}
}

/*
Kind returns the uint-based kind value associated with the
receiver instance.
*/
func (r ObjectClass) Kind() (kind uint) {
	kind = StructuralKind
	if !r.IsZero() {
		switch k := r.objectClass.Kind; k {
		case AbstractKind:
			kind = k
		case AuxiliaryKind:
			kind = k
		}
	}

	return
}

/*
SuperClasses returns an [ObjectClasses] stack instance containing zero
(0) or more superior [ObjectClass] instances from which the receiver
directly extends. This method does not walk any super class chains.
*/
func (r ObjectClass) SuperClasses() (sups ObjectClasses) {
	if !r.IsZero() {
		sups = r.objectClass.SuperClasses
	}

	return
}

/*
SubClasses returns an instance of [ObjectClasses] containing slices of
[ObjectClass] instances that are direct subordinates to the receiver
instance. As such, this method is essentially the inverse of the
[ObjectClass.SuperClasses] method.

The super chain is NOT traversed beyond immediate subordinate instances.

Note that the relevant [Schema] instance must have been set using the
[ObjectClass.SetSchema] method prior to invocation of this method.
Should this requirement remain unfulfilled, the return instance will
be a zero instance.
*/
func (r ObjectClass) SubClasses() (subs ObjectClasses) {
	if !r.IsZero() {
		subs = NewObjectClassOIDList()
		ocs := r.schema().ObjectClasses()
		for i := 0; i < ocs.Len(); i++ {
			typ := ocs.Index(i)
			supers := typ.SuperClasses()
			if got := supers.Get(r.NumericOID()); !got.IsZero() {
				subs.Push(typ)
			}
		}
	}

	return
}

func (r ObjectClass) schema() (s Schema) {
	if !r.IsZero() {
		s = r.objectClass.schema
	}

	return
}

/*
SuperChain returns an [ObjectClasses] stack of [ObjectClass] instances
which make up the super type chain of the receiver instance.
*/
func (r ObjectClass) SuperChain() (sups ObjectClasses) {
	if !r.IsZero() {
		sups = NewObjectClasses()
		_sups := r.objectClass.SuperClasses.superChain()
		for i := 0; i < _sups.Len(); i++ {
			sups.Push(_sups.Index(i))
		}
	}

	return
}

func (r ObjectClasses) superChain() (sups ObjectClasses) {
	sups = NewObjectClasses()
	for i := 0; i < r.Len(); i++ {
		sups.Push(r.Index(i))
		_sups := r.Index(i).SuperChain()
		for j := 0; j < _sups.Len(); j++ {
			sups.Push(_sups.Index(j))
		}
	}

	return
}

/*
Type returns the string literal "objectClass".
*/
func (r ObjectClass) Type() string {
	return `objectClass`
}

/*
Type returns the string literal "objectClasses".
*/
func (r ObjectClasses) Type() string {
	return `objectClasses`
}

/*
Inventory returns an instance of [Inventory] which represents the current
inventory of [ObjectClass] instances within the receiver.

The keys are numeric OIDs, while the values are zero (0) or more string
slices, each representing a name by which the definition is known.
*/
func (r ObjectClasses) Inventory() (inv Inventory) {
	inv = make(Inventory, 0)
	for i := 0; i < r.len(); i++ {
		def := r.index(i)
		inv[def.NumericOID()] = def.Names().List()

	}

	return
}

/*
SetSchema assigns an instance of [Schema] to the receiver instance.  This allows
internal verification of certain actions without the need for user input of
an instance of [Schema] manually at each juncture.

Note that the underlying [Schema] instance is automatically set when creating
instances of this type by way of parsing, as well as if the receiver instance
was initialized using the [Schema.NewObjectClass] method.

This is a fluent method.
*/
func (r ObjectClass) SetSchema(schema Schema) ObjectClass {
	if !r.IsZero() {
		r.objectClass.setSchema(schema)
	}

	return r
}

func (r *objectClass) setSchema(schema Schema) {
	r.schema = schema
}

/*
Schema returns the [Schema] instance associated with the receiver instance.
*/
func (r ObjectClass) Schema() (s Schema) {
	if !r.IsZero() {
		s = r.objectClass.getSchema()
	}

	return
}

func (r *objectClass) getSchema() (s Schema) {
	if r != nil {
		s = r.schema
	}

	return
}

/*
AllMust returns an [AttributeTypes] containing zero (0) or more required
[AttributeType] definitions for use with this class as well as those specified
by any and all applicable super classes. Duplicate references are silently
discarded.
*/
func (r ObjectClass) AllMust() (must AttributeTypes) {
	must = NewAttributeTypeOIDList()

	if sup := r.SuperClasses(); !sup.IsZero() {
		for i := 0; i < sup.len(); i++ {
			sm := sup.index(i)
			if sc := sm.AllMust(); !sc.IsZero() {
				for j := 0; j < sc.len(); j++ {
					must.push(sc.index(i))
				}
			}
		}
	}

	if lm := r.Must(); !lm.IsZero() {
		for i := 0; i < lm.len(); i++ {
			must.push(lm.index(i))
		}
	}

	return
}

/*
Must returns an [AttributeTypes] containing zero (0) or more required
[AttributeType] definitions for use with this class.
*/
func (r ObjectClass) Must() (must AttributeTypes) {
	if !r.IsZero() {
		must = r.objectClass.Must
	}

	return
}

/*
AllMay returns an [AttributeTypes] containing zero (0) or more allowed
[AttributeType] definitions for use with this class as well as those
specified by any and all applicable super classes.  Duplicate references
are silently discarded.
*/
func (r ObjectClass) AllMay() (may AttributeTypes) {
	may = NewAttributeTypeOIDList()

	if sup := r.SuperClasses(); !sup.IsZero() {
		for i := 0; i < sup.len(); i++ {
			sm := sup.index(i)
			if sc := sm.AllMay(); !sc.IsZero() {
				for j := 0; j < sc.len(); j++ {
					may.push(sc.index(i))
				}
			}
		}
	}

	if lm := r.May(); !lm.IsZero() {
		for i := 0; i < lm.len(); i++ {
			may.push(lm.index(i))
		}
	}

	return
}

/*
May returns an [AttributeTypes] containing zero (0) or more allowed
[AttributeType] definitions for use with this class.
*/
func (r ObjectClass) May() (may AttributeTypes) {
	if !r.IsZero() {
		may = r.objectClass.May
	}

	return
}

/*
OID returns the string representation of an OID -- which is either a
numeric OID or descriptor -- that is held by the receiver instance.
*/
func (r ObjectClass) OID() (oid string) {
	if !r.IsZero() {
		oid = r.NumericOID() // default
		if r.objectClass.Name.len() > 0 {
			oid = r.objectClass.Name.index(0)
		}
	}

	return
}

/*
IsIdentifiedAs returns a Boolean value indicative of whether id matches
either the numericOID or descriptor of the receiver instance.  Case is
not significant in the matching process.
*/
func (r ObjectClass) IsIdentifiedAs(id string) (ident bool) {
	if !r.IsZero() {
		ident = id == r.NumericOID() ||
			r.objectClass.Name.contains(id)
	}

	return
}

/*
Name returns the string form of the principal name of the receiver instance, if set.
*/
func (r ObjectClass) Name() (id string) {
	if !r.IsZero() {
		id = r.objectClass.Name.index(0)
	}

	return
}

/*
Names returns the underlying instance of [QuotedDescriptorList] from within
the receiver.
*/
func (r ObjectClass) Names() (names QuotedDescriptorList) {
	if !r.IsZero() {
		names = r.objectClass.Name
	}

	return
}

/*
Attributes returns an instance of [AttributeTypes] containing references
to [AttributeType] instances which reside within the receiver's MUST and MAY
clauses.

This is useful if a complete list of all immediate [AttributeType] instances
which are available for use and does not include those made available through
super classes.

Duplicate references are silently discarded.
*/
func (r ObjectClass) Attributes() (a AttributeTypes) {
	a = NewAttributeTypeOIDList()

	for _, list := range []AttributeTypes{
		r.objectClass.Must,
		r.objectClass.May,
	} {
		for i := 0; i < list.len(); i++ {
			if at := list.index(i); !at.IsZero() {
				a.push(at)
			}
		}
	}

	return
}

/*
AllAttributes returns an [AttributeTypes] instance containing zero (0) or
more permitted or required [AttributeType] definitions for use with this
class as well as those specified by any and all applicable super classes.

This is useful if a complete list of all [AttributeType] instances which
are available for use is needed.

Duplicate references are silently discarded.
*/
func (r ObjectClass) AllAttributes() (all AttributeTypes) {
	all = NewAttributeTypeOIDList()

	for _, attrs := range []AttributeTypes{
		r.AllMay(),
		r.AllMust(),
	} {
		for i := 0; i < attrs.len(); i++ {
			all.push(attrs.index(i))
		}
	}

	return
}

/*
Extensions returns the [Extensions] instance -- if set -- within
the receiver.
*/
func (r ObjectClass) Extensions() (e Extensions) {
	if !r.IsZero() {
		e = r.objectClass.Extensions
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
func (r ObjectClass) SetNumericOID(id string) ObjectClass {
	if !r.IsZero() {
		r.objectClass.setNumericOID(id)
	}

	return r
}

func (r *objectClass) setNumericOID(id string) {
	if isNumericOID(id) {
		if len(r.OID) == 0 {
			r.OID = id
		}
	}
}

/*
NumericOID returns the string representation of the numeric OID
held by the receiver instance.
*/
func (r ObjectClass) NumericOID() (noid string) {
	if !r.IsZero() {
		noid = r.objectClass.OID
	}

	return
}

/*
SetStringer allows the assignment of an individual [Stringer] function or
method to all [ObjectClass] slices within the receiver stack instance.

Input of zero (0) variadic values, or an explicit nil, will overwrite all
preexisting stringer functions with the internal closure default, which is
based upon a one-time use of the [text/template] package by all receiver
slice instances.

Input of a non-nil closure function value will overwrite all preexisting
stringers.

This is a fluent method and may be used multiple times.
*/
func (r ObjectClasses) SetStringer(function ...Stringer) ObjectClasses {
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
func (r ObjectClass) SetStringer(function ...Stringer) ObjectClass {
	if r.Compliant() {
		r.objectClass.setStringer(function...)
	}

	return r
}

func (r *objectClass) setStringer(function ...Stringer) {
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
XOrigin returns an instance of [ObjectClasses] containing only definitions
which bear the X-ORIGIN value of x. Case is not significant in the matching
process, nor is whitespace (e.g.: RFC 4517 vs. RFC4517).
*/
func (r ObjectClasses) XOrigin(x string) (defs ObjectClasses) {
	defs = NewObjectClasses()
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
func (r ObjectClass) String() (def string) {
	if !r.IsZero() {
		if r.objectClass.stringer != nil {
			def = r.objectClass.stringer()
		}
	}

	return
}

/*
ObjectClasses returns the [ObjectClasses] instance from within
the receiver instance.
*/
func (r Schema) ObjectClasses() (ocs ObjectClasses) {
	slice, _ := r.cast().Index(objectClassesIndex)
	ocs, _ = slice.(ObjectClasses)
	return
}

/*
prepareString returns a string an an error indicative of an attempt
to represent the receiver instance as a string using [text/template].
*/
func (r *objectClass) prepareString() (str string, err error) {
	buf := newBuf()
	var kind string = `STRUCTURAL`
	switch r.Kind {
	case AbstractKind:
		kind = `ABSTRACT`
	case AuxiliaryKind:
		kind = `AUXILIARY`
	}

	t := newTemplate(`objectClass`).
		Funcs(funcMap(map[string]any{
			`Kind`:         func() string { return kind },
			`ExtensionSet`: r.Extensions.tmplFunc,
			`MustLen`:      r.Must.len,
			`MayLen`:       r.May.len,
			`SuperLen`:     r.SuperClasses.len,
			`Obsolete`:     func() bool { return r.Obsolete },
		}))

	if t, err = t.Parse(objectClassTmpl); err == nil {
		if err = t.Execute(buf, struct {
			Definition *objectClass
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
Description returns the underlying (optional) descriptive text
assigned to the receiver instance.
*/
func (r ObjectClass) Description() (desc string) {
	if !r.IsZero() {
		desc = r.objectClass.Desc
	}

	return
}

/*
SetDescription parses desc into the underlying DESC clause within the
receiver instance.  Although a RFC 4512-compliant QuotedString is
required, the outer single-quotes need not be specified literally.
*/
func (r ObjectClass) SetDescription(desc string) ObjectClass {
	if !r.IsZero() {
		r.objectClass.setDescription(desc)
	}

	return r
}

func (r *objectClass) setDescription(desc string) {
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
func (r ObjectClasses) IsZero() bool {
	return r.cast().IsZero()
}

/*
toggle-based oid list stringer handler - do not exec directly.
*/
func (r ObjectClasses) oIDsStringer(_ ...any) (present string) {
	slice := r.index(0)
	hd := slice.Schema().Options().Positive(HangingIndents)
	id := r.cast().ID()
	if hd && id != `oc_oidlist` {
		return r.oIDsStringerPretty(len(id))
	}

	return r.oIDsStringerStd()
}

/*
prettified stackage closure func for oid lists - do not exec directly.

prepare a custom [stackage.PresentationPolicy] instance for our input
[QuotedDescriptorList] stack to convert the following:

	( top $ account $ engineeringEmployee )

... into ...

	( top
	$ account
	$ engineeringEmployee )

This has no effect if the stack has only one member, producing something
like:

	top
*/
func (r ObjectClasses) oIDsStringerPretty(lead int) string {
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

/*
factory default stackage closure func for oid lists - do not exec directly.
*/
func (r ObjectClasses) oIDsStringerStd(_ ...any) string {
	var _present []string
	for i := 0; i < r.len(); i++ {
		_present = append(_present, r.index(i).OID())
	}

	present := newStringBuilder()
	switch len(_present) {
	case 0:
		break
	case 1:
		present.WriteString(_present[0])
	default:
		padchar := string(rune(32))
		if !r.cast().IsPadded() {
			padchar = ``
		}

		joined := join(_present, padchar+`$`+padchar)
		present.WriteString("(" + padchar + joined + padchar + ")")
	}

	return present.String()
}

// stackage closure func - do not exec directly.
func (r ObjectClasses) canPush(x ...any) (err error) {
	if len(x) == 0 {
		return
	} else if !r.cast().IsInit() {
		err = ErrNilReceiver
		return
	}

	// scan receiver for push candidate
	for i := 0; i < len(x); i++ {
		oc, ok := x[i].(ObjectClass)
		if !ok {
			err = ErrTypeAssert
			break
		} else if oc.IsZero() {
			err = ErrNilInput
			break
		}

		if r.contains(oc.objectClass.OID) {
			err = mkerr(ErrNotUnique.Error() + ": " + oc.Type() + `, ` + oc.NumericOID())
			break
		}
	}

	return
}

/*
Len returns the current integer length of the receiver instance.
*/
func (r ObjectClasses) Len() int {
	return r.len()
}

func (r ObjectClasses) len() int {
	return r.cast().Len()
}

/*
String is a stringer method that returns the string representation
of the receiver instance.
*/
func (r ObjectClasses) String() string {
	return r.cast().String()
}

/*
Index returns the instance of [ObjectClass] found within the
receiver stack instance at index N.  If no instance is found at
the index specified, a zero [ObjectClass] instance is returned.
*/
func (r ObjectClasses) Index(idx int) ObjectClass {
	return r.index(idx)
}

func (r ObjectClasses) index(idx int) (oc ObjectClass) {
	slice, found := r.cast().Index(idx)
	if found {
		if _oc, ok := slice.(ObjectClass); ok {
			oc = _oc
		}
	}

	return
}

/*
Push returns an error following an attempt to push an [ObjectClass]
into the receiver instance.
*/
func (r ObjectClasses) Push(oc any) error {
	return r.push(oc)
}

func (r ObjectClasses) push(x any) (err error) {
	switch tv := x.(type) {
	case ObjectClass:
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
Contains calls [ObjectClasses.Get] to return a Boolean value indicative
of a successful, non-zero retrieval of an [ObjectClass] instance -- matching
the provided id -- from within the receiver stack instance.
*/
func (r ObjectClasses) Contains(id string) bool {
	return r.contains(id)
}

func (r ObjectClasses) contains(id string) bool {
	return !r.get(id).IsZero()
}

/*
Get returns an instance of [ObjectClass] based upon a search for id within
the receiver stack instance.

The return instance, if not nil, was retrieved based upon a textual match of
the principal identifier of an [ObjectClass] and the provided id.

The return instance is nil if no match was made.

Case is not significant in the matching process.
*/
func (r ObjectClasses) Get(id string) ObjectClass {
	return r.get(id)
}

func (r ObjectClasses) get(id string) (oc ObjectClass) {
	for i := 0; i < r.len() && oc.IsZero(); i++ {
		if _oc := r.index(i); !_oc.IsZero() {
			if _oc.objectClass.OID == id {
				oc = _oc
			} else if _oc.objectClass.Name.contains(id) {
				oc = _oc
			}
		}
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r ObjectClass) IsZero() bool {
	return r.objectClass == nil
}

/*
LoadObjectClasses returns an error following an attempt to load all
built-in [ObjectClass] slices into the receiver instance.
*/
func (r Schema) LoadObjectClasses() error {
	return r.loadObjectClasses()
}

func (r Schema) loadObjectClasses() (err error) {
	if !r.IsZero() {
		funks := []func() error{
			r.loadRFC4512ObjectClasses,
			r.loadRFC4519ObjectClasses,
			r.loadRFC4523ObjectClasses,
			r.loadRFC4524ObjectClasses,
			r.loadRFC2307ObjectClasses,
			r.loadRFC2079ObjectClasses,
			r.loadRFC2589ObjectClasses,
			r.loadRFC2798ObjectClasses,
			r.loadRFC3671ObjectClasses,
			r.loadRFC3672ObjectClasses,
			r.loadRFC4403ObjectClasses,
		}

		for i := 0; i < len(funks) && err == nil; i++ {
			err = funks[i]()
		}
	}

	return
}

/*
LoadRFC2079ObjectClasses returns an error following an attempt to
load all RFC 2079 [ObjectClass] slices into the receiver instance.
*/
func (r Schema) LoadRFC2079ObjectClasses() error {
	return r.loadRFC2079ObjectClasses()
}

func (r Schema) loadRFC2079ObjectClasses() (err error) {

	var i int
	for i = 0; i < len(rfc2079ObjectClasses) && err == nil; i++ {
		oc := rfc2079ObjectClasses[i]
		err = r.ParseObjectClass(string(oc))
	}

	if want := rfc2079ObjectClasses.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC2079 ObjectClasses parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC2798ObjectClasses returns an error following an attempt to
load all RFC 2798 [ObjectClass] slices into the receiver instance.
*/
func (r Schema) LoadRFC2798ObjectClasses() error {
	return r.loadRFC2798ObjectClasses()
}

func (r Schema) loadRFC2798ObjectClasses() (err error) {

	var i int
	for i = 0; i < len(rfc2798ObjectClasses) && err == nil; i++ {
		oc := rfc2798ObjectClasses[i]
		err = r.ParseObjectClass(string(oc))
	}

	if want := rfc2798ObjectClasses.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC2798 ObjectClasses parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC2307ObjectClasses returns an error following an attempt to
load all RFC 2307 [ObjectClass] slices into the receiver instance.
*/
func (r Schema) LoadRFC2307ObjectClasses() error {
	return r.loadRFC2307ObjectClasses()
}

func (r Schema) loadRFC2307ObjectClasses() (err error) {
	for k, v := range rfc2307Macros {
		r.Macros().Set(k, v)
	}

	var i int
	for i = 0; i < len(rfc2307ObjectClasses) && err == nil; i++ {
		oc := rfc2307ObjectClasses[i]
		err = r.ParseObjectClass(string(oc))
	}

	if want := rfc2307ObjectClasses.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC2307 ObjectClasses parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC3671ObjectClasses returns an error following an attempt to
load all RFC 3671 [ObjectClass] slices into the receiver instance.
*/
func (r Schema) LoadRFC3671ObjectClasses() error {
	return r.loadRFC3671ObjectClasses()
}

func (r Schema) loadRFC3671ObjectClasses() (err error) {

	var i int
	for i = 0; i < len(rfc3671ObjectClasses) && err == nil; i++ {
		oc := rfc3671ObjectClasses[i]
		err = r.ParseObjectClass(string(oc))
	}

	if want := rfc3671ObjectClasses.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC3671 ObjectClasses parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC3672ObjectClasses returns an error following an attempt to
load all RFC 3672 [ObjectClass] slices into the receiver instance.
*/
func (r Schema) LoadRFC3672ObjectClasses() error {
	return r.loadRFC3672ObjectClasses()
}

func (r Schema) loadRFC3672ObjectClasses() (err error) {

	var i int
	for i = 0; i < len(rfc3672ObjectClasses) && err == nil; i++ {
		oc := rfc3672ObjectClasses[i]
		err = r.ParseObjectClass(string(oc))
	}

	if want := rfc3672ObjectClasses.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC3672 ObjectClasses parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4512ObjectClasses returns an error following an attempt to
load all RFC 4512 [ObjectClass] slices into the receiver instance.
*/
func (r Schema) LoadRFC4512ObjectClasses() error {
	return r.loadRFC4512ObjectClasses()
}

/*
LoadRFC4512AttributeTypes returns an error following an attempt to
load all RFC 4512 [ObjectClass] slices into the receiver instance.
*/
func (r Schema) loadRFC4512ObjectClasses() (err error) {

	var i int
	for i = 0; i < len(rfc4512ObjectClasses) && err == nil; i++ {
		oc := rfc4512ObjectClasses[i]
		err = r.ParseObjectClass(string(oc))
	}

	if want := rfc4512ObjectClasses.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4512 ObjectClasses parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4519ObjectClasses returns an error following an attempt to
load all RFC 4519 [ObjectClass] slices into the receiver instance.
*/
func (r Schema) LoadRFC4519ObjectClasses() error {
	return r.loadRFC4519ObjectClasses()
}

func (r Schema) loadRFC4519ObjectClasses() (err error) {

	var i int
	for i = 0; i < len(rfc4519ObjectClasses) && err == nil; i++ {
		oc := rfc4519ObjectClasses[i]
		err = r.ParseObjectClass(string(oc))
	}

	if want := rfc4519ObjectClasses.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4519 ObjectClasses parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4523ObjectClasses returns an error following an attempt to
load all RFC 4523 [ObjectClass] slices into the receiver instance.
*/
func (r Schema) LoadRFC4523ObjectClasses() error {
	return r.loadRFC4523ObjectClasses()
}

func (r Schema) loadRFC4523ObjectClasses() (err error) {

	var i int
	for i = 0; i < len(rfc4523ObjectClasses) && err == nil; i++ {
		oc := rfc4523ObjectClasses[i]
		err = r.ParseObjectClass(string(oc))
	}

	if want := rfc4523ObjectClasses.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4523 ObjectClasses parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4524ObjectClasses returns an error following an attempt to
load all RFC 4524 [ObjectClass] slices into the receiver instance.
*/
func (r Schema) LoadRFC4524ObjectClasses() error {
	return r.loadRFC4524ObjectClasses()
}

func (r Schema) loadRFC4524ObjectClasses() (err error) {

	var i int
	for i = 0; i < len(rfc4524ObjectClasses) && err == nil; i++ {
		oc := rfc4524ObjectClasses[i]
		err = r.ParseObjectClass(string(oc))
	}

	if want := rfc4524ObjectClasses.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4524 ObjectClasses parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC2589ObjectClasses returns an error following an attempt to
load all RFC 2589 [ObjectClass] slices into the receiver instance.
*/
func (r Schema) LoadRFC2589ObjectClasses() error {
	return r.loadRFC2589ObjectClasses()
}

func (r Schema) loadRFC2589ObjectClasses() (err error) {

	var i int
	for i = 0; i < len(rfc2589ObjectClasses) && err == nil; i++ {
		oc := rfc2589ObjectClasses[i]
		err = r.ParseObjectClass(string(oc))
	}

	if want := rfc2589ObjectClasses.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC2589 ObjectClasses parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

/*
LoadRFC4403ObjectClasses returns an error following an attempt to
load all RFC 4403 [ObjectClass] slices into the receiver instance.
*/
func (r Schema) LoadRFC4403ObjectClasses() error {
	return r.loadRFC4403ObjectClasses()
}

func (r Schema) loadRFC4403ObjectClasses() (err error) {

	var i int
	for i = 0; i < len(rfc4403ObjectClasses) && err == nil; i++ {
		oc := rfc4403ObjectClasses[i]
		err = r.ParseObjectClass(string(oc))
	}

	if want := rfc4403ObjectClasses.Len(); i != want {
		if err == nil {
			err = mkerr("Unexpected number of RFC4403 ObjectClasses parsed: want " + itoa(want) + ", got " + itoa(i))
		}
	}

	return
}

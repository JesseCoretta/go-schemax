package schemax

import (
	"fmt"
	"sync"
)

/*
AttributeTypeCollection describes all of the following types:

- *AttributeTypes

- *RequiredAttributeTypes

- *PermittedAttributeTypes

- *ProhibitedAttributeTypes

- *ApplicableAttributeTypes
*/
type AttributeTypeCollection interface {
	// Contains returns the index number and presence boolean that
	// reflects the result of a term search within the receiver.
	Contains(interface{}) (int, bool)

	// Get returns the *AttributeType instance retrieved as a result
	// of a term search, based on Name or OID. If no match is found,
	// nil is returned.
	Get(interface{}) *AttributeType

	// Index returns the *AttributeType instance stored at the nth
	// index within the receiver, or nil.
	Index(int) *AttributeType

	// Equal performs a deep-equal between the receiver and the
	// interface AttributeTypeCollection provided.
	Equal(AttributeTypeCollection) bool

	// Set returns an error instance based on an attempt to add
	// the provided *AttributeType instance to the receiver.
	Set(*AttributeType) error

	// String returns a properly-delimited sequence of string
	// values, either as a Name or OID, for the receiver type.
	String() string

	// Label returns the field name associated with the interface
	// types, or a zero string if no label is appropriate.
	Label() string

	// IsZero returns a boolean value indicative of whether the
	// receiver is considered zero, or undefined.
	IsZero() bool

	// Len returns an integer value indicative of the current
	// number of elements stored within the receiver.
	Len() int
}

/*
Usage describes the intended usage of an AttributeType definition as a single text value.  This can be one of four constant values, the first of which (userApplication) is implied in the absence of any other value and is not necessary to reveal in such a case.
*/
type Usage uint8

const (
	UserApplication Usage = iota
	DirectoryOperation
	DistributedOperation
	DSAOperation
)

/*
AttributeType conforms to the specifications of RFC4512 Section 4.1.2. Boolean values, e.g: 'OBSOLETE', are supported internally and are not explicit fields.
*/
type AttributeType struct {
	OID         OID
	Name        Name
	Description Description
	SuperType   SuperiorAttributeType
	Equality    Equality
	Ordering    Ordering
	Substring   Substring
	Syntax      *LDAPSyntax
	Usage       Usage
	Extensions  Extensions
	bools       Boolean
	mub         uint
}

/*
String is an unsafe convenience wrapper for Unmarshal(r). If an error is encountered, an empty string definition is returned. If reliability and error handling are important, use Unmarshal.
*/
func (r AttributeType) String() (def string) {
	def, _ = Unmarshal(r)
	return
}

/*
SuperiorAttributeType contains an embedded instance of *AttributeType. This alias type reflects the SUP field of an attributeType definition.
*/
type SuperiorAttributeType struct {
	*AttributeType
}

/*
AttributeTypes is a thread-safe collection of *AttributeType slice instances.
*/
type AttributeTypes struct {
	mutex  *sync.Mutex
	slice  collection
	macros *Macros
}

/*
ApplicableAttributeTypes contains an embedded instance of *AttributeTypes. This alias type reflects the APPLIES field of a matchingRuleUse definition.
*/
type ApplicableAttributeTypes struct {
	*AttributeTypes
}

func (r ApplicableAttributeTypes) String() string {
        return r.slice.attrs_oids_string()
}

/*
RequiredAttributeTypes contains an embedded instance of *AttributeTypes. This alias type reflects the MUST fields of a dITContentRule or objectClass definitions.
*/
type RequiredAttributeTypes struct {
	*AttributeTypes
}

func (r RequiredAttributeTypes) String() string {
        return r.slice.attrs_oids_string()
}

/*
PermittedAttributeTypes contains an embedded instance of *AttributeTypes. This alias type reflects the MAY fields of a dITContentRule or objectClass definitions.
*/
type PermittedAttributeTypes struct {
	*AttributeTypes
}

func (r PermittedAttributeTypes) String() string {
        return r.slice.attrs_oids_string()
}

/*
ProhibitedAttributeTypes contains an embedded instance of *AttributeTypes. This alias type reflects the NOT field of a dITContentRule definition.
*/
type ProhibitedAttributeTypes struct {
	*AttributeTypes
}

func (r ProhibitedAttributeTypes) String() string {
        return r.slice.attrs_oids_string()
}

/*
SetMacros assigns the *Macros instance to the receiver, allowing subsequent OID resolution capabilities during the addition of new slice elements.
*/
func (r *AttributeTypes) SetMacros(macros *Macros) {
	r.macros = macros
}

/*
String is a stringer method that returns the string-form of the receiver instance.
*/
func (r Usage) String() string {
	switch r {
	case DirectoryOperation:
		return `directoryOperation`
	case DistributedOperation:
		return `distributedOperation`
	case DSAOperation:
		return `dSAOperation`
	}

	return `` // default is userApplication, but it need not be stated literally
}

/*
Contains is a thread-safe method that returns a collection slice element index integer and a presence-indicative boolean value based on a term search conducted within the receiver.
*/
func (r AttributeTypes) Contains(x interface{}) (int, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.macros.IsZero() {
		if oid, resolved := r.macros.Resolve(x); resolved {
			return r.slice.contains(oid)
		}
	}
	return r.slice.contains(x)
}

/*
Index is a thread-safe method that returns the nth collection slice element if defined, else nil. This method supports use of negative indices which should be used with special care.
*/
func (r AttributeTypes) Index(idx int) *AttributeType {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	assert, _ := r.slice.index(idx).(*AttributeType)
	return assert
}

/*
Get combines Contains and Index method executions to return an entry based on a term search conducted within the receiver.
*/
func (r AttributeTypes) Get(x interface{}) *AttributeType {
	idx, found := r.Contains(x)
	if !found {
		return nil
	}

	return r.Index(idx)
}

/*
Len is a thread-safe method that returns the effective length of the receiver slice collection.
*/
func (r AttributeTypes) Len() int {
	return r.slice.len()
}

/*
String is a stringer method used to return the properly-delimited and formatted series of attributeType name or OID definitions.
*/
func (r AttributeTypes) String() string {
	return ``
}

/*
IsZero returns a boolean value indicative of whether the receiver is considered empty or uninitialized.
*/
func (r *AttributeTypes) IsZero() bool {
        if r != nil {
                return r.slice.isZero()
        }
        return r == nil
}

/*
IsZero returns a boolean value indicative of whether the receiver is considered empty or uninitialized.
*/
func (r *AttributeType) IsZero() bool {
	return r == nil
}

/*
Set is a thread-safe append method that returns an error instance indicative of whether the append operation failed in some manner. Uniqueness is enforced for new elements based on Object Identifier and not the effective Name of the definition, if defined.
*/
func (r *AttributeTypes) Set(x *AttributeType) error {
	if _, exists := r.Contains(x.OID); exists {
		return fmt.Errorf("%T already contains %T:%s", r, x, x.OID)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.slice.append(x)
}

/*
Equal performs a deep-equal between the receiver and the provided collection type.
*/
func (r AttributeTypes) Equal(x AttributeTypeCollection) bool {
        return r.slice.equal(x.(*AttributeTypes).slice)
}

/*
Equal performs a deep-equal between the receiver and the provided definition type.

Description text is ignored.
*/
func (r *AttributeType) Equal(x interface{}) (equals bool) {
	z, ok := x.(*AttributeType)
	if !ok {
		return
	}

	if z.IsZero() && r.IsZero() {
		equals = true
		return
	} else if z.IsZero() || r.IsZero() {
		return
	}

	if !z.Name.Equal(r.Name) {
		return
	}

	if !r.OID.Equal(z.OID) {
		return
	}

	if z.Usage != r.Usage {
		return
	}

	if z.bools != r.bools {
		return
	}

	if !z.SuperType.IsZero() && !r.SuperType.IsZero() {
		if !z.SuperType.OID.Equal(r.SuperType.OID) {
			return
		}
	}

	if !r.Syntax.Equal(z.Syntax) {
		return
	}

	if !r.Equality.Equal(z.Equality) {
		return
	}

	if !r.Ordering.Equal(z.Ordering) {
		return
	}
	equals = r.Extensions.Equal(z.Extensions)

	return
}

/*
NewAttributeTypes initializes and returns a new AttributeTypeCollection interface object.
*/
func NewAttributeTypes() AttributeTypeCollection {
	var x interface{} = &AttributeTypes{
		mutex: &sync.Mutex{},
		slice: make(collection, 0, 0),
	}
	return x.(AttributeTypeCollection)
}

/*
NewApplicableAttributeTypes initializes an embedded instance of *AttributeTypes within the return value.
*/
func NewApplicableAttributeTypes() AttributeTypeCollection {
        var z *AttributeTypes = &AttributeTypes{
                mutex: &sync.Mutex{},
                slice: make(collection, 0, 0),
        }
        var x interface{} = &ApplicableAttributeTypes{z}
        return x.(AttributeTypeCollection)
}

/*
NewRequiredAttributeTypes initializes an embedded instance of *AttributeTypes within the return value.
*/
func NewRequiredAttributeTypes() AttributeTypeCollection {
        var z *AttributeTypes = &AttributeTypes{
                mutex: &sync.Mutex{},
                slice: make(collection, 0, 0),
        }
        var x interface{} = &RequiredAttributeTypes{z}
        return x.(AttributeTypeCollection)
}

/*
NewPermittedAttributeTypes initializes an embedded instance of *AttributeTypes within the return value.
*/
func NewPermittedAttributeTypes() AttributeTypeCollection {
        var z *AttributeTypes = &AttributeTypes{
                mutex: &sync.Mutex{},
                slice: make(collection, 0, 0),
        }
        var x interface{} = &PermittedAttributeTypes{z}
        return x.(AttributeTypeCollection)
}

/*
NewProhibitedAttributeTypes initializes an embedded instance of *AttributeTypes within the return value.
*/
func NewProhibitedAttributeTypes() AttributeTypeCollection {
        var z *AttributeTypes = &AttributeTypes{
                mutex: &sync.Mutex{},
                slice: make(collection, 0, 0),
        }
        var x interface{} = &ProhibitedAttributeTypes{z}
        return x.(AttributeTypeCollection)
}

func newUsage(x interface{}) Usage {
	switch tv := x.(type) {
	case string:
		switch toLower(tv) {
		case toLower(DirectoryOperation.String()):
			return DirectoryOperation
		case toLower(DistributedOperation.String()):
			return DistributedOperation
		case toLower(DSAOperation.String()):
			return DSAOperation
		}
	case uint:
		switch tv {
		case 0x1:
			return DirectoryOperation
		case 0x2:
			return DistributedOperation
		case 0x3:
			return DSAOperation
		}
	case int:
		if tv >= 0 {
			return newUsage(uint(tv))
		}
	}

	return UserApplication
}

/*
MaxLength returns the integer value, if one was specified, that defines the maximum acceptable value size supported by this *AttributeType per its associated *LDAPSyntax.  If not applicable, a 0 is returned.
*/
func (r AttributeType) MaxLength() int {
	return int(r.mub)
}

/*
SetMaxLength sets the minimum upper bounds, or maximum length, of the receiver instance. The argument must be a positive, non-zero integer.

This will only apply to *AttributeTypes that use a human-readable syntax.
*/
func (r *AttributeType) SetMaxLength(max int) {
	r.setMUB(max)
}

/*
setBoolean is a private method used by reflect to set the minimum upper bounds.
*/
func (r *AttributeType) setMUB(mub interface{}) {
	if r.IsZero() {
		return
	}

	if r.Syntax.IsZero() {
		return
	}

	if !r.Syntax.IsHumanReadable() {
		return
	}

	switch tv := mub.(type) {
	case string:
		n, err := atoi(tv)
		if err != nil || n < 0 {
			return
		}
		r.mub = uint(n)
	case int:
		if tv > 0 {
			r.mub = uint(tv)
		}
	case uint:
		r.mub = tv
	}

}

/*
is returns a boolean value indicative of whether the provided interface argument is either an enabled Boolean value, or an associated *MatchingRule or *LDAPSyntax.

In the case of an *LDAPSyntax argument, if the receiver is in fact a sub type of another *AttributeType instance, a reference to that super type is chased and analyzed accordingly.
*/
func (r AttributeType) is(b interface{}) bool {
	switch tv := b.(type) {
	case Boolean:
		return r.bools.is(tv)
	case *MatchingRule:
		switch {
		case tv.Equal(r.Equality.OID):
			return true
		case tv.Equal(r.Ordering.OID):
			return true
		case tv.Equal(r.Substring.OID):
			return true
		}
	case *LDAPSyntax:
		if r.Syntax != nil {
			return r.Syntax.OID.Equal(tv.OID)
		} else if !r.SuperType.IsZero() {
			return r.SuperType.is(tv)
		}
	}

	return false
}

/*
getSyntax will traverse the supertype chain upwards until it finds an explicit SYNTAX definition
*/
func (r *AttributeType) getSyntax() *LDAPSyntax {
	if r.IsZero() {
		return nil
	}
	if r.Syntax.IsZero() {
		return r.SuperType.getSyntax()
	}

	return r.Syntax
}

/*
Validate returns an error that reflects any fatal condition observed regarding the receiver configuration.
*/
func (r *AttributeType) Validate() (err error) {
	return r.validate()
}

func (r *AttributeType) validate() (err error) {
	if r.IsZero() {
		return raise(isZero, "%T.validate", r)
	}

	if err = validateBool(r.bools); err != nil {
		return
	}

	var ls *LDAPSyntax
	if ls, err = r.validateSyntax(); err != nil {
		return
	}

	// verify that no length was set with a
	// non-human-readable syntax in use.
	if !r.HumanReadable() && r.MaxLength() != 0 {
		err = raise(invalidUnmarshal,
			"%T.unmarshal: %d non-human-readable syntax:%t (%s) with minimum upper-bounds set",
			r, r.MaxLength(), r.HumanReadable(), ls.OID)
		return
	}

	if err = r.validateMatchingRules(ls); err != nil {
		return
	}

	if err = validateNames(r.Name.strings()...); err != nil {
		return
	}

	if err = validateDesc(r.Description); err != nil {
		return
	}

	if !r.SuperType.IsZero() {
		if r.SuperType.Syntax.IsZero() {
			err = raise(invalidUnmarshal, "%T.unmarshal: %T.%T: %s (sub-typed)",
				r.SuperType, r.SuperType, r.SuperType.Syntax, isZero.Error())
		}
	} else {
		if r.Syntax.IsZero() {
			err = raise(invalidUnmarshal, "%T.unmarshal: %T.%T: %s (not sub-typed)",
				r, r, r.Syntax, isZero.Error())
		}
	}

	return
}

func (r *AttributeType) validateSyntax() (ls *LDAPSyntax, err error) {
	ls = r.getSyntax()
	if ls.IsZero() {
		err = raise(invalidSyntax,
			"checkMatchingRules: %T is missing a syntax", r)
	}

	return
}

func (r *AttributeType) validateMatchingRules(ls *LDAPSyntax) (err error) {
	if err = r.validateEquality(ls); err != nil {
		return err
	}

	if err = r.validateOrdering(ls); err != nil {
		return err
	}

	if err = r.validateSubstr(ls); err != nil {
		return err
	}

	return
}

func (r *AttributeType) validateEquality(ls *LDAPSyntax) error {
	if !r.Equality.IsZero() {
		if contains(toLower(r.Equality.Name.Index(0)), `ordering`) ||
			contains(toLower(r.Equality.Name.Index(0)), `substring`) {
			return raise(invalidMatchingRule,
				"validateEquality: %T.Equality uses non-equality %T syntax (%s)",
				r, r.Equality, r.Equality.Syntax.OID.String())
		}
	}

	return nil
}

func (r *AttributeType) validateSubstr(ls *LDAPSyntax) error {
	if !r.Substring.IsZero() {
		if !contains(toLower(r.Substring.Name.Index(0)), `substring`) {
			return raise(invalidMatchingRule,
				"validateSubstr: %T.Substring uses non-substring %T syntax (%s)",
				r, r.Substring, r.Substring.Syntax.OID.String())
		}
	}

	return nil
}

func (r *AttributeType) validateOrdering(ls *LDAPSyntax) error {
	if !r.Ordering.IsZero() {
		if !contains(toLower(r.Ordering.Name.Index(0)), `ordering`) {
			return raise(invalidMatchingRule,
				"validateOrdering: %T.Ordering uses non-substring %T syntax (%s)",
				r, r.Ordering, r.Ordering.Syntax.OID.String())
		}
	}

	return nil
}

func (r *AttributeType) unmarshal(namesok bool) (def string, err error) {
	if err = r.validate(); err != nil {
		err = raise(invalidUnmarshal, err.Error())
		return
	}

	WHSP := ` `

	def += `(` + WHSP + r.OID.String() // will never be zero length

	if !r.Name.IsZero() {
		def += WHSP + r.Name.Label()
		def += WHSP + r.Name.String()
	}

	if !r.Description.IsZero() {
		def += WHSP + r.Description.Label()
		def += WHSP + r.Description.String()
	}

	if r.bools.is(Obsolete) {
		def += WHSP + Obsolete.String()
	}

	if !r.SuperType.IsZero() {
		def += WHSP + r.SuperType.Label()
		def += WHSP + r.SuperType.Name.Index(0)
	}

	if !r.Syntax.IsZero() {
		def += WHSP + r.Syntax.Label()
		def += WHSP + r.Syntax.OID.String()
	}

	if !r.Equality.IsZero() {
		def += WHSP + r.Equality.Label()
		def += WHSP + r.Equality.Name.Index(0)
	}

	if !r.Ordering.IsZero() {
		def += WHSP + r.Ordering.Label()
		def += WHSP + r.Ordering.Name.Index(0)
	}

	if !r.Substring.IsZero() {
		def += WHSP + r.Substring.Label()
		def += WHSP + r.Substring.Name.Index(0)
	}

	if r.Usage != UserApplication {
		def += WHSP + r.Usage.Label()
		def += WHSP + r.Usage.String()
	}

	if r.bools.is(NoUserModification) {
		def += WHSP + NoUserModification.String()
	}

	if !r.Extensions.IsZero() {
		def += WHSP + r.Extensions.String()
	}

	def += WHSP + `)`

	return
}

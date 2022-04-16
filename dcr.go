package schemax

import (
	"fmt"
	"sync"
)

/*
DITContentRuleCollection describes all DITContentRules-based types.
*/
type DITContentRuleCollection interface {
	// Contains returns the index number and presence boolean that
	// reflects the result of a term search within the receiver.
        Contains(interface{}) (int, bool)

	// Get returns the *DITContentRule instance retrieved as a result
	// of a term search, based on Name or OID. If no match is found,
	// nil is returned.
        Get(interface{}) *DITContentRule

	// Index returns the *DITContentRule instance stored at the nth
	// index within the receiver, or nil.
        Index(int) *DITContentRule

	// Equal performs a deep-equal between the receiver and the
	// interface DITContentRuleCollection provided.
        Equal(DITContentRuleCollection) bool

	// Set returns an error instance based on an attempt to add
	// the provided *DITContentRule instance to the receiver.
        Set(*DITContentRule) error

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
DITContentRule conforms to the specifications of RFC4512 Section 4.1.6. Boolean values, e.g: 'OBSOLETE', are supported internally and are not explicit fields.

The OID value of this type MUST match the OID of a known (and catalogued) STRUCTURAL *ObjectClass instance.
*/
type DITContentRule struct {
	OID         OID
	Name        Name
	Description Description
	Aux         ObjectClassCollection
	Must        AttributeTypeCollection
	May         AttributeTypeCollection
	Not         AttributeTypeCollection
	Extensions  Extensions
	bools       Boolean
}

/*
DITContentRules is a thread-safe collection of *DITContentRule slice instances.
*/
type DITContentRules struct {
	mutex  *sync.Mutex
	slice  collection
	macros *Macros
}

/*
Equal performs a deep-equal between the receiver and the provided collection type.
*/
func (r DITContentRules) Equal(x DITContentRuleCollection) bool {
        return r.slice.equal(x.(*DITContentRules).slice)
}

/*
SetMacros assigns the *Macros instance to the receiver, allowing subsequent OID resolution capabilities during the addition of new slice elements.
*/
func (r *DITContentRules) SetMacros(macros *Macros) {
	r.macros = macros
}

/*
Contains is a thread-safe method that returns a collection slice element index integer and a presence-indicative boolean value based on a term search conducted within the receiver.
*/
func (r DITContentRules) Contains(x interface{}) (int, bool) {
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
func (r DITContentRules) Index(idx int) *DITContentRule {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	assert, _ := r.slice.index(idx).(*DITContentRule)
	return assert
}

/*
Get combines Contains and Index method executions to return an entry based on a term search conducted within the receiver.
*/
func (r DITContentRules) Get(x interface{}) *DITContentRule {
	idx, found := r.Contains(x)
	if !found {
		return nil
	}

	return r.Index(idx)
}

/*
Len is a thread-safe method that returns the effective length of the receiver slice collection.
*/
func (r DITContentRules) Len() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.slice.len()
}

/*
String is a stringer method used to return the properly-delimited and formatted series of attributeType name or OID definitions.
*/
func (r DITContentRules) String() string { return `` }

/*
String is an unsafe convenience wrapper for Unmarshal(r). If an error is encountered, an empty string definition is returned. If reliability and error handling are important, use Unmarshal.
*/
func (r DITContentRule) String() (def string) {
	def, _ = Unmarshal(r)
	return
}

/*
IsZero returns a boolean value indicative of whether the receiver is considered empty or uninitialized.
*/
func (r DITContentRules) IsZero() bool {
	return r.slice.len() == 0
}

/*
IsZero returns a boolean value indicative of whether the receiver is considered empty or uninitialized.
*/
func (r *DITContentRule) IsZero() bool {
	return r == nil
}

/*
Set is a thread-safe append method that returns an error instance indicative of whether the append operation failed in some manner. Uniqueness is enforced for new elements based on Object Identifier and not the effective Name of the definition, if defined.
*/
func (r *DITContentRules) Set(x *DITContentRule) error {
	if _, exists := r.Contains(x.OID); exists {
		return fmt.Errorf("%T already contains %T:%s", r, x, x.OID)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.slice.append(x)
}

/*
Belongs returns a boolean value indicative of whether the provided AUXILIARY *ObjectClass belongs to the receiver instance of *DITContentRule.
*/
func (r DITContentRule) Belongs(aux *ObjectClass) (belongs bool) {
	if aux.IsZero() || !aux.Kind.is(Auxiliary) {
		return
	}

	_, belongs = r.Aux.Contains(aux)

	return
}

/*
Requires returns a boolean value indicative of whether the provided value is required per the receiver.
*/
func (r DITContentRule) Requires(x interface{}) (required bool) {
	switch tv := x.(type) {
	case *AttributeType:
		_, required = r.Must.Contains(tv)
	}

	return
}

/*
Permits returns a boolean value indicative of whether the provided value is allowed from use per the receiver.
*/
func (r DITContentRule) Permits(x interface{}) (permitted bool) {
	switch tv := x.(type) {
	case *AttributeType:
		_, permitted = r.May.Contains(tv)
	}

	return
}

/*
Prohibits returns a boolean value indicative of whether the provided value is prohibited from use per the receiver.
*/
func (r DITContentRule) Prohibits(x interface{}) (prohibited bool) {
	switch tv := x.(type) {
	case *AttributeType:
		_, prohibited = r.Not.Contains(tv)
	}

	return
}

/*
Equal performs a deep-equal between the receiver and the provided definition type.

Description text is ignored.
*/
func (r *DITContentRule) Equal(x interface{}) (equals bool) {

	z, ok := x.(*DITContentRule)
	if !ok {
		return
	}

	if z.IsZero() && r.IsZero() {
		equals = true
		return
	} else if z.IsZero() || r.IsZero() {
		return
	}

	if !r.OID.Equal(z.OID) {
		return
	}

	if !r.Name.Equal(z.Name) {
		return
	}

	if !r.Aux.Equal(z.Aux) {
		return
	}

	if !r.Must.Equal(z.Must) {
		return
	}

	if !r.May.Equal(z.May) {
		return
	}

	if !r.Not.Equal(z.Not) {
		return
	}

	equals = r.Extensions.Equal(z.Extensions)

	return
}

/*
NewDITContentRules initializes and returns a new DITContentRuleCollection interface object.
*/
func NewDITContentRules() DITContentRuleCollection {
	var x interface{} = &DITContentRules{
		mutex: &sync.Mutex{},
		slice: make(collection, 0, 0),
	}
	return x.(DITContentRuleCollection)
}

/*
Validate returns an error that reflects any fatal condition observed regarding the receiver configuration.
*/
func (r *DITContentRule) Validate() (err error) {
	return r.validate()
}

func (r *DITContentRule) validate() (err error) {
	if r.IsZero() {
		return raise(isZero, "%T.validate", r)
	}

	if err = validateDesc(r.Description); err != nil {
		return
	}

	return
}

func (r *DITContentRule) unmarshal(namesok bool) (def string, err error) {
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

	if r.bools.enabled(Obsolete) {
		def += WHSP + r.bools.Obsolete()
	}

	if !r.Aux.IsZero() {
		def += WHSP + r.Aux.Label()
		def += WHSP + r.Aux.String()
	}

	if !r.Must.IsZero() {
		def += WHSP + r.Must.Label()
		def += WHSP + r.Must.String()
	}

	if !r.May.IsZero() {
		def += WHSP + r.May.Label()
		def += WHSP + r.May.String()
	}

	if !r.Not.IsZero() {
		def += WHSP + r.Not.Label()
		def += WHSP + r.Not.String()
	}

	if !r.Extensions.IsZero() {
		def += WHSP + r.Extensions.String()
	}

	def += WHSP + `)`

	return
}
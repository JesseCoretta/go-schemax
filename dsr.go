package schemax

import (
	"fmt"
	"sync"
)

/*
DITStructureRuleCollection describes all of the following types:

- *DITStructureRules

- *SuperiorDITStructureRules
*/
type DITStructureRuleCollection interface {
	// Contains returns the index number and presence boolean that
	// reflects the result of a term search within the receiver.
	Contains(interface{}) (int, bool)

	// Get returns the *DITStructureRule instance retrieved as a result
	// of a term search, based on Name or ID. If no match is found,
	// nil is returned.
	Get(interface{}) *DITStructureRule

	// Index returns the *DITStructureRule instance stored at the nth
	// index within the receiver, or nil.
	Index(int) *DITStructureRule

	// Equal performs a deep-equal between the receiver and the
	// interface DITStructureRuleCollection provided.
	Equal(DITStructureRuleCollection) bool

	// Set returns an error instance based on an attempt to add
	// the provided *DITStructureRule instance to the receiver.
	Set(*DITStructureRule) error

	// String returns a properly-delimited sequence of string
	// values, either as a Name or ID, for the receiver type.
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
RuleID describes a numerical identifier for an instance of DITStructureRule.
*/
type RuleID uint

/*
DITStructureRule conforms to the specifications of RFC4512 Section 4.1.7.1. Boolean values, e.g: 'OBSOLETE', are supported internally and are not explicit fields.
*/
type DITStructureRule struct {
	ID            RuleID
	Name          Name
	Description   Description
	Form          *NameForm
	SuperiorRules DITStructureRuleCollection
	Extensions    Extensions
	bools         Boolean
}

/*
DITStructureRules is a thread-safe collection of *DITStructureRule slice instances.
*/
type DITStructureRules struct {
	mutex *sync.Mutex
	slice collection
}

/*
SuperiorDITStructureRules contains an embedded instance of *DITStructureRules. This alias type reflects the SUP field of an dITStructureRule definition.
*/
type SuperiorDITStructureRules struct {
	*DITStructureRules
}

/*
Equal performs a deep-equal between the receiver and the provided collection type.
*/
func (r DITStructureRules) Equal(x DITStructureRuleCollection) bool {
        return r.slice.equal(x.(*DITStructureRules).slice)
}

/*
String is a stringer method that returns the string-form of the receiver instance.
*/
func (r RuleID) String() string {
	return itoa(int(r))
}

/*
Equal returns a boolean value indicative of whether the provided value is numerically equal to the receiver.
*/
func (r RuleID) Equal(x interface{}) bool {
	rule := NewRuleID(x)
	return r.String() == rule.String()
}

/*
Contains is a thread-safe method that returns a collection slice element index integer and a presence-indicative boolean value based on a term search conducted within the receiver.
*/
func (r DITStructureRules) Contains(x interface{}) (int, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.slice.contains(x)
}

/*
Index is a thread-safe method that returns the nth collection slice element if defined, else nil. This method supports use of negative indices which should be used with special care.
*/
func (r DITStructureRules) Index(idx int) *DITStructureRule {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	assert, _ := r.slice.index(idx).(*DITStructureRule)
	return assert
}

/*
Get combines Contains and Index method executions to return an entry based on a term search conducted within the receiver.
*/
func (r DITStructureRules) Get(x interface{}) *DITStructureRule {
	idx, found := r.Contains(x)
	if !found {
		return nil
	}

	return r.Index(idx)
}

/*
Len is a thread-safe method that returns the effective length of the receiver slice collection.
*/
func (r DITStructureRules) Len() int {
	return r.slice.len()
}

func (r DITStructureRules) String() string {
	return ``
}

/*
String is an unsafe convenience wrapper for Unmarshal(r). If an error is encountered, an empty string definition is returned. If reliability and error handling are important, use Unmarshal.
*/
func (r DITStructureRule) String() (def string) {
	def, _ = Unmarshal(r)
	return
}

/*
String is a stringer method used to return the properly-delimited and formatted series of attributeType name or OID definitions.
*/
func (r SuperiorDITStructureRules) String() string {
	return r.slice.dsrs_ids_string()
}

/*
IsZero returns a boolean value indicative of whether the receiver is considered empty or uninitialized.
*/
func (r *DITStructureRules) IsZero() bool {
        if r != nil {
                return r.slice.isZero()
        }
        return r == nil
}

/*
IsZero returns a boolean value indicative of whether the receiver is considered empty or uninitialized.
*/
func (r *DITStructureRule) IsZero() bool {
	return r == nil
}

/*
Set is a thread-safe append method that returns an error instance indicative of whether the append operation failed in some manner. Uniqueness is enforced for new elements based on Object Identifier and not the effective Name of the definition, if defined.
*/
func (r *DITStructureRules) Set(x *DITStructureRule) error {
	if _, exists := r.Contains(x.ID); exists {
		return fmt.Errorf("%T already contains %T:%s", r, x, x.ID)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.slice.append(x)
}

/*
NewDITStructureRules initializes and returns a new DITStructureRuleCollection interface object.
*/
func NewDITStructureRules() DITStructureRuleCollection {
	var x interface{} = &DITStructureRules{
		mutex: &sync.Mutex{},
		slice: make(collection, 0, 0),
	}
	return x.(DITStructureRuleCollection)
}

/*
NewSuperiorDITStructureRules initializes an embedded instance of *DITStructureRules within the return value.
*/
func NewSuperiorDITStructureRules() DITStructureRuleCollection {
	var z *DITStructureRules = &DITStructureRules{
		mutex: &sync.Mutex{},
		slice: make(collection, 0, 0),
	}
	var x interface{} = &SuperiorDITStructureRules{z}
        return x.(DITStructureRuleCollection)
}

/*
Equal performs a deep-equal between the receiver and the provided definition type.

Description text is ignored.
*/
func (r *DITStructureRule) Equal(x interface{}) (equals bool) {

	z, ok := x.(*DITStructureRule)
	if !ok {
		return
	}

	if z.IsZero() && r.IsZero() {
		equals = true
		return
	} else if z.IsZero() || r.IsZero() {
		return
	}

	if !r.ID.Equal(z.ID) {
		return
	}

	if !r.Name.Equal(z.Name) {
		return
	}

	if !r.Form.Equal(z.Form) {
		return
	}

        if !z.SuperiorRules.IsZero() && !r.SuperiorRules.IsZero() {
                if !r.SuperiorRules.Equal(z.SuperiorRules) {
                        return
                }
        }

	equals = r.Extensions.Equal(z.Extensions)

	return
}

/*
NewRuleID returns a new instance of *RuleID, intended for assignment to an instance of *DITStructureRule.
*/
func NewRuleID(x interface{}) (rid RuleID) {
	switch tv := x.(type) {
	case int:
		if tv < 0 {
			return
		}
		x := uint(tv)
		rid = RuleID(x)
	case uint:
		rid = RuleID(tv)
	case string:
		if isDigit(tv) {
			if n, err := atoi(tv); err == nil && n > 0 {
				rid = RuleID(uint(n))
			}
		}
	}

	return
}

/*
Validate returns an error that reflects any fatal condition observed regarding the receiver configuration.
*/
func (r *DITStructureRule) Validate() (err error) {
	return r.validate()
}

func (r *DITStructureRule) validate() (err error) {
	if r.IsZero() {
		err = raise(isZero, "%T.validate", r, r)
		return
	}

	if r.Form.IsZero() {
		err = raise(invalidNameForm,
			"%T.validate: missing %T",
			r, r.Form)
		return
	}

	if err = validateNames(r.Name.strings()...); err != nil {
		return
	}

	if err = validateDesc(r.Description); err != nil {
		return
	}

	if err = validateBool(r.bools); err != nil {
		return
	}

	if r.SuperiorRules == nil {
		return
	}

	for i := 0; i < r.SuperiorRules.Len(); i++ {
		if err = r.SuperiorRules.Index(i).validate(); err != nil {
			return err
		}
	}

	return
}

func (r *DITStructureRule) unmarshal(namesok bool) (def string, err error) {
	if err = r.validate(); err != nil {
		err = raise(invalidUnmarshal, err.Error())
		return
	}

	WHSP := ` `

	def += `(` + WHSP + r.ID.String() // will never be zero length

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

	// Form will never be zero
	def += WHSP + r.Form.Label()
	def += WHSP + r.Form.String()

	if !r.SuperiorRules.IsZero() {
		def += WHSP + r.SuperiorRules.Label()
		def += WHSP + r.SuperiorRules.String()
	}

	if !r.Extensions.IsZero() {
		def += WHSP + r.Extensions.String()
	}

	def += WHSP + `)`

	return
}

package schemax

import (
	"fmt"
	"sync"
)

/*
ObjectClassCollection describes all ObjectClasses-based types:

- *SuperiorObjectClasses

- *AuxiliaryObjectClasses
*/
type ObjectClassCollection interface {
	// Contains returns the index number and presence boolean that
	// reflects the result of a term search within the receiver.
	Contains(interface{}) (int, bool)

	// Get returns the *ObjectClass instance retrieved as a result
	// of a term search, based on Name or OID. If no match is found,
	// nil is returned.
	Get(interface{}) *ObjectClass

	// Index returns the *ObjectClass instance stored at the nth
	// index within the receiver, or nil.
	Index(int) *ObjectClass

	// Equal performs a deep-equal between the receiver and the
	// interface ObjectClassCollection provided.
	Equal(ObjectClassCollection) bool

	// Set returns an error instance based on an attempt to add
	// the provided *ObjectClass instance to the receiver.
	Set(*ObjectClass) error

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
Kind is an unsigned 8-bit integer that describes the "kind" of ObjectClass definition bearing this type.  Only one distinct Kind value may be set for any given ObjectClass definition, and must be set explicitly (no default is implied).
*/
type Kind uint8

const (
	badKind Kind = iota
	Abstract
	Structural
	Auxiliary
)

/*
IsZero returns a boolean value indicative of whether the receiver is undefined.
*/
func (r Kind) IsZero() bool {
	return r == badKind
}

/*
ObjectClass conforms to the specifications of RFC4512 Section 4.1.1. Boolean values, e.g: 'OBSOLETE', are supported internally and are not explicit fields.
*/
type ObjectClass struct {
	OID         OID
	Name        Name
	Description Description
	SuperClass  ObjectClassCollection
	Kind        Kind
	Must        AttributeTypeCollection
	May         AttributeTypeCollection
	Extensions  Extensions
	bools       Boolean
}

/*
ObjectClasses is a thread-safe collection of *ObjectClass slice instances.
*/
type ObjectClasses struct {
	mutex  *sync.Mutex
	slice  collection
	macros *Macros
}

/*
StructuralObjectClass is a type alias of *ObjectClass intended for use solely within instances of NameForm within its "OC" field.
*/
type StructuralObjectClass struct {
	*ObjectClass
}

/*
SuperiorObjectClasses contains an embedded *ObjectClasses instance. This type alias is meant to reside within the SUP field of an objectClass definition.
*/
type SuperiorObjectClasses struct {
	*ObjectClasses
}

/*
AuxiliaryObjectClasses contains an embedded *ObjectClasses instance. This type alias is meant to reside within the AUX field of a dITContentRule definition.
*/
type AuxiliaryObjectClasses struct {
	*ObjectClasses
}

/*
SetMacros assigns the *Macros instance to the receiver, allowing subsequent OID resolution capabilities during the addition of new slice elements.
*/
func (r *ObjectClasses) SetMacros(macros *Macros) {
	r.macros = macros
}

/*
String is an unsafe convenience wrapper for Unmarshal(r). If an error is encountered, an empty string definition is returned. If reliability and error handling are important, use Unmarshal.
*/
func (r ObjectClass) String() (def string) {
	def, _ = Unmarshal(r)
	return
}

/*
String is a stringer method that returns the string-form of the receiver instance.
*/
func (r Kind) String() string {
	switch r {
	case Abstract:
		return `ABSTRACT`
	case Structural:
		return `STRUCTURAL`
	case Auxiliary:
		return `AUXILIARY`
	}

	return `` // no default
}

/*
Contains is a thread-safe method that returns a collection slice element index integer and a presence-indicative boolean value based on a term search conducted within the receiver.
*/
func (r ObjectClasses) Contains(x interface{}) (int, bool) {
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
func (r ObjectClasses) Index(idx int) *ObjectClass {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	assert, _ := r.slice.index(idx).(*ObjectClass)
	return assert
}

/*
Get combines Contains and Index method executions to return an entry based on a term search conducted within the receiver.
*/
func (r ObjectClasses) Get(x interface{}) *ObjectClass {
	idx, found := r.Contains(x)
	if !found {
		return nil
	}

	return r.Index(idx)
}

/*
Len is a thread-safe method that returns the effective length of the receiver slice collection.
*/
func (r ObjectClasses) Len() int {
	return r.slice.len()
}

/*
IsZero returns a boolean value indicative of whether the receiver is considered empty or uninitialized.
*/
func (r *ObjectClasses) IsZero() bool {
        if r != nil {
                return r.slice.isZero()
        }
        return r == nil
}

/*
IsZero returns a boolean value indicative of whether the receiver is considered empty or uninitialized.
*/
func (r *ObjectClass) IsZero() bool {
	return r == nil
}

/*
Set is a thread-safe append method that returns an error instance indicative of whether the append operation failed in some manner. Uniqueness is enforced for new elements based on Object Identifier and not the effective Name of the definition, if defined.
*/
func (r *ObjectClasses) Set(x *ObjectClass) error {
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
func (r ObjectClasses) Equal(x ObjectClassCollection) bool {
	return r.slice.equal(x.(*ObjectClasses).slice)
}

/*
Equal performs a deep-equal between the receiver and the provided definition type.

Description text is ignored.
*/
func (r *ObjectClass) Equal(x interface{}) (equals bool) {

	z, ok := x.(*ObjectClass)
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

	if r.Kind != z.Kind {
		return
	}

	if r.bools != z.bools {
		return
	}

	if !r.Must.Equal(z.Must) {
		return
	}

	if !r.May.Equal(z.May) {
		return
	}

	if !z.SuperClass.IsZero() && !r.SuperClass.IsZero() {
		if !r.SuperClass.Equal(z.SuperClass) {
			return
		}
	}

	equals = r.Extensions.Equal(z.Extensions)

	return
}

/*
NewObjectClasses returns an initialized instance of ObjectClasses cast as an ObjectClassCollection.
*/
func NewObjectClasses() ObjectClassCollection {
	var x interface{} = &ObjectClasses{
		mutex: &sync.Mutex{},
		slice: make(collection, 0, 0),
	}
	return x.(ObjectClassCollection)
}

/*
NewSuperiorObjectClasses returns an initialized instance of SuperiorObjectClasses cast as an ObjectClassCollection.
*/
func NewSuperiorObjectClasses() ObjectClassCollection {
        var z *ObjectClasses = &ObjectClasses{
                mutex: &sync.Mutex{},
                slice: make(collection, 0, 0),
        }
        var x interface{} = &SuperiorObjectClasses{z}
        return x.(ObjectClassCollection)
}

/*
NewAuxiliaryObjectClasses returns an initialized instance of AuxiliaryObjectClasses cast as an ObjectClassCollection.
*/
func NewAuxiliaryObjectClasses() ObjectClassCollection {
        var z *ObjectClasses = &ObjectClasses{
                mutex: &sync.Mutex{},
                slice: make(collection, 0, 0),
        }
        var x interface{} = &AuxiliaryObjectClasses{z}
        return x.(ObjectClassCollection)
}

func newKind(x interface{}) Kind {
	switch tv := x.(type) {
	case Kind:
		return newKind(tv.String())
	case string:
		switch toLower(tv) {
		case toLower(Abstract.String()):
			return Abstract
		case toLower(Structural.String()):
			return Structural
		case toLower(Auxiliary.String()):
			return Auxiliary
		}
	case uint:
		switch tv {
		case 0x1:
			return Abstract
		case 0x2:
			return Structural
		case 0x3:
			return Auxiliary
		}
	case int:
		if tv >= 0 {
			return newKind(uint(tv))
		}
	}

	return badKind
}

func (r Kind) is(x Kind) bool {
	return r == x
}

/*
is returns a boolean value indicative of whether the provided interface value is either a Kind or a Boolean AND is enabled within the receiver.
*/
func (r ObjectClass) is(b interface{}) bool {
	switch tv := b.(type) {
	case Boolean:
		return r.bools.is(tv)
	case Kind:
		return r.Kind.is(tv)
	}

	return false
}

func (r *ObjectClass) validateKind() (err error) {
	if newKind(r.Kind.String()) == badKind {
		err = invalidObjectClassKind
	}

	return
}

/*
Validate returns an error that reflects any fatal condition observed regarding the receiver configuration.
*/
func (r *ObjectClass) Validate() (err error) {
	return r.validate()
}

func (r *ObjectClass) validate() (err error) {
	if r.IsZero() {
		return raise(isZero, "%T.validate", r)
	}

	if err = validateBool(r.bools); err != nil {
		return
	}

	if err = r.validateKind(); err != nil {
		return
	}

	if err = validateNames(r.Name.strings()...); err != nil {
		return
	}

	if err = validateDesc(r.Description); err != nil {
		return
	}

	return
}

func (r *ObjectClass) getMay(m AttributeTypeCollection) (ok PermittedAttributeTypes) {
	for _, atr := range m.(*AttributeTypes).slice {
		at, assert := atr.(*AttributeType)
		if !assert {
			return
		}
		ok.Set(at)
	}

	if !r.SuperClass.IsZero() {
		for i := 0; i < r.SuperClass.Len(); i++ {
			oc := r.SuperClass.Index(0)
			if oc.IsZero() {
				continue
			}
			for j := 0; j < oc.May.Len(); j++ {
				may := oc.May.Index(j)
				if may.IsZero() {
					ok.Set(may)
				}
			}
		}
	}

	if !r.May.IsZero() {
                for i := 0; i < r.May.Len(); i++ {
                        may := r.May.Index(0)
                        if may.IsZero() {
                                continue
                        }
                        ok.Set(may)
                }
	}

	return
}

func (r *ObjectClass) getMust(m RequiredAttributeTypes) (req RequiredAttributeTypes) {
        for _, atr := range m.slice {
                at, ok := atr.(*AttributeType)
                if !ok {
                        return
                }
                req.Set(at)
        }

        if !r.SuperClass.IsZero() {
                for i := 0; i < r.SuperClass.Len(); i++ {
                        oc := r.SuperClass.Index(0)
                        if oc.IsZero() {
                                continue
                        }
                        for j := 0; j < oc.Must.Len(); j++ {
                                must := oc.Must.Index(j)
                                if must.IsZero() {
                                        req.Set(must)
                                }
                        }
                }
        }

        if !r.Must.IsZero() {
                for i := 0; i < r.Must.Len(); i++ {
                        must := r.Must.Index(0)
                        if must.IsZero() {
                                continue
                        }
                        req.Set(must)
                }
        }

	return
}

/*
String returns a properly-delimited sequence of string values, either as a Name or OID, for the receiver type.
*/
func (r ObjectClasses) String() string {
	return r.slice.ocs_oids_string()
}

/*
String returns a properly-delimited sequence of string values, either as a Name or OID, for the receiver type.
*/
func (r SuperiorObjectClasses) String() string {
	return r.slice.ocs_oids_string()
}

/*
String returns a properly-delimited sequence of string values, either as a Name or OID, for the receiver type.
*/
func (r AuxiliaryObjectClasses) String() string {
	return r.slice.ocs_oids_string()
}

func (r *ObjectClass) unmarshal(namesok bool) (def string, err error) {
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

	if r.is(Obsolete) {
		def += WHSP + r.bools.Obsolete()
	}

	if !r.SuperClass.IsZero() {
		def += WHSP + r.SuperClass.Label()
		def += WHSP + r.SuperClass.String()
	}

	// Kind will never be zero
	def += WHSP + r.Kind.String()

	if !r.Must.IsZero() {
		def += WHSP + r.Must.Label()
		def += WHSP + r.Must.String()
	}
	if !r.May.IsZero() {
		def += WHSP + r.May.Label()
		def += WHSP + r.May.String()
	}

	if !r.Extensions.IsZero() {
		def += WHSP + r.Extensions.String()
	}

	def += WHSP + `)`

	return
}
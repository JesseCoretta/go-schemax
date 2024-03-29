package schemax

import "sync"

/*
DITStructureRuleCollection describes all of the following types:

- *DITStructureRules

- *SuperiorDITStructureRules
*/
type DITStructureRuleCollection interface {
	// Get returns the *DITStructureRule instance retrieved as a result
	// of a term search, based on Name or ID. If no match is found,
	// nil is returned.
	Get(any) *DITStructureRule

	// Index returns the *DITStructureRule instance stored at the nth
	// index within the receiver, or nil.
	Index(int) *DITStructureRule

	// Equal performs a deep-equal between the receiver and the
	// interface DITStructureRuleCollection provided.
	Equal(DITStructureRuleCollection) bool

	// Set returns an error instance based on an attempt to add
	// the provided *DITStructureRule instance to the receiver.
	Set(*DITStructureRule) error

	// Contains returns the index number and presence boolean that
	// reflects the result of a term search within the receiver.
	Contains(any) (int, bool)

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

	// SetSpecifier assigns a string value to all definitions within
	// the receiver. This value is used in cases where a definition
	// type name (e.g.: attributetype, objectclass, etc.) is required.
	// This value will be displayed at the beginning of the definition
	// value during the unmarshal or unsafe stringification process.
	SetSpecifier(string)

	// SetUnmarshaler assigns the provided DefinitionUnmarshaler
	// signature to all definitions within the receiver. The provided
	// function shall be executed during the unmarshal or unsafe
	// stringification process.
	SetUnmarshaler(DefinitionUnmarshaler)
}

/*
RuleID describes a numerical identifier for an instance of DITStructureRule.
*/
type RuleID uint

/*
DITStructureRule conforms to the specifications of RFC4512 Section 4.1.7.1.
*/
type DITStructureRule struct {
	ID            RuleID
	Name          Name
	Description   Description
	Obsolete      bool
	Form          *NameForm
	SuperiorRules DITStructureRuleCollection
	Extensions    *Extensions
	ufn           DefinitionUnmarshaler
	spec          string
	info          []byte
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
Type returns the formal name of the receiver in order to satisfy signature requirements of the Definition interface type.
*/
func (r *DITStructureRule) Type() string {
	return `DITStructureRule`
}

/*
Equal performs a deep-equal between the receiver and the provided collection type.
*/
func (r DITStructureRules) Equal(x DITStructureRuleCollection) bool {
	return r.slice.equal(x.(*DITStructureRules).slice)
}

/*
SetSpecifier is a convenience method that executes the SetSpecifier method in iterative fashion for all definitions within the receiver.
*/
func (r *DITStructureRules) SetSpecifier(spec string) {
	for i := 0; i < r.Len(); i++ {
		r.Index(i).SetSpecifier(spec)
	}
}

/*
SetUnmarshaler is a convenience method that executes the SetUnmarshaler method in iterative fashion for all definitions within the receiver.
*/
func (r *DITStructureRules) SetUnmarshaler(fn DefinitionUnmarshaler) {
	for i := 0; i < r.Len(); i++ {
		r.Index(i).SetUnmarshaler(fn)
	}
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
func (r RuleID) Equal(x any) bool {
	rule := NewRuleID(x)
	return r.String() == rule.String()
}

/*
Contains is a thread-safe method that returns a collection slice element index integer and a presence-indicative boolean value based on a term search conducted within the receiver.
*/
func (r DITStructureRules) Contains(x any) (int, bool) {
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
func (r DITStructureRules) Get(x any) *DITStructureRule {
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
	if &r == nil {
		return 0
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.slice.len()
}

/*
String is a non-functional stringer method needed to satisfy interface type requirements and should not be used. See the SuperiorDITStructureRules.String() method instead.
*/
func (r DITStructureRules) String() string { return `` }

/*
String is an unsafe convenience wrapper for Unmarshal(r). If an error is encountered, an empty string definition is returned. If reliability and error handling are important, use Unmarshal.
*/
func (r DITStructureRule) String() (def string) {
	def, _ = r.unmarshal()
	return
}

/*
SetSpecifier assigns a string value to the receiver, useful for placement into configurations that require a type name (e.g.: ditstructurerule). This will be displayed at the beginning of the definition value during the unmarshal or unsafe stringification process.
*/
func (r *DITStructureRule) SetSpecifier(spec string) {
	r.spec = spec
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
		return nil //silent
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.slice.append(x)
}

/*
SetInfo assigns the byte slice to the receiver. This is a user-leveraged field intended to allow arbitrary information (documentation?) to be assigned to the definition.
*/
func (r *DITStructureRule) SetInfo(info []byte) {
	r.info = info
}

/*
Info returns the assigned informational byte slice instance stored within the receiver.
*/
func (r *DITStructureRule) Info() []byte {
	return r.info
}

/*
SetUnmarshaler assigns the provided DefinitionUnmarshaler signature value to the receiver. The provided function shall be executed during the unmarshal or unsafe stringification process.
*/
func (r *DITStructureRule) SetUnmarshaler(fn DefinitionUnmarshaler) {
	r.ufn = fn
}

/*
NewDITStructureRule returns a newly initialized, yet effectively nil, instance of *DITStructureRule.

Users generally do not need to execute this function unless an instance of the returned type will be manually populated (as opposed to parsing a raw text definition).
*/
func NewDITStructureRule() *DITStructureRule {
	dsr := new(DITStructureRule)
	dsr.SuperiorRules = NewSuperiorDITStructureRules()
	dsr.Extensions = NewExtensions()
	return dsr
}

/*
NewDITStructureRules initializes and returns a new DITStructureRuleCollection interface object.
*/
func NewDITStructureRules() DITStructureRuleCollection {
	var x any = &DITStructureRules{
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
	var x any = &SuperiorDITStructureRules{z}
	return x.(DITStructureRuleCollection)
}

/*
Equal performs a deep-equal between the receiver and the provided definition type.

Description text is ignored.
*/
func (r *DITStructureRule) Equal(x any) (eq bool) {

	z, ok := x.(*DITStructureRule)
	if !ok {
		return
	}

	if z.IsZero() && r.IsZero() {
		eq = true
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

	noexts := z.Extensions.IsZero() && r.Extensions.IsZero()
	if !noexts {
		eq = r.Extensions.Equal(z.Extensions)
	} else {
		eq = true
	}

	return
}

/*
NewRuleID returns a new instance of *RuleID, intended for assignment to an instance of *DITStructureRule.
*/
func NewRuleID(x any) (rid RuleID) {
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

func (r *DITStructureRule) unmarshal() (string, error) {
	if err := r.validate(); err != nil {
		err = raise(invalidUnmarshal, err.Error())
		return ``, err
	}

	if r.ufn != nil {
		return r.ufn(r)
	}
	return r.unmarshalBasic()
}

/*
Map is a convenience method that returns a map[string][]string instance containing the effective contents of the receiver.
*/
func (r *DITStructureRule) Map() (def map[string][]string) {
	if err := r.Validate(); err != nil {
		return
	}

	def = make(map[string][]string, 14)
	def[`RAW`] = []string{r.String()}
	def[`ID`] = []string{r.ID.String()}
	def[`TYPE`] = []string{r.Type()}

	if len(r.info) > 0 {
		def[`INFO`] = []string{string(r.info)}
	}

	if !r.Name.IsZero() {
		def[`NAME`] = make([]string, 0)
		for i := 0; i < r.Name.Len(); i++ {
			def[`NAME`] = append(def[`NAME`], r.Name.Index(i))
		}
	}

	if len(r.Description) > 0 {
		def[`DESC`] = []string{r.Description.String()}
	}

	if !r.Form.IsZero() {
		def[`FORM`] = []string{r.Form.OID.String(), r.Form.Name.Index(0)}
	}

	if !r.SuperiorRules.IsZero() {
		def[`SUP`] = make([]string, 0)
		for i := 0; i < r.SuperiorRules.Len(); i++ {
			sup := r.SuperiorRules.Index(i)
			term := sup.ID.String()
			def[`SUP`] = append(def[`SUP`], term)
		}
	}

	if !r.Extensions.IsZero() {
		for i := 0; i < r.Extensions.Len(); i++ {
			if ext := r.Extensions.Index(i); !ext.IsZero() {
				def[ext.Label] = ext.Value
			}
		}
	}

	if r.Obsolete {
		def[`OBSOLETE`] = []string{`TRUE`}
	}

	return def
}

/*
DITStructureRuleUnmarshaler is a package-included function that honors the signature of the first class (closure) DefinitionUnmarshaler type.

The purpose of this function, and similar user-devised ones, is to unmarshal a definition with specific formatting included, such as linebreaks, leading specifier declarations and indenting.
*/
func DITStructureRuleUnmarshaler(x any) (def string, err error) {
	var r *DITStructureRule
	switch tv := x.(type) {
	case *DITStructureRule:
		if tv.IsZero() {
			err = raise(isZero, "%T is nil", tv)
			return
		}
		r = tv
	default:
		err = raise(unexpectedType,
			"Bad type for unmarshal (%T)", tv)
		return
	}

	var (
		WHSP string = ` `
		idnt string = "\n\t"
		head string = `(`
		tail string = `)`
	)

	if len(r.spec) > 0 {
		head = r.spec + WHSP + head
	}

	def += head + WHSP + r.ID.String()

	if !r.Name.IsZero() {
		def += idnt + r.Name.Label()
		def += WHSP + r.Name.String()
	}

	if !r.Description.IsZero() {
		def += idnt + r.Description.Label()
		def += WHSP + r.Description.String()
	}

	if r.Obsolete {
		def += idnt + `OBSOLETE`
	}

	// Form will never be zero
	def += idnt + r.Form.Label()
	def += WHSP + r.Form.String()

	if !r.SuperiorRules.IsZero() {
		def += idnt + r.SuperiorRules.Label()
		def += WHSP + r.SuperiorRules.String()
	}

	if !r.Extensions.IsZero() {
		for i := 0; i < r.Extensions.Len(); i++ {
			if ext := r.Extensions.Index(i); !ext.IsZero() {
				def += idnt + ext.String()
			}
		}
	}

	def += WHSP + tail

	return
}

func (r *DITStructureRule) unmarshalBasic() (def string, err error) {
	var (
		WHSP string = ` `
		head string = `(`
		tail string = `)`
	)

	if len(r.spec) > 0 {
		head = r.spec + WHSP + head
	}

	def += head + WHSP + r.ID.String()

	if !r.Name.IsZero() {
		def += WHSP + r.Name.Label()
		def += WHSP + r.Name.String()
	}

	if !r.Description.IsZero() {
		def += WHSP + r.Description.Label()
		def += WHSP + r.Description.String()
	}

	if r.Obsolete {
		def += WHSP + `OBSOLETE`
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

	def += WHSP + tail

	return
}

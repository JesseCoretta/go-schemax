package schemax

/*
NewExtensions initializes and returns a new instance of [Extensions].
*/
func NewExtensions(o ...Option) (e Extensions) {
	e = newExtensions(o...)
	e.cast().SetPushPolicy(e.canPush)
	return
}

func (r Extensions) Push(x any) error {
	return r.push(x)
}

func (r Extensions) push(extn any) (err error) {
	if extn == nil {
		err = ErrNilInput
		return
	}

	r.cast().Push(extn)

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r Extension) IsZero() bool {
	return r.extension == nil
}

func (r Extensions) canPush(x ...any) (err error) {
	if len(x) == 0 {
		return
	}

	for i := 0; i < len(x) && err == nil; i++ {
		instance := x[i]
		if e, ok := instance.(Extension); !ok || e.IsZero() {
			err = ErrTypeAssert
		}
	}

	return
}

/*
Definition returns the [Definition] instance to which the receiver
instance is assigned.
*/
func (r Extensions) Definition() Definition {
	_m := r.cast().Auxiliary()[`def`]
	m, _ := _m.(Definition)
	return m
}

/*
setDefinition assigns input [Definition] x to the receiver instance.

This is a fluent method.
*/
func (r Extensions) setDefinition(x Definition) Extensions {
	r.cast().Auxiliary()[`def`] = x
	return r
}

/*
String is a stringer method that returns a properly formatted
quoted descriptor list based on the number of string values
within the receiver instance.
*/
func (r Extensions) String() (s string) {
	for i := 0; i < r.Len(); i++ {
		s += r.Index(i).String()
	}

	return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r Extensions) IsZero() bool {
	return r.cast().IsZero()
}

/*
Set assigns the provided key and values instances to the receiver instance,
thereby specifying a new extension value as described in RFC 4512.
*/
func (r Extensions) Set(key string, values ...string) {
	_key := uc(key)
	var hi, se bool
	if def := r.Definition(); def != nil {
		if sch := def.Schema(); !sch.IsZero() {
			opts := sch.Options()
			hi = opts.Positive(HangingIndents)
			se = opts.Positive(SortExtensions)
		}
	}
	if _values, found := r.get(_key); !found {
		_values = newQStringList(`extensions`)
		for v := 0; v < len(values); v++ {
			_values.cast().Push(values[v])
		}

		if _values.cast().Len() > 0 {
			ext := new(extension)
			ext.XString = key
			ext.Values = _values
			ext.hindent = hi
			ext.sortExts = se
			r.Push(Extension{ext})
		}
	} else {
		for i := 0; i < len(values); i++ {
			if !_values.Contains(values[i]) {
				_values.cast().Push(values[i])
			}
		}
	}
}

/*
Exists returns a Boolean value indicative of whether the specified key
exists within the receiver instance.  Case is not significant in the
matching process.
*/
func (r Extensions) Exists(key string) bool {
	return r.exists(key)
}

func (r Extensions) exists(key string) (found bool) {
	_, found = r.get(key)
	return
}

/*
Len returns the current integer length of the receiver instance.
*/
func (r Extensions) Len() int {
	return r.len()
}

func (r Extensions) len() int {
	return r.cast().Len()
}

/*
Get returns [QuotedStringList] and Boolean instances based
on a successful retrieval of values associated with key in
the receiver instance.  The Boolean value is indicative of
a positive match.  Case is not significant in the matching
process.
*/
func (r Extensions) Get(key string) (QuotedStringList, bool) {
	return r.get(key)
}

func (r Extensions) get(key string) (values QuotedStringList, found bool) {
	for i := 0; i < r.len() && values.cast().Len() == 0; i++ {
		if eq(r.index(i).XString, key) {
			values = r.index(i).Values
		}
	}

	found = values.cast().Len() > 0

	return
}

/*
Keys returns all keys found within the underlying map instance.
*/
func (r Extensions) Keys() (keys []string) {
	for i := 0; i < r.len(); i++ {
		keys = append(keys, r.index(i).XString)
	}

	return
}

func (r Extensions) tmplFunc() (e string) {
	if !r.IsZero() {
		e = r.String()
	}

	return
}

/*
String returns the string representation of the receiver instance.
*/
func (r Extension) String() (s string) {
	if !r.IsZero() {
		switch r.Values.Len() {
		case 0:
			break
		case 1:
			s = hindent(r.hindent) + r.XString + ` ` + `'` + r.Values.Index(0) + `'`
		default:
			s = hindent(r.hindent) + r.XString + ` ` + r.Values.String()
		}
	}

	return
}

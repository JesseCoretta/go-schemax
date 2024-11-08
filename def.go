package schemax

/*
defmap.go defines map related methods to allow simple interaction
with instances of types that have been marshaled into map form.
*/

func (r DefinitionMaps) Len() int {
	return len(r)
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r DefinitionMaps) IsZero() bool {
	return r.Len() == 0
}

func (r DefinitionMaps) Index(idx int) DefinitionMap {
	if r.Len() < idx {
		return DefinitionMap{}
	}

	return r[idx]
}

func (r DefinitionMap) Get(id string) (val []string) {
	var key string
	for key, val = range r {
		if eq(key, id) {
			break
		}
	}

	return
}

func (r DefinitionMap) Contains(id string) bool {
	return len(r.Get(id)) > 0
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r DefinitionMap) IsZero() bool {
	return r.Len() == 0
}

func (r DefinitionMap) Len() int {
	return len(r)
}

/*
Keys returns slices of names, each representing a key residing within
the receiver instance. Note that the order of keys cannot be guaranteed.
*/
func (r DefinitionMap) Keys() (keys []string) {
	for key := range r {
		keys = append(keys, key)
	}

	return
}

/*
clean removes valueless (or zero-length valued) keys/value pairs
from the receiver instance.
*/
func (r *DefinitionMap) clean() {
	for k, v := range *r {
		if len(v) == 0 {
			delete(*r, k)
		} else if len(v[0]) == 0 {
			delete(*r, k)
		} else if k == `MUB` && v[0] == `0` {
			delete(*r, k)
		}
	}
}

/*
Type returns the first value assigned to the TYPE key, if defined,
else `unknown` is returned.
*/
func (r DefinitionMap) Type() (t string) {
	t = `unknown`
	if _type := r.Get(`TYPE`); len(_type) > 0 {
		if _t := _type[0]; len(_t) > 0 {
			t = _t
		}
	}

	return
}

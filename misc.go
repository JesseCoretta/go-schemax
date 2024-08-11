package schemax

/*
IsNumericOID wraps [objectid.NewDotNotation] to parse input value id in
order to assess its validity as an ASN.1 OBJECT IDENTIFIER in dot form,
e.g.:

	1.3.6.1.4.1.56521.999

The qualifications are as follows:

  - Must only consist of digits (arcs) and dot (ASCII FULL STOP) delimiters
  - Must begin with a root arc: 0, 1 or 2
  - Second-level arcs below root arcs 0 and 1 cannot be greater than 39
  - Cannot end with a dot
  - Dots cannot be contiguous
  - Though arcs are unbounded, no arc may ever be negative
  - OID must consist of at least two (2) arcs
*/
func IsNumericOID(id string) bool {
	return isNumericOID(id)
}

func isNumericOID(id string) bool {
	if !isValidOIDPrefix(id) {
		return false
	}

	var last rune
	for i, c := range id {
		switch {
		case c == '.':
			if last == c {
				return false
			} else if i == len(id)-1 {
				return false
			}
			last = '.'
		case ('0' <= c && c <= '9'):
			last = c
			continue
		}
	}

	return true
}

func isValidOIDPrefix(id string) bool {
	slices := split(id, `.`)
	if len(slices) < 2 {
		return false
	}

	root, err := atoi(slices[0])
	if err != nil {
		return false
	}
	if !(0 <= root && root <= 2) {
		return false
	}

	var sub int
	if sub, err = atoi(slices[1]); err != nil {
		return false
	} else if !(0 <= sub && sub <= 39) && root != 2 {
		return false
	}

	return true
}

func isAlnum(x rune) bool {
	return isDigit(x) || isAlpha(x)
}

func bool2str(x bool) string {
	if x {
		return `true`
	}

	return `false`
}

/*
IsDescriptor scans the input string val and judges whether it
qualifies as an RFC 4512 "descr", in that all of the following
evaluate as true:

  - Non-zero in length
  - Begins with an alphabetical character
  - Ends in an alphanumeric character
  - Contains only alphanumeric characters or hyphens
  - No contiguous hyphens

This function is an alternative to engaging the [antlr4512]
parsing subsystem.
*/
func IsDescriptor(val string) bool {
	return isDescriptor(val)
}

func isDescriptor(val string) bool {
	if len(val) == 0 {
		return false
	}

	// must begin with an alpha.
	if !isAlpha(rune(val[0])) {
		return false
	}

	// can only end in alnum.
	if !isAlnum(rune(val[len(val)-1])) {
		return false
	}

	// watch hyphens to avoid contiguous use
	var lastHyphen bool

	// iterate all characters in val, checking
	// each one for "descr" validity.
	for i := 0; i < len(val); i++ {
		ch := rune(val[i])
		switch {
		case isAlnum(ch):
			lastHyphen = false
		case ch == '-':
			if lastHyphen {
				// cannot use consecutive hyphens
				return false
			}
			lastHyphen = true
		default:
			// invalid character (none of [a-zA-Z0-9\-])
			return false
		}
	}

	return true
}

/*
mapTransferExtensions returns the provided dest instance of DefinitionMap,
following an attempt to copy all extensions found within src into dest.

This is mainly used to keep cyclomatics low during presentation and marshaling
procedures and may be used for any Definition qualifier.

The dest input value must be initialized else go will panic.
*/
func mapTransferExtensions(src Definition, dest DefinitionMap) DefinitionMap {
	exts := src.Extensions()
	for _, k := range exts.Keys() {
		if ext, found := exts.get(k); found {
			dest[k] = ext.List()
		}
	}

	return dest
}

/*
condenseWHSP returns input string b with all contiguous
WHSP characters condensed into single space characters.

WHSP is qualified through space or TAB chars (ASCII #32
and #9 respectively).
*/
func condenseWHSP(b string) (a string) {
	// remove leading and trailing
	// WHSP characters ...
	b = trimS(b)
	b = repAll(b, string(rune(10)), string(rune(32)))

	var last bool
	for i := 0; i < len(b); i++ {
		c := rune(b[i])
		switch c {
		// match space (32) or tab (9)
		case rune(9), rune(10), rune(32):
			if !last {
				last = true
				a += string(rune(32))
			}
		default:
			if last {
				last = false
			}
			a += string(c)
		}
	}

	a = trimS(a) //once more
	return
}

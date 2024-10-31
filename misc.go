package schemax

/*
IsNumericOID returns a Boolean value indicative of the outcome of an
attempt to parse input value id in the context of an unencoded ASN.1
OBJECT IDENTIFIER in dot form, e.g.:

	1.3.6.1.4.1.56521.999

The qualifications are as follows:

  - Must only consist of digits (arcs) and dot (ASCII FULL STOP) delimiters
  - Must begin with a root arc: 0, 1 or 2
  - Second-level arcs below root arcs 0 and 1 cannot be greater than 39
  - Cannot end with a dot
  - Dots cannot be contiguous
  - Though arcs are unbounded, no arc may ever be negative
  - OID must consist of at least two (2) arcs

Note: poached from JesseCoretta/go-objectid.
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

/*
governedDistinguishedName contains the components of a distinguished
name and the integer length of those components that are not distinct.

For example, a root suffix of "dc=example,dc=com" has two (2) components,
meaning its flattened integer length is one (1), as "dc=example" and
"dc=com" are not separate and distinct.

An easier, though less descriptive, explanation of the integer length
is simply "the number of comma characters at the far right (root) of
the DN which do NOT describe separate entries.  Again, the comma in the
"dc=example,dc=com" suffix equals a length of one (1).
*/
type governedDistinguishedName struct {
	components [][][]string
	flat       int
	length     int
}

func (r *governedDistinguishedName) isZero() bool {
	if r != nil {
		return len(r.components) == 0
	}

	return true
}

/*
tokenizeDN will attempt to tokenize the input dn value.

Through the act of tokenization, the following occurs:

An LDAP distinguished name, such as "uid=jesse+gidNumber=5042,ou=People,dc=example,dc=com,

... is returned as:

	[][][]string{
	  [][]string{
	    []string{`uid`,`jesse`},
	    []string{`gidNumber`,`5042`},
	  },
	  [][]string{
	    []string{`ou`,`People`},
	  },
	  [][]string{
	    []string{`dc`,`example`},
	  },
	  [][]string{
	    []string{`dc`,`com`},
	  },
	}

Please note that this function is NOT considered a true parser. If actual
parsing of component attribute values within a given DN is either desired
or required, consider use of a proper parser such as [go-ldap/v3's ParseDN]
function.

flat is an integer value that describes the flattened root suffix "length".
For instance, given the root suffix of "dc=example,dc=com" -- which is a
single entry and not two separate entries -- the input value should be the
integer 1.

[go-ldap/v3's ParseDN]: https://pkg.go.dev/github.com/go-ldap/ldap/v3#ParseDN
*/
func tokenizeDN(d string, flat ...int) (x *governedDistinguishedName) {
	if len(d) == 0 {
		return
	}

	x = &governedDistinguishedName{
		components: make([][][]string, 0),
	}

	if len(flat) > 0 {
		x.flat = flat[0]
	}

	rdns := splitUnescaped(d, `,`, `\`)
	lr := len(rdns)

	if lr == x.flat || x.flat < 0 {
		// bogus depth
		return
	}

	for i := 0; i < lr; i++ {
		var atvs [][]string = make([][]string, 0)
		srdns := splitUnescaped(rdns[i], `+`, `\`)
		for j := 0; j < len(srdns); j++ {
			if atv := split(srdns[j], `=`); len(atv) == 2 {
				atvs = append(atvs, atv)
			} else {
				atvs = append(atvs, []string{})
			}
		}

		x.components = append(x.components, atvs)
	}

	if x.flat > 0 {
		e := lr - 1
		f := e - x.flat
		x.components[f] = append(x.components[f], x.components[e]...)
		x.components = x.components[:e]
	}
	x.length = len(x.components)

	return
}

func detokenizeDN(x *governedDistinguishedName) (dtkz string) {
	if x.isZero() {
		return
	}

	var rdns []string
	for i := 0; i < x.length; i++ {
		rdn := x.components[i]
		char := `,`
		if i < x.length-x.flat && len(rdn) > 1 {
			char = `+`
		}

		//var r []string
		var atv []string
		for j := 0; j < len(rdn); j++ {
			atv = append(atv, rdn[j][0]+`=`+rdn[j][1])
		}
		rdns = append(rdns, join(atv, char))
	}

	dtkz = join(rdns, `,`)
	return
}

func splitUnescaped(str, sep, esc string) (slice []string) {
	slice = split(str, sep)
	for i := len(slice) - 2; i >= 0; i-- {
		if hasSfx(slice[i], esc) {
			slice[i] = slice[i][:len(slice[i])-len(esc)] + sep + slice[i+1]
			slice = append(slice[:i+1], slice[i+2:]...)
		}
	}

	return
}

/*
strInSlice returns a Boolean value indicative of whether the
specified string (str) is present within slice. Please note
that case is a significant element in the matching process.
*/
func strInSlice(str string, slice []string) bool {
	for i := 0; i < len(slice); i++ {
		if str == slice[i] {
			return true
		}
	}
	return false
}

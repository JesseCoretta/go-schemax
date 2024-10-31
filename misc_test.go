package schemax

import (
	"fmt"
	"testing"
)

func ExampleIsNumericOID() {
	fmt.Printf("%t", IsNumericOID(`1.3.6.1.4.1`))
	// Output: true
}

func ExampleIsDescriptor() {
	fmt.Printf("%t", IsDescriptor(`telephoneNumber`))
	// Output: true
}

func TestMisc_codecov(t *testing.T) {
	_ = IsDescriptor(`test_`)
	_ = IsDescriptor(`tes--t`)
	_ = IsDescriptor(`te?st`)
	_ = IsDescriptor(``)
	_ = IsDescriptor(`_`)
	_ = IsNumericOID(`1.3.6.1.4.1`)
	_ = IsNumericOID(`^73`)
	_ = IsNumericOID(`l`)
	_ = condenseWHSP(`Lorem     ipsum 
dolor sit amet, 
           consectetur adipiscing elit,
  sed do eiusmod tempor

incididunt ut labore et dolore magna aliqua.`)
}

func TestGovernedDistinguishedName(t *testing.T) {
	for _, strukt := range []struct {
		DN string
		L  int
	}{
		{`dc=example,dc=com`, 0},
		{`dc=example,dc=com`, 1},
		{`ou=People,dc=example,dc=com`, 1},
		{`ou=People+ou=Employees,dc=example,dc=com`, 1},
		{`o=example`, 0},
	} {
		if cdn := tokenizeDN(strukt.DN, strukt.L); cdn.isZero() {
			t.Errorf("%s failed: no content", t.Name())
			return
		} else {
			if dtkx := detokenizeDN(cdn); dtkx != strukt.DN {
				t.Errorf("%s failed: want %s, got %s [raw:%#v]",
					t.Name(), strukt.DN, dtkx, cdn)
				return
			}
		}
	}
}

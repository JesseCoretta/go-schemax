package rfc2377

type NameFormDefinitions []NameFormDefinition
type NameFormDefinition string

func (r NameFormDefinitions) Len() int {
	return len(r)
}

var AllNameForms NameFormDefinitions
var (
	DomainNameForm,
	DCOrgNameForm,
	DCOrgUnitNameForm,
	DCLNameForm,
	UIDOrgPersonNameForm NameFormDefinition
)

func (r NameFormDefinition) String() string {
	return string(r)
}

func init() {

	// Note: these definitions include corrections suggested
	// in Errata IDs 8114 through 8118 (reporter: me)
	// See: https://www.rfc-editor.org/errata_search.php?rfc=2377&rec_status=15&submitter_name=Jesse+Coretta&presentation=records
	DomainNameForm = NameFormDefinition(`( 1.3.6.1.1.2.1 NAME 'domainNameForm' OC domain MUST dc X-ORIGIN 'RFC2377' )`)
	DCOrgNameForm = NameFormDefinition(`( 1.3.6.1.1.2.2 NAME 'dcOrganizationNameForm' OC organization MUST dc X-ORIGIN 'RFC2377' )`)
	DCOrgUnitNameForm = NameFormDefinition(`( 1.3.6.1.1.2.3 NAME 'dcOrganizationalUnitNameForm' OC organizationalUnit MUST dc X-ORIGIN 'RFC2377' )`)
	DCLNameForm = NameFormDefinition(`( 1.3.6.1.1.2.4 NAME 'dcLocalityNameForm' OC locality MUST dc X-ORIGIN 'RFC2377' )`)
	UIDOrgPersonNameForm = NameFormDefinition(`( 1.3.6.1.1.2.5 NAME 'uidOrganizationalPersonNameForm' OC organizationalPerson MUST uid X-ORIGIN 'RFC2377' )`)

	AllNameForms = NameFormDefinitions{
		DomainNameForm,
		DCOrgNameForm,
		DCOrgUnitNameForm,
		DCLNameForm,
		UIDOrgPersonNameForm,
	}

}

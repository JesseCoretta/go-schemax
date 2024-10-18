package x501

type NameFormDefinitions []NameFormDefinition
type NameFormDefinition string

func (r NameFormDefinitions) Len() int {
	return len(r)
}

var AllNameForms NameFormDefinitions
var SubentryNameForm NameFormDefinition

func (r NameFormDefinition) String() string {
	return string(r)
}

func init() {

	SubentryNameForm = NameFormDefinition(`( 2.5.15.16 NAME 'subentryNameForm' OC subentry MUST cn X-ORIGIN 'X.501' )`)
	AllNameForms = NameFormDefinitions{
		SubentryNameForm,
	}

}

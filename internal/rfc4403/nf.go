package rfc4403

type NameFormDefinitions []NameFormDefinition
type NameFormDefinition string

func (r NameFormDefinitions) Len() int {
	return len(r)
}

var AllNameForms NameFormDefinitions
var (
	UDDIBusinessEntityNameForm,
	UDDIContactNameForm,
	UDDIAddressNameForm,
	UDDIBusinessServiceNameForm,
	UDDIBindingTemplateNameForm,
	UDDITModelInstanceInfoNameForm,
	UDDITModelNameForm,
	UDDIPublisherAssertionNameForm,
	UDDIv3SubscriptionNameForm,
	UDDIv3EntityObituaryNameForm NameFormDefinition
)

func (r NameFormDefinition) String() string {
	return string(r)
}

func init() {

	UDDIBusinessEntityNameForm = NameFormDefinition(`( 1.3.6.1.1.10.15.1 NAME 'uddiBusinessEntityNameForm' OC uddiBusinessEntity MUST ( uddiBusinessKey ) X-ORIGIN 'RFC4403' )`)
	UDDIContactNameForm = NameFormDefinition(`( 1.3.6.1.1.10.15.2 NAME 'uddiContactNameForm' OC uddiContact MUST ( uddiUUID ) X-ORIGIN 'RFC4403' )`)
	UDDIAddressNameForm = NameFormDefinition(`( 1.3.6.1.1.10.15.3 NAME 'uddiAddressNameForm' OC uddiAddress MUST ( uddiUUID ) X-ORIGIN 'RFC4403' )`)
	UDDIBusinessServiceNameForm = NameFormDefinition(`( 1.3.6.1.1.10.15.4 NAME 'uddiBusinessServiceNameForm' OC uddiBusinessService MUST ( uddiServiceKey ) X-ORIGIN 'RFC4403' )`)
	UDDIBindingTemplateNameForm = NameFormDefinition(`( 1.3.6.1.1.10.15.5 NAME 'uddiBindingTemplateNameForm' OC uddiBindingTemplate MUST ( uddiBindingKey ) X-ORIGIN 'RFC4403' )`)
	UDDITModelInstanceInfoNameForm = NameFormDefinition(`( 1.3.6.1.1.10.15.6 NAME 'uddiTModelInstanceInfoNameForm' OC uddiTModelInstanceInfo MUST ( uddiTModelKey ) X-ORIGIN 'RFC4403' )`)
	UDDITModelNameForm = NameFormDefinition(`( 1.3.6.1.1.10.15.7 NAME 'uddiTModelNameForm' OC uddiTModel MUST ( uddiTModelKey ) X-ORIGIN 'RFC4403' )`)
	UDDIPublisherAssertionNameForm = NameFormDefinition(`( 1.3.6.1.1.10.15.8 NAME 'uddiPublisherAssertionNameForm' OC uddiPublisherAssertion MUST ( uddiUUID ) X-ORIGIN 'RFC4403' )`)
	UDDIv3SubscriptionNameForm = NameFormDefinition(`( 1.3.6.1.1.10.15.9 NAME 'uddiv3SubscriptionNameForm' OC uddiv3Subscription MUST ( uddiUUID ) X-ORIGIN 'RFC4403' )`)
	UDDIv3EntityObituaryNameForm = NameFormDefinition(`( 1.3.6.1.1.10.15.10 NAME 'uddiv3EntityObituaryNameForm' OC uddiv3EntityObituary MUST ( uddiUUID ) X-ORIGIN 'RFC4403' )`)

	AllNameForms = NameFormDefinitions{
		UDDIBusinessEntityNameForm,
		UDDIContactNameForm,
		UDDIAddressNameForm,
		UDDIBusinessServiceNameForm,
		UDDIBindingTemplateNameForm,
		UDDITModelInstanceInfoNameForm,
		UDDITModelNameForm,
		UDDIPublisherAssertionNameForm,
		UDDIv3SubscriptionNameForm,
		UDDIv3EntityObituaryNameForm,
	}

}

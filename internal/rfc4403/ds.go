package rfc4403

type DITStructureRuleDefinitions []DITStructureRuleDefinition
type DITStructureRuleDefinition string

func (r DITStructureRuleDefinitions) Len() int {
	return len(r)
}

var AllDITStructureRules DITStructureRuleDefinitions
var (
	UDDIBusinessEntityStructureRule,
	UDDIContactStructureRule,
	UDDIAddressStructureRule,
	UDDIBusinessServiceStructureRule,
	UDDIBindingTemplateStructureRule,
	UDDITModelInstanceInfoStructureRule,
	UDDITModelStructureRule,
	UDDIPublisherAssertionStructureRule,
	UDDIv3SubscriptionStructureRule,
	UDDIv3EntityObituaryStructureRule DITStructureRuleDefinition
)

func (r DITStructureRuleDefinition) String() string {
	return string(r)
}

func init() {

	UDDIBusinessEntityStructureRule = DITStructureRuleDefinition(`( 1 NAME 'uddiBusinessEntityStructureRule' FORM uddiBusinessEntityNameForm X-ORIGIN 'RFC4403' )`)
	UDDIContactStructureRule = DITStructureRuleDefinition(`( 2 NAME 'uddiContactStructureRule' FORM uddiContactNameForm SUP ( 1 ) X-ORIGIN 'RFC4403' )`)
	UDDIAddressStructureRule = DITStructureRuleDefinition(`( 3 NAME 'uddiAddressStructureRule' FORM uddiAddressNameForm SUP ( 2 ) X-ORIGIN 'RFC4403' )`)
	UDDIBusinessServiceStructureRule = DITStructureRuleDefinition(`( 4 NAME 'uddiBusinessServiceStructureRule' FORM uddiBusinessServiceNameForm SUP ( 1 ) X-ORIGIN 'RFC4403' )`)
	UDDIBindingTemplateStructureRule = DITStructureRuleDefinition(`( 5 NAME 'uddiBindingTemplateStructureRule' FORM uddiBindingTemplateNameForm SUP ( 4 ) X-ORIGIN 'RFC4403' )`)
	UDDITModelInstanceInfoStructureRule = DITStructureRuleDefinition(`( 6 NAME 'uddiTModelInstanceInfoStructureRule' FORM uddiTModelInstanceInfoNameForm SUP ( 5 ) X-ORIGIN 'RFC4403' )`)
	UDDITModelStructureRule = DITStructureRuleDefinition(`( 7 NAME 'uddiTModelStructureRule' FORM uddiTModelNameForm X-ORIGIN 'RFC4403' )`)
	UDDIPublisherAssertionStructureRule = DITStructureRuleDefinition(`( 8 NAME 'uddiPublisherAssertionStructureRule' FORM uddiPublisherAssertionNameForm X-ORIGIN 'RFC4403' )`)
	UDDIv3SubscriptionStructureRule = DITStructureRuleDefinition(`( 9 NAME 'uddiv3SubscriptionStructureRule' FORM uddiv3SubscriptionNameForm X-ORIGIN 'RFC4403' )`)
	UDDIv3EntityObituaryStructureRule = DITStructureRuleDefinition(`( 10 NAME 'uddiv3EntityObituaryStructureRule' FORM uddiv3EntityObituaryNameForm X-ORIGIN 'RFC4403' )`)

	AllDITStructureRules = DITStructureRuleDefinitions{
		UDDIBusinessEntityStructureRule,
		UDDIContactStructureRule,
		UDDIAddressStructureRule,
		UDDIBusinessServiceStructureRule,
		UDDIBindingTemplateStructureRule,
		UDDITModelInstanceInfoStructureRule,
		UDDITModelStructureRule,
		UDDIPublisherAssertionStructureRule,
		UDDIv3SubscriptionStructureRule,
		UDDIv3EntityObituaryStructureRule,
	}
}

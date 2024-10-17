package rfc4403

type ObjectClassDefinitions []ObjectClassDefinition
type ObjectClassDefinition string

func (r ObjectClassDefinitions) Len() int {
	return len(r)
}

var AllObjectClasses ObjectClassDefinitions
var (
	UDDIBusinessEntity,
	UDDIContact,
	UDDIAddress,
	UDDIBusinessService,
	UDDIBindingTemplate,
	UDDITModelInstanceInfo,
	UDDITModel,
	UDDIPublisherAssertion,
	UDDIv3Subscription,
	UDDIv3EntityObituary ObjectClassDefinition
)

func (r ObjectClassDefinition) String() string {
	return string(r)
}

func init() {

	UDDIBusinessEntity = ObjectClassDefinition(`( 1.3.6.1.1.10.6.1 NAME 'uddiBusinessEntity'  SUP top STRUCTURAL MUST ( uddiBusinessKey $ uddiName) MAY ( uddiAuthorizedName $ uddiOperator $ uddiDiscoveryURLs $ uddiDescription $ uddiIdentifierBag $ uddiCategoryBag $ uddiv3BusinessKey $ uddiv3DigitalSignature $ uddiv3EntityModificationTime $ uddiv3NodeId) X-ORIGIN 'RFC4403' )`)
	UDDIContact = ObjectClassDefinition(`( 1.3.6.1.1.10.6.2 NAME 'uddiContact'  SUP top STRUCTURAL MUST ( uddiPersonName $ uddiUUID ) MAY ( uddiUseType $ uddiDescription $ uddiPhone $ uddiEMail ) X-ORIGIN 'RFC4403' )`)
	UDDIAddress = ObjectClassDefinition(`( 1.3.6.1.1.10.6.3 NAME 'uddiAddress'  SUP top STRUCTURAL MUST ( uddiUUID ) MAY ( uddiUseType $ uddiSortCode $ uddiTModelKey $ uddiv3TmodelKey $ uddiAddressLine $ uddiLang) X-ORIGIN 'RFC4403' )`)
	UDDIBusinessService = ObjectClassDefinition(`( 1.3.6.1.1.10.6.4 NAME 'uddiBusinessService'  SUP top STRUCTURAL MUST ( uddiServiceKey ) MAY ( uddiName $ uddiBusinessKey $ uddiDescription $ uddiCategoryBag $ uddiIsProjection $ uddiv3ServiceKey $ uddiv3BusinessKey $ uddiv3DigitalSignature $ uddiv3EntityCreationTime $ uddiv3EntityModificationTime $ uddiv3NodeId) X-ORIGIN 'RFC4403' )`)
	UDDIBindingTemplate = ObjectClassDefinition(`( 1.3.6.1.1.10.6.5 NAME 'uddiBindingTemplate'  SUP top STRUCTURAL MUST ( uddiBindingKey ) MAY ( uddiServiceKey $ uddiDescription $ uddiAccessPoint $ uddiHostingRedirector $ uddiCategoryBag $ uddiv3BindingKey $ uddiv3ServiceKey $ uddiv3DigitalSignature $ uddiv3EntityCreationTime $ uddiv3NodeId) X-ORIGIN 'RFC4403' )`)
	UDDITModelInstanceInfo = ObjectClassDefinition(`( 1.3.6.1.1.10.6.6 NAME 'uddiTModelInstanceInfo'  SUP top STRUCTURAL MUST ( uddiTModelKey ) MAY ( uddiDescription $ uddiInstanceDescription $ uddiInstanceParms $ uddiOverviewDescription $ uddiOverviewURL $ uddiv3TmodelKey) X-ORIGIN 'RFC4403' )`)
	UDDITModel = ObjectClassDefinition(`( 1.3.6.1.1.10.6.7 NAME 'uddiTModel'  SUP top STRUCTURAL MUST ( uddiTModelKey $ uddiName ) MAY ( uddiAuthorizedName $ uddiOperator $ uddiDescription $ uddiOverviewDescription $ uddiOverviewURL $ uddiIdentifierBag $ uddiCategoryBag $ uddiIsHidden $ uddiv3TModelKey $ uddiv3DigitalSignature $ uddiv3NodeId) X-ORIGIN 'RFC4403' )`)
	UDDIPublisherAssertion = ObjectClassDefinition(`( 1.3.6.1.1.10.6.8 NAME 'uddiPublisherAssertion'  SUP top STRUCTURAL MUST ( uddiFromKey $ uddiToKey $ uddiKeyedReference $ uddiUUID ) MAY ( uddiv3DigitalSignature $ uddiv3NodeId) X-ORIGIN 'RFC4403' )`)
	UDDIv3Subscription = ObjectClassDefinition(`( 1.3.6.1.1.10.6.9 NAME 'uddiv3Subscription'  SUP top STRUCTURAL MUST ( uddiv3SubscriptionFilter $ uddiUUID) MAY ( uddiAuthorizedName $ uddiv3SubscriptionKey $ uddiv3BindingKey $ uddiv3NotificationInterval $ uddiv3MaxEntities $ uddiv3ExpiresAfter $ uddiv3BriefResponse $ uddiv3NodeId) X-ORIGIN 'RFC4403' )`)
	UDDIv3EntityObituary = ObjectClassDefinition(`( 1.3.6.1.1.10.6.10 NAME 'uddiv3EntityObituary'  SUP top STRUCTURAL MUST ( uddiv3EntityKey $ uddiUUID) MAY ( uddiAuthorizedName $ uddiv3EntityCreationTime $ uddiv3EntityDeletionTime $ uddiv3NodeId) X-ORIGIN 'RFC4403' )`)

	AllObjectClasses = ObjectClassDefinitions{
		UDDIBusinessEntity,
		UDDIContact,
		UDDIAddress,
		UDDIBusinessService,
		UDDIBindingTemplate,
		UDDITModelInstanceInfo,
		UDDITModel,
		UDDIPublisherAssertion,
		UDDIv3Subscription,
		UDDIv3EntityObituary,
	}

}

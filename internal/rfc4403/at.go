package rfc4403

type AttributeTypeDefinitions []AttributeTypeDefinition
type AttributeTypeDefinition string

func (r AttributeTypeDefinitions) Len() int {
	return len(r)
}

var AllAttributeTypes AttributeTypeDefinitions

var (
	UDDIBusinessKey,
	UDDIAuthorizedName,
	UDDIOperator,
	UDDIName,
	UDDIDescription,
	UDDIDiscoveryURLs,
	UDDIUseType,
	UDDIPersonName,
	UDDIPhone,
	UDDIEMail,
	UDDISortCode,
	UDDITModelKey,
	UDDIAddressLine,
	UDDIIdentifierBag,
	UDDICategoryBag,
	UDDIKeyedReference,
	UDDIServiceKey,
	UDDIBindingKey,
	UDDIAccessPoint,
	UDDIHostingRedirector,
	UDDIInstanceDescription,
	UDDIInstanceParms,
	UDDIOverviewDescription,
	UDDIOverviewURL,
	UDDIFromKey,
	UDDIToKey,
	UDDIUUID,
	UDDIIsHidden,
	UDDIIsProjection,
	UDDILang,
	UDDIv3BusinessKey,
	UDDIv3ServiceKey,
	UDDIv3BindingKey,
	UDDIv3TModelKey,
	UDDIv3DigitalSignature,
	UDDIv3NodeId,
	UDDIv3EntityModificationTime,
	UDDIv3SubscriptionKey,
	UDDIv3SubscriptionFilter,
	UDDIv3NotificationInterval,
	UDDIv3MaxEntities,
	UDDIv3ExpiresAfter,
	UDDIv3BriefResponse,
	UDDIv3EntityKey,
	UDDIv3EntityCreationTime,
	UDDIv3EntityDeletionTime AttributeTypeDefinition
)

func (r AttributeTypeDefinition) String() string {
	return string(r)
}

func init() {

	UDDIBusinessKey = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.1 NAME 'uddiBusinessKey' DESC 'businessEntity unique identifier' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIAuthorizedName = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.2 NAME 'uddiAuthorizedName' DESC 'businessEntity publisher name' EQUALITY distinguishedNameMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.12 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIOperator = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.3 NAME 'uddiOperator' DESC 'registry site operator of businessEntitys master copy' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIName = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.4 NAME 'uddiName' DESC 'human readable name' EQUALITY caseIgnoreMatch ORDERING caseIgnoreOrderingMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4403' )`)
	UDDIDescription = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.5 NAME 'uddiDescription' DESC 'short description' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4403' )`)
	UDDIDiscoveryURLs = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.6 NAME 'uddiDiscoveryURLs' DESC 'URL to retrieve a businessEntity instance' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4403' )`)
	UDDIUseType = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.7 NAME 'uddiUseType' DESC 'name of convention the referenced document follows' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIPersonName = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.8 NAME 'uddiPersonName' DESC 'name of person or job role available for contact' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIPhone = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.9 NAME 'uddiPhone' DESC 'telephone number for contact' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4403' )`)
	UDDIEMail = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.10 NAME 'uddiEMail' DESC 'e-mail address for contact' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4403' )`)
	UDDISortCode = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.11 NAME 'uddiSortCode' DESC 'specifies an external display mechanism' EQUALITY caseIgnoreMatch ORDERING caseIgnoreOrderingMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDITModelKey = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.12 NAME 'uddiTModelKey' DESC 'tModel unique identifier' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIAddressLine = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.13 NAME 'uddiAddressLine' DESC 'address' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4403' )`)
	UDDIIdentifierBag = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.14 NAME 'uddiIdentifierBag' DESC 'identification information' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4403' )`)
	UDDICategoryBag = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.15 NAME 'uddiCategoryBag' DESC 'categorization information' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4403' )`)
	UDDIKeyedReference = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.16 NAME 'uddiKeyedReference' DESC 'categorization information' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4403' )`)
	UDDIServiceKey = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.17 NAME 'uddiServiceKey' DESC 'businessService unique identifier' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIBindingKey = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.18 NAME 'uddiBindingKey' DESC 'bindingTemplate unique identifier' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIAccessPoint = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.19 NAME 'uddiAccessPoint' DESC 'entry point address to call a web service' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIHostingRedirector = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.20 NAME 'uddiHostingRedirector' DESC 'designates a pointer to another bindingTemplate' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIInstanceDescription = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.21 NAME 'uddiInstanceDescription' DESC 'instance details description' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4403' )`)
	UDDIInstanceParms = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.22 NAME 'uddiInstanceParms' DESC 'URL reference to required settings' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIOverviewDescription = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.23 NAME 'uddiOverviewDescription' DESC 'outlines tModel usage' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4403' )`)
	UDDIOverviewURL = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.24 NAME 'uddiOverviewURL' DESC 'URL reference to overview document' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIFromKey = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.25 NAME 'uddiFromKey' DESC 'unique businessEntity key reference' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIToKey = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.26 NAME 'uddiToKey' DESC 'unique businessEntity key reference' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIUUID = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.27 NAME 'uddiUUID' DESC 'unique attribute' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIIsHidden = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.28 NAME 'uddiIsHidden' DESC 'isHidden attribute' EQUALITY booleanMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.7 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIIsProjection = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.29 NAME 'uddiIsProjection' DESC 'isServiceProjection attribute' EQUALITY booleanMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.7 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDILang = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.30 NAME 'uddiLang' DESC 'xml:lang value in v3 Address structure' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3BusinessKey = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.31 NAME 'uddiv3BusinessKey' DESC 'UDDIv3 businessEntity unique identifier' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3ServiceKey = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.32 NAME 'uddiv3ServiceKey' DESC 'UDDIv3 businessService unique identifier' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3BindingKey = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.33 NAME 'uddiv3BindingKey' DESC 'UDDIv3 BindingTemplate unique identifier' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3TModelKey = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.34 NAME 'uddiv3TModelKey' DESC 'UDDIv3 TModel unique identifier' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3DigitalSignature = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.35 NAME 'uddiv3DigitalSignature' DESC 'UDDIv3 entity digital signature' EQUALITY caseExactMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4403' )`)
	UDDIv3NodeId = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.36 NAME 'uddiv3NodeId' DESC 'UDDIv3 Node Identifier' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3EntityModificationTime = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.37 NAME 'uddiv3EntityModificationTime' DESC 'UDDIv3 Last Modified Time for Entity' EQUALITY generalizedTimeMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.24 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3SubscriptionKey = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.38 NAME 'uddiv3SubscriptionKey' DESC 'UDDIv3 Subscription unique identifier' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3SubscriptionFilter = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.39 NAME 'uddiv3SubscriptionFilter' DESC 'UDDIv3 Subscription Filter' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3NotificationInterval = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.40 NAME 'uddiv3NotificationInterval' DESC 'UDDIv3 Notification Interval' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3MaxEntities = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.41 NAME 'uddiv3MaxEntities' DESC 'UDDIv3 Subscription maxEntities field' EQUALITY integerMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.27 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3ExpiresAfter = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.42 NAME 'uddiv3ExpiresAfter' DESC 'UDDIv3 Subscription ExpiresAfter field' EQUALITY generalizedTimeMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.24 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3BriefResponse = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.43 NAME 'uddiv3BriefResponse' DESC 'UDDIv3 Subscription ExpiresAfter field' EQUALITY booleanMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.7 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3EntityKey = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.44 NAME 'uddiv3EntityKey' DESC 'UDDIv3 Entity unique identifier' EQUALITY caseIgnoreMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3EntityCreationTime = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.45 NAME 'uddiv3EntityCreationTime' DESC 'UDDIv3 Entity Creation Time' EQUALITY generalizedTimeMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.24 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)
	UDDIv3EntityDeletionTime = AttributeTypeDefinition(`( 1.3.6.1.1.10.4.46 NAME 'uddiv3EntityDeletionTime' DESC 'UDDIv3 Entity Deletion Time' EQUALITY generalizedTimeMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.24 SINGLE-VALUE X-ORIGIN 'RFC4403' )`)

	AllAttributeTypes = AttributeTypeDefinitions{
		UDDIBusinessKey,
		UDDIAuthorizedName,
		UDDIOperator,
		UDDIName,
		UDDIDescription,
		UDDIDiscoveryURLs,
		UDDIUseType,
		UDDIPersonName,
		UDDIPhone,
		UDDIEMail,
		UDDISortCode,
		UDDITModelKey,
		UDDIAddressLine,
		UDDIIdentifierBag,
		UDDICategoryBag,
		UDDIKeyedReference,
		UDDIServiceKey,
		UDDIBindingKey,
		UDDIAccessPoint,
		UDDIHostingRedirector,
		UDDIInstanceDescription,
		UDDIInstanceParms,
		UDDIOverviewDescription,
		UDDIOverviewURL,
		UDDIFromKey,
		UDDIToKey,
		UDDIUUID,
		UDDIIsHidden,
		UDDIIsProjection,
		UDDILang,
		UDDIv3BusinessKey,
		UDDIv3ServiceKey,
		UDDIv3BindingKey,
		UDDIv3TModelKey,
		UDDIv3DigitalSignature,
		UDDIv3NodeId,
		UDDIv3EntityModificationTime,
		UDDIv3SubscriptionKey,
		UDDIv3SubscriptionFilter,
		UDDIv3NotificationInterval,
		UDDIv3MaxEntities,
		UDDIv3ExpiresAfter,
		UDDIv3BriefResponse,
		UDDIv3EntityKey,
		UDDIv3EntityCreationTime,
		UDDIv3EntityDeletionTime,
	}
}

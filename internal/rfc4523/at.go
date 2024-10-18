package rfc4523

type AttributeTypeDefinitions []AttributeTypeDefinition
type AttributeTypeDefinition string

func (r AttributeTypeDefinitions) Len() int {
	return len(r)
}

var (
	AllAttributeTypes AttributeTypeDefinitions
)

var (
	UserCertificate           AttributeTypeDefinition
	CACertificate             AttributeTypeDefinition
	CrossCertificatePair      AttributeTypeDefinition
	CertificateRevocationList AttributeTypeDefinition
	AuthorityRevocationList   AttributeTypeDefinition
	DeltaRevocationList       AttributeTypeDefinition
	SupportedAlgorithms       AttributeTypeDefinition
)

func (r AttributeTypeDefinition) String() string {
	return string(r)
}

func init() {

	UserCertificate = AttributeTypeDefinition(`( 2.5.4.36 NAME 'userCertificate' DESC 'X.509 user certificate' EQUALITY certificateExactMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.8 X-ORIGIN 'RFC4523' )`)
	CACertificate = AttributeTypeDefinition(`( 2.5.4.37 NAME 'cACertificate' DESC 'X.509 CA certificate' EQUALITY certificateExactMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.8 X-ORIGIN 'RFC4523' )`)
	CrossCertificatePair = AttributeTypeDefinition(`( 2.5.4.40 NAME 'crossCertificatePair' DESC 'X.509 cross certificate pair' EQUALITY certificatePairExactMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.10 X-ORIGIN 'RFC4523' )`)
	CertificateRevocationList = AttributeTypeDefinition(`( 2.5.4.39 NAME 'certificateRevocationList' DESC 'X.509 certificate revocation list' EQUALITY certificateListExactMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.9 X-ORIGIN 'RFC4523' )`)
	AuthorityRevocationList = AttributeTypeDefinition(`( 2.5.4.38 NAME 'authorityRevocationList' DESC 'X.509 authority revocation list' EQUALITY certificateListExactMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.9 X-ORIGIN 'RFC4523' )`)
	DeltaRevocationList = AttributeTypeDefinition(`( 2.5.4.53 NAME 'deltaRevocationList' DESC 'X.509 delta revocation list' EQUALITY certificateListExactMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.9 X-ORIGIN 'RFC4523' )`)
	SupportedAlgorithms = AttributeTypeDefinition(`( 2.5.4.52 NAME 'supportedAlgorithms' DESC 'X.509 supported algorithms' EQUALITY algorithmIdentifierMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.49 X-ORIGIN 'RFC4523' )`)

	AllAttributeTypes = AttributeTypeDefinitions{
		UserCertificate,
		CACertificate,
		CrossCertificatePair,
		CertificateRevocationList,
		AuthorityRevocationList,
		DeltaRevocationList,
		SupportedAlgorithms,
	}
}

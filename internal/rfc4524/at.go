package rfc4524

type AttributeTypeDefinitions []AttributeTypeDefinition
type AttributeTypeDefinition string

func (r AttributeTypeDefinitions) Len() int {
	return len(r)
}

var (
	AllAttributeTypes AttributeTypeDefinitions
)

var (
	AssociatedDomain,
	AssociatedName,
	BuildingName,
	CO,
	DocumentAuthor,
	DocumentIdentifier,
	DocumentLocation,
	DocumentPublisher,
	DocumentTitle,
	DocumentVersion,
	Drink,
	HomePhone,
	HomePostalAddress,
	Host,
	Info,
	Mail,
	Manager,
	Mobile,
	OrganizationalStatus,
	Pager,
	PersonalTitle,
	RoomNumber,
	Secretary,
	UniqueIdentifier,
	UserClass,
	SingleLevelQuality AttributeTypeDefinition
)

func (r AttributeTypeDefinition) String() string {
	return string(r)
}

func init() {
	AssociatedDomain = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.37 NAME 'associatedDomain' EQUALITY caseIgnoreIA5Match SUBSTR caseIgnoreIA5SubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.26 X-ORIGIN 'RFC4524' )`)
	AssociatedName = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.38 NAME 'associatedName' EQUALITY distinguishedNameMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.12 X-ORIGIN 'RFC4524' )`)
	BuildingName = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.48 NAME 'buildingName' EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{256} X-ORIGIN 'RFC4524' )`)
	CO = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.43 NAME ( 'co' 'friendlyCountryName' ) EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4524' )`)
	DocumentAuthor = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.14 NAME 'documentAuthor' EQUALITY distinguishedNameMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.12 X-ORIGIN 'RFC4524' )`)
	DocumentIdentifier = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.11 NAME 'documentIdentifier' EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{256} X-ORIGIN 'RFC4524' )`)
	DocumentLocation = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.15 NAME 'documentLocation' EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{256} X-ORIGIN 'RFC4524' )`)
	DocumentPublisher = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.56 NAME 'documentPublisher' EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15 X-ORIGIN 'RFC4524' )`)
	DocumentTitle = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.12 NAME 'documentTitle' EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{256} X-ORIGIN 'RFC4524' )`)
	DocumentVersion = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.13 NAME 'documentVersion' EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{256} X-ORIGIN 'RFC4524' )`)
	Drink = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.5 NAME ( 'drink' 'favouriteDrink' ) EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{256} X-ORIGIN 'RFC4524' )`)
	HomePhone = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.20 NAME ( 'homePhone' 'homeTelephoneNumber' ) EQUALITY telephoneNumberMatch SUBSTR telephoneNumberSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.50 X-ORIGIN 'RFC4524' )`)
	HomePostalAddress = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.39 NAME 'homePostalAddress' EQUALITY caseIgnoreListMatch SUBSTR caseIgnoreListSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.41 X-ORIGIN 'RFC4524' )`)
	Host = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.9 NAME 'host' EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{256} X-ORIGIN 'RFC4524' )`)
	Info = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.4 NAME 'info' EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{2048} X-ORIGIN 'RFC4524' )`)
	Mail = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.3 NAME ( 'mail' 'rfc822Mailbox' ) EQUALITY caseIgnoreIA5Match SUBSTR caseIgnoreIA5SubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.26{256} X-ORIGIN 'RFC4524' )`)
	Manager = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.10 NAME 'manager' EQUALITY distinguishedNameMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.12 X-ORIGIN 'RFC4524' )`)
	Mobile = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.41 NAME ( 'mobile' 'mobileTelephoneNumber' ) EQUALITY telephoneNumberMatch SUBSTR telephoneNumberSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.50 X-ORIGIN 'RFC4524' )`)
	OrganizationalStatus = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.45 NAME 'organizationalStatus' EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{256} X-ORIGIN 'RFC4524' )`)
	Pager = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.42 NAME ( 'pager' 'pagerTelephoneNumber' ) EQUALITY telephoneNumberMatch SUBSTR telephoneNumberSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.50 X-ORIGIN 'RFC4524' )`)
	PersonalTitle = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.40 NAME 'personalTitle' EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{256} X-ORIGIN 'RFC4524' )`)
	RoomNumber = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.6 NAME 'roomNumber' EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{256} X-ORIGIN 'RFC4524' )`)
	Secretary = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.21 NAME 'secretary' EQUALITY distinguishedNameMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.12 X-ORIGIN 'RFC4524' )`)
	SingleLevelQuality = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.50 NAME 'singleLevelQuality' DESC 'Single Level Quality' SYNTAX 1.3.6.1.4.1.1466.115.121.1.13 SINGLE-VALUE X-ORIGIN 'RFC4524' )`)
	UniqueIdentifier = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.44 NAME 'uniqueIdentifier' EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{256} X-ORIGIN 'RFC4524' )`)
	UserClass = AttributeTypeDefinition(`( 0.9.2342.19200300.100.1.8 NAME 'userClass' EQUALITY caseIgnoreMatch SUBSTR caseIgnoreSubstringsMatch SYNTAX 1.3.6.1.4.1.1466.115.121.1.15{256} X-ORIGIN 'RFC4524' )`)

	AllAttributeTypes = AttributeTypeDefinitions{
		AssociatedDomain,
		AssociatedName,
		BuildingName,
		CO,
		DocumentAuthor,
		DocumentIdentifier,
		DocumentLocation,
		DocumentPublisher,
		DocumentTitle,
		DocumentVersion,
		Drink,
		HomePhone,
		HomePostalAddress,
		Host,
		Info,
		Mail,
		Manager,
		Mobile,
		OrganizationalStatus,
		Pager,
		PersonalTitle,
		RoomNumber,
		Secretary,
		SingleLevelQuality,
		UniqueIdentifier,
		UserClass,
	}
}

package rfc2307

type ObjectClassDefinitions []ObjectClassDefinition
type ObjectClassDefinition string

func (r ObjectClassDefinitions) Len() int {
	return len(r)
}

var (
	AllObjectClasses ObjectClassDefinitions
)

var (
	POSIXAccount   ObjectClassDefinition
	ShadowAccount  ObjectClassDefinition
	POSIXGroup     ObjectClassDefinition
	IPService      ObjectClassDefinition
	IPProtocol     ObjectClassDefinition
	ONCRPC         ObjectClassDefinition
	IPHost         ObjectClassDefinition
	IPNetwork      ObjectClassDefinition
	NISNetgroup    ObjectClassDefinition
	NISMap         ObjectClassDefinition
	NISObject      ObjectClassDefinition
	IEEE802Device  ObjectClassDefinition
	BootableDevice ObjectClassDefinition
)

func init() {

	POSIXAccount = ObjectClassDefinition(`( nisSchema.2.0 NAME 'posixAccount' SUP top AUXILIARY DESC 'Abstraction of an account with POSIX attributes' MUST ( cn $ uid $ uidNumber $ gidNumber $ homeDirectory ) MAY ( userPassword $ loginShell $ gecos $ description ) X-ORIGIN 'RFC2307' )`)
	ShadowAccount = ObjectClassDefinition(`( nisSchema.2.1 NAME 'shadowAccount' SUP top AUXILIARY DESC 'Additional attributes for shadow passwords' MUST uid MAY ( userPassword $ shadowLastChange $ shadowMin $ shadowMax $ shadowWarning $ shadowInactive $ shadowExpire $ shadowFlag $ description ) X-ORIGIN 'RFC2307' )`)
	POSIXGroup = ObjectClassDefinition(`( nisSchema.2.2 NAME 'posixGroup' SUP top STRUCTURAL DESC 'Abstraction of a group of accounts' MUST ( cn $ gidNumber ) MAY ( userPassword $ memberUid $ description ) X-ORIGIN 'RFC2307' )`)
	IPService = ObjectClassDefinition(`( nisSchema.2.3 NAME 'ipService' SUP top STRUCTURAL DESC 'Abstraction an Internet Protocol service. Maps an IP port and protocol (such as tcp or udp) to one or more names; the distinguished value of the cn attribute denotes the service\'s canonical name' MUST ( cn $ ipServicePort $ ipServiceProtocol ) MAY ( description ) X-ORIGIN 'RFC2307' )`)
	IPProtocol = ObjectClassDefinition(`( nisSchema.2.4 NAME 'ipProtocol' SUP top STRUCTURAL DESC 'Abstraction of an IP protocol. Maps a protocol number to one or more names. The distinguished value of the cn attribute denotes the protocol\'s canonical name' MUST ( cn $ ipProtocolNumber $ description ) MAY description X-ORIGIN 'RFC2307' )`)
	ONCRPC = ObjectClassDefinition(`( nisSchema.2.5 NAME 'oncRpc' SUP top STRUCTURAL DESC 'Abstraction of an Open Network Computing (ONC) [RFC1057] Remote Procedure Call (RPC) binding. This class maps an ONC RPC number to a name. The distinguished value of the cn attribute denotes the RPC service\'s canonical name' MUST ( cn $ oncRpcNumber $ description ) MAY description X-ORIGIN 'RFC2307' )`)
	IPHost = ObjectClassDefinition(`( nisSchema.2.6 NAME 'ipHost' SUP top AUXILIARY DESC 'Abstraction of a host, an IP device. The distinguished value of the cn attribute denotes the hosts canonical name. Device SHOULD be used as a structural class' MUST ( cn $ ipHostNumber ) MAY ( l $ description $ manager ) X-ORIGIN 'RFC2307' )`)
	IPNetwork = ObjectClassDefinition(`( nisSchema.2.7 NAME 'ipNetwork' SUP top STRUCTURAL DESC 'Abstraction of a network. The distinguished value of the cn attribute denotes the networks canonical name' MUST ( cn $ ipNetworkNumber ) MAY ( ipNetmaskNumber $ l $ description $ manager ) X-ORIGIN 'RFC2307' )`)
	NISNetgroup = ObjectClassDefinition(`( nisSchema.2.8 NAME 'nisNetgroup' SUP top STRUCTURAL DESC 'Abstraction of a netgroup. May refer to other netgroups' MUST cn MAY ( nisNetgroupTriple $ memberNisNetgroup $ description ) X-ORIGIN 'RFC2307' )`)
	NISMap = ObjectClassDefinition(`( nisSchema.2.9 NAME 'nisMap' SUP top STRUCTURAL DESC 'A generic abstraction of a NIS map' MUST nisMapName MAY description X-ORIGIN 'RFC2307' )`)
	NISObject = ObjectClassDefinition(`( nisSchema.2.10 NAME 'nisObject' SUP top STRUCTURAL DESC 'An entry in a NIS map' MUST ( cn $ nisMapEntry $ nisMapName ) MAY description X-ORIGIN 'RFC2307' )`)
	IEEE802Device = ObjectClassDefinition(`( nisSchema.2.11 NAME 'ieee802Device' SUP top AUXILIARY DESC 'A device with a MAC address; device SHOULD be used as a structural class' MAY macAddress X-ORIGIN 'RFC2307' )`)
	BootableDevice = ObjectClassDefinition(`( nisSchema.2.12 NAME 'bootableDevice' DESC 'A device with boot parameters; device SHOULD be used as a structural class' SUP top AUXILIARY MAY ( bootFile $ bootParameter ) X-ORIGIN 'RFC2307' )`)

	AllObjectClasses = ObjectClassDefinitions{
		POSIXAccount,
		ShadowAccount,
		POSIXGroup,
		IPService,
		IPProtocol,
		ONCRPC,
		IPHost,
		IPNetwork,
		NISNetgroup,
		NISMap,
		NISObject,
		IEEE802Device,
		BootableDevice,
	}

}

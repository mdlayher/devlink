// WARNING: This file has automatically been generated on Thu, 31 May 2018 16:53:43 EDT.
// By https://git.io/c-for-go. DO NOT EDIT.

package dlh

const (
	// GenlName as defined in dlh/devlink.h:16
	GenlName = "devlink"
	// GenlVersion as defined in dlh/devlink.h:17
	GenlVersion = 0x1
	// GenlMcgrpConfigName as defined in dlh/devlink.h:18
	GenlMcgrpConfigName = "config"
	// CmdEswitchModeGet as defined in dlh/devlink.h:62
	CmdEswitchModeGet = CmdEswitchGet
	// CmdEswitchModeSet as defined in dlh/devlink.h:66
	CmdEswitchModeSet = CmdEswitchSet
	// SbThresholdToAlphaMax as defined in dlh/devlink.h:116
	SbThresholdToAlphaMax = 20
)

// devlinkCommand as declared in dlh/devlink.h:20
type devlinkCommand int32

// devlinkCommand enumeration from dlh/devlink.h:20
const (
	CmdUnspec                = iota
	CmdGet                   = 1
	CmdSet                   = 2
	CmdNew                   = 3
	CmdDel                   = 4
	CmdPortGet               = 5
	CmdPortSet               = 6
	CmdPortNew               = 7
	CmdPortDel               = 8
	CmdPortSplit             = 9
	CmdPortUnsplit           = 10
	CmdSbGet                 = 11
	CmdSbSet                 = 12
	CmdSbNew                 = 13
	CmdSbDel                 = 14
	CmdSbPoolGet             = 15
	CmdSbPoolSet             = 16
	CmdSbPoolNew             = 17
	CmdSbPoolDel             = 18
	CmdSbPortPoolGet         = 19
	CmdSbPortPoolSet         = 20
	CmdSbPortPoolNew         = 21
	CmdSbPortPoolDel         = 22
	CmdSbTcPoolBindGet       = 23
	CmdSbTcPoolBindSet       = 24
	CmdSbTcPoolBindNew       = 25
	CmdSbTcPoolBindDel       = 26
	CmdSbOccSnapshot         = 27
	CmdSbOccMaxClear         = 28
	CmdEswitchGet            = 29
	CmdEswitchSet            = 30
	CmdDpipeTableGet         = 31
	CmdDpipeEntriesGet       = 32
	CmdDpipeHeadersGet       = 33
	CmdDpipeTableCountersSet = 34
	CmdResourceSet           = 35
	CmdResourceDump          = 36
	CmdReload                = 37
	__CmdMax                 = 38
	CmdMax                   = __CmdMax - 1
)

// devlinkPortType as declared in dlh/devlink.h:86
type devlinkPortType int32

// devlinkPortType enumeration from dlh/devlink.h:86
const (
	PortTypeNotset = iota
	PortTypeAuto   = 1
	PortTypeEth    = 2
	PortTypeIb     = 3
)

// devlinkSbPoolType as declared in dlh/devlink.h:93
type devlinkSbPoolType int32

// devlinkSbPoolType enumeration from dlh/devlink.h:93
const (
	SbPoolTypeIngress = iota
	SbPoolTypeEgress  = 1
)

// devlinkSbThresholdType as declared in dlh/devlink.h:111
type devlinkSbThresholdType int32

// devlinkSbThresholdType enumeration from dlh/devlink.h:111
const (
	SbThresholdTypeStatic  = iota
	SbThresholdTypeDynamic = 1
)

// devlinkEswitchMode as declared in dlh/devlink.h:118
type devlinkEswitchMode int32

// devlinkEswitchMode enumeration from dlh/devlink.h:118
const (
	EswitchModeLegacy    = iota
	EswitchModeSwitchdev = 1
)

// devlinkEswitchInlineMode as declared in dlh/devlink.h:123
type devlinkEswitchInlineMode int32

// devlinkEswitchInlineMode enumeration from dlh/devlink.h:123
const (
	EswitchInlineModeNone      = iota
	EswitchInlineModeLink      = 1
	EswitchInlineModeNetwork   = 2
	EswitchInlineModeTransport = 3
)

// devlinkEswitchEncapMode as declared in dlh/devlink.h:130
type devlinkEswitchEncapMode int32

// devlinkEswitchEncapMode enumeration from dlh/devlink.h:130
const (
	EswitchEncapModeNone  = iota
	EswitchEncapModeBasic = 1
)

// devlinkAttr as declared in dlh/devlink.h:135
type devlinkAttr int32

// devlinkAttr enumeration from dlh/devlink.h:135
const (
	AttrUnspec                    = iota
	AttrBusName                   = 1
	AttrDevName                   = 2
	AttrPortIndex                 = 3
	AttrPortType                  = 4
	AttrPortDesiredType           = 5
	AttrPortNetdevIfindex         = 6
	AttrPortNetdevName            = 7
	AttrPortIbdevName             = 8
	AttrPortSplitCount            = 9
	AttrPortSplitGroup            = 10
	AttrSbIndex                   = 11
	AttrSbSize                    = 12
	AttrSbIngressPoolCount        = 13
	AttrSbEgressPoolCount         = 14
	AttrSbIngressTcCount          = 15
	AttrSbEgressTcCount           = 16
	AttrSbPoolIndex               = 17
	AttrSbPoolType                = 18
	AttrSbPoolSize                = 19
	AttrSbPoolThresholdType       = 20
	AttrSbThreshold               = 21
	AttrSbTcIndex                 = 22
	AttrSbOccCur                  = 23
	AttrSbOccMax                  = 24
	AttrEswitchMode               = 25
	AttrEswitchInlineMode         = 26
	AttrDpipeTables               = 27
	AttrDpipeTable                = 28
	AttrDpipeTableName            = 29
	AttrDpipeTableSize            = 30
	AttrDpipeTableMatches         = 31
	AttrDpipeTableActions         = 32
	AttrDpipeTableCountersEnabled = 33
	AttrDpipeEntries              = 34
	AttrDpipeEntry                = 35
	AttrDpipeEntryIndex           = 36
	AttrDpipeEntryMatchValues     = 37
	AttrDpipeEntryActionValues    = 38
	AttrDpipeEntryCounter         = 39
	AttrDpipeMatch                = 40
	AttrDpipeMatchValue           = 41
	AttrDpipeMatchType            = 42
	AttrDpipeAction               = 43
	AttrDpipeActionValue          = 44
	AttrDpipeActionType           = 45
	AttrDpipeValue                = 46
	AttrDpipeValueMask            = 47
	AttrDpipeValueMapping         = 48
	AttrDpipeHeaders              = 49
	AttrDpipeHeader               = 50
	AttrDpipeHeaderName           = 51
	AttrDpipeHeaderId             = 52
	AttrDpipeHeaderFields         = 53
	AttrDpipeHeaderGlobal         = 54
	AttrDpipeHeaderIndex          = 55
	AttrDpipeField                = 56
	AttrDpipeFieldName            = 57
	AttrDpipeFieldId              = 58
	AttrDpipeFieldBitwidth        = 59
	AttrDpipeFieldMappingType     = 60
	AttrPad                       = 61
	AttrEswitchEncapMode          = 62
	AttrResourceList              = 63
	AttrResource                  = 64
	AttrResourceName              = 65
	AttrResourceId                = 66
	AttrResourceSize              = 67
	AttrResourceSizeNew           = 68
	AttrResourceSizeValid         = 69
	AttrResourceSizeMin           = 70
	AttrResourceSizeMax           = 71
	AttrResourceSizeGran          = 72
	AttrResourceUnit              = 73
	AttrResourceOcc               = 74
	AttrDpipeTableResourceId      = 75
	AttrDpipeTableResourceUnits   = 76
	__AttrMax                     = 77
	AttrMax                       = __AttrMax - 1
)

// devlinkDpipeFieldMappingType as declared in dlh/devlink.h:236
type devlinkDpipeFieldMappingType int32

// devlinkDpipeFieldMappingType enumeration from dlh/devlink.h:236
const (
	DpipeFieldMappingTypeNone    = iota
	DpipeFieldMappingTypeIfindex = 1
)

// devlinkDpipeMatchType as declared in dlh/devlink.h:242
type devlinkDpipeMatchType int32

// devlinkDpipeMatchType enumeration from dlh/devlink.h:242
const (
	DpipeMatchTypeFieldExact = iota
)

// devlinkDpipeActionType as declared in dlh/devlink.h:247
type devlinkDpipeActionType int32

// devlinkDpipeActionType enumeration from dlh/devlink.h:247
const (
	DpipeActionTypeFieldModify = iota
)

// devlinkDpipeFieldEthernetId as declared in dlh/devlink.h:251
type devlinkDpipeFieldEthernetId int32

// devlinkDpipeFieldEthernetId enumeration from dlh/devlink.h:251
const (
	DpipeFieldEthernetDstMac = iota
)

// devlinkDpipeFieldIpv4Id as declared in dlh/devlink.h:255
type devlinkDpipeFieldIpv4Id int32

// devlinkDpipeFieldIpv4Id enumeration from dlh/devlink.h:255
const (
	DpipeFieldIpv4DstIp = iota
)

// devlinkDpipeFieldIpv6Id as declared in dlh/devlink.h:259
type devlinkDpipeFieldIpv6Id int32

// devlinkDpipeFieldIpv6Id enumeration from dlh/devlink.h:259
const (
	DpipeFieldIpv6DstIp = iota
)

// devlinkDpipeHeaderId as declared in dlh/devlink.h:263
type devlinkDpipeHeaderId int32

// devlinkDpipeHeaderId enumeration from dlh/devlink.h:263
const (
	DpipeHeaderEthernet = iota
	DpipeHeaderIpv4     = 1
	DpipeHeaderIpv6     = 2
)

// devlinkResourceUnit as declared in dlh/devlink.h:269
type devlinkResourceUnit int32

// devlinkResourceUnit enumeration from dlh/devlink.h:269
const (
	ResourceUnitEntry = iota
)

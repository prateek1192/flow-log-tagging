package types

type FlowLog struct {
	Version     int32
	AccountId   string
	InterfaceId string
	SrcAddr     string
	DstAddr     string
	SrcPort     string
	DstPort     string
	Protocol    string
	Packets     int64
	Bytes       int64
	Start       int64
	End         int64
	Action      string
	LogStatus   string
}

type LookupEntry struct {
	DstPort  string `csv:"dstport"`
	Protocol string `csv:"protocol"`
	Tag      string `csv:"tag"`
}

type LookupTable map[LookupKey][]string

type LookupKey struct {
	DstPort  string
	Protocol string
}

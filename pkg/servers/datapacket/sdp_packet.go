package datapacket

type SdpItem struct {
	AVType             string
	Codec              string
	TimeScale          int
	Control            string
	Rtpmap             int
	Config             []byte
	SpropParameterSets [][]byte
	PayloadType        int
	SizeLength         int
	IndexLength        int
}

type SdpData struct {
	Items map[string]*SdpItem
}

func MakeSdpData() *SdpData {

	tmp := &SdpData{}
	tmp.Items = make(map[string]*SdpItem)

	return tmp
}

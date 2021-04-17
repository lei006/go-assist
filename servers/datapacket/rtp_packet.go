package datapacket

type RtpPacket struct {
	Version        int
	Padding        bool
	Extension      bool
	CSRCCnt        int
	Marker         bool
	PayloadType    int
	SequenceNumber int
	Timestamp      int
	SSRC           int
	Payload        []byte
	PayloadOffset  int
}

func (rtp *RtpPacket) GetType() PacketType {
	return PacketType_RTP
}

func (rtp *RtpPacket) GetDataSize() int {
	return len(rtp.Payload)
}
func (rtp *RtpPacket) GetData() []byte {
	return rtp.Payload
}

func (rtp *RtpPacket) PacketInfo() string {

	return "this is rtpPacket!"
}

package datapacket

type RtpHeader struct {
	Version        int
	Padding        bool
	Extension      bool
	CSRCCnt        int
	Marker         bool
	PayloadType    int
	SequenceNumber int
	Timestamp      int
	SSRC           int //同步标识， 是一个随机数，在同一个RTP会话中只有一个同步标识

	PayloadStart int
	PayloadEnd   int
}

type RtpPacket struct {
	Header RtpHeader
	Raw    []byte //有效负荷
}

func (rtp *RtpPacket) GetType() PacketType {
	return PacketType_RTP
}

func (rtp *RtpPacket) GetDataSize() int {
	return len(rtp.Raw)
}

func (rtp *RtpPacket) Payload() []byte {

	if rtp.Header.PayloadEnd-rtp.Header.PayloadStart < 1 {
		return nil
	}

	return rtp.Raw[rtp.Header.PayloadStart:rtp.Header.PayloadEnd]
}

func (rtp *RtpPacket) PacketInfo() string {

	return "this is rtpPacket!"
}

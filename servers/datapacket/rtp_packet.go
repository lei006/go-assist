package datapacket

import (
	pion_rtp "github.com/pion/rtp"
)

type RtpPacket struct {
	pion_rtp.Packet
}

func (rtp *RtpPacket) GetType() PacketType {
	return PacketType_RTP
}

func (rtp *RtpPacket) GetDataSize() int {
	return len(rtp.Raw)
}

func (rtp *RtpPacket) PacketInfo() string {

	return "this is rtpPacket!"
}

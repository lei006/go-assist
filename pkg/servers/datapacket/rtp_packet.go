package datapacket

import (
	"errors"

	"github.com/lei006/go-assist/utils"
	pion_rtp "github.com/pion/rtp"
)

const (
	// 注意，一般情况下，AVC使用96，AAC使用97，HEVC使用98
	// 但是我还遇到过：
	// HEVC使用96
	// AVC使用105
	RTPPacketTypeAVCOrHEVC = 96
	RTPPacketTypeAAC       = 97
	RTPPacketTypeHEVC      = 98
)

type RtpPacket struct {
	pion_rtp.Packet

	h264Packet *H264Packet
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

func (packet *RtpPacket) Unmarshal(payload []byte) error {
	return packet.Packet.Unmarshal(payload)
}

//转成 datapacket 包.
func (packet *RtpPacket) ToDataPacket() (DataPacket, error) {

	if packet.Header.PayloadType == RTPPacketTypeAVCOrHEVC {

		if packet.h264Packet == nil {
			// 是一个新H264包
			packet.h264Packet = &H264Packet{}
		}

		is_success, err := packet.h264Packet.UnmarshalRTP(packet.Payload)
		if err != nil {
			packet.h264Packet = nil
			return nil, errors.New("rtp packet: parsing the h264 failure--" + err.Error())
		}
		if !is_success {
			//长度为 0 说明包未结束, 所以不返回错，也不返回包
			return nil, nil
		}
		//fmt.Println("len =", len(packet.h264Packet.Payload))

		utils.PrintBinStr(packet.h264Packet.Payload, 12, "xx + "+string(packet.h264Packet.GetType()))

		ret := packet.h264Packet
		packet.h264Packet = nil //清包，下次是一个新包...
		return ret, nil
	}

	return nil, errors.New("rtp ToDataPacket: Untreated load package type")
}

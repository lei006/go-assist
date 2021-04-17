package server_rtsp

import (
	"bytes"

	"github.com/lei006/go-assist/servers/datapacket"
)

type RtspPacket struct {
	datapacket.DataPacket

	Buffer *bytes.Buffer
}

func MakeRtspPacket(packet *RTPPack) *RtspPacket {
	tmp := &RtspPacket{}

	tmp.Buffer = packet.Buffer

	//buf := packet.Buffer.Bytes()

	return tmp
}

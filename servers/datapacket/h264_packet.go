package datapacket

import (
	pion_codecs "github.com/pion/rtp/codecs"
)

type H264Packet struct {
	pion_codecs.H264Packet

	DataPacket
	IsKey bool
	SPS   []byte
	PPS   []byte
}

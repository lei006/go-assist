package datapacket

type H264Packet struct {
	DataPacket
	IsKey bool
	SPS   []byte
	PPS   []byte
}

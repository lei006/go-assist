package datapacket

type PacketType string

const (
	PacketType_H264  PacketType = "H264"
	PacketType_H265  PacketType = "H265"
	PacketType_FLV   PacketType = "FLV"
	PacketType_AAC   PacketType = "AAC"
	PacketType_G711A PacketType = "G711A"
	PacketType_G711U PacketType = "G711U"
)

type DataPacket interface {
	GetType() PacketType
	GetDataSize() int
	GetData() []byte
	PacketInfo() string
}

type AudioPacket interface {
	DataPacket
}

type VideoPacket interface {
	DataPacket
	IsKeyFrame() bool //是否是关键帧
}

type PacketCallback func(DataPacket)

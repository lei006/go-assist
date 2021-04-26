package datapacket

import (
	pion_codecs "github.com/pion/rtp/codecs"
)

type H264Packet struct {
	packet pion_codecs.H264Packet

	Payload []byte
	SPS     []byte
	PPS     []byte
}

func (p *H264Packet) GetType() PacketType {
	return PacketType_H264
}

func (p *H264Packet) UnmarshalRTP(payload []byte) (bool, error) {

	if len(payload) > 3 && payload[0] == 0 && payload[1] == 0 && payload[2] == 1 {
		p.Payload = payload
		return true, nil
	}

	data, err := p.packet.Unmarshal(payload)
	if err != nil {
		//解析出错
		return false, err
	}

	if len(data) == 0 {
		//包未结束
		return false, nil
	}

	//解析出数据
	p.Payload = data

	return true, err
}

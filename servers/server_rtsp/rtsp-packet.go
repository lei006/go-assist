package server_rtsp

import "github.com/lei006/go-assist/servers/datapacket"

type RtspPacket interface {
	datapacket.DataPacket
}

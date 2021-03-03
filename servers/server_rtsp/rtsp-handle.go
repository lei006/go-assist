package server_rtsp

import "net"

type RTSPHandle struct {
	conn net.Conn
}

func MakeRtspHandle(conn net.Conn) *RTSPHandle {
	handle := RTSPHandle{
		conn: conn,
	}

	return &handle
}

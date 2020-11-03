package main

import (
	"go-assist/protocol/rtmp_server"
)

func main() {

	rtmpServer := rtmp_server.MakeRtmpServer()
	rtmpServer.Start()

}

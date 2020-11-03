package main

import (
	"fmt"
	"go-assist/protocol/rtmp_server"
)

func main() {
	fmt.Println("1111111111")
	rtmpServer := rtmp_server.MakeRtmpServer()
	rtmpServer.Start()

}

package main

import (
	"fmt"
	"go-assist/protocol/rtmp_server"
)

func main() {
	fmt.Println("============>")
	rtmpServer := rtmp_server.MakeRtmpServer()
	rtmpServer.Start()

}

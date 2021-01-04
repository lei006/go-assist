package main

import (
	"fmt"
	"github.com/lei006/go-assist/protocol/rtmp_server"
)

func main() {
	fmt.Println("============>")
	rtmpServer := rtmp_server.MakeRtmpServer()
	rtmpServer.Start()

}

package main

import (
	"bytes"
	"fmt"

	"github.com/lei006/go-assist/servers/server_rtsp"
)

func main() {

	rtspServer := server_rtsp.GetServer()
	var index int
	rtspServer.PacketCallback(func(url string, buffer *bytes.Buffer) {
		fmt.Println("url =", url, index, buffer.Bytes()[:20])
		index++
	})

	rtspServer.ClientAsk(func(req *server_rtsp.Request) error {
		fmt.Println("url =========================  start", req.URL)
		fmt.Println("sdp :", req.Body)
		fmt.Println("url =========================  end")

		return nil
	})

	rtspServer.Start()

}

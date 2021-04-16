package main

import "github.com/lei006/go-assist/servers/server_rtsp"

func main() {

	rtspServer := server_rtsp.GetServer()

	rtspServer.Start()

}

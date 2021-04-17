package main

import "github.com/lei006/go-assist/servers/server_cvr"

func main() {

	rtspServer := server_cvr.ServerCVR{}

	rtspServer.Start()

}

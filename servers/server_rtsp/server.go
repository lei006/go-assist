package server_rtsp

import "fmt"

type RtspServer struct {
}

func MakeServer() *RtspServer {

	tmp := &RtspServer{}

	return tmp
}

func (this *RtspServer) Start() {
	fmt.Println("xxxxxxxxxxx")
}

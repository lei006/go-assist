package rtsp_server

import (
	"fmt"
)

type RtspServer struct {
	Name string
}

// MakeNewServer 创建新的webrtc 服务
func MakeNewServer() *RtspServer {
	return &RtspServer{Name: "RtspServer"}
}

func (server *RtspServer) GetName() string {
	fmt.Println("RtspServer ==> GetName", server.Name)
	return server.Name
}

func (server *RtspServer) Start() {

	fmt.Println("RtspServer => Start")

}

func (server *RtspServer) Stop() {

	fmt.Println("RtspServer ==> Stop")

}

func (server *RtspServer) GetInfo() {

	fmt.Println("RtspServer ==> GetInfo")

}

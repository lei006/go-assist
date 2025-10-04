package server_rtmp

import (
	"fmt"

	"github.com/lei006/go-assist/servers/livertc"
	"github.com/q191201771/lal/pkg/base"
)

type HandlePubProxy struct {
	session *livertc.Session
}

func MakeHandlePubProxy(session *livertc.Session) *HandlePubProxy {
	return &HandlePubProxy{
		session: session,
	}
}

// 注意，回调结束后，内部会复用Payload内存块
func (handle *HandlePubProxy) OnReadRTMPAVMsg(msg base.RTMPMsg) {
	fmt.Println("====================== OnReadRTMPAVMsg ")
}

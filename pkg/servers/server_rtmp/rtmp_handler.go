package server_rtmp

import (
	"fmt"

	"github.com/lei006/go-assist/servers/livertc"
	cmap "github.com/orcaman/concurrent-map"

	lal_rtmp "github.com/q191201771/lal/pkg/rtmp"
)

type RtmpHandler struct {
	engine *livertc.EngineMedia

	pub_cmap cmap.ConcurrentMap //存放发布者..
}

func (handler *RtmpHandler) OnRTMPConnect(session *lal_rtmp.ServerSession, opa lal_rtmp.ObjectPairArray) {

}

// 返回true则允许推流，返回false则强制关闭这个连接
func (handler *RtmpHandler) OnNewRTMPPubSession(session *lal_rtmp.ServerSession) bool {

	fmt.Println("======================  OnNewRTMPPubSession path =", session.URL())
	/////////////////////////////////////// 申请推流
	engion_session, err := handler.engine.AddPubSession(session.URL())
	if err != nil {
		// 创建失败.
		return false
	}

	handlerSession := MakeHandlePubProxy(engion_session)
	session.SetPubSessionObserver(handlerSession)
	session.SetPubSessionObserver(handlerSession)
	return true

}

func (handler *RtmpHandler) OnDelRTMPPubSession(session *lal_rtmp.ServerSession) {

}

// 返回true则允许拉流，返回false则强制关闭这个连接
func (handler *RtmpHandler) OnNewRTMPSubSession(session *lal_rtmp.ServerSession) bool {
	return true
}

func (handler *RtmpHandler) OnDelRTMPSubSession(session *lal_rtmp.ServerSession) {

}

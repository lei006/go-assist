package server_rtsp

import (
	"fmt"

	"github.com/lei006/go-assist/servers/livertc"
	cmap "github.com/orcaman/concurrent-map"

	lal_rtsp "github.com/q191201771/lal/pkg/rtsp"
)

type RtspHandler struct {
	engine *livertc.EngineMedia

	pub_cmap cmap.ConcurrentMap //存放发布者..
}

// @brief 使得上层有能力管理未进化到Pub、Sub阶段的Session
func (handle *RtspHandler) OnNewRTSPSessionConnect(session *lal_rtsp.ServerCommandSession) {
	fmt.Println("======================  OnNewRTSPSessionConnect")
}

// @brief 注意，对于已经进化到了Pub、Sub阶段的Session，该回调依然会被调用
func (handle *RtspHandler) OnDelRTSPSession(session *lal_rtsp.ServerCommandSession) {
	fmt.Println("======================  OnDelRTSPSession")

}

///////////////////////////////////////////////////////////////////////////

// @brief  Announce阶段回调
// @return 如果返回false，则表示上层要强制关闭这个推流请求
func (handle *RtspHandler) OnNewRTSPPubSession(session *lal_rtsp.PubSession) bool {

	fmt.Println("======================  OnNewRTSPPubSession path =", session.URL())
	/////////////////////////////////////// 申请推流
	engion_session, err := handle.engine.AddPubSession(session.URL())
	if err != nil {
		// 创建失败.
		return false
	}

	handlerSession := MakeHandlePubProxy(engion_session)
	session.SetObserver(handlerSession)
	return true
}

func (handle *RtspHandler) OnDelRTSPPubSession(session *lal_rtsp.PubSession) {
	fmt.Println("======================  OnDelRTSPPubSession")
	//handle.engine.RemoveSession()
}

///////////////////////////////////////////////////////////////////////////

// @return 如果返回false，则表示上层要强制关闭这个拉流请求
// @return sdp
func (handle *RtspHandler) OnNewRTSPSubSessionDescribe(session *lal_rtsp.SubSession) (ok bool, sdp []byte) {
	fmt.Println("======================  OnNewRTSPSubSessionDescribe")

	return true, nil
}

// @brief Describe阶段回调
// @return ok  如果返回false，则表示上层要强制关闭这个拉流请求
func (handle *RtspHandler) OnNewRTSPSubSessionPlay(session *lal_rtsp.SubSession) bool {
	fmt.Println("======================  OnNewRTSPSubSessionPlay")
	/////////////////////////////////////// 申请播放..

	return true
}

func (handle *RtspHandler) OnDelRTSPSubSession(session *lal_rtsp.SubSession) {
	fmt.Println("======================  OnDelRTSPSubSession")
}

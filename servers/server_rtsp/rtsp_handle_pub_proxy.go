package server_rtsp

import (
	"fmt"

	"github.com/lei006/go-assist/servers/datapacket"
	"github.com/lei006/go-assist/servers/livertc"
	"github.com/q191201771/lal/pkg/base"
	"github.com/q191201771/lal/pkg/rtprtcp"
)

type HandlePubProxy struct {
	session *livertc.Session
}

func MakeHandlePubProxy(session *livertc.Session) *HandlePubProxy {
	return &HandlePubProxy{
		session: session,
	}
}

var rtp_packet datapacket.RtpPacket

func (handle *HandlePubProxy) OnRTPPacket(pkt rtprtcp.RTPPacket) {
	//utils.PrintBin(pkt.Raw, 18)

	err := rtp_packet.Unmarshal(pkt.Raw)
	if err != nil {
		fmt.Println(" ", err.Error())
		return
	}

	packet, err := rtp_packet.ToDataPacket()
	if err != nil {
		fmt.Println(" error: ", err.Error())
		return
	}
	if packet == nil {
		fmt.Println(" package not end ")
		return
	}

	//fmt.Println(" packet.GetType() =", packet.GetType())

}

// @param asc: AAC AudioSpecificConfig，注意，如果不存在音频或音频不为AAC，则为nil
// @param vps, sps, pps 如果都为nil，则没有视频，如果sps, pps不为nil，则vps不为nil是H265，vps为nil是H264
//
// 注意，4个参数可能同时为nil
func (handle *HandlePubProxy) OnAVConfig(asc, vps, sps, pps []byte) {
	fmt.Println("====================== OnAVConfig ")

}

// @param pkt: pkt结构体中字段含义见rtprtcp.OnAVPacket
func (handle *HandlePubProxy) OnAVPacket(pkt base.AVPacket) {

}

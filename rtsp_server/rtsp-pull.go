package rtsp_server

import (
	
	"log"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/codec/h264parser"
	"github.com/deepch/vdk/format/rtsp"
)

type RtspPuller struct {
	core.PublishContext
	RtspUrl string `bson:"rtsp_url" json:"rtsp_url"`
}

func MakeNewRtspPull(id string, rtsp_url string, title string) (*RtspPuller, error) {

	//id := utils.Md5V2(rtsp_url)

	rtspPuller := &RtspPuller{}
	rtspPuller.Title = title
	rtspPuller.Tag = rtsp_url

	err := rtspPuller.Init(id, "rtsp")
	if err != nil {
		return nil, err
	}

	//处理数据包函数...
	handleWritePacket := func(codecs []av.CodecData, pck av.Packet) {

		sps := codecs[0].(h264parser.CodecData).SPS()
		pps := codecs[0].(h264parser.CodecData).PPS()

		//utils.PrintBin(pck.Data[4:], 30)
		/*
			if pck.IsKeyFrame {
				pck.Data = append([]byte("\000\000\001"+string(sps)+"\000\000\001"+string(pps)+"\000\000\001"), pck.Data[4:]...)

			} else {
				pck.Data = pck.Data[4:]
			}
		*/
		if pck.Idx == 0 {

			//packet := core.MakeNewDataPacket(pck.Data, samples)
			// 把长度去掉...
			packet := core.MakeNewDataPacket(pck.Data[4:])
			packet.IsKeyFrame = pck.IsKeyFrame
			packet.SPS = sps
			packet.PPS = pps
			rtspPuller.Publish(packet)

		} else {
		}
	}

	go func(name, url string) {

		const KeepSleepTime = 5
		for {
			rtsp.DebugRtsp = false
			session, err := rtsp.Dial(url)
			if err != nil {
				if rtspPuller.IsExit() {
					break
				}
				time.Sleep(2 * time.Second)
				continue
			}
			//连接成功....
			rtspPuller.SetState("connected")

			session.RtpKeepAliveTimeout = KeepSleepTime * time.Second
			if err != nil {
				log.Println(name, err)
				time.Sleep(KeepSleepTime * time.Second)
				continue
			}
			codec, err := session.Streams()
			if err != nil {
				log.Println(name, err)
				time.Sleep(KeepSleepTime * time.Second)
				continue
			}

			for {
				pkt, err := session.ReadPacket()
				if err != nil {
					log.Println(name, err)
					break
				}

				handleWritePacket(codec, pkt)
				if rtspPuller.IsExit() {
					rtspPuller.Close()
					break
				}
			}
			err = session.Close()
			if err != nil {
				log.Println("session Close error", err)
			}

			if rtspPuller.IsExit() {
				rtspPuller.Close()
				break
			}

			log.Println(name, "reconnect wait 5s")

			time.Sleep(KeepSleepTime * time.Second)
		}

	}(id, rtsp_url)

	return nil, nil
}

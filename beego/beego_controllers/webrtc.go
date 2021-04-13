package beego_controllers

import (
	"livertc/core/beego_assert"
	"livertc/core/intfs"
	"livertc/core/tools"
	"livertc/core/types"

	beego "github.com/beego/beego/v2/server/web"
)

type WebrtcController struct {
	beego_assert.BaseController
}

func (this *WebrtcController) Router(app intfs.IApplication, prefix string) {
	this.App = app

	beego.Router(prefix+"webrtc/subscribers", this, "get:GetSubscribers")      //定阅者列表
	beego.Router(prefix+"webrtc/subscriber/:id", this, "get:GetSubscriber")    //取得定阅者信息
	beego.Router(prefix+"webrtc/subscriber/:id", this, "delete:DelSubscriber") //发布者
	beego.Router(prefix+"webrtc/subscriber", this, "post:AddSubscriber")       //发布者

	beego.Router(prefix+"webrtc/publishers", this, "get:GetPublishers")      //发布者列表
	beego.Router(prefix+"webrtc/publisher", this, "post:AddPublisher")       //发布者
	beego.Router(prefix+"webrtc/publisher/:id", this, "delete:DelPublisher") //发布者

}

/**
 * @api {get} /publisher 1. 发布列表
 * @apiGroup 3.发布媒体
 * @apiName publisher list
 * @apiHeader {String} x-token 授权码
 * @apiSuccessExample {json} 成功响应(三份数据.):
{
  code: 20000
  data: [,…]
  0: {id: "Gz39D7rmoPUA1xXi", title: "Web发布者", type: "webrtc", state: "connected", tag: "sdp", startAt: 0,…}
  1: {id: "11111111133", title: "RTSP发布者", type: "rtsp", state: "connected",…}
  2: {id: "9yTtuGxn1XHv9NLF", title: "Web发布者", type: "webrtc", state: "connected", tag: "sdp", startAt: 0,…}
  message: "success"
}
*/

func (this *WebrtcController) GetPublishers() {
	//list := application.GetDefaultApp().GetEngine().GetPublisherList()
	//this.SuccReturnList(list, int64(len(list)))
}

/**
 * @api {get} /subscriber 1. 定阅列表
 * @apiGroup 4.定阅媒体
 * @apiName subscriber list
 * @apiHeader {String} x-token 授权码
 * @apiSuccessExample {json} 成功响应:
{
  code: 20000
  data: [{id: "mgno4ibtv53bkck", title: "", type: "", state: "connected", tag: "sdp", startAt: 1599555562,…}]
  message: "success"
}

*/

func (this *WebrtcController) GetSubscribers() {

	//list := application.GetDefaultApp().GetEngine().GetSublisherList()
	//this.SuccReturnList(list, int64(len(list)))
}

func (this *WebrtcController) AddPublisher() {
	/*
		str_body := string(this.Ctx.Input.RequestBody)
		pub_id := gjson.Get(str_body, "pub_id").String()
		sdp_bate64 := gjson.Get(str_body, "sdp").String()

		old_pub := application.GetDefaultApp().GetEngine().GetPublisher(pub_id)
		if old_pub != nil {
			old_pub.Close("add publisher close old")
		}

		new_pub, err := application.GetDefaultApp().GetEngine().MakePublisher(pub_id, "webrtc")
		if err != nil {
			logs.Debug("MakePublisher error : " + err.Error())
			this.ReturnFail("make new publisher fail:" + err.Error())
		}

		logs.Debug("add new  public ", pub_id, "    sublisher count = ", new_pub.GetSublisherCount())

		iceConfig, err := this.getIceServers()
		if err != nil {
			this.ReturnFail("make new publisher fail:" + err.Error())
		}

		receiver, err := p2p_webrtc.MakeReceiver(sdp_bate64, iceConfig.ToWebrtcConfig())
		if err != nil {
			logs.Debug("p2p_webrtc.MakeSender error : " + err.Error())
			this.ReturnFail("new sender:" + err.Error())
		}

		go func() {
			for {

				// 读视频...
				rtpPacket := receiver.ReadVideoPacket()

				if new_pub.IsAskExit() {
					receiver.Close()
					fmt.Println("exit  video---2-----------", new_pub.GetID())
					break
				}

				if receiver.IsAskExit() == true {
					new_pub.Close("remote close")

					fmt.Println("exit  video---1-----------", new_pub.GetID())

					break
				}

				newPack := datapacket.MakeRtpVideoPacket(*rtpPacket)
				new_pub.SendPacketP2P(newPack)

			}

		}()

		go func() {
			// 读音频...

			// 读视频...
			for {
				rtpPacket := receiver.ReadAudioPacket()
				//logs.Debug("read audio packet ......")
				if new_pub.IsAskExit() {
					receiver.Close()
					fmt.Println("exit  audio---2-----------", new_pub.GetID())

					break
				}

				if receiver.IsAskExit() == true {
					new_pub.Close("remote close")
					fmt.Println("exit  audio---1-----------", new_pub.GetID())

					break
				}

				newPack := datapacket.MakeRtpAudioPacket(*rtpPacket)
				new_pub.SendPacketP2P(newPack)
			}

		}()

		this.ReturnSuccess(receiver)
	*/
}

/**
 * @api {Post} /subscriber 2. 定阅媒体
 * @apiGroup 4.定阅媒体
 * @apiName subscriber add
 * @apiHeader {String} x-token 授权码
 * @apiParam {String} pub_id 定阅者id
 * @apiParam {string} helper 辅助数据(webrtc为bate64的sdp)
 * @apiSuccessExample {json} 成功响应:
{
  "message": "success",
  "code": 20000,
  "data": [
    {
      "id": "l9xhg7lblf285ao",
      "source_type": "rtsp",
      "bufferLen": 10,
      "timeout": 100000000,
      "timelen": 6371,
      "startAt": "2020-06-19 14:49:52",
      "state": "play",
      "inBytes": 0,
      "onlines": 0
    }
  ]
}

*/

func (this *WebrtcController) AddSubscriber() {
	/*
		str_body := string(this.Ctx.Input.RequestBody)
		pub_id := gjson.Get(str_body, "pub_id").String()
		sdp_bate64 := gjson.Get(str_body, "sdp").String()

		new_sub, err := application.GetDefaultApp().GetEngine().MakeSubscriber(pub_id, "webrtc")
		if err != nil {
			logs.Debug("MakeSubscriber error : " + err.Error())

			this.ReturnFailCode(40010, "make new subscriber fail:"+err.Error())
		}

		iceConfig, err := this.getIceServers()
		if err != nil {
			this.ReturnFailCode(40010, "make new publisher fail:"+err.Error())
		}

		sender, err := p2p_webrtc.MakeSender(sdp_bate64, iceConfig.ToWebrtcConfig())
		if err != nil {
			logs.Debug("p2p_webrtc.MakeSender error : " + err.Error())
			this.ReturnFailCode(40010, "new sender:"+err.Error())
		}
		sender.Id = new_sub.GetID()
		sender.PubID = new_sub.GetPublisherID()

		// 处理对讲视频，语音...
		new_sub.RecvPacketP2P(func(packet_data intfs.DataPacket) {
			switch packet_data.(type) {
			case *datapacket.PacketG711:
				sender.SendPacketH264(packet_data)

			case *datapacket.H264Packet:
				sender.SendPacketH264(packet_data)
			default:
				logs.Warn("error packet type ")
			}
		})

		go func() {

			var err error

			for {

				packet_data := new_sub.RecvPacket()
				switch packet_data.(type) {
				case *datapacket.PacketG711:
					sender.SendPacketH264(packet_data)

				case *datapacket.H264Packet:
					err = sender.SendPacketH264(packet_data)
				case *datapacket.AudioAlawPacket:
					err = sender.SendPacketAudioAlaw(packet_data.(*datapacket.AudioAlawPacket))
				default:
					logs.Warn("error packet type ")
				}

				if err != nil {
					sender.Close()
					new_sub.Close("sender webrtc SendPacketH264 error:" + err.Error())
					break
				}

				if new_sub.IsAskExit() {
					sender.Close()
					break
				}

				if sender.IsAskExit() {
					new_sub.Close("sender webrtc remote peer close")
					break
				}

			}
		}()

		logs.Notice(" new subscriber id=", sender.Id, "  pub_id=", pub_id, "   from ", this.Ctx.Input.IP())
		this.SuccReturn(sender)
	*/
}

/**
 * @api {delete} /subscriber/:id 3. 移除定阅
 * @apiGroup 4.定阅媒体
 * @apiName Remove Subscriber
 * @apiHeader {String} x-token 授权码
 * @apiParam {String} id 定阅者id
 * @apiSuccessExample {json} 成功响应:
{
  "message": "success",
  "code": 20000,
  "data": "ok"
}
*/

func (this *WebrtcController) DelSubscriber() {
	/*
		id := this.GetString(":id")

		sub := application.GetDefaultApp().GetEngine().GetSubscriber(id)
		if sub == nil {
			this.ReturnFailCode(40010, "没有找到这个定阅!")
		}
		sub.Close("通过 system/subscriber api 主动关闭定阅")
	*/

	this.SuccReturn("ok")
}

func (this *WebrtcController) DelPublisher() {

	/*
		id := this.GetString(":id")

		pub := application.GetDefaultApp().GetEngine().GetPublisher(id)
		if pub == nil {
			this.ReturnFail("删除发布出错: 没有找到!")
		}
		pub.Close("通过 system/publisher api 主动关闭定阅")
	*/
	this.SuccReturn("ok")
}

func (this *WebrtcController) GetSubscriber() {
	/*
		id := this.GetString(":id")

		sub := application.GetDefaultApp().GetEngine().GetSubscriber(id)
		if sub == nil {
			this.ReturnFail("没有找到这个定阅!")
		}
		this.SuccReturn(sub)
	*/
}

func (this *WebrtcController) getIceServers() (*types.IceServers, error) {

	iceServer := types.IceServers{}

	{

		val, err := tools.GetConfigItem("StunServers")
		if err != nil {
			return nil, err
		}
		if val != nil {
			iceServer.Stun = val.(string)
		}
	}

	{
		val, err := tools.GetConfigItem("TurnServers")
		if err != nil {
			return nil, err
		}
		if val != nil {
			iceServer.Turn = val.(string)
		}
	}

	{
		val, err := tools.GetConfigItem("TurnAuth")
		if err != nil {
			return nil, err
		}
		if val != nil {
			iceServer.Auth = val.(string)
		}
	}
	return &iceServer, nil
}

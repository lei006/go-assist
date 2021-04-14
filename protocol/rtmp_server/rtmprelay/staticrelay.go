package rtmprelay

import (
	"fmt"
	"sync"

	"github.com/beego/beego/v2/core/logs"
	"github.com/lei006/go-assist/protocol/intfs"
	"github.com/lei006/go-assist/servers/server_livego/configure"
	"github.com/lei006/go-assist/servers/server_livego/protocol/rtmp/core"
)

type StaticPush struct {
	RtmpUrl       string
	packet_chan   chan *intfs.Packet
	sndctrl_chan  chan string
	connectClient *core.ConnClient
	startflag     bool
}

var G_StaticPushMap = make(map[string](*StaticPush))
var g_MapLock = new(sync.RWMutex)

var (
	STATIC_RELAY_STOP_CTRL = "STATIC_RTMPRELAY_STOP"
)

func GetStaticPushList(appname string) ([]string, error) {
	pushurlList, ok := configure.GetStaticPushUrlList(appname)

	if !ok {
		return nil, fmt.Errorf("no static push url")
	}

	return pushurlList, nil
}

func GetAndCreateStaticPushObject(rtmpurl string) *StaticPush {
	g_MapLock.RLock()
	staticpush, ok := G_StaticPushMap[rtmpurl]
	logs.Debug("GetAndCreateStaticPushObject: %s, return %v", rtmpurl, ok)
	if !ok {
		g_MapLock.RUnlock()
		newStaticpush := NewStaticPush(rtmpurl)

		g_MapLock.Lock()
		G_StaticPushMap[rtmpurl] = newStaticpush
		g_MapLock.Unlock()

		return newStaticpush
	}
	g_MapLock.RUnlock()

	return staticpush
}

func GetStaticPushObject(rtmpurl string) (*StaticPush, error) {
	g_MapLock.RLock()
	if staticpush, ok := G_StaticPushMap[rtmpurl]; ok {
		g_MapLock.RUnlock()
		return staticpush, nil
	}
	g_MapLock.RUnlock()

	return nil, fmt.Errorf("G_StaticPushMap[%s] not exist....", rtmpurl)
}

func ReleaseStaticPushObject(rtmpurl string) {
	g_MapLock.RLock()
	if _, ok := G_StaticPushMap[rtmpurl]; ok {
		g_MapLock.RUnlock()

		logs.Debug("ReleaseStaticPushObject %s ok", rtmpurl)
		g_MapLock.Lock()
		delete(G_StaticPushMap, rtmpurl)
		g_MapLock.Unlock()
	} else {
		g_MapLock.RUnlock()
		logs.Debug("ReleaseStaticPushObject: not find %s", rtmpurl)
	}
}

func NewStaticPush(rtmpurl string) *StaticPush {
	return &StaticPush{
		RtmpUrl:       rtmpurl,
		packet_chan:   make(chan *intfs.Packet, 500),
		sndctrl_chan:  make(chan string),
		connectClient: nil,
		startflag:     false,
	}
}

func (self *StaticPush) Start() error {
	if self.startflag {
		return fmt.Errorf("StaticPush already start %s", self.RtmpUrl)
	}

	self.connectClient = core.NewConnClient()

	logs.Debug("static publish server addr:%v starting....", self.RtmpUrl)
	err := self.connectClient.Start(self.RtmpUrl, "publish")
	if err != nil {
		logs.Debug("connectClient.Start url=%v error", self.RtmpUrl)
		return err
	}
	logs.Debug("static publish server addr:%v started, streamid=%d", self.RtmpUrl, self.connectClient.GetStreamId())
	go self.HandleAvPacket()

	self.startflag = true
	return nil
}

func (self *StaticPush) Stop() {
	if !self.startflag {
		return
	}

	logs.Debug("StaticPush Stop: %s", self.RtmpUrl)
	self.sndctrl_chan <- STATIC_RELAY_STOP_CTRL
	self.startflag = false
}

func (self *StaticPush) WriteAvPacket(packet *intfs.Packet) {
	if !self.startflag {
		return
	}

	self.packet_chan <- packet
}

func (self *StaticPush) sendPacket(p *intfs.Packet) {
	if !self.startflag {
		return
	}
	var cs core.ChunkStream

	cs.Data = p.Data
	cs.Length = uint32(len(p.Data))
	cs.StreamID = self.connectClient.GetStreamId()
	cs.Timestamp = p.TimeStamp
	//cs.Timestamp += v.BaseTimeStamp()

	//logs.Printf("Static sendPacket: rtmpurl=%s, length=%d, streamid=%d",
	//	self.RtmpUrl, len(p.Data), cs.StreamID)
	if p.IsVideo {
		cs.TypeID = intfs.TAG_VIDEO
	} else {
		if p.IsMetadata {
			cs.TypeID = intfs.TAG_SCRIPTDATAAMF0
		} else {
			cs.TypeID = intfs.TAG_AUDIO
		}
	}

	self.connectClient.Write(cs)
}

func (self *StaticPush) HandleAvPacket() {
	if !self.IsStart() {
		logs.Debug("static push %s not started", self.RtmpUrl)
		return
	}

	for {
		select {
		case packet := <-self.packet_chan:
			self.sendPacket(packet)
		case ctrlcmd := <-self.sndctrl_chan:
			if ctrlcmd == STATIC_RELAY_STOP_CTRL {
				self.connectClient.Close(nil)
				logs.Debug("Static HandleAvPacket close: publishurl=%s", self.RtmpUrl)
				return
			}
		}
	}
}

func (self *StaticPush) IsStart() bool {
	return self.startflag
}

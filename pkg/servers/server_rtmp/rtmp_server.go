package server_rtmp

import (
	"errors"

	"github.com/lei006/go-assist/servers/livertc"
	"github.com/sirupsen/logrus"

	lal_rtmp "github.com/q191201771/lal/pkg/rtmp"
)

type ServerRtmp struct {
	engine *livertc.EngineMedia

	handler *RtmpHandler

	rtmpServer *lal_rtmp.Server
	debug      bool

	logger *logrus.Logger
}

func MakeRtmpServer(engine *livertc.EngineMedia) *ServerRtmp {

	handler := &RtmpHandler{
		engine: engine, //尽量不入侵原始代码...
	}

	return &ServerRtmp{
		engine:  engine,
		handler: handler,
		logger:  logrus.New(),
	}
}

func (rtmp *ServerRtmp) ListenAt(addr string) error {

	if rtmp.rtmpServer == nil {
		rtmpServer := lal_rtmp.NewServer(rtmp.handler, addr)
		err := rtmpServer.Listen()
		if err != nil {
			return err
		}
		rtmp.rtmpServer = rtmpServer
		return nil
	}

	return errors.New("rtmp is started")
}

func (rtmp *ServerRtmp) RunLoop() error {

	return rtmp.rtmpServer.RunLoop()
}

func (rtmp *ServerRtmp) SetDebug() {
	rtmp.debug = true
}

package server_rtsp

import (
	"errors"

	"github.com/lei006/go-assist/servers/livertc"

	"github.com/sirupsen/logrus"

	lal_rtsp "github.com/q191201771/lal/pkg/rtsp"
)

type ServerRtsp struct {
	engine *livertc.EngineMedia

	handler *RtspHandler

	rtspServer *lal_rtsp.Server
	debug      bool

	logger *logrus.Logger
}

func MakeRtspServer(engine *livertc.EngineMedia) *ServerRtsp {

	handler := &RtspHandler{
		engine: engine, //尽量不入侵原始代码...
	}

	return &ServerRtsp{
		engine:  engine,
		handler: handler,
		logger:  logrus.New(),
	}
}
func (rtsp *ServerRtsp) ListenAt(addr string) error {

	if rtsp.rtspServer == nil {
		rtspServer := lal_rtsp.NewServer(addr, rtsp.handler)
		err := rtspServer.Listen()
		if err != nil {
			return err
		}
		rtsp.rtspServer = rtspServer
		return nil
	}

	return errors.New("Rtsp is started")
}

func (rtsp *ServerRtsp) RunLoop() error {
	return rtsp.rtspServer.RunLoop()
}

func (rtsp *ServerRtsp) SetDebug() {
	rtsp.debug = true
}

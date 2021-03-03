package server_rtsp

import (
	"fmt"
	"livertc/core/application"
	"livertc/core/intfs"
	"net"
	"runtime/debug"
	"sync"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type RtspServer struct {
	SessionLogger
	TCPListener *net.TCPListener
	TCPPort     int
	Stoped      bool
	pushers     map[string]*Pusher // Path <-> Pusher
	pushersLock sync.RWMutex
	Name        string

	handlePlayFilter func(url string) (string, error)
	handleRecvFilter func(url string) (string, error)
}

const (
	Default_RtspPort = 554
)

func MakeServer() *RtspServer {

	//port := appContext.Config.GetConfigInt32("RtspPort", Default_RtspPort)

	port := beego.AppConfig.DefaultInt("rtspport", Default_RtspPort)

	srv := &RtspServer{
		//SessionLogger:  SessionLogger{log.New(os.Stdout, "[LiveRTC-RTSP]", log.LstdFlags|log.Lshortfile)},
		SessionLogger: SessionLogger{logger: logs.GetBeeLogger()},
		Stoped:        true,
		TCPPort:       port,
		pushers:       make(map[string]*Pusher),
		Name:          "rtsp-server",
	}

	return srv
}

func (server *RtspServer) GetName() string {
	return server.Name
}

func (this *RtspServer) GetContext() *intfs.ServerContext {
	context := &intfs.ServerContext{
		ServerName: this.GetName(),
		Ip:         application.GetDefaultApp().RegeditIP,
		Port:       this.TCPPort,
		Weight:     1.0,
		Running:    !this.Stoped,
	}
	return context
}

func (server *RtspServer) Start() (err error) {

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", server.TCPPort))
	if err != nil {
		return
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}

	go server.ListenStart(listener)

	return err
}

func (server *RtspServer) ListenStart(listener *net.TCPListener) (err error) {

	server.Stoped = false
	server.TCPListener = listener

	for !server.Stoped {
		conn, err := server.TCPListener.Accept()
		if err != nil {
			beego.Debug(err.Error())
			continue
		}

		handle := MakeRtspHandle(conn)

		go server.handleConn(handle)

	}
	return
}

func (this *RtspServer) handleConn(handle *RTSPHandle) error {
	defer func() {
		if r := recover(); r != nil {
			beego.Error("rtmp server handleConn panic: ", r, string(debug.Stack()))
		}
	}()

	networkBuffer := 1048576

	if tcpConn, ok := handle.conn.(*net.TCPConn); ok {
		if err := tcpConn.SetReadBuffer(networkBuffer); err != nil {
			beego.Debug("rtsp server conn set read buffer error, %v", err)
		}
		if err := tcpConn.SetWriteBuffer(networkBuffer); err != nil {
			beego.Debug("rtsp server conn set write buffer error, %v", err)
		}
	}

	session := NewSession(this, handle.conn)
	session.Start()

	return nil
}

func (server *RtspServer) Stop() {

	beego.Debug("rtsp server stop on", server.TCPPort)
	server.Stoped = true
	if server.TCPListener != nil {
		server.TCPListener.Close()
		server.TCPListener = nil
	}
	server.pushersLock.Lock()
	server.pushers = make(map[string]*Pusher)
	server.pushersLock.Unlock()

}

func (server *RtspServer) AddPusher(pusher *Pusher) bool {

	added := false
	server.pushersLock.Lock()
	_, ok := server.pushers[pusher.Path()]
	if !ok {
		server.pushers[pusher.Path()] = pusher
		beego.Debug("pusher start:", pusher, "     now pusher size:", len(server.pushers))
		added = true
	} else {
		added = false
	}
	server.pushersLock.Unlock()
	if added {
		go pusher.Start()

	}
	return added
}

func (server *RtspServer) TryAttachToPusher(session *Session) (int, *Pusher) {
	server.pushersLock.Lock()
	attached := 0
	var pusher *Pusher = nil
	if _pusher, ok := server.pushers[session.Path]; ok {
		if _pusher.RebindSession(session) {
			beego.Debug("Attached to a pusher")
			attached = 1
			pusher = _pusher
		} else {
			attached = -1
		}
	}
	server.pushersLock.Unlock()
	return attached, pusher
}

func (server *RtspServer) RemovePusher(pusher *Pusher) {

	pusher.publisher.Close("rtsp pusher close 2")
	server.pushersLock.Lock()
	if _pusher, ok := server.pushers[pusher.Path()]; ok && pusher.ID() == _pusher.ID() {
		delete(server.pushers, pusher.Path())
		beego.Debug("pusher end:", pusher, "     now pusher size:", len(server.pushers))
	}
	server.pushersLock.Unlock()
}

func (server *RtspServer) GetPusher(path string) (pusher *Pusher) {
	server.pushersLock.RLock()
	pusher = server.pushers[path]
	server.pushersLock.RUnlock()
	return
}

func (server *RtspServer) GetPushers() (pushers map[string]*Pusher) {
	pushers = make(map[string]*Pusher)
	server.pushersLock.RLock()
	for k, v := range server.pushers {
		pushers[k] = v
	}
	server.pushersLock.RUnlock()
	return
}

func (server *RtspServer) GetPusherSize() (size int) {
	server.pushersLock.RLock()
	size = len(server.pushers)
	server.pushersLock.RUnlock()
	return
}

func (server *RtspServer) SetRecvFilter(cb_handle func(url string) (string, error)) {

	server.handleRecvFilter = cb_handle

}

func (server *RtspServer) SetPlayFilter(cb_handle func(url string) (string, error)) {

	server.handlePlayFilter = cb_handle

}

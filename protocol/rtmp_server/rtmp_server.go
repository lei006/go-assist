package rtmp_server

import (
	"fmt"
	"github.com/lei006/go-assist/protocol/intfs"
	"github.com/lei006/go-assist/protocol/rtmp_server/core"
	"net"

	"github.com/astaxie/beego/logs"
	"github.com/gwuhaolin/livego/configure"
)

type RtmpServer struct {
	address string
}

func MakeRtmpServer() intfs.Server {

	fmt.Println("xxxxxxxx111xxx")

	srv := &RtmpServer{}
	srv.address = ":1556"
	return srv
}

func (this *RtmpServer) Start() (bool, error) {
	fmt.Println("xxxxxxxxxxx")

	listener, err := net.Listen("tcp", this.address)
	if err != nil {
		logs.Error(err)
		return false, err
	}

	defer func() {
		if r := recover(); r != nil {
			//logs.Error("rtmp serve panic: ", r)
			logs.Error("rtmp serve panic: ", r)
		}
	}()

	logs.Debug(listener, this.address)
	/*
		for {
			var netconn net.Conn
			netconn, err = listener.Accept()
			if err != nil {
				return false, err
			}
			fmt.Println(netconn)

			//andler := MakeRtmpHandler(netconn)
			//go this.Handle()

			conn := core.NewConn(netconn, 4*1024)
			//logs.Debug("new client, connect remote: ", conn.RemoteAddr().String(),"local:", conn.LocalAddr().String())
			go this.handleConn(conn)
		}
	*/
	return true, err
}
func (this *RtmpServer) Stop() {

}

func (this *RtmpServer) handleConn(conn *core.Conn) error {

	// 与服务器握手
	if err := conn.HandshakeServer(); err != nil {
		conn.Close()
		logs.Error("handleConn HandshakeServer err: ", err)
		return err
	}

	connServer := core.NewConnServer(conn)

	// 读消息
	if err := connServer.ReadMsg(); err != nil {
		conn.Close()
		logs.Error("handleConn read msg err: ", err)
		return err
	}

	//取得信息..--> 检查appname
	appname, name, _ := connServer.GetInfo()

	if ret := configure.CheckAppName(appname); !ret {
		err := fmt.Errorf("application name=%s is not configured", appname)
		conn.Close()
		logs.Error("CheckAppName err: ", err)
		return err
	}

	logs.Debug("handleConn: IsPublisher=%v", connServer.IsPublisher())
	if connServer.IsPublisher() {
		//是发布
		channel, err := configure.RoomKeys.GetChannel(name)
		if err != nil {
			err := fmt.Errorf("invalid key")
			conn.Close()
			logs.Error("CheckKey err: ", err)
			return err
		}
		connServer.PublishInfo.Name = channel
		if pushlist, ret := configure.GetStaticPushUrlList(appname); ret && (pushlist != nil) {
			logs.Debug("GetStaticPushUrlList: %v", pushlist)
		}

		// 创建 Reader
		/*
			reader := core.NewVirReader(connServer)
			s.handler.HandleReader(reader)
			logs.Debug("new publisher: %+v", reader.Info())

				if s.getter != nil {
					writeType := reflect.TypeOf(s.getter)
					logs.Debug("handleConn:writeType=%v", writeType)
					writer := s.getter.GetWriter(reader.Info())
					s.handler.HandleWriter(writer)
				}

			flvWriter := new(flv.FlvDvr)
			s.handler.HandleWriter(flvWriter.GetWriter(reader.Info()))
		*/
	} else {
		// 定阅.
		writer := core.NewVirWriter(connServer)
		logs.Debug("new player: %+v", writer.Info())
		//s.handler.HandleWriter(writer)

	}

	return nil
}

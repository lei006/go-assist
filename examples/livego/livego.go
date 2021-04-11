package main

import (
	"fmt"
	"net"

	"github.com/lei006/go-assist/livego/configure"
	"github.com/lei006/go-assist/livego/protocol/hls"
	"github.com/lei006/go-assist/livego/protocol/httpflv"
	"github.com/lei006/go-assist/livego/protocol/rtmp"

	log "github.com/sirupsen/logrus"
)

type LiveServer struct {
	HlsPort     int
	HttpFlvPort int
	RtmpPort    int
}

func (this *LiveServer) Start(rtmp_port, httpflv_port, hls_port int) error {
	this.HlsPort = hls_port
	this.HttpFlvPort = httpflv_port
	this.RtmpPort = rtmp_port

	stream := rtmp.NewRtmpStream()
	hlsServer, err := this.startHls(hls_port)
	if err != nil {
		return err
	}

	err = this.startHTTPFlv(httpflv_port, stream)
	if err != nil {
		return err
	}

	err = this.startRtmp(rtmp_port, stream, hlsServer)
	if err != nil {
		return err
	}

	return nil
}

func (this *LiveServer) Stop() {

}

func (this *LiveServer) startHls(port int) (*hls.Server, error) {
	hlsAddr := configure.Config.GetString("hls_addr")
	hlsListen, err := net.Listen("tcp", hlsAddr)
	if err != nil {
		return nil, err
	}

	hlsServer := hls.NewServer()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("HLS server panic: ", r)
			}
		}()
		log.Info("HLS listen On ", hlsAddr)
		hlsServer.Serve(hlsListen)
	}()
	return hlsServer, nil
}

func (this *LiveServer) startRtmp(port int, stream *rtmp.RtmpStream, hlsServer *hls.Server) error {
	rtmpAddr := fmt.Sprintf(":%d", port)

	rtmpListen, err := net.Listen("tcp", rtmpAddr)
	if err != nil {
		return err
	}

	var rtmpServer *rtmp.Server

	if hlsServer == nil {
		rtmpServer = rtmp.NewRtmpServer(stream, nil)
		log.Info("HLS server disable....")
	} else {
		rtmpServer = rtmp.NewRtmpServer(stream, hlsServer)
		log.Info("HLS server enable....")
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("RTMP server panic: ", r)
		}
	}()
	log.Info("RTMP Listen On ", rtmpAddr)

	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	rtmpServer.Serve(rtmpListen)

	return nil
}

func (this *LiveServer) startHTTPFlv(port int, stream *rtmp.RtmpStream) error {
	httpflvAddr := fmt.Sprintf(":%d", port)

	flvListen, err := net.Listen("tcp", httpflvAddr)
	if err != nil {
		return err
	}

	hdlServer := httpflv.NewServer(stream)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("HTTP-FLV server panic: ", r)
			}
		}()
		log.Info("HTTP-FLV listen On ", httpflvAddr)
		hdlServer.Serve(flvListen)
	}()

	return nil
}

func main() {
	fmt.Println("xxxxxx")

	srv := LiveServer{}
	err := srv.Start(1370, 7001, 7002)
	fmt.Println("xx33xxxx", err)

}

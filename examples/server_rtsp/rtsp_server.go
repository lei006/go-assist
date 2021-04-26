package main

import (
	"github.com/lei006/go-assist/servers/livertc"
	"github.com/lei006/go-assist/servers/server_rtsp"
)

// ffmpeg.exe -f dshow -i video="USB2.0 PC CAMERA" -preset veryfast -deinterlace -rtsp_transport tcp -f rtsp rtsp://127.0.0.1/3.sdp
// ffmpeg.exe -f dshow -i video="USB2.0 PC CAMERA" -preset veryfast -deinterlace -rtsp_transport tcp -f rtsp rtsp://127.0.0.1/3.sdp
// ffmpeg.exe -r 25 -f dshow -i video="Video (00 Pro Capture Mini HDMI)" -vcodec libx264 -pix_fmt nv12 -s 1280x720 -bf 0 -preset veryfast -deinterlace -rtsp_transport tcp -f rtsp rtsp://127.0.0.1/3.sdp
// ffmpeg.exe -r 25 -f dshow -i video="Video (00 Pro Capture Mini HDMI)" -vcodec libx264 -pix_fmt nv12 -s 1280x720  -r 25 -bf 0  -preset veryfast -deinterlace -rtsp_transport tcp -f rtsp rtsp://127.0.0.1/3.sdp
func main() {
	engine := livertc.Default()

	//ffmpeg -re -i  rtmp://58.200.131.2:1935/livetv/cctv16 -vcodec libx264 -acodec aac  -f rtsp -rtsp_transport tcp rtsp://127.0.0.1:554/test.sdp
	rtsp := server_rtsp.MakeRtspServer(engine)
	rtsp.ListenAt(":554")
	rtsp.RunLoop()
}

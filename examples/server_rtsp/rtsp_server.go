package main

import (
	"bytes"
	"fmt"

	"github.com/lei006/go-assist/servers/datapacket"
	"github.com/lei006/go-assist/servers/server_rtsp"
	"github.com/lei006/go-assist/utils"
)

// ffmpeg.exe -f dshow -i video="USB2.0 PC CAMERA" -preset veryfast -deinterlace -rtsp_transport tcp -f rtsp rtsp://127.0.0.1/3.sdp
// ffmpeg.exe -f dshow -i video="USB2.0 PC CAMERA" -preset veryfast -deinterlace -rtsp_transport tcp -f rtsp rtsp://127.0.0.1/3.sdp
// ffmpeg.exe -r 25 -f dshow -i video="Video (00 Pro Capture Mini HDMI)" -vcodec libx264 -pix_fmt nv12 -s 1280x720  -r 25 -bf 0  -preset veryfast -deinterlace -rtsp_transport tcp -f rtsp rtsp://127.0.0.1/3.sdp
func main() {

	rtp_packet := datapacket.RtpPacket{}

	rtspServer := server_rtsp.GetServer()
	var index int
	rtspServer.PacketCallback(func(url string, buffer *bytes.Buffer) {

		data_buf := buffer.Bytes()

		/*
			rtp_packet, err := rtp_decoder.ParseToRTP(data_buf)
			if err != nil {
				fmt.Println(" ", err.Error())
				return
			}
		*/
		utils.PrintBin(data_buf, 18)

		err := rtp_packet.Unmarshal(data_buf)
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
			fmt.Println(" package not over ")
			return
		}

		//fmt.Println(" packet.GetType() =", packet.GetType())

		/*
			//fmt.Println("buf type =", rtp_packet.Payload[0]&0x1f)
			//utils.PrintBin(rtp_packet.Payload, 18)

			//fmt.Println("Marker =", packet.Marker, packet.Payload[0]&0x1f, packet.Payload[1], packet.Payload[2])
			//fmt.Println("PayloadType =", packet.PayloadType)

			data_len := len(rtp_packet.Payload)
			if data_len > 36 {
				//fmt.Println("PayloadType = ", packet.PayloadType, packet.SequenceNumber, "       len=", data_len, packet.Payload[:8])
				//fmt.Println("buffer = ", buffer.Bytes()[:36])
			}
		*/
		index++
	})

	rtspServer.ClientAsk(func(req *server_rtsp.Request) error {
		fmt.Println("url =========================  start", req.URL)
		fmt.Println("sdp :", req.Body)
		fmt.Println("url =========================  end")

		return nil
	})

	rtspServer.Start()

}

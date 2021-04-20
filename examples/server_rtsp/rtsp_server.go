package main

import (
	"bytes"
	"fmt"

	"github.com/lei006/go-assist/utils"
	"github.com/lei006/go-assist/servers/encoder"
	"github.com/lei006/go-assist/servers/server_rtsp"
)

// ffmpeg.exe -f dshow -i video="USB2.0 PC CAMERA" -preset veryfast -deinterlace -rtsp_transport tcp -f rtsp rtsp://127.0.0.1/3.sdp
// ffmpeg.exe -f dshow -i video="USB2.0 PC CAMERA" -preset veryfast -deinterlace -rtsp_transport tcp -f rtsp rtsp://127.0.0.1/3.sdp
func main() {

	rtspServer := server_rtsp.GetServer()
	var index int
	rtspServer.PacketCallback(func(url string, buffer *bytes.Buffer) {
		//fmt.Println("url =", url, index, buffer.Bytes()[:8])

		data_buf := buffer.Bytes();
		utils.PrintBin(data_buf, 18)



		decoder := encoder.RtpDecoder{}
		packet := decoder.ParseRTP(data_buf)
		if packet != nil {

			fmt.Println("Marker =", packet.Marker, packet.Payload[0]&0x1f, packet.Payload[1], packet.Payload[2])
			//fmt.Println("PayloadType =", packet.PayloadType)



			data_len := len(packet.Payload)
			if data_len > 36 {
				//fmt.Println("PayloadType = ", packet.PayloadType, packet.SequenceNumber, "       len=", data_len, packet.Payload[:8])
				//fmt.Println("buffer = ", buffer.Bytes()[:36])
			}
		}

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

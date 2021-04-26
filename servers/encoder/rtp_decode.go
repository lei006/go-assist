package encoder

import (
	"errors"

	"github.com/lei006/go-assist/servers/datapacket"
)

const (
	// 注意，一般情况下，AVC使用96，AAC使用97，HEVC使用98
	// 但是我还遇到过：
	// HEVC使用96
	// AVC使用105
	RTPPacketTypeAVCOrHEVC = 96
	RTPPacketTypeAAC       = 97
	RTPPacketTypeHEVC      = 98
)

type RtpDecoder struct {
	h264Parser H264Parser
}

//拆包
func (decoder *RtpDecoder) ParsePacket(rtpBytes []byte) (*datapacket.DataPacket, error) {

	//只拆包头..
	rtp_packet := &datapacket.RtpPacket{}
	err := rtp_packet.Header.Unmarshal(rtpBytes)
	if err != nil {
		return nil, err
	}

	//H264的PT值为96
	// PayloadType 可以去这儿查: https://blog.csdn.net/qq_40732350/article/details/88374707

	if rtp_packet.Header.PayloadType == RTPPacketTypeAVCOrHEVC {
		return decoder.parsePacket_H264(rtpBytes, rtp_packet.Header.Marker)
	} else if rtp_packet.Header.PayloadType == RTPPacketTypeAAC {
		return decoder.parsePacket_H265(rtpBytes)
	} else if rtp_packet.Header.PayloadType == RTPPacketTypeHEVC {
		return decoder.parsePacket_AAC(rtpBytes)
	}

	return nil, errors.New("rtp: Unsupported data types")
}

//  H264 在 rtp 中的打包形式
//  RTP header (12bytes)+ FU Indicator (1byte)  +  FU header(1 byte) + NALU payload

// 1、单个NAL包单元
// 12字节的RTP头后面的就是音视频数据，比较简单。一个封装单个NAL单元包到RTP的NAL单元流的RTP序号必须符合NAL单元的解码顺序。

// 2、FU-A的分片格式
//2、FU-A的分片格式
//数据比较大的H264视频包，被RTP分片发送。12字节的RTP头后面跟随的就是FU-A分片：
//FU indicator有以下格式：
//      +---------------+
//      |0|1|2|3|4|5|6|7|
//      +-+-+-+-+-+-+-+-+
//      |F|NRI|  Type   |
//      +---------------+
//   FU指示字节的类型域 Type=28表示FU-A。。NRI域的值必须根据分片NAL单元的NRI域的值设置。
//
//   FU header的格式如下：
//      +---------------+
//      |0|1|2|3|4|5|6|7|
//      +-+-+-+-+-+-+-+-+
//      |S|E|R|  Type   |
//      +---------------+
//   S: 1 bit
//   当设置成1,开始位指示分片NAL单元的开始。当跟随的FU荷载不是分片NAL单元荷载的开始，开始位设为0。
//   E: 1 bit
//   当设置成1, 结束位指示分片NAL单元的结束，即, 荷载的最后字节也是分片NAL单元的最后一个字节。当跟随的FU荷载不是分片NAL单元的最后分片,结束位设置为0。
//   R: 1 bit
//   保留位必须设置为0，接收者必须忽略该位。
//   Type: 5 bits
//   NAL单元荷载类型定义见下表
//
//表1.  单元类型以及荷载结构总结
//      Type   Packet      Type name
//      ---------------------------------------------------------
//      0      undefined                                    -
//      1-23   NAL unit    Single NAL unit packet per H.264
//      24     STAP-A     Single-time aggregation packet
//      25     STAP-B     Single-time aggregation packet
//      26     MTAP16    Multi-time aggregation packet
//      27     MTAP24    Multi-time aggregation packet
//      28     FU-A      Fragmentation unit
//      29     FU-B      Fragmentation unit
//      30-31  undefined

func (decoder *RtpDecoder) parsePacket_H264(payload []byte, is_first bool) (*datapacket.DataPacket, error) {

	//解码去H264数据包...

	decoder.h264Parser.UnmarshalRTP(payload)

	return nil, nil
}

func (decoder *RtpDecoder) parsePacket_H265(rtpBytes []byte) (*datapacket.DataPacket, error) {

	return nil, nil
}

func (decoder *RtpDecoder) parsePacket_AAC(rtpBytes []byte) (*datapacket.DataPacket, error) {

	return nil, nil
}

func (decoder *RtpDecoder) UnmarshalRTP(payload []byte) (*datapacket.DataPacket, error) {

	return nil, nil
}

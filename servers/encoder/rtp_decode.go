package encoder

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/lei006/go-assist/servers/datapacket"
)

const (
	RTP_FIXED_HEADER_LENGTH = 12 // rtp 固定头长度...
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
}

// RTP header (12bytes)+ FU Indicator (1byte)  +  FU header(1 byte) + NALU payload

func (decoder *RtpDecoder) ParseHead(rtpBytes []byte) (datapacket.RtpHeader, error) {

	firstByte := rtpBytes[0]  //第一个字节
	secondByte := rtpBytes[1] //第二个字节

	header := datapacket.RtpHeader{}
	header.Version = int(firstByte >> 6)                               //版本号
	header.Padding = (firstByte>>5)&1 == 1                             //填充标识，如果为1，则在尾部有额外的1个或多个8位组
	header.Extension = (firstByte>>4)&1 == 1                           //头部是否有扩展位
	header.CSRCCnt = int(firstByte & 0x0f)                             //CSRC计数器
	header.Marker = secondByte>>7 == 1                                 //标识：  视频的结束   音频的开始
	header.PayloadType = int(secondByte & 0x7f)                        // 负载的数据类型..区分音频流和视频流
	header.SequenceNumber = int(binary.BigEndian.Uint16(rtpBytes[2:])) //序号-随机
	//h264的采样率为90000HZ，因此时间戳的单位为1(秒)/90000，因此如果当前视频帧率为25fps，那时间戳间隔或者说增量应该为3600，如果帧率为30fps，则增量为3000，以此类推
	header.Timestamp = int(binary.BigEndian.Uint32(rtpBytes[4:])) //时间戳..采样率为90000HZ
	header.SSRC = int(binary.BigEndian.Uint32(rtpBytes[8:]))

	if header.Version != 2 {
		return header, errors.New("Version number required to 2")
	}

	offset := RTP_FIXED_HEADER_LENGTH
	end := len(rtpBytes)

	//跳过 csrc 区域
	if end-offset >= 4*header.CSRCCnt {
		offset += 4 * header.CSRCCnt
	}

	// 跳过头部扩展区
	if header.Extension && end-offset >= 4 {
		// 不知道为什么要加 + 2 明明是四个字节..
		extLen := 4 * int(binary.BigEndian.Uint16(rtpBytes[offset+2:]))
		offset += 4 //跳过长度区
		if end-offset >= extLen {
			offset += extLen //跳过 扩展区.
		}
	}

	// 忽略尾部填充区
	if header.Padding && end-offset > 0 {
		paddingLen := int(rtpBytes[end-1])
		if end-offset >= paddingLen {
			end -= paddingLen
		}
	}

	header.PayloadStart = offset
	header.PayloadEnd = end

	if header.PayloadEnd-header.PayloadStart < 1 {
		return header, errors.New("Valid data can not be empty")
	}

	return header, nil
}

// H264
// 解包数据--拆成 rtp包..
func (decoder *RtpDecoder) ParseToRTP(rtpBytes []byte) (*datapacket.RtpPacket, error) {
	if len(rtpBytes) < RTP_FIXED_HEADER_LENGTH {
		return nil, errors.New("RTP packets is too small")
	}

	header, err := decoder.ParseHead(rtpBytes)
	if err != nil {
		return nil, err
	}

	rtp_packet := &datapacket.RtpPacket{
		Header: header,
		Raw:    rtpBytes,
	}

	rtp_packet.Header = header

	return rtp_packet, nil
}

//拆包
func (decoder *RtpDecoder) ParsePacket(rtpBytes []byte) (*datapacket.DataPacket, error) {
	//拆包头..
	rtp_packet, err := decoder.ParseToRTP(rtpBytes)
	if err != nil {
		return nil, err
	}

	payload := rtp_packet.Payload()

	//H264的PT值为96
	// PayloadType 可以去这儿查: https://blog.csdn.net/qq_40732350/article/details/88374707

	if rtp_packet.Header.PayloadType == RTPPacketTypeAVCOrHEVC {
		return decoder.parsePacket_H264(payload, rtp_packet.Header.Marker)
	} else if rtp_packet.Header.PayloadType == RTPPacketTypeAAC {
		return decoder.parsePacket_H265(payload)
	} else if rtp_packet.Header.PayloadType == RTPPacketTypeHEVC {
		return decoder.parsePacket_AAC(payload)
	}

	return nil, errors.New("RTP Unsupported data types")
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

	//注意长度
	if len(payload) < 2 {
		return nil, errors.New("RTP H264 packets is too small")
	}

	//NALU头的计算..
	FU_Indicator := payload[0]
	FU_Header := payload[1]

	_type := FU_Indicator & 0x1f
	if _type >= 1 && _type <= 23 {
		_type = 1
	}

	switch _type {
	case 0, 1:
		//单一NAL单元包/
		fmt.Println("1-23   NAL unit    Single NAL unit packet per H.264")
	case 24:
		//一次聚合包
		fmt.Println("24     STAP-A     Single-time aggregation packet")
	case 28:

		packet_start := (FU_Header&0x80 != 0)
		packet_end := (FU_Header&0x40 != 0)

		//碎片化单元
		fmt.Println("28     FU-A      Fragmentation unit", packet_start, packet_end)
	default:
		// 其它情况不处理....
		return nil, errors.New("RTP H264 packets is too small")
	}

	if is_first {
		//fmt.Println("h264 first nalu_type =", _type, FU_Header, len(payload))
	} else {
		//fmt.Println("h264       nalu_type =", _type, FU_Header, len(payload))
	}

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

package encoder

import (
	"encoding/binary"

	"github.com/lei006/go-assist/servers/datapacket"
)

const (
	RTP_FIXED_HEADER_LENGTH = 12 // rtp 固定头长度...
)

type RtpDecoder struct {

}


// 解包数据--拆成 rtp包..
func (decoder *RtpDecoder) ParseRTP(rtpBytes []byte) *datapacket.RtpPacket {
	if len(rtpBytes) < RTP_FIXED_HEADER_LENGTH {
		return nil
	}
	
	firstByte := rtpBytes[0]
	secondByte := rtpBytes[1]
	info := &datapacket.RtpPacket{
		Version:   int(firstByte >> 6),		//版本号
		Padding:   (firstByte>>5)&1 == 1,	//填充标识，如果为1，则在尾部有额外的1个或多个8位组
		Extension: (firstByte>>4)&1 == 1,	//头部是否有扩展位
		CSRCCnt:   int(firstByte & 0x0f),   //CSRC计数器
		Marker:         secondByte>>7 == 1,				//标识：  视频的结束   音频的开始
		PayloadType:    int(secondByte & 0x7f),	// 负载的数据类型..区分音频流和视频流
		SequenceNumber: int(binary.BigEndian.Uint16(rtpBytes[2:])), //序号
		Timestamp:      int(binary.BigEndian.Uint32(rtpBytes[4:])),
		SSRC:           int(binary.BigEndian.Uint32(rtpBytes[8:])),
	}

	offset := RTP_FIXED_HEADER_LENGTH
	end := len(rtpBytes)

	//跳过 csrc 区域
	if end-offset >= 4*info.CSRCCnt {
		offset += 4 * info.CSRCCnt
	}

	// 跳过头部扩展区
	if info.Extension && end-offset >= 4 {
		// 不知道为什么要加 + 2 明明是四个字节..
		extLen := 4 * int(binary.BigEndian.Uint16(rtpBytes[offset+2:]))
		offset += 4	//跳过长度区
		if end-offset >= extLen {
			offset += extLen	//跳过 扩展区.
		}
	}

	// 忽略尾部填充区
	if info.Padding && end-offset > 0 {
		paddingLen := int(rtpBytes[end-1])
		if end-offset >= paddingLen {
			end -= paddingLen
		}
	}
	// 这才是真正的有效数据...
	info.Payload = rtpBytes[offset:end]
	info.PayloadOffset = offset
	if end-offset < 1 {
		return nil
	}

	return info
}


func (decoder *RtpDecoder)UnmarshalRTP(payload []byte) (*DataPacket, error) {
		

	return nil,nil;
}







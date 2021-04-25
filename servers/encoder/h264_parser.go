package encoder

import (
	"bytes"
	"fmt"
	"io"

	h264_parser "github.com/deepch/vdk/codec/h264parser"

	"github.com/deepch/vdk/utils/bits/pio"
)

const (
	i_frame byte = 0
	p_frame byte = 1
	b_frame byte = 2
)

const (
	nalu_type_not_define byte = 0
	nalu_type_slice      byte = 1  //slice_layer_without_partioning_rbsp() sliceheader
	nalu_type_dpa        byte = 2  // slice_data_partition_a_layer_rbsp( ), slice_header
	nalu_type_dpb        byte = 3  // slice_data_partition_b_layer_rbsp( )
	nalu_type_dpc        byte = 4  // slice_data_partition_c_layer_rbsp( )
	nalu_type_idr        byte = 5  // slice_layer_without_partitioning_rbsp( ),sliceheader
	nalu_type_sei        byte = 6  //sei_rbsp( )
	nalu_type_sps        byte = 7  //seq_parameter_set_rbsp( )
	nalu_type_pps        byte = 8  //pic_parameter_set_rbsp( )
	nalu_type_aud        byte = 9  // access_unit_delimiter_rbsp( )
	nalu_type_eoesq      byte = 10 //end_of_seq_rbsp( )
	nalu_type_eostream   byte = 11 //end_of_stream_rbsp( )
	nalu_type_filler     byte = 12 //filler_data_rbsp( )
)

const (
	naluBytesLen int = 4
	maxSpsPpsLen int = 2 * 1024
)

var (
	decDataNil        = fmt.Errorf("dec buf is nil")
	spsDataError      = fmt.Errorf("sps data error")
	ppsHeaderError    = fmt.Errorf("pps header error")
	ppsDataError      = fmt.Errorf("pps data error")
	naluHeaderInvalid = fmt.Errorf("nalu header invalid")
	videoDataInvalid  = fmt.Errorf("video data not match")
	dataSizeNotMatch  = fmt.Errorf("data size not match")
	naluBodyLenError  = fmt.Errorf("nalu body len error")
)

var startCode = []byte{0x00, 0x00, 0x00, 0x01}
var naluAud = []byte{0x00, 0x00, 0x00, 0x01, 0x09, 0xf0}

type H264Parser struct {
	frameType    byte
	specificInfo []byte
	pps          *bytes.Buffer

	IsKeyFrame bool
	Payload    []byte
	fuStarted  bool
	gotpkt     bool

	fuBuffer []byte

	SPS []byte
	PPS []byte
}

type sequenceHeader struct {
	configVersion        byte //8bits
	avcProfileIndication byte //8bits
	profileCompatility   byte //8bits
	avcLevelIndication   byte //8bits
	reserved1            byte //6bits
	naluLen              byte //2bits
	reserved2            byte //3bits
	spsNum               byte //5bits
	ppsNum               byte //8bits
	spsLen               int
	ppsLen               int
}

func NewH264Parser() *H264Parser {
	return &H264Parser{
		pps: bytes.NewBuffer(make([]byte, maxSpsPpsLen)),
	}
}

//return value 1:sps, value2 :pps
func (H264Parser *H264Parser) parseSpecificInfo(src []byte) error {
	if len(src) < 9 {
		return decDataNil
	}
	sps := []byte{}
	pps := []byte{}

	var seq sequenceHeader
	seq.configVersion = src[0]
	seq.avcProfileIndication = src[1]
	seq.profileCompatility = src[2]
	seq.avcLevelIndication = src[3]
	seq.reserved1 = src[4] & 0xfc
	seq.naluLen = src[4]&0x03 + 1
	seq.reserved2 = src[5] >> 5

	//get sps
	seq.spsNum = src[5] & 0x1f
	seq.spsLen = int(src[6])<<8 | int(src[7])

	if len(src[8:]) < seq.spsLen || seq.spsLen <= 0 {
		return spsDataError
	}
	sps = append(sps, startCode...)
	sps = append(sps, src[8:(8+seq.spsLen)]...)

	//get pps
	tmpBuf := src[(8 + seq.spsLen):]
	if len(tmpBuf) < 4 {
		return ppsHeaderError
	}
	seq.ppsNum = tmpBuf[0]
	seq.ppsLen = int(0)<<16 | int(tmpBuf[1])<<8 | int(tmpBuf[2])
	if len(tmpBuf[3:]) < seq.ppsLen || seq.ppsLen <= 0 {
		return ppsDataError
	}

	pps = append(pps, startCode...)
	pps = append(pps, tmpBuf[3:]...)

	H264Parser.specificInfo = append(H264Parser.specificInfo, sps...)
	H264Parser.specificInfo = append(H264Parser.specificInfo, pps...)

	return nil
}

func (H264Parser *H264Parser) isNaluHeader(src []byte) bool {
	if len(src) < naluBytesLen {
		return false
	}
	return src[0] == 0x00 &&
		src[1] == 0x00 &&
		src[2] == 0x00 &&
		src[3] == 0x01
}

func (H264Parser *H264Parser) naluSize(src []byte) (int, error) {
	if len(src) < naluBytesLen {
		return 0, fmt.Errorf("nalusizedata invalid")
	}
	buf := src[:naluBytesLen]
	size := int(0)
	for i := 0; i < len(buf); i++ {
		size = size<<8 + int(buf[i])
	}
	return size, nil
}

func (H264Parser *H264Parser) getAnnexbH264(src []byte, w io.Writer) error {
	dataSize := len(src)
	if dataSize < naluBytesLen {
		return videoDataInvalid
	}
	H264Parser.pps.Reset()
	_, err := w.Write(naluAud)
	if err != nil {
		return err
	}

	index := 0
	nalLen := 0
	hasSpsPps := false
	hasWriteSpsPps := false

	for dataSize > 0 {
		nalLen, err = H264Parser.naluSize(src[index:])
		if err != nil {
			return dataSizeNotMatch
		}
		index += naluBytesLen
		dataSize -= naluBytesLen
		if dataSize >= nalLen && len(src[index:]) >= nalLen && nalLen > 0 {
			nalType := src[index] & 0x1f
			switch nalType {
			case nalu_type_aud:
			case nalu_type_idr:
				if !hasWriteSpsPps {
					hasWriteSpsPps = true
					if !hasSpsPps {
						if _, err := w.Write(H264Parser.specificInfo); err != nil {
							return err
						}
					} else {
						if _, err := w.Write(H264Parser.pps.Bytes()); err != nil {
							return err
						}
					}
				}
				fallthrough
			case nalu_type_slice:
				fallthrough
			case nalu_type_sei:
				_, err := w.Write(startCode)
				if err != nil {
					return err
				}
				_, err = w.Write(src[index : index+nalLen])
				if err != nil {
					return err
				}
			case nalu_type_sps:
				fallthrough
			case nalu_type_pps:
				hasSpsPps = true
				_, err := H264Parser.pps.Write(startCode)
				if err != nil {
					return err
				}
				_, err = H264Parser.pps.Write(src[index : index+nalLen])
				if err != nil {
					return err
				}
			}
			index += nalLen
			dataSize -= nalLen
		} else {
			return naluBodyLenError
		}
	}
	return nil
}

func (H264Parser *H264Parser) Parse(b []byte, isSeq bool, w io.Writer) (err error) {
	switch isSeq {
	case true:
		err = H264Parser.parseSpecificInfo(b)
	case false:
		// is annexb
		if H264Parser.isNaluHeader(b) {
			_, err = w.Write(b)
		} else {
			err = H264Parser.getAnnexbH264(b, w)
		}
	}
	return
}

/*
// 把多个 nal 粘成一个 nal
//func (H264Parser *H264Parser) UnmarshalRTP(rtpPacket *pion_rtp.Packet) (*core.DataPacket, error) {
func (H264Parser *H264Parser) UnmarshalRTP(payload []byte) (*datapacket.H264Packet, error) {

	rtpPacket := &pion_rtp.Packet{}
	if err := rtpPacket.Unmarshal(payload); err != nil {
		//出错了.........
		return nil, errors.New("不是有效Nal")
	}

	err := H264Parser.handleH264Payload(0, rtpPacket.Payload)
	if H264Parser.gotpkt == true || H264Parser.IsKeyFrame == true {
		//utils.PrintBinStr(H264Parser.Payload, 40, "par ")
		newPack := datapacket.MakeH264Packet(H264Parser.Payload[4:])

		H264Parser.gotpkt = false
		H264Parser.IsKeyFrame = false
		H264Parser.Payload = []byte{}

		return newPack, nil
	} else {
		return nil, err
	}
}
*/
func (H264Parser *H264Parser) PrintLog() {
	if H264Parser.gotpkt {

		if H264Parser.IsKeyFrame {
			fmt.Println("gotpkt = ", H264Parser.gotpkt, H264Parser.IsKeyFrame, len(H264Parser.Payload), "=======-------------")
		} else {
			fmt.Println("gotpkt = ", H264Parser.gotpkt, H264Parser.IsKeyFrame, len(H264Parser.Payload))
		}
	}

}

func (self *H264Parser) handleBuggyAnnexbH264Packet(timestamp uint32, packet []byte) (isBuggy bool, err error) {
	if len(packet) >= 4 && packet[0] == 0 && packet[1] == 0 && packet[2] == 0 && packet[3] == 1 {
		isBuggy = true
		if nalus, typ := h264_parser.SplitNALUs(packet); typ != h264_parser.NALU_RAW {
			for _, nalu := range nalus {
				if len(nalu) > 0 {
					if err = self.handleH264Payload(timestamp, nalu); err != nil {
						return
					}
				}
			}
		}
	}
	return
}

func (self *H264Parser) handleH264Payload(timestamp uint32, packet []byte) (err error) {
	if len(packet) < 2 {
		err = fmt.Errorf("rtp: h264 packet too short")
		return
	}

	var isBuggy bool
	if isBuggy, err = self.handleBuggyAnnexbH264Packet(timestamp, packet); isBuggy {
		return
	}

	naluType := packet[0] & 0x1f

	switch {
	case naluType >= 1 && naluType <= 5:
		if naluType == 5 {
			self.IsKeyFrame = true
		}
		self.gotpkt = true
		// raw nalu to avcc
		b := make([]byte, 4+len(packet))
		pio.PutU32BE(b[0:4], uint32(len(packet)))
		copy(b[4:], packet)
		self.Payload = b

	case naluType == 7: // sps
		//utils.PrintBinStr(packet, 30, "RTSP-SPS")
		self.SPS = packet
		break

	case naluType == 8: // pps
		//utils.PrintBinStr(packet, 30, "RTSP-PPS")
		self.PPS = packet
		break
	case naluType == 28: // FU-A 处理 h264碎片

		fuIndicator := packet[0]
		fuHeader := packet[1]
		isStart := fuHeader&0x80 != 0
		isEnd := fuHeader&0x40 != 0
		if isStart {
			self.fuStarted = true
			self.fuBuffer = []byte{fuIndicator&0xe0 | fuHeader&0x1f}
		}
		if self.fuStarted {
			self.fuBuffer = append(self.fuBuffer, packet[2:]...)
			if isEnd {
				self.fuStarted = false
				if err = self.handleH264Payload(timestamp, self.fuBuffer); err != nil {
					return
				}
			}
		}

	case naluType == 24: // STAP-A   一次聚合包
		/*
			0                   1                   2                   3
			0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
			+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			|                          RTP Header                           |
			+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			|STAP-A NAL HDR |         NALU 1 Size           | NALU 1 HDR    |
			+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			|                         NALU 1 Data                           |
			:                                                               :
			+               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			|               | NALU 2 Size                   | NALU 2 HDR    |
			+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			|                         NALU 2 Data                           |
			:                                                               :
			|                               +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			|                               :...OPTIONAL RTP padding        |
			+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

			Figure 7.  An example of an RTP packet including an STAP-A
			containing two single-time aggregation units
		*/
		packet = packet[1:]
		for len(packet) >= 2 {
			size := int(packet[0])<<8 | int(packet[1])
			if size+2 > len(packet) {
				break
			}
			if err = self.handleH264Payload(timestamp, packet[2:size+2]); err != nil {
				return
			}
			packet = packet[size+2:]
		}
		return

	case naluType >= 6 && naluType <= 23: // other single NALU packet
		/*
			fuIndicator := packet[0]
			fuHeader := packet[1]
			self.fuBuffer = []byte{fuIndicator&0xe0 | fuHeader&0x1f}
			//packet[1] = []byte{fuIndicator&0xe0 | fuHeader&0x1f}
			self.Payload = packet[1:]

			return
		*/

	case naluType == 25: // STAB-B
	case naluType == 26: // MTAP-16
	case naluType == 27: // MTAP-24
	default:
		err = fmt.Errorf("rtsp: unsupported H264 naluType=%d", naluType)
		return
	}

	return
}

package flv

import (
	"fmt"
	"go-assist/protocol/intfs"
)

var (
	ErrAvcEndSEQ = fmt.Errorf("avc end sequence")
)

type Demuxer struct {
}

func NewDemuxer() *Demuxer {
	return &Demuxer{}
}

func (d *Demuxer) DemuxH(p *intfs.Packet) error {
	var tag Tag
	_, err := tag.ParseMediaTagHeader(p.Data, p.IsVideo)
	if err != nil {
		return err
	}
	p.Header = &tag

	return nil
}

func (d *Demuxer) Demux(p *intfs.Packet) error {
	var tag Tag
	n, err := tag.ParseMediaTagHeader(p.Data, p.IsVideo)
	if err != nil {
		return err
	}
	if tag.CodecID() == intfs.VIDEO_H264 &&
		p.Data[0] == 0x17 && p.Data[1] == 0x02 {
		return ErrAvcEndSEQ
	}
	p.Header = &tag
	p.Data = p.Data[n:]

	return nil
}

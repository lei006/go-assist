package cache

import (
	"bytes"

	"github.com/lei006/go-assist/protocol/intfs"

	"github.com/beego/beego/v2/core/logs"

	"github.com/lei006/go-assist/livego/protocol/amf"
)

const (
	SetDataFrame string = "@setDataFrame"
	OnMetaData   string = "onMetaData"
)

var setFrameFrame []byte

func init() {

	b := bytes.NewBuffer(nil)
	encoder := &amf.Encoder{}
	if _, err := encoder.Encode(b, SetDataFrame, amf.AMF0); err != nil {
		logs.Critical(err)
	}
	setFrameFrame = b.Bytes()
}

type SpecialCache struct {
	full bool
	p    *intfs.Packet
}

func NewSpecialCache() *SpecialCache {
	return &SpecialCache{}
}

func (specialCache *SpecialCache) Write(p *intfs.Packet) {
	specialCache.p = p
	specialCache.full = true
}

func (specialCache *SpecialCache) Send(w intfs.WriteCloser) error {
	if !specialCache.full {
		return nil
	}
	return w.Write(specialCache.p)
}

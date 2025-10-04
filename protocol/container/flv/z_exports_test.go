package flv

import "testing"

func TestExports(t *testing.T) {
  _ = NewDemuxer
  _ = NewFlv
  _ = NewFLVWriter
}

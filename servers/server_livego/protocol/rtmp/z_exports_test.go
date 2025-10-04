package rtmp

import "testing"

func TestExports(t *testing.T) {
  _ = NewRtmpClient
  _ = NewRtmpServer
  _ = NewVirWriter
  _ = NewVirReader
  _ = NewRtmpStream
  _ = NewStream
}

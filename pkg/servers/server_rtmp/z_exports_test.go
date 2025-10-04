package server_rtmp

import "testing"

func TestExports(t *testing.T) {
  _ = MakeHandlePubProxy
  _ = MakeRtmpServer
}

package server_rtsp

import "testing"

func TestExports(t *testing.T) {
  _ = MakeHandlePubProxy
  _ = MakeRtspServer
}

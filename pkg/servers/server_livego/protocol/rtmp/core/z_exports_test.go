package core

import "testing"

func TestExports(t *testing.T) {
  _ = NewConn
  _ = initControlMsg
  _ = NewConnClient
  _ = NewConnServer
  _ = hsMakeDigest
  _ = hsCalcDigestPos
  _ = hsFindDigest
  _ = hsParse1
  _ = hsCreate01
  _ = hsCreate2
  _ = NewReadWriter
}

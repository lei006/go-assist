package httpflv

import "testing"

func TestExports(t *testing.T) {
  _ = NewServer
  _ = NewFLVWriter
}

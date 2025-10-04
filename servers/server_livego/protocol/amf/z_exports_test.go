package amf

import "testing"

func TestExports(t *testing.T) {
  _ = NewDecoder
  _ = NewTrait
  _ = NewTypedObject
  _ = readFlags
  _ = init
  _ = MetaDataReform
  _ = DumpBytes
  _ = Dump
  _ = WriteByte
  _ = WriteBytes
  _ = ReadByte
  _ = ReadBytes
  _ = WriteMarker
  _ = ReadMarker
  _ = AssertMarker
}

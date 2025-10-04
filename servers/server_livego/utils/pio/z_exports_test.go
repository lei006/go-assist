package pio

import "testing"

func TestExports(t *testing.T) {
  _ = U8
  _ = U16BE
  _ = I16BE
  _ = I24BE
  _ = U24BE
  _ = I32BE
  _ = U32LE
  _ = U32BE
  _ = U40BE
  _ = U64BE
  _ = I64BE
  _ = PutU8
  _ = PutI16BE
  _ = PutU16BE
  _ = PutI24BE
  _ = PutU24BE
  _ = PutI32BE
  _ = PutU32BE
  _ = PutU32LE
  _ = PutU40BE
  _ = PutU48BE
  _ = PutU64BE
  _ = PutI64BE
}

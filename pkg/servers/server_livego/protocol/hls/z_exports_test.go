package hls

import "testing"

func TestExports(t *testing.T) {
  _ = newAudioCache
  _ = NewTSCacheItem
  _ = NewServer
  _ = NewTSItem
  _ = NewSource
  _ = newStatus
}

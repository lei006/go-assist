package rtmprelay

import "testing"

func TestExports(t *testing.T) {
  _ = NewRtmpRelay
  _ = GetStaticPushList
  _ = GetAndCreateStaticPushObject
  _ = GetStaticPushObject
  _ = ReleaseStaticPushObject
  _ = NewStaticPush
}

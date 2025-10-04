package configure

import "testing"

func TestExports(t *testing.T) {
  _ = Init
  _ = initLog
  _ = init
  _ = CheckAppName
  _ = GetStaticPushUrlList
}

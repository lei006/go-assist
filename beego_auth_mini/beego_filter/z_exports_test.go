package beego_filter

import "testing"

func TestExports(t *testing.T) {
  _ = FilterCheckDomain
  _ = Add_CheckAdmin
  _ = FilterPremissions
  _ = FilterCheckToken
  _ = Add_NoCheckToken
}

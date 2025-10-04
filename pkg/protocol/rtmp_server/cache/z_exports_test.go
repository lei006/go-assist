package cache

import "testing"

func TestExports(t *testing.T) {
  _ = NewCache
  _ = newArray
  _ = NewGopCache
  _ = init
  _ = NewSpecialCache
}

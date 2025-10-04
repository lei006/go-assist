package api

import "testing"

func TestExports(t *testing.T) {
  _ = NewServer
  _ = JWTMiddleware
}

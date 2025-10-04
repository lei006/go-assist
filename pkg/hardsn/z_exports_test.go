package hardsn

import "testing"

func TestExports(t *testing.T) {
  _ = GetPhysicalID
  _ = getMACAddress
  _ = getCPUInfo
  _ = wmiQuery
  _ = getMachineGuid
  _ = getMd5String
}

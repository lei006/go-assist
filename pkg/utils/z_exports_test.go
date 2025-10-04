package utils

import "testing"

func TestExports(t *testing.T) {
  _ = ToJsonString
  _ = GetCurrentDirectory
  _ = PathMustExist
  _ = Is_Go_Run_Mod
  _ = Is_RunAtVs
  _ = GetBinPath
  _ = GetWorkPath
  _ = GetExePath
  _ = PrintRTP
  _ = PrintH264
  _ = PrintBin
  _ = PrintBinStr
  _ = PrintfBin
  _ = init
  _ = RandNumber
  _ = RandomString
  _ = RandomString1
}

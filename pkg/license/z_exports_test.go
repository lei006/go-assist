package license

import "testing"

func TestExports(t *testing.T) {
  _ = NewValidationError
  _ = verifyUser
  _ = verifyExp
  _ = verifyIat
  _ = verifyNbf
  _ = verifyIss
  _ = TestLicense
  _ = KeyEncryption
  _ = KeyDecryption
  _ = Decryption
  _ = MakeLicenseKey
  _ = MakEccP521Key
  _ = generateECCKeyString
  _ = init
  _ = RandString
}

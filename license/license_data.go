package license

import (
	jwt "github.com/golang-jwt/jwt/v4"
)

type LicenseData struct {
	jwt.StandardClaims
}

//签名一个数具
func (data *LicenseData) Sign(key string) (string, error) {

	//data.si

	return "", nil
}

//检验
func (data *LicenseData) Verify(sign string, key string) (bool, error) {

	return false, nil
}

package license

import (
	"encoding/json"

	jwt "github.com/golang-jwt/jwt/v4"
)

type LicenseClaims struct {
	jwt.StandardClaims
	Number int64 `json:"num,omitempty"`
}

func (lic_claims *LicenseClaims) ToCompare(tmp_lic_claims *LicenseClaims) bool {

	if lic_claims.Number != tmp_lic_claims.Number {
		return false
	}
	if lic_claims.Audience != tmp_lic_claims.Audience {
		return false
	}
	if lic_claims.ExpiresAt != tmp_lic_claims.ExpiresAt {
		return false
	}
	if lic_claims.Id != tmp_lic_claims.Id {
		return false
	}
	if lic_claims.IssuedAt != tmp_lic_claims.IssuedAt {
		return false
	}
	if lic_claims.Issuer != tmp_lic_claims.Issuer {
		return false
	}
	if lic_claims.NotBefore != tmp_lic_claims.NotBefore {
		return false
	}
	if lic_claims.Subject != tmp_lic_claims.Subject {
		return false
	}

	return true
}

//签名一个数具
func (data *LicenseClaims) ToJson() (string, error) {

	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (lic_data *LicenseClaims) FromJson(data_str string) error {

	data := []byte(data_str)
	err := json.Unmarshal(data, lic_data)
	return err

}

//签名一个数具
func (data *LicenseClaims) Sign(key LicenseKey) (string, error) {

	//data.si

	return "", nil
}

//检验
func (data *LicenseClaims) Verify(sign string, key string) (bool, error) {

	return false, nil
}

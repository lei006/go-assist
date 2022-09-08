package license

import (
	"encoding/json"

	jwt "github.com/golang-jwt/jwt/v4"
)

type LicenseClaims struct {
	jwt.StandardClaims
	Number int64 `json:"num,omitempty"`
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

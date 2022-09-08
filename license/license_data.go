package license

import (
	"encoding/json"
)

type licenseData struct {
	Claims string `json:"claims"`
	Sign   string `json:"sign"`
	PubKey string `json:"pub_key"`
}

//签名一个数具
func (data *licenseData) ToJson() (string, error) {

	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (lic_data *licenseData) FromJson(data_str string) error {

	data := []byte(data_str)
	err := json.Unmarshal(data, lic_data)
	return err

}

func (lic_data *licenseData) Verify() (bool, error) {

	//
	key := &LicenseKey{
		PubKey: lic_data.PubKey,
	}

	verify, err := key.Verify(lic_data.Claims, lic_data.Sign)

	return verify, err
}

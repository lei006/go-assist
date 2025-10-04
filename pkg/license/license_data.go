package license

import (
	"encoding/base64"
	"encoding/json"
)

type licenseData struct {
	Claims string `json:"claims"`
	Sign   string `json:"sign"`
	PubKey string `json:"pub_key"`
}

//签名一个数具
func (lic_data *licenseData) toJson() (string, error) {

	b, err := json.Marshal(lic_data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (lic_data *licenseData) fromJson(data_str string) error {

	data := []byte(data_str)
	err := json.Unmarshal(data, lic_data)
	return err

}

func (lic_data *licenseData) ToString() (string, error) {

	str_json, err := lic_data.toJson()
	if err != nil {
		return "", err
	}

	base64_str := base64.StdEncoding.EncodeToString([]byte(str_json))

	return base64_str, nil
}

func (lic_data *licenseData) FromString(data_str string) error {

	byte_data, err := base64.StdEncoding.DecodeString(data_str)
	if err != nil {
		return err
	}
	lic_data.fromJson(string(byte_data))
	return nil
}

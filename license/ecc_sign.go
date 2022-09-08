package license

import "encoding/json"

type EccSign struct {
	Rtext string `json:"r_text"`
	Stext string `json:"s_text"`
}

//签名一个数具
func (sign *EccSign) ToJson() (string, error) {

	b, err := json.Marshal(sign)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (sign *EccSign) FromJson(data_str string) error {

	data := []byte(data_str)
	err := json.Unmarshal(data, sign)
	return err
}

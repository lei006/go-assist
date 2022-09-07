package license

import "encoding/json"

type LicenseKey struct {
	PubKey string `json:"pub_key"` //公钥
	PriKey string `json:"pri_key"` //私钥
}

func (key *LicenseKey) ToString() string {
	return "pub:" + key.PubKey + ",pri:" + key.PriKey
}

func (key *LicenseKey) ToJson() (string, error) {

	b, err := json.Marshal(key)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

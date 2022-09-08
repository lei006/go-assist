package license

import (
	"errors"
)

func Sign(data *licenseData, key *LicenseKey) {

}

func TestLicense() (bool, error) {

	return false, nil
}

func Encryption(data LicenseClaims, key *LicenseKey) (string, error) {

	json_str, err := data.ToJson()
	if err != nil {
		return "", err
	}

	//生成数字签名
	sign_str, err := key.Sign(json_str)
	if err != nil {
		return "", err
	}

	lic_data := &licenseData{}
	lic_data.Claims = json_str
	//lic_data.Sign = base64.StdEncoding.EncodeToString([]byte(sign_str))
	lic_data.Sign = sign_str
	lic_data.PubKey = key.PubKey

	lic_data_json, err := lic_data.ToJson()
	if err != nil {
		return "", err
	}

	return lic_data_json, nil
}

func Decryption(enc_data string, key *LicenseKey) (*LicenseClaims, error) {

	lic_data := &licenseData{}
	err := lic_data.FromJson(enc_data)
	if err != nil {
		return nil, err
	}

	verify, err := lic_data.Verify()
	if err != nil {
		return nil, err
	}

	if verify == false {
		return nil, errors.New("verify sign error")
	}

	lic_claims := &LicenseClaims{}
	err = lic_claims.FromJson(lic_data.Claims)
	if err != nil {
		return nil, errors.New("json to claims error:" + err.Error())
	}

	return lic_claims, nil

}

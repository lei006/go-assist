package license

import (
	"errors"
	"fmt"
)

func TestLicense() (bool, error) {

	key, err := MakEccP521Key()
	if err != nil {
		return false, errors.New("生成KEY错误:" + err.Error())
	}

	enc_data := &LicenseClaims{}
	enc_data.StandardClaims.Id = RandString(10)

	lic_data, err := KeyEncryption(enc_data, key)
	if err != nil {
		return false, errors.New("license Encryption error:" + err.Error())
	}
	//fmt.Println("lic_data=", lic_data)
	dec_data, err := Decryption(lic_data)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return false, errors.New("license Decryption error:" + err.Error())
	}
	if enc_data.StandardClaims.Id != dec_data.StandardClaims.Id {
		return false, nil
	}
	return true, nil

	/*
		dec_data, err := KeyDecryption(lic_data, key)
		if err != nil {
			fmt.Println("Decryption error:", err)
			return false, errors.New("license Decryption error:" + err.Error())
		}
		if enc_data.StandardClaims.Id != dec_data.StandardClaims.Id {
			return false, nil
		}
	*/
	return true, nil
}

func KeyEncryption(data *LicenseClaims, key *LicenseKey) (string, error) {

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
	lic_data.PubKey = key.GetBase64PubKey()

	lic_data_json, err := lic_data.ToString()
	if err != nil {
		return "", err
	}

	return lic_data_json, nil
}

func KeyDecryption(enc_data string, key *LicenseKey) (*LicenseClaims, error) {

	lic_data := &licenseData{}
	err := lic_data.FromString(enc_data)
	if err != nil {
		return nil, err
	}

	verify, err := key.VerifySign(lic_data.Claims, lic_data.Sign)
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

func Decryption(license_data string) (*LicenseClaims, error) {

	lic_data := &licenseData{}
	err := lic_data.FromString(license_data)
	if err != nil {
		return nil, errors.New("json data error:" + err.Error())
	}

	key := &LicenseKey{}
	key.SetBase64PubKey(lic_data.PubKey)

	claims, err := KeyDecryption(license_data, key)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

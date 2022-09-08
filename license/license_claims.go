package license

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"time"
)

var TimeFunc = time.Now

type LicenseClaims struct {
	User      string `json:"user,omitempty"`       //使用者
	ExpiresAt int64  `json:"expires_at,omitempty"` //过期时间
	AppName   string `json:"app_name,omitempty"`   //程序名
	IssuedAt  int64  `json:"issuer_at,omitempty"`  //发行时间
	Issuer    string `json:"issuer,omitempty"`     //发行人
	Hardsn    string `json:"hardsn,omitempty"`     //硬件ID
	NotBefore int64  `json:"not_before,omitempty"` //生效时间
	Subject   string `json:"subject,omitempty"`    //标题
	Number    int64  `json:"number,omitempty"`     //发行数量
}

func (lic_claims *LicenseClaims) ToCompare(tmp_lic_claims *LicenseClaims) bool {

	if lic_claims.Number != tmp_lic_claims.Number {
		return false
	}
	if lic_claims.User != tmp_lic_claims.User {
		return false
	}
	if lic_claims.ExpiresAt != tmp_lic_claims.ExpiresAt {
		return false
	}
	if lic_claims.AppName != tmp_lic_claims.AppName {
		return false
	}
	if lic_claims.IssuedAt != tmp_lic_claims.IssuedAt {
		return false
	}
	if lic_claims.Issuer != tmp_lic_claims.Issuer {
		return false
	}
	if lic_claims.Hardsn != tmp_lic_claims.Hardsn {
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

// Valid validates time based claims "exp, iat, nbf". There is no accounting for clock skew.
// As well, if any of the above claims are not in the token, it will still
// be considered a valid claim.
func (c LicenseClaims) Valid() error {
	vErr := new(ValidationError)
	now := TimeFunc().Unix()

	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if !c.VerifyExpiresAt(now, false) {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("%s by %s", ErrTokenExpired, delta)
		vErr.Errors |= ValidationErrorExpired
	}

	if !c.VerifyIssuedAt(now, false) {
		vErr.Inner = ErrTokenUsedBeforeIssued
		vErr.Errors |= ValidationErrorIssuedAt
	}

	if !c.VerifyNotBefore(now, false) {
		vErr.Inner = ErrTokenNotValidYet
		vErr.Errors |= ValidationErrorNotValidYet
	}

	if vErr.valid() {
		return nil
	}

	return vErr
}

// VerifyAudience compares the aud claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *LicenseClaims) VerifyUser(cmp string, req bool) bool {
	return verifyUser([]string{c.User}, cmp, req)
}

// VerifyExpiresAt compares the exp claim against cmp (cmp < exp).
// If req is false, it will return true, if exp is unset.
func (c *LicenseClaims) VerifyExpiresAt(cmp int64, req bool) bool {
	if c.ExpiresAt == 0 {
		return verifyExp(nil, time.Unix(cmp, 0), req)
	}

	t := time.Unix(c.ExpiresAt, 0)
	return verifyExp(&t, time.Unix(cmp, 0), req)
}

// VerifyIssuedAt compares the iat claim against cmp (cmp >= iat).
// If req is false, it will return true, if iat is unset.
func (c *LicenseClaims) VerifyIssuedAt(cmp int64, req bool) bool {
	if c.IssuedAt == 0 {
		return verifyIat(nil, time.Unix(cmp, 0), req)
	}

	t := time.Unix(c.IssuedAt, 0)
	return verifyIat(&t, time.Unix(cmp, 0), req)
}

// VerifyNotBefore compares the nbf claim against cmp (cmp >= nbf).
// If req is false, it will return true, if nbf is unset.
func (c *LicenseClaims) VerifyNotBefore(cmp int64, req bool) bool {
	if c.NotBefore == 0 {
		return verifyNbf(nil, time.Unix(cmp, 0), req)
	}

	t := time.Unix(c.NotBefore, 0)
	return verifyNbf(&t, time.Unix(cmp, 0), req)
}

// VerifyIssuer compares the iss claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *LicenseClaims) VerifyIssuer(cmp string, req bool) bool {
	return verifyIss(c.Issuer, cmp, req)
}

// ----- helpers

func verifyUser(aud []string, cmp string, required bool) bool {
	if len(aud) == 0 {
		return !required
	}
	// use a var here to keep constant time compare when looping over a number of claims
	result := false

	var stringClaims string
	for _, a := range aud {
		if subtle.ConstantTimeCompare([]byte(a), []byte(cmp)) != 0 {
			result = true
		}
		stringClaims = stringClaims + a
	}

	// case where "" is sent in one or many aud claims
	if len(stringClaims) == 0 {
		return !required
	}

	return result
}

func verifyExp(exp *time.Time, now time.Time, required bool) bool {
	if exp == nil {
		return !required
	}
	return now.Before(*exp)
}

func verifyIat(iat *time.Time, now time.Time, required bool) bool {
	if iat == nil {
		return !required
	}
	return now.After(*iat) || now.Equal(*iat)
}

func verifyNbf(nbf *time.Time, now time.Time, required bool) bool {
	if nbf == nil {
		return !required
	}
	return now.After(*nbf) || now.Equal(*nbf)
}

func verifyIss(iss string, cmp string, required bool) bool {
	if iss == "" {
		return !required
	}
	if subtle.ConstantTimeCompare([]byte(iss), []byte(cmp)) != 0 {
		return true
	} else {
		return false
	}
}

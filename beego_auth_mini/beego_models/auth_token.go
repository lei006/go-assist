package beego_models

import (
	"errors"
	"fmt"
	"livertc/core/utils"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Token struct {
	Token     string    `orm:"column(token);pk" bson:"token"  json:"token"`
	Username  string    `orm:"index" bson:"username" json:"username"`                                 //用户名
	Avatar    string    `bson:"avatar" json:"avatar"`                                                 //角色（admin,user）
	SectionId int64     `orm:"index" bson:"section_id" json:"section_id"`                             //用来存用户部门
	CreatedAt time.Time `orm:"index;auto_now_add;type(datetime)" bson:"created_at" json:"created_at"` //创建时间
	UpdatedAt time.Time `orm:"index;auto_now;type(datetime)" bson:"updated_at" json:"updated_at"`     //修改时间
	ExpiredAt time.Time `orm:"index;type(datetime)" bson:"expired_at" json:"expired_at"`              //过期时间
}

func init() {
	orm.RegisterModel(new(Token))
}

//管理员
func (this *Token) IsAdmin() bool {
	return this.Avatar == "admin"
}

//操作员
func (this *Token) IsUser() bool {
	return this.Avatar == "user"
}

func (token *Token) ToString() string {
	return fmt.Sprintf("%+v", token)
}

type TokenModel struct {
}

var ModToken TokenModel

func (this *TokenModel) AddNewToken(userinfo *User) (*Token, error) {

	token, err := this.AddOne(userinfo)
	if err != nil {
		return nil, err
	}

	tokenInfo, err := this.GetOne(token)
	if err != nil {
		return nil, errors.New("GetOne : " + token + "   " + err.Error())
	}

	return tokenInfo, nil
}

func (this *TokenModel) AddOne(userinfo *User) (string, error) {

	new_token := utils.RandomString(25)

	tmp := &Token{
		Token:     new_token,
		Username:  userinfo.Username,
		Avatar:    userinfo.Avatar,
		SectionId: userinfo.SectionId,
		ExpiredAt: time.Now().Add(3 * time.Minute), //默认过期设置在3分钟后，
	}
	_, err := orm.NewOrm().Insert(tmp)
	if err != nil {
		return "", err
	}

	return new_token, nil

}

func (this *TokenModel) DeleteOne(token string) error {
	tmp := Token{Token: token}
	_, err := orm.NewOrm().Delete(&tmp)
	return err
}

func (this *TokenModel) DeleteOverTime(over_time time.Duration) error {
	overTime := time.Now().Add(-over_time)
	_, err := orm.NewOrm().QueryTable(Token{}).Filter("updated_at__lte", overTime).Delete()
	return err
}

func (this *TokenModel) UpdateOne(token string) error {
	tmp := Token{Token: token}
	err := orm.NewOrm().Read(&tmp)
	if err != nil {
		return err
	}

	//省略的是一个int类型的返回值，代表的是更新了多少条数据
	_, err = orm.NewOrm().Update(&tmp)
	if err != nil {
		return err
	}

	return nil
}

func (this *TokenModel) UpdateExpiredAt(token string, expired_time time.Duration) error {
	tmp := Token{
		Token:     token,
		ExpiredAt: time.Now().Add(expired_time),
	}

	//省略的是一个int类型的返回值，代表的是更新了多少条数据
	_, err := orm.NewOrm().Update(&tmp, "expired_at")
	if err != nil {
		return err
	}

	return nil
}

func (this *TokenModel) GetOne(token string) (*Token, error) {

	count, err := orm.NewOrm().QueryTable(Token{}).Filter("token", token).Count()
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, errors.New("no find " + token)
	}

	tmp := Token{Token: token}
	err = orm.NewOrm().Read(&tmp)
	if err != nil {
		return nil, err
	}

	return &tmp, nil
}

func (this *TokenModel) CheckToken(token string) error {
	tmp := Token{Token: token}
	err := orm.NewOrm().Read(&tmp)
	if err != nil {
		return err
	}
	return nil
}

func (this *TokenModel) IsExpired(token string) (bool, error) {

	// 过期时间，还要旧
	count, err := orm.NewOrm().QueryTable(Token{}).Filter("token", token).Filter("expired_at__lte", time.Now()).Count()
	if err != nil {
		return true, err
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}

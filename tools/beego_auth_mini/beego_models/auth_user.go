package beego_models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id        int64  `pk:"auto" bson:"id" json:"id"`
	Username  string `orm:"unique;index" bson:"username" json:"username"`
	Nickname  string `bson:"nickname" json:"nickname"`
	Password  string `bson:"password" json:"password"`
	Avatar    string `orm:"index" bson:"avatar" json:"avatar"` //admin,option,user
	SectionId int64  `orm:"index" bson:"section_id" json:"section_id"`
	IsEnable  bool   `bson:"is_enable" json:"is_enable"`
	Desc      string `orm:"size(8196)" bson:"desc" json:"desc"` //

	CreatedAt time.Time `orm:"index;auto_now_add;type(datetime)" bson:"created_at" json:"created_at"` //创建时间
	UpdatedAt time.Time `orm:"index;auto_now;type(datetime)" bson:"updated_at" json:"updated_at"`     //修改时间
}

func init() {
	orm.RegisterModel(new(User))
}

type UserModel struct {
}

var ModUser UserModel

func (this *UserModel) GetAll() ([]User, int, error) {

	var users []User
	_, err := orm.NewOrm().QueryTable(User{}).All(&users)

	return users, len(users), err
}

func (this *UserModel) AddOne(u User) (int64, error) {
	tmp_orm := orm.NewOrm()

	return tmp_orm.Insert(&u)
}

func (this *UserModel) DeleteOne(user_id int64) error {
	tmp := User{Id: user_id}
	_, err := orm.NewOrm().Delete(&tmp)
	return err
}

func (this *UserModel) UpdateOne(user_id int64, user *User) error {

	user.Id = user_id
	_, err := orm.NewOrm().Update(user, "nickname", "password", "avatar", "section_id", "desc", "is_enable")
	if err != nil {
		return err
	}

	return nil
}

func (this *UserModel) EnableOne(user_id int64, enable bool) error {

	tmp := &User{Id: user_id, IsEnable: enable}
	_, err := orm.NewOrm().Update(tmp, "is_enable")
	if err != nil {
		return err
	}

	return nil
}

func (this *UserModel) GetOne(user_id int64) (*User, error) {
	tmp := User{Id: user_id}
	err := orm.NewOrm().Read(&tmp)
	if err != nil {
		return nil, err
	}
	return &tmp, nil
}

func (this *UserModel) GetOneByUsername(username string) (*User, error) {
	tmp := User{}
	err := orm.NewOrm().QueryTable(User{}).Filter("username", username).One(&tmp)
	if err != nil {
		return nil, err
	}
	return &tmp, nil
}

//用户是否存在..
func (this *UserModel) IsExist(username string) (bool, error) {

	count, err := orm.NewOrm().QueryTable(User{}).Filter("username", username).Count()
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

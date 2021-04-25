package beego_assert

import (
	"github.com/beego/beego/v2/client/orm"

	_ "github.com/go-sql-driver/mysql"
)

func InitMySql(dataSource string) error {

	orm.RegisterDriver("mysql", orm.DRMySQL)
	if err := orm.RegisterDataBase("default", "mysql", dataSource); err != nil {
		return err
	}

	orm.RunSyncdb("default", false, true)

	return nil
}

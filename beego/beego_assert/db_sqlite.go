package beego_assert

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/mattn/go-sqlite3"
)

func InitSqlite(dataSource string) error {

	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", dataSource)
	orm.RunSyncdb("default", false, true)

	return nil
}

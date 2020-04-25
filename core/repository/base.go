package repository

import (
	"fmt"
	"sync"

	"github.com/ilovelili/dongfeng/util"
	"github.com/jinzhu/gorm"

	// mysql driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	instance *gorm.DB
	once     sync.Once
)

// db singleton db client
func db() *gorm.DB {
	config := util.LoadConfig()
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	connectionString := fmt.Sprintf("%v:%v@tcp(%v)/%v?parseTime=true", config.DataBase.User, config.DataBase.Password, config.DataBase.Host, config.DataBase.DataBase)
	once.Do(func() {
		db, err := gorm.Open("mysql", connectionString)
		if err == nil {
			if config.DataBase.Debug {
				instance = db.Set("gorm:auto_preload", true).Debug().LogMode(true)
			} else {
				instance = db.Set("gorm:auto_preload", true).LogMode(false)
			}
		}
	})

	return instance
}

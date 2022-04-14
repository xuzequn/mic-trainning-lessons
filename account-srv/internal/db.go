package internal

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mic-trainning-lessons/account-srv/model"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	newlogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), //io.writer
		logger.Config{
			SlowThreshold:             time.Second, //slowsql
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        //忽略ErrorRecordNotFound(记录未找到)报错
			Colorful:                  true,        // 禁用彩色打印
		},
	)
	conn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "123456", "127.0.0.1", 3306, "happy_account_mic_training")
	DB, err := gorm.Open(mysql.Open(conn), &gorm.Config{
		Logger: newlogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //表使用英文单数形式
		},
	})
	if err != nil {
		panic("数据库连接失败" + err.Error())
	}
	err = DB.AutoMigrate(&model.Account{})
	if err != nil {
		fmt.Println(err)
	}
}

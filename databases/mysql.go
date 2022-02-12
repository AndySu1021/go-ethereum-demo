package databases

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var MySqlClient *gorm.DB

// InitMySql 初始化連線資料庫，生成可操作基本增刪改查結構的變數
func InitMySql() (err error)  {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("databases.mysql.username"),
		viper.GetString("databases.mysql.password"),
		viper.GetString("databases.mysql.host"),
		viper.GetString("databases.mysql.port"),
		viper.GetString("databases.mysql.database"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	MySqlClient, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: newLogger,
	})
	if err != nil{
		return err
	}

	sqlDB, _ := MySqlClient.DB()
	return sqlDB.Ping()
}

func Close() (err error) {
	sqlDB, _ := MySqlClient.DB()
	err = sqlDB.Close()
	return
}

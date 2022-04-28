package dao

import (
	"RedRock-Work/lesson_6/lv1/user/server/model"
	"RedRock-Work/lesson_6/lv1/user/server/tool"
	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var db *gorm.DB
var pool *redis.Pool

func InitMysql() {
	config := tool.GetConfig().Mysql
	dB, err := gorm.Open(mysql.Open(config.Gorm), &gorm.Config{
		SkipDefaultTransaction:                   false,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := dB.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)           //连接池中最大的空闲连接数
	sqlDB.SetMaxOpenConns(10)           //连接池最多容纳的链接数量
	sqlDB.SetConnMaxLifetime(time.Hour) //连接池中链接的最大可复用时间

	db = dB

	err = db.AutoMigrate(&model.UserInfo{}) //todo
	if err != nil {
		panic(err)
	}
}

func InitRedis() {
	config := tool.GetConfig().Redis

	p := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Address)
			if err != nil {
				return nil, err
			}
			_, err = c.Do("AUTH", config.Password)
			if err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
	pool = p
}

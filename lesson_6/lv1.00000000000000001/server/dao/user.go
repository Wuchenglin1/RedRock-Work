package dao

import (
	"RedRock-Work/lesson_6/lv1.00000000000000001/server/model"
	"github.com/gomodule/redigo/redis"
)

func SearchUserInfoByUserName(u *model.User) error {
	conn := pool.Get()
	defer conn.Close()
	password, err := redis.String(conn.Do("get", u.Username))
	if err != nil {
		if err.Error()[8:] == "nil returned" {
			res := db.Where("username = ?", u.Username).First(&u)
			if res.Error != nil {
				return res.Error
			}
		} else {
			return err
		}
	}
	u.Password = password
	return nil
}

func InsertUser(u model.User) (error, error) {
	tx := db.Begin()

	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("set", u.Username, u.Password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	res := tx.Create(&u)
	if res.Error != nil {
		tx.Rollback()
		return res.Error, nil
	}
	tx.Commit()
	return nil, nil
}

func Login(u *model.User) error {
	conn := pool.Get()
	defer conn.Close()
	password, err := redis.String(conn.Do("get", u.Username))
	if err != nil {
		if err.Error()[8:] == "nil returned" {
			res := db.Where("username = ?", u.Username).First(&u)
			if res.Error != nil {
				return res.Error
			}
			return nil
		}
		return err
	}
	u.Password = password
	return nil
}

func ChangePassword(u model.User) error {
	tx := db.Begin()

	res := tx.Where("username = ?", u.Username).Updates(&u)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("set", u.Username, u.Password)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

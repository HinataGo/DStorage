package service

import (
	"context"
	"testing"

	"user/dao"
	"user/redis"
)

func TestUserServiceImpl_Login(t *testing.T) {

	err := dao.InitDB("127.0.0.1", "3306", "root", "123456", "test")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("127.0.0.1", "6379", "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	userService := &UserServiceImpl{
		userDao: &dao.UserDaoImpl{},
	}

	user, err := userService.Login(context.Background(), "aoho@mail.com", "aoho")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.Id)

}

func TestUserServiceImpl_Register(t *testing.T) {

	err := dao.InitDB("127.0.0.1", "3306", "root", "123456", "test")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("127.0.0.1", "6379", "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	userService := &UserServiceImpl{
		userDao: &dao.UserDaoImpl{},
	}

	user, err := userService.Register(context.Background(),
		&RegisterUserVO{
			Username: "aoho",
			Password: "aoho",
			Email:    "aoho@mail.com",
		})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.Id)

}

package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"user/dao"
	"user/redis"
)

// 定义用户数据转换 结构体
type UserInfoDTO struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// 定义用户注册结构体
type RegisterUserVO struct {
	Username string
	Password string
	Email    string
}

var (
	ErrUserExisted = errors.New("user is existed")
	ErrPassword    = errors.New("email and password are not match")
	ErrRegistering = errors.New("email is registering")
)

// 定义接口
type UserService interface {
	Login(ctx context.Context, email, password string) (*UserInfoDTO, error)
	Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO, error)
}
type UserServiceImpl struct {
	userDao dao.UserDao
}

func MakeUserServiceImpl(userDao dao.UserDao) UserService {
	return &UserServiceImpl{
		userDao: userDao,
	}
}

func (userService *UserServiceImpl) Login(ctx context.Context, email, password string) (*UserInfoDTO, error) {

	user, err := userService.userDao.SelectByEmail(email)
	if err == nil {
		if user.Password == password {
			return &UserInfoDTO{
				Id:       user.Id,
				Username: user.Username,
				Email:    user.Email,
			}, nil
		} else {
			return nil, ErrPassword
		}
	} else {
		log.Printf("err : %s", err)
	}
	return nil, err
}

func (userService UserServiceImpl) Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO, error) {

	lock := redis.GetRedisLock(vo.Email, time.Duration(5)*time.Second)
	err := lock.Lock()
	if err != nil {
		log.Printf("err : %s", err)
		return nil, ErrRegistering
	}
	defer lock.Unlock()

	existUser, err := userService.userDao.SelectByEmail(vo.Email)

	if (err == nil && existUser == nil) || err == gorm.ErrRecordNotFound {
		newUser := &dao.UserEntity{
			Username: vo.Username,
			Password: vo.Password,
			Email:    vo.Email,
		}
		err = userService.userDao.Save(newUser)
		if err == nil {
			return &UserInfoDTO{
				Id:       newUser.Id,
				Username: newUser.Username,
				Email:    newUser.Email,
			}, nil
		}
	}
	if err == nil {
		err = ErrUserExisted
	}
	return nil, err

}

package dao

import "time"

type UserEntity struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserEntity) TableName() string {
	return "user"
}

type UserDao interface {
	SelectByEmail(email string) (*UserEntity, error)
	Save(user *UserEntity) error
}
type UserDaoImpl struct {
}

func (UserDao *UserDaoImpl) SelectByEmail(email string) (*UserEntity, error) {
	user := &UserEntity{}
	err := db.Where("email = ?", email).First(user).Error
	return user, err
}

func (UserDao *UserDaoImpl) Save(user *UserEntity) error {
	return db.Create(user).Error
}

package service

import (
	db "DStorage/db/mysql"

	"fmt"
)

// User : 用户表model
type User struct {
	Username       string
	Email          string
	Phone          string
	SignupDate     string
	LastActiveDate string
	Status         int
}

// UserSignUp : 通过用户名及密码完成user表的注册操作
func UserSignUp(username string, passwd string) bool {
	stmt, err := db.Conn().Prepare(
		"insert ignore into storage_user (`user_name`,`user_pwd`) values (?,?)")
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

// UserSignIn : 判断密码是否一致
func UserSignIn(username string, verifyPwd string) bool {
	stmt, err := db.Conn().Prepare("select * from storage_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found: " + username)
		return false
	}
	pRows := db.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == verifyPwd {
		return true
	}
	return false
}

// UpdateToken : 刷新用户登录的token
func UpdateToken(username string, token string) bool {
	stmt, err := db.Conn().Prepare(
		"replace into storage_user_token (`user_name`,`user_token`) values (?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// GetUserInfo : 查询用户信息
func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := db.Conn().Prepare(
		"select user_name,signup_date from storage_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	// 执行查询的操作
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupDate)
	if err != nil {
		return user, err
	}
	return user, nil
}

// UserExist : 查询用户是否存在
func UserExist(username string) (bool, error) {

	stmt, err := db.Conn().Prepare(
		"select 1 from storage_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		return false, err
	}
	return rows.Next(), nil
}

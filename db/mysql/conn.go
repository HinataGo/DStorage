package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"DStorage/config"
)

var mySQL *sql.DB

// init : 初始化
func init() {
	mySQL, _ := sql.Open("mysql", config.MySQLSource)
	mySQL.SetMaxOpenConns(1000)
	err := mySQL.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
}

// Conn : 连接数据库
func Conn() *sql.DB {
	return mySQL
}

// ParseRows : 解析行数据
func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		// 将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		// TODO 处理异常错误数据
	}
}

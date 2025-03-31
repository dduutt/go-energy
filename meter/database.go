package meter

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

// 初始化数据库连接
func InitDB() (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_CONF["USER"], DB_CONF["PASSWORD"], DB_CONF["HOST"], DB_CONF["PORT"], DB_CONF["DATABASE"])
	db, err = sql.Open("mysql", dsn)
	dbOnce.Do(func() {
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal(err)
		}

		// 配置连接池
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		// 验证连接
		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}
	})

	return db, err
}

// 插入 Energy 数据到数据库
func InsertEnergy(energy *Energy) error {

	// 插入 SQL 语句
	query := `
		INSERT INTO energy (
			code, workshop, room, name, protocol, ip, port, slave_or_area, 
			start, size, data_type,byte_order, bytes, value, magnification
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	// 执行插入操作
	_, err = stmt.Exec(
		energy.Code,
		energy.WorkShop,
		energy.Room,
		energy.Name,
		energy.Protocol,
		energy.IP,
		energy.Port,
		energy.SlaveOrArea,
		energy.Start,
		energy.Size,
		energy.DataType,
		energy.ByteOrder,
		hex.EncodeToString(energy.Bytes),
		energy.Value,
		energy.Magnification,
	)
	return err
}

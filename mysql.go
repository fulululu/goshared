package goshared

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlRawConfig struct {
	Host              string `envconfig:"MYSQL_HOST,default=127.0.0.1"`
	Port              string `envconfig:"MYSQL_PORT,default=3306"`
	TLP               string `envconfig:"MYSQL_TLP,optional"`
	User              string `envconfig:"MYSQL_USER"`
	Password          string `envconfig:"MYSQL_PASSWORD"`
	Database          string `envconfig:"MYSQL_DATABASE"`
	ConnMaxLifeMinute int    `envconfig:"MYSQL_CONN_MAX_LIFE_MINUTE,default=3"`
	MaxOpenConns      int    `envconfig:"MYSQL_MAX_OPEN_CONNS,default=10"`
	MaxIdleConns      int    `envconfig:"MYSQL_MAX_IDLE_CONNS,default=10"`
}

type MysqlORMConfig struct {
	MySQL mysql.Config
	GORM  gorm.Config
}

func NewMysqlClient(rawCfg MysqlRawConfig, ormCfg MysqlORMConfig) (*gorm.DB, error) {
	// Raw layer
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=True&loc=Local", rawCfg.User, rawCfg.Password, "tcp", rawCfg.Host, rawCfg.Port, rawCfg.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * time.Duration(rawCfg.ConnMaxLifeMinute))
	db.SetMaxOpenConns(rawCfg.MaxOpenConns)
	db.SetMaxIdleConns(rawCfg.MaxIdleConns)

	// ORM layer
	ormCfg.MySQL.Conn = db
	gormDialector := mysql.New(ormCfg.MySQL)
	gormDB, err := gorm.Open(gormDialector, &ormCfg.GORM)
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

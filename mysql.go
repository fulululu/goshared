package goshared

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConfig struct {
	IP                string `envconfig:"MYSQL_IP"`
	Port              string `envconfig:"default=3306,MYSQL_PORT"`
	TLP               string `envconfig:"optional,MYSQL_TLP"`
	User              string `envconfig:"MYSQL_USER"`
	Password          string `envconfig:"MYSQL_PASSWORD"`
	Database          string `envconfig:"MYSQL_DATABASE"`
	ConnMaxLifeMinute int    `envconfig:"default=3,MYSQL_CONN_MAX_LIFE_MINUTE"`
	MaxOpenConns      int    `envconfig:"default=10,MYSQL_MAX_OPEN_CONNS"`
	MaxIdleConns      int    `envconfig:"default=10,MYSQL_MAX_IDLE_CONNS"`
}

type MySQLManager struct {
	Client    *sql.DB
	ORMClient *gorm.DB
}

func (ins *MySQLManager) Init(cfg MySQLConfig) error {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=True&loc=Local", cfg.User, cfg.Password, "tcp", cfg.IP, cfg.Port, cfg.Database)
	db, e := sql.Open("mysql", dsn)
	if e != nil {
		return e
	}

	e = db.Ping()
	if e != nil {
		return e
	}
	db.SetConnMaxLifetime(time.Minute * time.Duration(cfg.ConnMaxLifeMinute))
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	ins.Client = db
	gormDB, e := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if e != nil {
		return e
	}
	ins.ORMClient = gormDB
	return nil
}

package goshared

import (
	"database/sql"
	"fmt"
	"time"
)

type MySQLConfig struct {
	IP                string
	Port              string `envconfig:"default=3306"`
	TLP               string `envconfig:"optional"`
	User              string
	Password          string
	Database          string
	ConnMaxRetryTimes int `envconfig:"default=10"`
	ConnMaxLifeMinute int `envconfig:"default=3"`
	MaxOpenConns      int `envconfig:"default=10"`
	MaxIdleConns      int `envconfig:"default=10"`
}

type MySQLManager struct {
	Client *sql.DB
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

	return nil
}

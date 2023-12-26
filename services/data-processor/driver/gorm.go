package driver

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Database string

type GormOptions struct {
	Dialector             gorm.Dialector
	MaxOpenConnection     int
	MaxIdleConnection     int
	MaxConnectionLifetime time.Duration
}

type GormConnection struct {
	Options      *Options
	Instance     *sql.DB
	GormInstance *gorm.DB
}

const (
	MYSQL    Database = "MYSQL"
	POSTGRES Database = "POSTGRES"
)

func (cp *GormConnection) Init() error {
	if gormdb, e := gorm.Open(cp.Options.Gorm.Dialector, &gorm.Config{}); e != nil {
		return e
	} else {
		cp.GormInstance = gormdb
		if sqlDB, e := gormdb.DB(); e != nil {
			return e
		} else {
			cp.Instance = sqlDB
			sqlDB.SetMaxOpenConns(cp.Options.Gorm.MaxOpenConnection)
			sqlDB.SetMaxIdleConns(cp.Options.Gorm.MaxIdleConnection)
			sqlDB.SetConnMaxLifetime(cp.Options.Gorm.MaxConnectionLifetime)

			return nil
		}
	}
}

func (cp *GormConnection) GetInstance() any {
	return cp.Instance
}

func (cp *GormConnection) GetWrapperInstance() any {
	return cp.GormInstance
}

func (cp *GormConnection) GetOptions() *Options {
	return cp.Options
}

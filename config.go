package sql

import (
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

const (
	defaultDriverName = "mysql"
)

type Config struct {
	DriverName string

	User       string
	Passwd     string
	ServerName string
	DBName     string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (cfg Config) FormatDSN() string {
	net := "tcp"
	if strings.HasPrefix(cfg.ServerName, "/") {
		net = "unix"
	}
	mysqlcfg := mysql.NewConfig()
	mysqlcfg.User = cfg.User
	mysqlcfg.Passwd = cfg.Passwd
	mysqlcfg.Net = net
	mysqlcfg.Addr = cfg.ServerName
	mysqlcfg.DBName = cfg.DBName
	mysqlcfg.Collation = "utf8mb4_bin"
	mysqlcfg.InterpolateParams = true
	return mysqlcfg.FormatDSN()
}

func (cfg Config) driverName() string {
	if len(cfg.DriverName) > 0 {
		return cfg.DriverName
	}
	return defaultDriverName
}

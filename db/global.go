package db

import "github.com/ordishs/gocore"

var (
	dbHost, _     = gocore.Config().Get("db_host")
	dbPort, _     = gocore.Config().GetInt("db_port")
	dbName, _     = gocore.Config().Get("db_name")
	dbUser, _     = gocore.Config().Get("db_user")
	dbPassword, _ = gocore.Config().Get("db_password")
)

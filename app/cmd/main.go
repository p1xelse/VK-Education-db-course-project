package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	elog "github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var cfgPostgres = postgres.Config{DSN: "host=ws_pg user=postgres password=postgres port=5432"}

func main() {
	db, err := gorm.Open(postgres.New(cfgPostgres),
		&gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	userDB := userRep.New(db)
	forumDB := forumRep.New(db)
	postDB := postRep.New(db)
	threadDB := threadRep.New(db)
	serviceDB := serviceRep.New(db)

	userUC := userUsecase.New(userDB)
	forumUC := forumUsecase.New(forumDB, userDB)
	postUC := postUsecase.New(postDB, userDB, threadDB, forumDB)
	threadUC := threadUsecase.New(threadDB, userDB, forumDB)
	serviceUC := serviceUsecase.New(serviceDB)

	e := echo.New()

	e.Logger.SetHeader(`time=${time_rfc3339} level=${level} prefix=${prefix} ` +
		`file=${short_file} line=${line} message:`)
	e.Logger.SetLevel(elog.INFO)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `time=${time_custom} remote_ip=${remote_ip} ` +
			`host=${host} method=${method} uri=${uri} user_agent=${user_agent} ` +
			`status=${status} error="${error}" ` +
			`bytes_in=${bytes_in} bytes_out=${bytes_out}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

}

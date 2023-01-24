package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	elog "github.com/labstack/gommon/log"
	"github.com/p1xelse/VK_DB_course_project/app/cmd/server"
	_forumDelivery "github.com/p1xelse/VK_DB_course_project/app/internal/forum/delivery"
	forumRep "github.com/p1xelse/VK_DB_course_project/app/internal/forum/repository"
	forumUsecase "github.com/p1xelse/VK_DB_course_project/app/internal/forum/usecase"
	_postDelivery "github.com/p1xelse/VK_DB_course_project/app/internal/post/delivery"
	postRep "github.com/p1xelse/VK_DB_course_project/app/internal/post/repository"
	postUsecase "github.com/p1xelse/VK_DB_course_project/app/internal/post/usecase"
	_serviceDelivery "github.com/p1xelse/VK_DB_course_project/app/internal/service/delivery"
	serviceRep "github.com/p1xelse/VK_DB_course_project/app/internal/service/repository"
	serviceUsecase "github.com/p1xelse/VK_DB_course_project/app/internal/service/usecase"
	_threadDelivery "github.com/p1xelse/VK_DB_course_project/app/internal/thread/delivery"
	threadRep "github.com/p1xelse/VK_DB_course_project/app/internal/thread/repository"
	threadUsecase "github.com/p1xelse/VK_DB_course_project/app/internal/thread/usecase"
	_userDelivery "github.com/p1xelse/VK_DB_course_project/app/internal/user/delivery"
	userRep "github.com/p1xelse/VK_DB_course_project/app/internal/user/repository"
	userUsecase "github.com/p1xelse/VK_DB_course_project/app/internal/user/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var cfgPostgres = postgres.Config{DSN: "host=localhost user=db_pg password=db_postgres database=db_forum port=5432"}

//var cfgPostgres = postgres.Config{DSN: "host=localhost user=postgres password=postgres database=postgres port=8080"}

func main() {
	db, err := gorm.Open(postgres.New(cfgPostgres),
		&gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	userDB := userRep.NewUserRepository(db)
	forumDB := forumRep.NewForumRepository(db)
	postDB := postRep.NewPostRepository(db)
	threadDB := threadRep.NewThreadRepository(db)
	serviceDB := serviceRep.NewServiceRepository(db)

	userUC := userUsecase.NewServiceUsecase(userDB)
	forumUC := forumUsecase.NewForumUsecase(forumDB, userDB)
	postUC := postUsecase.NewPostUsecase(postDB, userDB, threadDB, forumDB)
	threadUC := threadUsecase.NewThreadUsecase(threadDB, userDB, forumDB)
	serviceUC := serviceUsecase.NewServiceUsecase(serviceDB)

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

	_userDelivery.NewDelivery(e, userUC)
	_forumDelivery.NewDelivery(e, forumUC)
	_postDelivery.NewDelivery(e, postUC)
	_threadDelivery.NewDelivery(e, threadUC)
	_serviceDelivery.NewDelivery(e, serviceUC)

	s := server.NewServer(e)
	if err := s.Start(); err != nil {
		e.Logger.Fatal(err)
	}
}

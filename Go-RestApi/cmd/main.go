package main

import (
	"go-restapi/internal/api"
	"go-restapi/internal/repository"
	"go-restapi/internal/service"
	"go-restapi/pkg/common"
	"go-restapi/pkg/db"
	"log"
)

func main() {
	config := config{
		addr: common.GetString("APP_PORT", ":8080"),
		dbCfg: databaseConfig{
			dbUrl:       common.GetString("DB_URL", "postgres://root:mypassword@127.0.0.1:5432/go_restapi?sslmode=disable"),
			maxOpenCons: common.GetInt("DB_MAX_CONS", 30),
			maxIdleCons: common.GetInt("DB_MAX_IDLE", 30),
			maxIdleTime: common.GetString("DB_TIME_IDLE", "15m"),
		},
	}

	db, err := db.New(
		config.dbCfg.dbUrl,
		config.dbCfg.maxOpenCons,
		config.dbCfg.maxIdleCons,
		config.dbCfg.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}

	repo := repository.New(db)
	service := service.New(repo)
	api := api.NewAPI(service)

	app := &application{
		config:  config,
		db:      db,
		service: service,
		api:     api,
	}

	log.Fatal(app.serve())
}

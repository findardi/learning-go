package main

import (
	"go-restapi/internal/api"
	"go-restapi/internal/repository"
	"go-restapi/internal/service"
	"go-restapi/pkg/common"
	"go-restapi/pkg/common/logger"
	"go-restapi/pkg/config/limiter"
	"go-restapi/pkg/db"
	"log"
	"time"
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
		rateLimiter: limiter.Config{
			RequestPerTimeFrame: 100,
			TimeFrame:           time.Minute,
			Enabled:             true,
		},
	}

	defer logger.Sync()

	db, err := db.New(
		config.dbCfg.dbUrl,
		config.dbCfg.maxOpenCons,
		config.dbCfg.maxIdleCons,
		config.dbCfg.maxIdleTime)
	if err != nil {
		logger.Fatal(err.Error())
	}

	repo := repository.New(db)
	service := service.New(repo)
	api := api.NewAPI(service)
	ratelimiter := limiter.NewFixedWindowRateLimiter(config.rateLimiter.RequestPerTimeFrame, config.rateLimiter.TimeFrame)

	app := &application{
		config:  config,
		api:     api,
		limiter: ratelimiter,
	}

	log.Fatal(app.serve())
}

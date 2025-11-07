package db

import (
	"context"
	"database/sql"
	"go-restapi/pkg/common/logger"
	"time"

	_ "github.com/lib/pq"
)

func New(addr string, maxOpenCons, maxIdleCons int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenCons)
	db.SetMaxIdleConns(maxIdleCons)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	logger.Info("database connected...")

	return db, nil
}

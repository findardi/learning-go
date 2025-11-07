package main

import (
	"context"
	"errors"
	"go-restapi/internal/api"
	"go-restapi/pkg/common/appmiddleware"
	"go-restapi/pkg/common/logger"
	"go-restapi/pkg/config/limiter"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type application struct {
	config  config
	api     *api.API
	limiter limiter.Limiter
}

type config struct {
	addr        string
	dbCfg       databaseConfig
	rateLimiter limiter.Config
}

type databaseConfig struct {
	dbUrl       string
	maxOpenCons int
	maxIdleCons int
	maxIdleTime string
}

func (app *application) mount() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	// router.Use(middleware.Logger)
	router.Use(appmiddleware.LoggerMiddleware)
	router.Use(middleware.Timeout(60 * time.Second))

	if app.config.rateLimiter.Enabled {
		router.Use(appmiddleware.RatelimiterMiddleware(app.limiter))
	}

	router.Route("/api", func(r chi.Router) {
		r.Get("/health", app.api.HealthCheck)

		r.Route("/v1", func(r chi.Router) {
			r.Route("/auth", func(r chi.Router) {
				r.Post("/register", app.api.Auth.Register)
				r.Post("/login", app.api.Auth.Login)
			})

			r.Group(func(r chi.Router) {
				r.Use(appmiddleware.AuthenticateMiddleware)

				r.Route("/user", func(r chi.Router) {
					r.Get("/profile", app.api.User.Profile)
				})
			})
		})
	})

	return router
}

func (app *application) serve() error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      app.mount(),
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Minute * 2,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		logger.Info("signal caugh", zap.String("signal", s.String()))

		shutdown <- srv.Shutdown(ctx)
	}()

	logger.Info("server has started", zap.String("addr", app.config.addr))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	logger.Info("server has shutdown", zap.String("addr", app.config.addr))
	return nil
}

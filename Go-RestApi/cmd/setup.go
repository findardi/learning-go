package main

import (
	"database/sql"
	"go-restapi/internal/api"
	"go-restapi/internal/repository"
	"go-restapi/internal/service"
	"go-restapi/pkg/common/appmiddleware"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config     config
	db         *sql.DB
	repository *repository.Queries
	service    *service.Service
	api        *api.API
}

type config struct {
	addr  string
	dbCfg databaseConfig
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
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(60 * time.Second))

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

	log.Printf("server has started at %s", app.config.addr)
	return srv.ListenAndServe()
}

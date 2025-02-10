package internal

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"
	_ "github.com/marcboeker/go-duckdb"
	"go.uber.org/zap"
)

func NewApp(ctx context.Context) (*App, error) {
	db, err := sql.Open("duckdb", "classified.duckdb")
	if err != nil {
		return nil, err
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	a := &App{
		db:     db,
		logger: logger,
	}

	a.initHandlers()

	return a, nil
}

type App struct {
	db     *sql.DB
	logger *zap.Logger
	router *echo.Echo
}

func (a *App) initHandlers() {
	e := echo.New()

	e.GET("/", a.indexPageHandler)

	e.GET("/signup", a.signupPageHandler)
	e.POST("/signup", a.signupHandler)

	e.GET("/login", a.loginPageHandler)
	e.POST("/login", a.loginHandler)

	e.GET("/post", a.postPageHandler)
	e.POST("/post", a.postCreateHandler)

	e.GET("/post/:post_id", a.getPostByIDHandler)

	e.Static("/static", "static")

	a.router = e
}

func (a *App) Run(ctx context.Context) error {
	return a.router.Start(":8000")
}

func (a *App) GetLogger() *zap.Logger {
	return a.logger
}

func (a *App) Close() error {
	defer func() {
		if a.logger != nil {
			err := a.logger.Sync()
			if err != nil {
				a.logger.Error("call logger sync", zap.Error(err))
			}
		}

	}()

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			a.logger.Error("close db", zap.Error(err))
		}
	}

	return nil
}

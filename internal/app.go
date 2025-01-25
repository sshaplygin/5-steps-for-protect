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

	logger, err := zap.NewProduction()
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

func (a *App) InitTables() error {
	_, err := a.db.Exec(`
		CREATE SEQUENCE users_id_sequence START 1;

		CREATE TABLE users (
			id INTEGER PRIMARY KEY DEFAULT nextval('users_id_sequence'), 
			login TEXT, 
			password TEXT, 
			is_admin BOOLEAN,
		);
	`)
	if err != nil {
		return err
	}

	_, err = a.db.Exec(`
		CREATE SEQUENCE posts_id_sequence START 1;

		CREATE TABLE posts (
			id INTEGER PRIMARY KEY DEFAULT nextval('posts_id_sequence'), 
			title TEXT, 
			content TEXT, 
			user_id INTEGER, 
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initHandlers() {
	e := echo.New()

	e.GET("/", a.indexPageHandler)

	e.GET("/signup", a.indexPageHandler)
	e.POST("/signup", a.signupHandler)

	e.GET("/login", a.loginPageHandler)
	e.POST("/login", a.loginHandler)

	e.GET("/post", a.postPageHandler)
	e.POST("/post", a.postCreateHandler)

	e.GET("/post/{{user_id}}/{{post_id}}", a.postPageHandler)

	a.router = e
}

func (a *App) Run(ctx context.Context) error {
	return a.router.Start(":8000")
}

func (a *App) GetLogger() *zap.Logger {
	return a.logger
}

func (a *App) Close() error {
	if a.logger != nil {
		err := a.logger.Sync()
		if err != nil {
			// TODO:
		}
	}

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			// TODO:
		}
	}

	return nil
}

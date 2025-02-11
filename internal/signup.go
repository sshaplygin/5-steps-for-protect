package internal

import (
	"backend/templates"
	"net/http"

	"github.com/alexedwards/argon2id"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func (a *App) signupHandler(c echo.Context) error {
	login := c.FormValue("username")
	password := c.FormValue("password")

	const adminLogin = "admin_login"

	isAdmin := login == adminLogin

	const query = `
		INSERT INTO users (login, password, is_admin) 
		VALUES (?, ?, ?) 
		RETURNING id
	`

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		a.logger.Error("create new hash by user password", zap.Error(err))

		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	row := a.db.QueryRow(query, login, hash, isAdmin)
	var userID int
	if err := row.Scan(&userID); err != nil {
		a.logger.Error("create new user", zap.Error(err))

		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	if row.Err() != nil {
		a.logger.Error("create new user row", zap.Error(row.Err()))

		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	c.SetCookie(makeUserCookie(userID))

	return c.Redirect(http.StatusSeeOther, "/")
}

func (a *App) signupPageHandler(c echo.Context) error {
	token, ok := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	if !ok {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	templates.Get().ExecuteTemplate(c.Response().Writer, "signup.html", map[string]string{
		middleware.DefaultCSRFConfig.ContextKey: token,
	})

	return nil
}

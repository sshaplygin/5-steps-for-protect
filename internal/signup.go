package internal

import (
	"backend/templates"
	"net/http"

	"github.com/labstack/echo/v4"
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

	row := a.db.QueryRow(query, login, password, isAdmin)
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
	templates.Get().ExecuteTemplate(c.Response().Writer, "signup.html", nil)

	return nil
}

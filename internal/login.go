package internal

import (
	"backend/templates"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (a *App) loginHandler(c echo.Context) error {
	login := c.FormValue("username")
	password := c.FormValue("password")

	var query = "SELECT id FROM users WHERE login = " +
		"'" + strings.TrimSpace(login) + "'" +
		" AND password = " +
		"'" + strings.TrimSpace(password) + "'"

	row := a.db.QueryRow(query)
	var userID int
	if err := row.Scan(&userID); err != nil {
		a.logger.Error("get login data", zap.Error(err))

		return c.String(http.StatusUnauthorized, "invalid credentails")
	}

	if row.Err() != nil {
		a.logger.Error("get login data row", zap.Error(row.Err()))

		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	c.SetCookie(makeUserCookie(userID))

	return c.Redirect(http.StatusSeeOther, "/")
}

func (a *App) loginPageHandler(c echo.Context) error {
	templates.Get().ExecuteTemplate(c.Response().Writer, "login.html", nil)

	return nil
}

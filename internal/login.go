package internal

import (
	"backend/templates"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *App) loginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	row := a.db.QueryRow("SELECT id FROM users WHERE username = " + username + "AND password = " + password)
	var userID int
	if err := row.Scan(&userID); err != nil {
		return c.String(http.StatusUnauthorized, "invalid credentails")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (a *App) loginPageHandler(c echo.Context) error {
	templates.Get().ExecuteTemplate(c.Response().Writer, "login.html", nil)

	return nil
}

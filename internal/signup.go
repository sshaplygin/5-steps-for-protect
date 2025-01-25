package internal

import (
	"backend/templates"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *App) signupHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	const adminLogin = "admin_login"

	isAdmin := username == adminLogin

	_, err := a.db.Exec("INSERT INTO users (username, password, is_admin) VALUES (?, ?, ?)", username, password, isAdmin)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}

func (a *App) signupPageHandler(c echo.Context) error {
	templates.Get().ExecuteTemplate(c.Response().Writer, "signup.html", nil)

	return nil
}

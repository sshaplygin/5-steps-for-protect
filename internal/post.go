package internal

import (
	"backend/templates"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *App) postCreateHandler(c echo.Context) error {
	title := c.FormValue("title")
	content := c.FormValue("content")

	_, err := a.db.Exec("INSERT INTO posts (title, content) VALUES (?, ?)", title, content)
	if err != nil {
		return c.String(http.StatusInternalServerError, "internal server error")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (a *App) postPageHandler(c echo.Context) error {
	templates.Get().ExecuteTemplate(c.Response().Writer, "post.html", nil)

	return nil
}

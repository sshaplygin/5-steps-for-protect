package internal

import (
	"backend/templates"
	"net/http"

	"github.com/alexedwards/argon2id"
	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (a *App) loginHandler(c echo.Context) error {
	login := c.FormValue("username")
	password := c.FormValue("password")

	const query = "SELECT id, password FROM users WHERE login = ?"

	row := a.db.QueryRow(query, login)
	var (
		userID int
		hash   string
	)
	if err := row.Scan(&userID, &hash); err != nil {
		a.logger.Error("get login data", zap.Error(err))

		return c.String(http.StatusUnauthorized, "invalid credentails")
	}

	if row.Err() != nil {
		a.logger.Error("get login data row", zap.Error(row.Err()))

		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil || !match {
		return c.String(http.StatusForbidden, "Invalid Credentails")
	}

	c.SetCookie(makeUserCookie(userID))

	return c.Redirect(http.StatusSeeOther, "/")
}

func (a *App) loginPageHandler(c echo.Context) error {
	templates.Get().ExecuteTemplate(c.Response().Writer, "login.html", nil)

	return nil
}

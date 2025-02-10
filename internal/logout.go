package internal

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (a *App) logoutHandler(c echo.Context) error {
	val, err := c.Cookie(UD_COOKIE)
	if err != nil || val == nil || val.Value == "" {
		a.logger.Error("get user data cookie", zap.Error(err))

		return c.String(http.StatusBadRequest, "Bad request")
	}

	userID, err := strconv.Atoi(val.Value)
	if err != nil {
		a.logger.Error("create post query", zap.Error(err))

		return c.String(http.StatusInternalServerError, "Internal error")
	}

	c.SetCookie(makeExipedUserCookie(userID))

	return c.Redirect(http.StatusSeeOther, "/")
}

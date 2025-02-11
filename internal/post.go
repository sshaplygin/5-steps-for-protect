package internal

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"backend/templates"
)

func (a *App) postCreateHandler(c echo.Context) error {
	title := c.FormValue("title")
	content := c.FormValue("content")

	val, err := c.Cookie(UD_COOKIE)
	if err != nil {
		a.logger.Error("get user data cookie", zap.Error(err))

		return c.String(http.StatusBadRequest, "Bad request")
	}

	if val.Value == "" {
		return c.String(http.StatusUnauthorized, "unauthorized request")
	}

	const query = `
		INSERT INTO posts (title, content, user_id) 
		VALUES (?, ?, ?)
	`

	_, err = a.db.Exec(query, title, content, val.Value)
	if err != nil {
		a.logger.Error("create post query", zap.Error(err))

		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (a *App) getPostByIDHandler(c echo.Context) error {
	val, err := c.Cookie(UD_COOKIE)
	if err != nil {
		a.logger.Error("get cookie value", zap.Error(err))

		return c.String(http.StatusUnauthorized, "Unauthorized request")
	}

	if val.Value == "" {
		return c.String(http.StatusUnauthorized, "Unauthorized request")
	}

	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		a.logger.Error("get post_id from request path", zap.Error(err))

		return c.String(http.StatusBadRequest, "Bad request")
	}

	userID, err := strconv.Atoi(val.Value)
	if err != nil {
		a.logger.Error("get user_id from cookie value", zap.Error(err))

		return c.String(http.StatusBadRequest, "Bad request")
	}

	const query = `
		SELECT title, content, created_at 
		FROM posts 
		WHERE id = ? AND user_id = ?
	`
	row := a.db.QueryRow(query, postID, userID)

	var (
		title, content string
		createdAt      time.Time
	)

	post := map[string]string{
		"id": strconv.Itoa(postID),
	}

	if err := row.Scan(&title, &content, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			templates.Get().ExecuteTemplate(c.Response().Writer, "user_post.html", post)

			return nil
		}

		a.logger.Error("get post by id", zap.Error(err))

		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	post["title"] = title
	post["content"] = content
	post["created_at"] = createdAt.Format(time.RFC822)

	templates.Get().ExecuteTemplate(c.Response().Writer, "user_post.html", post)

	return nil
}

func (a *App) postPageHandler(c echo.Context) error {
	token, ok := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	if !ok {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	templates.Get().ExecuteTemplate(c.Response().Writer, "post.html", map[string]string{
		middleware.DefaultCSRFConfig.ContextKey: token,
	})

	return nil
}

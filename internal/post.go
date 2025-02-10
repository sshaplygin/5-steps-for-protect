package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"backend/templates"
)

func (a *App) postCreateHandler(c echo.Context) error {
	title := c.FormValue("title")
	content := c.FormValue("content")

	val, err := c.Cookie(UD_COOKIE)
	if err != nil {
		a.logger.Error("create post query", zap.Error(err))

		return c.String(http.StatusBadRequest, "Bad request")
	}

	if val.Value == "" {
		return c.String(http.StatusUnauthorized, "unauthorized request")
	}

	const query = `
		INSERT INTO posts (title, content, user_id) 
		VALUES ('%s', '%s', %s)
	`

	_, err = a.db.Exec(fmt.Sprintf(query, title, content, val.Value))
	if err != nil {
		a.logger.Error("create post query", zap.Error(err))

		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func (a *App) getPostByIDHandler(c echo.Context) error {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		a.logger.Error("get post_id from request path", zap.Error(err))

		return c.String(http.StatusBadRequest, "Bad request")
	}

	const query = `
		SELECT title, content, created_at 
		FROM posts 
		WHERE id = $1
	`
	row := a.db.QueryRow(query, postID)

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
	templates.Get().ExecuteTemplate(c.Response().Writer, "post.html", nil)

	return nil
}

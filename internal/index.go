package internal

import (
	"net/http"
	"strconv"
	"time"

	"backend/templates"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type pageData struct {
	Logged bool
	Posts  []map[string]string
}

func (a *App) indexPageHandler(c echo.Context) error {
	const query = `
		SELECT user_id, title, content, created_at 
		FROM posts 
		ORDER BY created_at DESC
	`

	rows, err := a.db.Query(query)
	if err != nil {
		a.logger.Error("get posts feed", zap.Error(err))

		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	defer rows.Close()

	var (
		posts          []map[string]string
		title, content string
		createdAt      time.Time
		userID         int64
	)

	for rows.Next() {
		rows.Scan(&userID, &title, &content, &createdAt)

		posts = append(posts, map[string]string{
			"user_id":    strconv.FormatInt(userID, 10),
			"title":      title,
			"content":    content,
			"created_at": createdAt.Format(time.RFC822),
		})
	}

	if err = rows.Err(); err != nil {
		a.logger.Error("get post_id from request path", zap.Error(err))

		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	val, _ := c.Cookie(UD_COOKIE)
	logged := val != nil && val.Value != ""

	templates.Get().ExecuteTemplate(c.Response().Writer, "index.html", pageData{logged, posts})

	return nil
}

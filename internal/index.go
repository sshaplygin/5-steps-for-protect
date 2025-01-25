package internal

import (
	"net/http"

	"backend/templates"

	"github.com/labstack/echo/v4"
)

func (a *App) indexPageHandler(c echo.Context) error {
	rows, err := a.db.Query("SELECT title, content, created_at FROM posts ORDER BY created_at ASC")
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	defer rows.Close()

	var posts []map[string]string
	var title, content, createdAt string
	for rows.Next() {

		rows.Scan(&title, &content, &createdAt)

		posts = append(posts, map[string]string{
			"title":      title,
			"content":    content,
			"created_at": createdAt,
		})
	}

	templates.Get().ExecuteTemplate(c.Response().Writer, "index.html", posts)

	return nil
}

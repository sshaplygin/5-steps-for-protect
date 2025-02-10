package internal

import (
	"net/http"
	"strconv"
	"time"
)

const UD_COOKIE = "ud"

func makeUserCookie(userID int) *http.Cookie {
	maxAge := 24 * time.Hour

	return &http.Cookie{
		Name:     UD_COOKIE, // user_data
		Value:    strconv.Itoa(userID),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   int(maxAge.Seconds()),
	}
}

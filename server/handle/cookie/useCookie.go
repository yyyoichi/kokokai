package cke

import (
	"errors"
	"net/http"
	"os"
	"time"
)

func NewUserCookie(jwtToken string) *http.Cookie {
	ev := os.Getenv("ENV")
	c := &http.Cookie{
		Name:    "token",
		Value:   jwtToken,
		Expires: time.Now().AddDate(0, 0, 7),
	}
	switch ev {
	case "DEV":
		c.HttpOnly = false
		c.Secure = false
	case "STG":
		c.HttpOnly = true
		c.Secure = true
		c.SameSite = http.SameSiteNoneMode
		c.Domain = "collokaistg.yyyoichi.com"
	case "PRO":
		c.HttpOnly = true
		c.Secure = true
		c.SameSite = http.SameSiteNoneMode
		c.Domain = "collokai.yyyoichi.com"
	}
	return c
}

func FromUserCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, errors.New("no user cookie: " + err.Error())
	}
	return cookie, nil
}

func UpdateUserCookie(r *http.Request, jwtToken string) (*http.Cookie, error) {
	c, err := FromUserCookie(r)
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:    "token",
		Value:   jwtToken,
		Expires: c.Expires,
	}, nil
}

package sess

import (
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte(os.Getenv("SSS_SECRET"))
	store = sessions.NewCookieStore(key)
)

func NewUserCSRFToken(r *http.Request) *sessions.Session {
	s, _ := store.Get(r, "user")
	s.Values["csrftoken"] = newId()
	ev := os.Getenv("ENV")
	switch ev {
	case "STG":
		s.Options.HttpOnly = true
		s.Options.Secure = true
		s.Options.SameSite = http.SameSiteNoneMode
	case "PRO":
		s.Options.HttpOnly = true
		s.Options.Secure = true
		s.Options.SameSite = http.SameSiteNoneMode
		s.Options.Domain = "https://collokai.yyyoichi.com"
	}
	return s
}

func FromUserCSRFToken(r *http.Request) (*string, bool) {
	session, _ := store.Get(r, "user")
	token, ok := session.Values["csrftoken"].(string)
	if !ok || token == "" {
		return nil, false
	}
	return &token, true
}

func newId() string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var b strings.Builder
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 20; i++ {
		b.WriteByte(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

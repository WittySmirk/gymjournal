package auth

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"os"
)

var sessionSecret = []byte(os.Getenv("SESSION_SECRET"))
var isProduction = os.Getenv("ENV") == "production"

func CreateAuth() {
	store := sessions.NewCookieStore(sessionSecret)
	store.Options.Path = "/"
	store.Options.MaxAge = 86400 * 30
	store.Options.HttpOnly = true
	store.Options.Secure = isProduction
	gothic.Store = store
	goth.UseProviders(google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_SECRET"), os.Getenv("BACKEND_URL")+"/auth/google/callback"))
}

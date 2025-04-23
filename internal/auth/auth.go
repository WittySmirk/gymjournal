package auth

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"net/http"
	"os"
)

var isProduction = os.Getenv("ENV") == "production"

var store *sessions.CookieStore = nil

func CreateStore() {
	sessionSecret := []byte(os.Getenv("SESSION_SECRET"))
	store = sessions.NewCookieStore(sessionSecret)
}

func CreateAuth() {
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	}
	gothic.Store = store
	goth.UseProviders(google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_SECRET"), os.Getenv("BACKEND_URL")+"/auth/google/callback"))
}

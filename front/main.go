package main

import (
	"fmt"
	"os"

	"github.com/go-martini/martini"
	_ "github.com/joho/godotenv/autoload"
	"github.com/martini-contrib/gzip"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	goauth2 "golang.org/x/oauth2"
)

func main() {
	m := App()
	m.Run()
}

func App() *martini.ClassicMartini {
	m := martini.Classic()

	m.Use(gzip.All())
	m.Use(sessions.Sessions("my_session",
		sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))))
	m.Use(martini.Static("assets"))

	m.Use(oauth2.Github(
		&goauth2.Config{
			ClientID:     os.Getenv("GH_CLIENT_ID"),
			ClientSecret: os.Getenv("GH_CLIENT_SECRET"),
			Scopes:       []string{""},
		},
	))

	m.Use(render.Renderer(render.Options{
		Directory: "templates",
		Layout:    "layout",
	}))

	m.Get("", Index)
	m.Get("/info", func(tokens oauth2.Tokens) string {
		if tokens.Expired() {
			return "not logged in, or the access token is expired"
		}
		return "logged in"
	})
	m.Get("/restrict", oauth2.LoginRequired, func(tokens oauth2.Tokens) string {
		return tokens.Access()
	})

	// m.Group("/actions", func(r martini.Router) {
	// 	r.Get("", ActionsIndex)
	// 	r.Get("/:id", ActionsShow)
	// 	r.Post("", ActionsNew)
	// 	r.Put("/:id", ActionsUpdate)
	// 	r.Delete("/:id", ActionsDelete)
	// })

	return m
}

func Index(r render.Render, tokens oauth2.Tokens) {
	if !tokens.Expired() {
		r.HTML(200, "hello", "true")
		fmt.Println(tokens.Access())
	} else {
		r.HTML(200, "hello", nil)
	}
}

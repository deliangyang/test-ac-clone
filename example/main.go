package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	docgen "github.com/deliangyang/chi-api-doc"
	"github.com/go-chi/chi/v5"
)

func adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", HelloWorld)
	r.Get("/accounts", HelloWorld)
	return r
}

func main() {
	r := chi.NewRouter()
	r.Get("/test/log", HelloWorld)
	r.Get("/test/2222", HelloWorld)

	r.Route("/articles", func(r chi.Router) {
		r.Get("/", HelloWorld)                     // GET /articles
		r.Get("/{month}-{day}-{year}", HelloWorld) // GET /articles/01-16-2017

		r.Post("/", HelloWorld)      // POST /articles
		r.Get("/search", HelloWorld) // GET /articles/search

		// Regexp url parameters:
		r.Get("/{articleSlug:[a-z-]+}", HelloWorld) // GET /articles/home-is-toronto

		// Subrouters:
		r.Route("/{articleID}", func(r chi.Router) {
			r.Get("/", HelloWorld)    // GET /articles/123
			r.Put("/", HelloWorld)    // PUT /articles/123
			r.Delete("/", HelloWorld) // DELETE /articles/123
		})
	})

	// Mount the admin sub-router
	r.Mount("/admin", adminRouter())

	fmt.Println(toMarkdown(docgen.Parse(r)))
	docs, err := json.MarshalIndent(docgen.Parse(r), "", "\t")

	if err != nil {
		panic(err)
	}
	fmt.Println(string(docs))
}

// @module Hello world
//
// @intro
// 用户检测（可获取用户昵称）
//
// @request
// <code>
//  {
//      userId: (integer)                  // 用户ID
//  }
// </code>
// @response
// <code>
//  {
//      user: {
//          userId: (integer),             // 用户ID
//          nickname: (string),            // 用户昵称 @since 2.1.3
//      }
//  }
// </code>
// @version 2.1.0
// @reviewed
// @deprecated true
//
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func toMarkdown(ad docgen.Docs) string {
	var lastModule string
	var markdown string
	for _, route := range ad {
		if lastModule != route.Module {
			markdown += "### " + route.Module + "\n"
		}

		template := `
#### %s

##### URL
%s %s

##### Request
%s

##### Response
%s
---
`
		markdown += fmt.Sprintf(template, route.Name, route.Method,
			route.Path, route.Request, route.Response)

		lastModule = route.Module
	}

	return markdown
}

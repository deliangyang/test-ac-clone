package docgen

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
)

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
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func TestParse(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/test/log", HelloWorld)
	routes := Parse(r)

	for _, value := range routes {
		if value.Path != "/test/log" {
			t.Fail()
		}

		// unit test无法获取comment
		if value.Version != "2.1.0" {
			t.Fail()
		}
	}
}

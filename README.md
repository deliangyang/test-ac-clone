

### Usage

```go
package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    
    docgen "github.com/deliangyang/api-docgen"
    "github.com/go-chi/chi"
)

func main() {
    r := chi.NewRouter()
    r.Get("/test/log", HelloWorld)
    r.Get("/test/2222", HelloWorld)
    
    docs, err := json.MarshalIndent(docgen.Parse(r), "", "\t")
    
    if err != nil {
        panic(err)
    }
    fmt.Println(string(docs))
}

/**
 * @module
 */
func HelloWorld(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("hello world"))
}
```

### 示例

```
example/main.go
```

### Happy API注释约束

```php
// @module
// 模块
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
//          nickname: (string),            // 用户昵称 @since 2.1.2
//      }
//  }
// </code>
// @version 2.1.0
// @reviewed
// @deprecated
```

### 生成输出格式
```json
[
  {
    "name": "日志记录",
    "module": "xxxx",
    "path": "/test/log",
    "method": "GET",
    "request": "{}",
    "response": "{}",
    "version": "2.1.1",
    "reviewed": false,
    "deprecated": false
  }
]
```

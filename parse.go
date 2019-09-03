package docgen

import (
	"regexp"
	"sort"
	"strings"

	"github.com/go-chi/chi"
	chiDocgen "github.com/go-chi/docgen"
)

// APIDoc 文档json序列化
type APIDoc struct {
	Name       string `json:"name"`
	Module     string `json:"module"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	Request    string `json:"request"`
	Response   string `json:"response"`
	Version    string `json:"version"`
	Reviewed   bool   `json:"reviewed"`
	Deprecated bool   `json:"deprecated"`
}

var (
	codeReg    = regexp.MustCompile(`</?code>`)
	keyWordReg = regexp.MustCompile(`^@([a-zA-z]+)`)
)

// Docs 文档数组
type Docs []APIDoc

func (d Docs) Len() int           { return len(d) }
func (d Docs) Less(i, j int) bool { return d[i].Path < d[j].Path && d[i].Method > d[j].Method }
func (d Docs) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

// Parse parse the comment
func Parse(r chi.Routes) Docs {
	rs, err := chiDocgen.BuildDoc(r)
	if err != nil {
		panic(err)
	}
	routes := subRoutes(rs.Router.Routes, "")
	sort.Sort(routes)
	return routes
}

// parseKeywords 获取关键词^@request, etc.
func parseKeywords(code string) map[string]string {
	meta := make(map[string]string)
	lines := strings.Split(code, "\n")
	var keyword string
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "@") {
			keywords := keyWordReg.FindStringSubmatch(line)
			keyword = keywords[1]
			meta[keyword] = strings.Replace(line, "@"+keyword, "", -1) + "\n"
		} else if keyword != "" {
			meta[keyword] += line + "\n"
		}
	}

	return meta
}

// reIndent 重新格式化缩进
func reIndent(code string) string {
	var indent int
	var currentIndent int
	var newCode string
	for _, line := range strings.Split(code, "\n") {
		line = strings.TrimSpace(line)
		if len(line) <= 0 {
			continue
		} else if strings.HasSuffix(line, "{") || strings.HasSuffix(line, "[") {
			currentIndent = indent
			indent++
		} else if strings.HasPrefix(line, "}") || strings.HasPrefix(line, "]") {
			indent--
			currentIndent = indent
		} else {
			currentIndent = indent
		}
		for i := 0; i < currentIndent; i++ {
			line = "\t" + line
		}
		newCode += line + "\n"
	}
	return newCode
}

// paresComment 解析route handle的注释
func parseComment(path string, method string, c string) APIDoc {
	apiDoc := APIDoc{
		Path:       path,
		Method:     method,
		Reviewed:   false,
		Deprecated: false,
	}
	for k, v := range parseKeywords(c) {
		v = strings.TrimSpace(v)
		switch strings.ToLower(string(k)) {
		case "module":
			apiDoc.Module = v
		case "request":
			apiDoc.Request = reIndent(codeReg.ReplaceAllString(v, ""))
			break
		case "response":
			apiDoc.Response = reIndent(codeReg.ReplaceAllString(v, ""))
		case "intro":
			apiDoc.Name = v
		case "version":
			apiDoc.Version = v
		case "reviewed":
			if v == "true" || len(v) <= 0 {
				apiDoc.Reviewed = true
			}
		case "deprecated":
			if v == "true" || len(v) <= 0 {
				apiDoc.Deprecated = true
			}
		}
	}
	return apiDoc
}

func subRoutes(drs chiDocgen.DocRoutes, p string) Docs {
	var docs Docs

	for path, rs := range drs {
		if len(p) > 0 {
			path = strings.TrimRight(p, "/*") + path
		}

		if rs.Router != nil {
			docs = append(docs, subRoutes(rs.Router.Routes, path)...)
		}
		for method, info := range rs.Handlers {
			path = strings.TrimRight(path, "/")
			ad := parseComment(path, method, info.Comment)
			docs = append(docs, ad)
		}
	}
	return docs
}

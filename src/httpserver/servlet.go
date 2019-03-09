package httpserver

import (
	"net/http"
	"bytes"
	"regexp"
)

type routeConf struct {
	// 是否允许EncodePath
	useEncodedPath bool

	// 严格分割线
	strictSlash bool
}

// matcher URL匹配器
type matcher interface {
	Match(*http.Request) bool
}

type Router struct {
	// 请求处理器
	handler http.Handler

	// 路径表达式
	pathPattern string

	// 编译好的正则表达式
	pathReg *regexp.Regexp

	// 路径参数Key
	pathParamKeys []string

	// 匹配器
	matchers []matcher
}

func (r *Router) Path(path string) {
	r.pathPattern = path

	// 解析正则表达式
	r.pathPattern, r.pathParamKeys = parsePathPattern(path)

	r.pathReg = regexp.Compile(r.pathPattern)
}

// Match URL是否匹配了
func (r *Router) Match(request *http.Request) bool {
	url := request.RequestURI

	matchers := r.pathReg.FindStringSubmatch(url)
	if matchers != nil && len(matchers) > 0 {
		for index, key := range r.pathParamKeys {
			request.Form.Add(key, matchers[index + 1])
		}
		return true
	}
	return false
}

func parsePathPattern(tpl string) ([]string, string) {
	idxs, err := braceIndices(tpl)
	if err != nil {
		return
	}

	keys := make([]string, len(idxs)/2)
	pattern := bytes.NewBufferString("")
	pattern.WriteByte('^')
	
	for i:=0; i < len(idxs); i++ {
		if i == 0 {
			pattern.WriteString(tpl[0:idxs[i]])	
		} else if i % 2 != 0 {
			keys[i / 2] = tpl[idxs[i-1]:idxs[i]]
			pattern.WriteString(`([a-zA-Z0-9.%]+)`)
			if i + 1 < len(idxs) && idxs[i] != idxs[i+1] {
				pattern.WriteString(tpl[idxs[i]:idxs[i+1]])
			} else if i + 1 > len(idxs) {
				pattern.WriteString(tpl[idxs[i]:])
			}
		}
		
	}

	return keys, pattern.String()
}

func braceIndices(s string) ([]int, error) {
	var level, idx int
	var idxs []int

	for i:=0; i < len(s); i++ {
		switch s[i] {
		case '{': 
			level++
			if level == 1 {
				idx = i
			}
		case '}':
			level--
			if level == 0 {
				idxs = append(idxs, idx, i + 1)
			} else if level < 0 {
				return nil, fmt.Errorf("unbalanced braces in %q", s)
			}
	}
	if level != 0 {
		return nil, fmt.Errorf("unbalanced braces in %q", s)
	}

	return idxs, nil
}

type Servlet struct {
	// 路由数组
	routers []Router
}
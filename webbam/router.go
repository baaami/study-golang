package main

import (
	"net/http"
	"strings"
)

type router struct {
	// key: http 메서드
	// value: URL 패턴별로 실행할 HandlerFunc
	handlers map[string]map[string]HandlerFunc
}

func (r *router) HandleFunc(method, pattern string, h HandlerFunc) {
	// http 메서드로 등록된 맵이 있는지 확인
	m, ok := r.handlers[method]
	if !ok {
		// 등록된 맵이 없으면 새 맵을 생성
		m = make(map[string]HandlerFunc)
		r.handlers[method] = m
	}

	// http 메서드로 등록된 맵에 URL 패턴과 핸들러 함수 등록
	m[pattern] = h
}

func (r *router) handler() HandlerFunc {
	return func(c *Context) {
		// http 메서드에 맞는 모든 hanlders를 반복하며 요청 URL에 해당하는 handler를 찾음
		for pattern, handler := range r.handlers[c.Request.Method] {
			if ok, params := match(pattern, c.Request.URL.Path); ok {
				for k, v := range params {
					c.Params[k] = v
				}

				// 요청 URL에 해당하는 handler 수행
				handler(c)
				return
			}
		}

		http.NotFound(c.ResponseWriter, c.Request)
		return
	}
}

func match(pattern, path string) (bool, map[string]string) {
	// 패턴과 패스가 정확히 일치하면 바로 true를 반환
	if pattern == path {
		return true, nil
	}

	// 패턴과 패스를 "/" 단위로 구분
	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	// 패턴과 패스를 "/"로 구분한 후 부분 문자열 집합의 개수가 다르면 false를 반환
	if len(patterns) != len(paths) {
		return false, nil
	}

	// 패턴에 일치하는 URL 매개변수를 담기 위한 parrams 맵 생성
	params := make(map[string]string)

	// "/"로 구분된 패턴/패스의 각 문자열을 하나씩 비교
	for i := 0; i < len(patterns); i++ {
		switch {
		case patterns[i] == paths[i]:
			// 패턴과 패스의 부분 문자열이 일치하면 바로 다음 루프 수행
		case len(patterns[i]) > 0 && patterns[i][0] == ':':
			// 패턴이 ':' 문자로 시작하면 params에 URL params를 담은 후 다음 루프 수행
			params[patterns[i][1:]] = paths[i]
		default:
			// 일치하는 경우가 없으면 false를 반환
			return false, nil
		}
	}

	// true와 params를 반환
	return true, params
}

package webhook

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Action func(name string, params url.Values)

type WebHook struct {
	Host       string
	Port       int
	BasePath   string
	ActionPath map[string]string
	Actions    map[string]Action
}

func New(host string, basePath string, port int) *WebHook {
	path := basePath
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	return &WebHook{
		Host:       host,
		BasePath:   path,
		Port:       port,
		ActionPath: map[string]string{},
	}
}

func (that WebHook) Register(id string, accessToken string, action Action) {

	path := that.BasePath + id + "/"
	that.ActionPath[path] = id
	that.webhook(path, func(ctx *Context) {
		logDebug(fmt.Sprintf("action trigger: %s, token=%s", path, ctx.Query.Get("AccessToken")))
		if ctx.Query.Get("AccessToken") != accessToken {
			ctx.Status(http.StatusForbidden)
			return
		}
		actionId := strings.TrimLeft(ctx.Path, that.BasePath)
		action(actionId, ctx.Query)
		ctx.Status(http.StatusOK)
	})
}

func (that *WebHook) Listen() {
	addr := fmt.Sprintf("%s:%d", that.Host, that.Port)

	logDebug("start server: " + addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Cannot start server due to: " + err.Error())
	}
}

func (that WebHook) webhook(path string, handler func(ctx *Context)) {
	logDebug("register: " + path)
	http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		logDebug("http " + request.URL.Path)
		p := request.URL.Path //strings.TrimLeft(request.URL.Path, "/")
		if that.ActionPath[p] == "" {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		ctx := Context{
			Path:    request.URL.Path,
			Query:   request.URL.Query(),
			WebHook: that,
			Request: request,
			Writer:  writer,
		}
		handler(&ctx)
	})
}

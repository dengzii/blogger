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
	return &WebHook{
		Host:       host,
		BasePath:   basePath,
		Port:       port,
		ActionPath: map[string]string{},
	}
}

func (that WebHook) Register(id string, accessToken string, action Action) {

	path := that.BasePath + id
	that.ActionPath[path] = id
	that.webhook(path, func(ctx *Context) {
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

	err := http.ListenAndServe(addr, &handler{
		allowPath: that.ActionPath,
	})
	if err != nil {
		fmt.Println("Cannot start server due to: " + err.Error())
	}
}

type handler struct {
	allowPath map[string]string
}

func (that *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (that WebHook) webhook(path string, handler func(ctx *Context)) {
	http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		p := strings.TrimLeft(request.URL.Path, "/")
		if that.ActionPath[p] == "" {
			writer.WriteHeader(http.StatusForbidden)
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

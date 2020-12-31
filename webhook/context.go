package webhook

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Context struct {
	Path    string
	Query   url.Values
	WebHook WebHook
	Request *http.Request
	Writer  http.ResponseWriter
}

func (that *Context) Status(code int) {
	that.Writer.WriteHeader(code)
}

func (that *Context) Header(key string, value string) {
	that.Writer.Header().Add(key, value)
}

func (that *Context) writeString(str string) (err error) {
	if that.Writer != nil {
		_, err = that.Writer.Write(str2ByteArr(str))
	}
	return
}

func (that *Context) writeJson(model interface{}) error {
	j, err := json.Marshal(model)
	if err == nil {
		_, err = that.Writer.Write(j)
	}
	return err
}

func str2ByteArr(str string) []byte {
	return []byte(str)
}

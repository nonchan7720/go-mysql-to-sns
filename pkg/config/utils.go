package config

import (
	"fmt"
	"io"

	"github.com/valyala/fasttemplate"
)

func templateExecute(template string, mp map[string]interface{}) string {
	tpl := fasttemplate.New(template, "{{", "}}")
	return tpl.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		v := mp[tag]
		return w.Write([]byte(fmt.Sprintf("%v", v)))
	})
}

package config

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/valyala/fasttemplate"
)

type TemplateType string

const (
	GoTemplate   TemplateType = TemplateType("Go")
	FastTemplate TemplateType = TemplateType("Fast")
)

func (typ TemplateType) String() string {
	return string(typ)
}

func templateExecute(typ TemplateType, str string, mp map[string]interface{}) string {
	switch typ {
	case FastTemplate:
		return fastTemplate(str, mp)
	case GoTemplate:
		return goTemplate(str, mp)
	}
	return ""
}

func builtins() template.FuncMap {
	return template.FuncMap{
		"print":      fmt.Sprint,
		"printf":     fmt.Sprintf,
		"println":    fmt.Sprintln,
		"html":       template.HTMLEscaper,
		"js":         template.JSEscaper,
		"urlquery":   template.URLQueryEscaper,
		"replace":    strings.Replace,
		"replaceAll": strings.ReplaceAll,
		"newUUID":    uuid.NewString,
	}
}

type mapper map[string]interface{}

func replaceQuote(key string) string {
	return strings.ReplaceAll(key, `"`, "")
}

func (m mapper) find(key string) (interface{}, bool) {
	keys := strings.Split(replaceQuote(key), ".")
	return nestedMapLookup(m, keys...)
}

func (m mapper) findWithDefault(key string, default_ interface{}) interface{} {
	if v, ok := m.find(key); ok {
		return v
	}
	if v, ok := default_.(string); ok {
		return replaceQuote(v)
	}
	return default_
}

func safeCall(fun reflect.Value, args []reflect.Value) (val reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	ret := fun.Call(args)
	if len(ret) == 2 && !ret[1].IsNil() {
		return ret[0], ret[1].Interface().(error)
	}
	return ret[0], nil
}

var (
	reflectValueType = reflect.TypeOf((*reflect.Value)(nil)).Elem()
)

func unwrap(v reflect.Value) reflect.Value {
	if v.Type() == reflectValueType {
		v = v.Interface().(reflect.Value)
	}
	return v
}

func fastTemplate(str string, mp map[string]interface{}) string {
	mapper := mapper(mp)
	fnMap := builtins()
	tpl := fasttemplate.New(str, "{{", "}}")
	runFunc := func(tag string) (interface{}, error) {
		values := strings.Split(tag, " ")
		fn, ok := fnMap[values[0]]
		if ok {
			fv := reflect.ValueOf(fn)
			argv := []reflect.Value{}
			for _, value := range values[1:] {
				argv = append(argv, reflect.ValueOf(mapper.findWithDefault(value, value)))
			}
			v, err := safeCall(fv, argv)
			if err != nil {
				return "", err
			}
			return unwrap(v).Interface(), nil
		}
		return "", fmt.Errorf("Missing tag: %s", tag)
	}
	return tpl.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		v, ok := mapper.find(tag)
		if !ok {
			var err error
			v, err = runFunc(tag)
			if err != nil {
				return 0, err
			}
		}
		return w.Write([]byte(fmt.Sprintf("%v", v)))
	})
}

func goTemplate(str string, mp map[string]interface{}) string {
	tpl, err := template.New("").Funcs(builtins()).Parse(str)
	if err != nil {
		slog.Warn("Unparsable template syntax.", "template", str)
		return ""
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, mp); err != nil {
		slog.Warn(err.Error(), "template", str)
		return ""
	}
	return buf.String()
}

func nestedMapLookup(m map[string]interface{}, ks ...string) (interface{}, bool) {
	var ok bool
	var val interface{}

	// 入力の検証
	if len(ks) == 0 {
		return nil, false
	}
	if ks[0] == "" {
		return nil, false
	}

	// 最初のキーで値を取得する。
	val, ok = m[ks[0]]
	if !ok {
		return nil, false
	}

	// 最後のキーの場合、値を返す。
	if len(ks) == 1 {
		return val, true
	}

	// 値がマップの場合、再帰的にネストされたマップを探索する。
	if m, ok := val.(map[string]interface{}); ok {
		return nestedMapLookup(m, ks[1:]...)
	}

	// 値がマップでない場合、エラーを返す。
	return nil, false
}

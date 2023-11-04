package config

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strings"
	"text/template"

	"github.com/creasty/defaults"
	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/nonchan7720/go-storage-to-messenger/pkg/utils"
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
	mapper := utils.Mapper(mp)
	fnMap := builtins()
	tpl := fasttemplate.New(str, "{{", "}}")
	runFunc := func(tag string) (interface{}, error) {
		values := strings.Split(tag, " ")
		fn, ok := fnMap[values[0]]
		if ok {
			fv := reflect.ValueOf(fn)
			argv := []reflect.Value{}
			for _, value := range values[1:] {
				argv = append(argv, reflect.ValueOf(mapper.FindWithDefault(value, value)))
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
		v, ok := mapper.Find(tag)
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

type TConfig interface {
	Config | Outbox | OutboxPolling
}

func loadConfig[T TConfig](filePath string) (*T, error) {
	f, err := NewExpandEnv(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var config T
	if err := loadYaml(f, &config); err != nil {
		return nil, err
	}
	if err := validate.Struct(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func loadYaml(r io.Reader, v any) error {
	if err := yaml.NewDecoder(r).Decode(v); err != nil {
		return err
	}
	if err := defaults.Set(v); err != nil {
		return err
	}
	return nil
}

func WriteConfig[T TConfig](w io.Writer) error {
	var config T
	if err := yaml.NewEncoder(w).Encode(&config); err != nil {
		return err
	}
	return nil
}

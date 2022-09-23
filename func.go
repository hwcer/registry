package registry

import (
	"path"
	"reflect"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"
)

func Clean(paths ...string) (r string) {
	paths = append(RouterPrefix, paths...)
	p := path.Join(paths...)
	if r == "/" {
		r = ""
	} else {
		r = strings.ToLower(p)
	}
	return
}
func FuncName(i interface{}) (fname string) {
	fn := ValueOf(i)
	fname = runtime.FuncForPC(reflect.Indirect(fn).Pointer()).Name()
	if fname != "" {
		lastIndex := strings.LastIndex(fname, ".")
		if lastIndex >= 0 {
			fname = fname[lastIndex+1:]
		}
	}
	return strings.TrimSuffix(fname, "-fm")
}

func RouteName(name string) string {
	if strings.HasPrefix(name, PathMatchParam) {
		return PathMatchParam
	} else if strings.HasPrefix(name, PathMatchVague) {
		return PathMatchVague
	}
	return name
}

func ValueOf(i interface{}) reflect.Value {
	v, ok := i.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(i)
	}
	return v
}

func IsExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

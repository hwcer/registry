package registry

import (
	"errors"
	"fmt"
	"github.com/hwcer/cosgo/logger"
	"reflect"
	"strings"
)

/*
所有接口都必须已经登录
使用updater时必须使用playerHandle.Data()来获取updater
*/

//NewService name: /x/y
//文件加载init()中调用

func NewService(name string, opts *Options) *Service {
	r := &Service{
		Options: opts,
		nodes:   make(map[string]*Node),
		//method:  make(map[string]reflect.Value),
	}
	r.prefix = r.Clean(name)
	if len(r.prefix) > 1 {
		r.name = r.prefix[1:]
	}
	return r
}

type Service struct {
	*Options
	name   string
	prefix string
	nodes  map[string]*Node
	//method map[string]reflect.Value
}

func (this *Service) Name() string {
	return this.name
}
func (this *Service) Prefix() string {
	return this.prefix
}

// Register 服务注册
func (this *Service) Register(i interface{}, prefix ...string) error {
	v := reflect.ValueOf(i)
	var kind reflect.Kind
	if v.Kind() == reflect.Ptr {
		kind = v.Elem().Kind()
	} else {
		kind = v.Kind()
	}
	switch kind {
	case reflect.Func:
		return this.RegisterFun(v, prefix...)
	case reflect.Struct:
		return this.RegisterStruct(v, prefix...)
	default:
		return fmt.Errorf("未知的注册类型：%v", v.Kind())
	}
}

func (this *Service) format(name string, prefix ...string) string {
	if len(prefix) == 0 {
		return this.Clean(name)
	}
	s := this.Clean(prefix...)
	s = strings.Replace(s, "%v", strings.ToLower(name), -1)
	return s
}

func (this *Service) RegisterFun(i interface{}, prefix ...string) error {
	v := ValueOf(i)
	if v.Kind() != reflect.Func {
		return errors.New("RegisterFun fn type must be reflect.Func")
	}

	name := this.format(FuncName(v), prefix...)
	if name == "" {
		return errors.New("RegisterFun name empty")
	}
	node := &Node{name: name, value: v}
	if this.Options.Filter != nil && !this.Options.Filter(this, node) {
		return fmt.Errorf("RegisterFun filter return false:%v", name)
	}

	if _, ok := this.nodes[name]; ok {
		return fmt.Errorf("RegisterFun exist:%v", name)
	}
	this.nodes[name] = node
	//this.method[fname] = v
	this.Options.addRoutePath(this, name)
	return nil
}

// Register 注册一组handle
func (this *Service) RegisterStruct(i interface{}, prefix ...string) error {
	v := ValueOf(i)
	if v.Kind() != reflect.Ptr {
		return errors.New("RegisterStruct handle type must be reflect.Struct")
	}
	if v.Elem().Kind() != reflect.Struct {
		return errors.New("RegisterStruct handle type must be reflect.Struct")
	}
	handleType := v.Type()
	serviceName := handleType.Elem().Name()
	//sname := this.format(handleType.Elem().Name(), prefix...)
	//if sname == "" {
	//	return errors.New("RegisterStruct name empty")
	//}
	//if _, ok := this.nodes[sname]; ok {
	//	return fmt.Errorf("RegisterStruct name exist:%v", sname)
	//}
	//node := NewNode(sname, v)
	//this.nodes[sname] = node
	//logger.Debug("Watch:%v\n", sname)
	for m := 0; m < handleType.NumMethod(); m++ {
		method := handleType.Method(m)
		//methodType := method.Type
		methodName := method.Name
		//logger.Debug("Watch,sname:%v,type:%v", fname, methodType)
		// value must be exported.
		if method.PkgPath != "" {
			logger.Debug("Watch value PkgPath Not End,value:%v.%v(),PkgPath:%v", serviceName, methodName, method.PkgPath)
			continue
		}
		if !IsExported(methodName) {
			logger.Debug("Watch value Can't Exported,value:%v.%v()", serviceName, methodName)
			continue
		}
		name := this.format(serviceName+"/"+methodName, prefix...)
		node := &Node{name: name, binder: v, value: method.Func}
		if this.Options.Filter != nil && !this.Options.Filter(this, node) {
			continue
		}
		this.nodes[name] = node
		this.Options.addRoutePath(this, name)
	}
	return nil
}

// Match 匹配一个路径
// path : $prefix/$methodName
// path : $prefix/$nodeName/$methodName
func (this *Service) Match(path string) (node *Node, ok bool) {
	index := len(this.prefix)
	if index > 0 && !strings.HasPrefix(path, this.prefix) {
		return
	}
	name := path[index:]
	node, ok = this.nodes[name]
	return
}

func (this *Service) Paths() (r []string) {
	for k, _ := range this.nodes {
		r = append(r, k)
	}
	return
}

package registry

import (
	"path"
	"strings"
)

type Filter func(node *Node) bool

type filterHandle interface {
	Filter(node *Node) bool //用于判断struct中的方法是否合法接口
}

func NewOptions() *Options {
	return &Options{route: make(map[string]*Node)}
}

type Options struct {
	route  map[string]*Node    //ALL PATH
	Format func(string) string //格式化路径
}

// Clean 将所有path 格式化成 /a/b 模式
func (this *Options) Clean(paths ...string) (r string) {
	paths = append([]string{"/"}, paths...)
	p := path.Join(paths...)
	if this.Format != nil {
		r = this.Format(p)
	} else {
		r = strings.ToLower(p)
	}
	if r == "/" {
		r = ""
	}
	return
}

func (this *Options) addNode(node *Node) {
	k := path.Join(node.service.prefix, node.name)
	this.route[k] = node
}

package registry

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

const (
	PathMatchParam string = ":"
	PathMatchVague string = "*"
)

var RouterPrefix = []string{"/"}

func PathName(paths ...string) (r string) {
	paths = append(RouterPrefix, paths...)
	p := path.Join(paths...)
	if r == "/" {
		r = ""
	} else {
		r = strings.ToLower(p)
	}
	return
}

func NodeName(name string) string {
	if strings.HasPrefix(name, PathMatchParam) {
		return PathMatchParam
	} else if strings.HasPrefix(name, PathMatchVague) {
		return PathMatchVague
	}
	return name
}

func NewRouter() *Router {
	return newRouter("", []string{""}, 0)
}

func newRouter(name string, arr []string, step int) *Router {
	node := &Router{
		step:   step,
		name:   name,
		child:  make(map[string]*Router),
		route:  arr[0 : step+1],
		static: make(map[string]*Router),
	}
	return node
}

// 静态路由
func newStatic(arr []string, handle interface{}) *Router {
	l := len(arr) - 1
	node := newRouter(NodeName(arr[l]), arr, l)
	node.handle = handle
	return node
}

type Router struct {
	step   int                //层级
	name   string             // string,:,*,当前层匹配规则
	child  map[string]*Router //子路径
	route  []string           //当前路由绝对路径
	handle interface{}        //handle入口
	static map[string]*Router //静态路由,不含糊任何匹配参数
	//middleware []MiddlewareFunc //中间件
}

// Router is the registry of all registered Routes for an `engine` instance for
// Request matching and URL path parameter parsing.
//type Router struct {
//	root   map[string]*Router //method->Router
//	static map[string]*Router //静态路由,不含糊任何匹配参数
//}

/*
/s/:id/update
/s/123
*/
func (this *Router) Match(paths ...string) (nodes []*Router) {
	route := PathName(paths...)
	//静态路由
	if v, ok := this.static[route]; ok {
		nodes = append(nodes, v)
		return
	}
	//模糊匹配
	arr := strings.Split(route, "/")
	lastPathIndex := len(arr) - 1

	var spareNode []*Router
	var selectNode *Router

	for _, k := range []string{PathMatchVague, PathMatchParam, arr[1]} {
		if node := this.child[k]; node != nil {
			spareNode = append(spareNode, node)
		}
	}
	//
	//for _, node := range spareNode {
	//	fmt.Printf("spareNode:%v\n", strings.Join(node.Route(), "/"))
	//}
	n := len(spareNode)
	for selectNode != nil || n > 0 {
		fmt.Printf("==========%v====%v====\n", n, len(spareNode))
		if selectNode == nil {
			n -= 1
			selectNode = spareNode[n]
			spareNode = spareNode[0:n]
		}

		fmt.Printf("selectNode step:%v  PATH:%v  \n", selectNode.step, strings.Join(selectNode.Route(), "/"))
		if selectNode.name == PathMatchVague || selectNode.step == lastPathIndex {
			if selectNode.handle != nil {
				nodes = append(nodes, selectNode)
			}
			selectNode = nil
			fmt.Printf("匹配成功\n")
		} else {
			fmt.Printf("查询子节点:%v \n", selectNode.childes())
			//查询子节点
			i := selectNode.step + 1
			k := arr[i]
			if node := selectNode.child[PathMatchVague]; node != nil {
				n += 1
				spareNode = append(spareNode, node)
				fmt.Printf("添加候选节点 step:%v  PATH:%v  \n", node.step, strings.Join(node.Route(), "/"))
			}
			if node := selectNode.child[PathMatchParam]; node != nil {
				n += 1
				spareNode = append(spareNode, node)
				fmt.Printf("添加候选节点 step:%v  PATH:%v  \n", node.step, strings.Join(node.Route(), "/"))
			}
			if node := selectNode.child[k]; node != nil {
				selectNode = node
				fmt.Printf("添加候选节点 step:%v  PATH:%v  \n", node.step, strings.Join(node.Route(), "/"))
			} else {
				selectNode = nil
			}
		}
	}
	return
}

func (this *Router) Register(route string, handle interface{}) (err error) {
	if route == "" {
		return errors.New("Router.Watch method or route empty")
	}
	route = PathName(route)
	arr := strings.Split(route, "/")
	//静态路径
	if !strings.Contains(route, PathMatchParam) && !strings.Contains(route, PathMatchVague) {
		if _, ok := this.static[route]; ok {
			err = fmt.Errorf("route exist:%v", route)
		} else {
			this.static[route] = newStatic(arr, handle)
		}
		return
	}
	//匹配路由
	var node *Router
	for i := 1; i < len(arr); i++ {
		if node == nil {
			node = newRouter(NodeName(arr[i]), arr, i)
			this.child[node.name] = node
		} else {
			node, err = node.addChild(arr, i)
		}
		if err != nil {
			fmt.Printf("路由冲突: step:%v  PATH:%v  \n", node.step, strings.Join(node.Route(), "/"))
			return
		} else {
			//fmt.Printf("添加节点 step:%v  PATH:%v  \n", node.step, strings.Join(node.Route(), "/"))
		}
	}
	if node != nil {
		node.handle = handle
	}
	return
}

func (this *Router) Route() (r []string) {
	r = append(r, this.route...)
	return r
}

func (this *Router) Handle() interface{} {
	return this.handle
}

func (this *Router) Params(paths ...string) map[string]string {
	r := make(map[string]string)
	arr := strings.Split(PathName(paths...), "/")

	m := len(arr)
	if m > len(this.route) {
		m = len(this.route)
	}
	for i := 0; i < m; i++ {
		s := this.route[i]
		if strings.HasPrefix(s, PathMatchParam) {
			k := strings.TrimPrefix(s, PathMatchParam)
			r[k] = arr[i]
		} else if strings.HasPrefix(s, PathMatchVague) {
			k := strings.TrimPrefix(s, PathMatchVague)
			if k == "" {
				k = PathMatchVague
			}
			r[k] = strings.Join(arr[i:], "/")
		}

	}
	return r
}

func (this *Router) addChild(arr []string, step int) (node *Router, err error) {
	name := NodeName(arr[step])
	index := len(arr) - 1
	//(*)必须放在结尾
	if name == PathMatchVague && index != step {
		return nil, fmt.Errorf("router(*) must be at the end:%v", strings.Join(arr, "/"))
	}
	//路由重复
	node = this.child[name]
	if node != nil && node.handle != nil && index == step {
		err = fmt.Errorf("route exist:%v", strings.Join(arr, "/"))
		return
	}
	if node == nil {
		node = newRouter(name, arr, step)
		this.child[name] = node
		fmt.Printf("创建节点 step:%v  PATH:%v  \n", node.step, strings.Join(node.Route(), "/"))
	} else {
		fmt.Printf("存在节点 step:%v  PATH:%v  \n", node.step, strings.Join(node.Route(), "/"))
	}
	return
}

func (this *Router) String() string {
	if len(this.route) < 2 {
		return "/"
	}
	return strings.Join(this.route, "/")
}

func (this *Router) childes() (r []string) {
	for _, c := range this.child {
		r = append(r, c.String())
	}
	return
}

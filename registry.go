package registry

type Registry struct {
	dict   map[string]*Service
	router *Router
}

func New(router *Router) *Registry {
	if router == nil {
		router = NewRouter()
	}
	return &Registry{
		dict:   make(map[string]*Service),
		router: router,
	}
}

func (this *Registry) Len() int {
	return len(this.dict)
}

func (this *Registry) Get(name string) (srv *Service, ok bool) {
	prefix := Clean(name)
	srv, ok = this.dict[prefix]
	return
}
func (this *Registry) Has(name string) (ok bool) {
	prefix := Clean(name)
	_, ok = this.dict[prefix]
	return
}

func (this *Registry) Clean(paths ...string) string {
	return Clean(paths...)
}

func (this *Registry) Merge(r *Registry) {
	for _, s := range r.Services() {
		prefix := s.prefix
		if _, ok := this.dict[prefix]; !ok {
			this.dict[prefix] = NewService(prefix, this.router)
		}
		this.dict[prefix].Merge(s)
	}
}

func (this *Registry) Router() *Router {
	return this.router
}

// Match 通过路径匹配Route,path必须是使用 Registry.Clean()处理后的
func (this *Registry) Match(paths ...string) (node *Node, ok bool) {
	nodes := this.router.Match(paths...)
	for i := len(nodes) - 1; i >= 0; i-- {
		if node, ok = nodes[i].Handle().(*Node); ok {
			return
		}
	}
	return
}

// Service GET OR CREATE
func (this *Registry) Service(name string) *Service {
	prefix := Clean(name)
	if r, ok := this.dict[prefix]; ok {
		return r
	}
	srv := NewService(prefix, this.router)
	this.dict[prefix] = srv
	return srv
}

// Services 获取所有ServicePath
func (this *Registry) Services() (r []*Service) {
	for _, s := range this.dict {
		r = append(r, s)
	}
	return
}

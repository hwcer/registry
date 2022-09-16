package registry

type Registry struct {
	*Options
	dict map[string]*Service
}

func New(opts *Options) *Registry {
	if opts == nil {
		opts = NewOptions()
	}
	return &Registry{
		dict:    make(map[string]*Service),
		Options: opts,
	}
}

func (this *Registry) Len() int {
	return len(this.dict)
}

func (this *Registry) Get(name string) (srv *Service, ok bool) {
	prefix := this.Clean(name)
	srv, ok = this.dict[prefix]
	return
}
func (this *Registry) Has(name string) (ok bool) {
	prefix := this.Clean(name)
	_, ok = this.dict[prefix]
	return
}
func (this *Registry) Merge(r *Registry) {
	for _, s := range r.Services() {
		prefix := s.prefix
		if _, ok := this.dict[prefix]; !ok {
			this.dict[prefix] = &Service{name: s.name, prefix: s.prefix, nodes: make(map[string]*Node), Options: this.Options}
		}
		this.dict[prefix].Merge(s)
	}
}

// Match 通过路径匹配Route,path必须是使用 Registry.Clean()处理后的
func (this *Registry) Match(path string) (node *Node, ok bool) {
	path = this.Clean(path)
	node, ok = this.Options.route[path]
	return
}

// Service GET OR CREATE
func (this *Registry) Service(name string) *Service {
	prefix := this.Clean(name)
	if r, ok := this.dict[prefix]; ok {
		return r
	}
	srv := NewService(prefix, this.Options)
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

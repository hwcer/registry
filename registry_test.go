package registry

import (
	"testing"
)

type services struct {
}

func (s *services) Test() {

}

func ABC() {

}

func TestRoute(t *testing.T) {
	registry := New(nil)

	service := registry.Service("srv")
	service.Register(&services{})

	if err := service.Register(ABC); err != nil {
		t.Logf("ERROR:%v", err)
	}

	registry.Range(func(name string, service *Service) bool {
		t.Logf("servicePath:%v,serviceMethod:%v", name, service.Paths())
		return true
	})

}

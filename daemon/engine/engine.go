package engine

import (
	"github.com/songxinjianqwe/scheduler/common"
	"github.com/songxinjianqwe/scheduler/daemon/engine/standalone"
	"sync"
)



type Engine interface {
	List() ([]common.Task, error)
	Submit(task common.Task) error
	Get(id string, watch bool, version int64) (common.Task, error)
	Stop(id string) error
	Delete(id string) error
}

var instantiated Engine
var once sync.Once

func NewEngine() Engine {
	once.Do(func() {
		instantiated = &standalone.StandAloneEngine{}
	})
	return instantiated
}




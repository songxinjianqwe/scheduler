package engine

import (
	"github.com/songxinjianqwe/scheduler/common"
)

type Engine interface {
	List() ([]common.Task, error)
	Submit(task common.Task) error
	Get(id string, watch bool, version int64) (common.Task, error)
	Stop(id string) error
	Delete(id string) error
}

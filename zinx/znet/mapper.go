package znet

import (
	"fmt"
	"my-zinx/zinx/ziface"
)

// Controller Tags
const (
	// 0-99是给用户预留的自定义tag

	TAG_HEARTBEAT = iota + 100
)

type ControllerMapper struct {
	apis map[uint16]ziface.IController
}

func NewControllerMapper() *ControllerMapper {
	return &ControllerMapper{
		apis: make(map[uint16]ziface.IController),
	}
}

func (r *ControllerMapper) GetController(tag uint16) ziface.IController {
	return r.apis[tag]
}

func (r *ControllerMapper) AddController(tag uint16, controller ziface.IController) ziface.IControllerMapper {
	r.apis[tag] = controller
	return r
}

func (r *ControllerMapper) ExecController(tag uint16, request ziface.IRequest) error {
	if controller, ok := r.apis[tag]; ok {
		if err := controller.PreHandle(request); err != nil {
			return fmt.Errorf("PreHandle error: %v", err)
		}
		if err := controller.Handle(request); err != nil {
			return fmt.Errorf("Handle error: %v", err)
		}
		if err := controller.PostHandle(request); err != nil {
			return fmt.Errorf("PostHandle error: %v", err)
		}
	}
	return fmt.Errorf("no controller for tag[%d]", tag)
}

var _ ziface.IControllerMapper = (*ControllerMapper)(nil)

package hosts

import (
	"gogo-cmdb/commons"
)

var RegisterObj *Register

type Register struct {
	RegisterChan chan *commons.RegisterMsg
}

func NewRegister() *Register {
	return &Register{
		RegisterChan: make(chan *commons.RegisterMsg, 1000),
	}
}

func (r *Register) Run() {
	go func() {
		for {
			rgMsg := <-r.RegisterChan
			DefaultHostManager.Register(rgMsg)
		}
	}()
}

func init() {
	RegisterObj = NewRegister()
	RegisterObj.Run()
}

package hosts

import (
	"apiserver/models"
)

var RegisterObj *Register

type Register struct {
	RegisterChan chan *models.RegisterMsg
}

func NewRegister() *Register {
	return &Register{
		RegisterChan: make(chan *models.RegisterMsg, 1000),
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

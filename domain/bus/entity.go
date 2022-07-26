package bus

import (
	"encoding/json"
	"github.com/sunist-c/CeobeBotFramework/infrastructure/authenticator"
)

type IScheduledEntity interface {
	Options(interface{}) error
	Component() interface{}
	Rights() []authenticator.IRight
}

type ScheduledEntity struct {
	options   []byte
	component interface{}
	rights    []authenticator.IRight
}

func (s ScheduledEntity) Options(opt interface{}) (err error) {
	err = json.Unmarshal(s.options, opt)
	return err
}

func (s ScheduledEntity) Component() interface{} {
	return s.component
}

func (s ScheduledEntity) Rights() []authenticator.IRight {
	return s.rights
}

func NewScheduledEntity(opts, component interface{}, rights ...authenticator.IRight) (e IScheduledEntity, err error) {
	s := ScheduledEntity{
		options:   nil,
		component: component,
		rights:    rights,
	}
	s.options, err = json.Marshal(opts)
	return s, err
}

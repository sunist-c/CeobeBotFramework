package bus

import "github.com/sunist-c/CeobeBotFramework/infrastructure/authenticator"

type IScheduler interface {
	GetComponent(key string, terminal authenticator.ITerminal) (result IScheduledEntity, err error)
	AddComponent(key string, entity IScheduledEntity) (err error)
	ImplComponent(key string, component interface{}, opt interface{}) (err error)
	AddTerminal(terminals []authenticator.ITerminal, permissions []authenticator.IRight)
}

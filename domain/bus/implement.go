package bus

import (
	"github.com/sunist-c/CeobeBotFramework/infrastructure/authenticator"
	"sync"
)

type SimpleScheduler struct {
	entities      sync.Map // map[string(entity.key)]IScheduledEntity
	authenticator authenticator.IAuthenticator
}

func (s *SimpleScheduler) GetComponent(key string, terminal authenticator.ITerminal) (result IScheduledEntity, err error) {
	required, ok := s.entities.Load(key)
	if !ok {
		err = ComponentNotFoundError{
			key:       key,
			scheduler: s,
		}
		return nil, err
	}

	entity := required.(IScheduledEntity)
	ok, err = s.authenticator.Introspect(terminal, entity)

	if !ok {
		return nil, PermissionDeniedError{
			key: key,
		}
	} else if err != nil {
		return nil, err
	} else {
		return entity, nil
	}
}

func (s *SimpleScheduler) AddComponent(key string, entity IScheduledEntity) (err error) {
	if _, ok := s.entities.Load(key); ok {
		return ComponentAlreadyExistError{
			key:       key,
			scheduler: s,
		}
	} else {
		s.entities.Store(key, entity)
		return nil
	}
}

func (s *SimpleScheduler) ImplComponent(key string, component interface{}, opt interface{}) (err error) {
	if e, ok := s.entities.Load(key); !ok {
		return ComponentNotFoundError{
			key:       key,
			scheduler: s,
		}
	} else if e.(IScheduledEntity).Component() != nil {
		return ComponentAlreadyExistError{
			key:       key,
			scheduler: s,
		}
	} else {
		se, err := NewScheduledEntity(opt, component, e.(IScheduledEntity).Rights()...)
		if err != nil {
			return err
		} else {
			s.entities.Store(key, se)
			return nil
		}
	}
}

func (s *SimpleScheduler) AddTerminal(terminals []authenticator.ITerminal, permissions []authenticator.IRight) {
	var users []authenticator.IUser
	for _, terminal := range terminals {
		users = append(users, authenticator.NewUser(terminal.Key()))
	}
	s.authenticator.AddScopes(authenticator.NewScopeFields(users, permissions))
}

func DefaultScheduler() IScheduler {
	return &SimpleScheduler{
		entities:      sync.Map{},
		authenticator: authenticator.NewSimpleAuthenticator(),
	}
}

func DefaultSchedulerWithConf(conf Config) IScheduler {
	s := &SimpleScheduler{
		entities:      sync.Map{},
		authenticator: authenticator.NewSimpleAuthenticator(),
	}

	for _, option := range conf.EntityOptions {
		var user []authenticator.IUser
		se := ScheduledEntity{
			options:   nil,
			component: nil,
			rights:    []authenticator.IRight{},
		}

		for key, right := range option.TerminalsRight {
			se.rights = append(se.rights, right)
			user = append(user, authenticator.NewUser(key))
		}
		s.authenticator.AddResources(se)
		s.authenticator.AddScopes(authenticator.NewScopeFields(user, se.rights))
	}

	return s
}

func NewScheduler(authenticator authenticator.IAuthenticator) IScheduler {
	return &SimpleScheduler{
		entities:      sync.Map{},
		authenticator: authenticator,
	}
}

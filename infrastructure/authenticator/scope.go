package authenticator

import "sync"

type Permission byte

const (
	NonePermission          Permission = iota // NonePermission 不需要权限，即所有人可查看，但不可更改
	ReadPermission                            // ReadPermission 只读权限，拥有此权限的人可以查看，但不可更改
	WritePermission                           // WritePermission 可写权限，拥有此权限的人可以更改，但不可进行管理
	AdministratorPermission                   // AdministratorPermission 管理权限，最高权限
)

type IScope interface {
	Introspect(user IUser, rights ...IRight) bool
}

type IRight interface {
	Name() string
	Permission() Permission
}

type Right struct {
	name       string
	permission Permission
}

func (r Right) Name() string {
	return r.name
}

func (r Right) Permission() Permission {
	return r.permission
}

func NewRightFields(resName string, permission Permission) IRight {
	return Right{
		name:       resName,
		permission: permission,
	}
}

func EmptyRight(resName string) IRight {
	return Right{
		name:       resName,
		permission: NonePermission,
	}
}

type Scope struct {
	owners sync.Map // map[string(owner.key)]IUser
	scopes sync.Map // map[string(right.name)]Permission
}

func (s *Scope) Introspect(user IUser, rights ...IRight) bool {
	if _, ok := s.owners.Load(user.Key()); !ok {
		for _, right := range rights {
			if right.Permission() > NonePermission {
				return false
			}
		}
	}

	for _, right := range rights {
		if p, ok := s.scopes.Load(right.Name()); !ok && right.Permission() > NonePermission {
			return false
		} else if _, ok := p.(Permission); !ok {
			return false
		} else if right.Permission() > p.(Permission) {
			return false
		}
	}

	return true
}

func NewScopeFields(users []IUser, rights []IRight) IScope {
	scope := Scope{
		owners: sync.Map{},
		scopes: sync.Map{},
	}

	for _, user := range users {
		scope.owners.Store(user.Key(), nil)
	}

	for _, right := range rights {
		scope.scopes.Store(right.Name(), right.Permission())
	}

	return &scope
}

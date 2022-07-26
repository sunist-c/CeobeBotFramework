package authenticator

type NonAuthenticator struct {
}

func (n NonAuthenticator) Introspect(terminal ITerminal, resources ...IResource) (ok bool, err error) {
	return true, nil
}

func (n NonAuthenticator) AddResources(resources ...IResource) {
	return
}

func (n NonAuthenticator) AddScopes(scopes ...IScope) {
	return
}

func NewNonAuthenticator() IAuthenticator {
	return NonAuthenticator{}
}

type SimpleAuthenticator struct {
	scopes    []IScope
	resources []IResource
}

func (s *SimpleAuthenticator) Introspect(terminal ITerminal, resources ...IResource) (ok bool, err error) {
	authed := 0
	for _, resource := range resources {
		for i := 0; i < len(s.scopes); i++ {
			status := s.scopes[i].Introspect(terminal, resource.Rights()...)
			if status {
				authed += 1
				break
			}
		}
	}

	return authed == len(resources), err
}

func (s *SimpleAuthenticator) AddResources(resources ...IResource) {
	s.resources = append(s.resources, resources...)
}

func (s *SimpleAuthenticator) AddScopes(scopes ...IScope) {
	s.scopes = append(s.scopes, scopes...)
}

func NewSimpleAuthenticator() IAuthenticator {
	return &SimpleAuthenticator{
		scopes:    []IScope{},
		resources: []IResource{},
	}
}

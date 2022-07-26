package authenticator

// IAuthenticator the interface definition of Authenticator
type IAuthenticator interface {
	Introspect(terminal ITerminal, resources ...IResource) (ok bool, err error)
	AddResources(resources ...IResource)
	AddScopes(scopes ...IScope)
}

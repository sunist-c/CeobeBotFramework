package authenticator

type IResource interface {
	Rights() []IRight
}

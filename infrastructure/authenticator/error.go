package authenticator

import "fmt"

type OwnerNameAlreadyExistError struct {
	key string
}

func (o OwnerNameAlreadyExistError) Error() string {
	return fmt.Sprintf("%v as owner-name is already existed", o.key)
}

package bus

import (
	"fmt"
)

type ComponentAlreadyExistError struct {
	key       string
	scheduler IScheduler
}

func (c ComponentAlreadyExistError) Error() string {
	return fmt.Sprintf("key %v in scheduler %v is already exist, cannot replace", c.key, c.scheduler)
}

type ComponentNotFoundError struct {
	key       string
	scheduler IScheduler
}

func (c ComponentNotFoundError) Error() string {
	return fmt.Sprintf("secheduler %v cannot find component named %v", c.scheduler, c.key)
}

type PermissionDeniedError struct {
	key string
}

func (p PermissionDeniedError) Error() string {
	return fmt.Sprintf("permission denied when access %v", p.key)
}

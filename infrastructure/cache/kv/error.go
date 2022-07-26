package kv

import "fmt"

type OutofRangeError struct {
	upperIndex int
	lowerIndex int
}

func (o OutofRangeError) Error() string {
	return fmt.Sprintf("kv-db out of range with [%v:%v]", o.lowerIndex, o.upperIndex)
}

type NotFoundError struct {
	key interface{}
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("kv-db cannot find key: %#v", n.key)
}

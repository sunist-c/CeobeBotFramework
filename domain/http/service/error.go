package service

import "fmt"

type UnknownBodyTypeError struct {
	t string
}

func (u UnknownBodyTypeError) Error() string {
	return fmt.Sprintf("unknown %v as body type", u.t)
}

type BodyTypeNotMatchedError struct {
	got  string
	want string
}

func (b BodyTypeNotMatchedError) Error() string {
	return fmt.Sprintf("want body type as %v, but got %v", b.want, b.got)
}

type UnMarshalError struct {
	mapTo interface{}
	info  string
}

func (u UnMarshalError) Error() string {
	return fmt.Sprintf("error occured when unmarshal to %#v: %v", u.mapTo, u.info)
}

type NecessaryFieldNotFoundError struct {
	field string
	where string
}

func (n NecessaryFieldNotFoundError) Error() string {
	return fmt.Sprintf("field %v in %v is taged as necessary, but not found", n.field, n.where)
}

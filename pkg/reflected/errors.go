package reflected

import "fmt"

type ContainerContext interface{}
type ReflectedError struct {
	Msg        string
	Container  ContainerContext
	Underlying error
}

func (e ReflectedError) Error() string {
	return fmt.Sprintf("Reflected error [%s]: %s %+v", e.Msg, e.Underlying.Error(), e.Container)
}

func NewReflectedError(origin ContainerContext, msg string, underlying error) error {
	return ReflectedError{
		Msg:        msg,
		Container:  origin,
		Underlying: underlying,
	}
}

package reflected

import "fmt"

type ContainerContext interface{}
type ReflectedError struct {
	Msg        string
	Container  ContainerContext
	Underlying error
}

func (e ReflectedError) Error() string {
	var underlyingErrorMsg string
	if e.Underlying != nil {
		underlyingErrorMsg = e.Underlying.Error()
	}
	return fmt.Sprintf("Reflected error [%s]: %s %+v", e.Msg, underlyingErrorMsg, e.Container)
}

func NewReflectedError(origin ContainerContext, msg string, underlying error) error {
	return ReflectedError{
		Msg:        msg,
		Container:  origin,
		Underlying: underlying,
	}
}

package reflected

import "fmt"

type ReflectedError struct {
	Msg        string
	Origin     *TypedContainer
	Underlying error
}

func (e ReflectedError) Error() string {
	return fmt.Sprintf("Reflected error [%s]: %s %+v", e.Msg, e.Underlying.Error(), e.Origin)
}

func NewReflectedError(origin *TypedContainer, msg string, underlying error) error {
	return ReflectedError{
		Msg:        msg,
		Origin:     origin,
		Underlying: underlying,
	}
}

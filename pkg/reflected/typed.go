package reflected

var typeHandlers = make(map[string]Typed)

func RegisterHandler(t Typed) {
	typeHandlers[t.TypeName()] = t
}

type Typed interface {
	TypeName() string
	Cast(interface{}) Typed
}

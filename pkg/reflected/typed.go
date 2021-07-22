package reflected

var typeHandlers = make(map[string]Payload)

func RegisterHandler(t Payload) {
	typeHandlers[t.TypeName()] = t
}

type Payload interface {
	TypeName() string
	Cast(interface{}) Payload
}

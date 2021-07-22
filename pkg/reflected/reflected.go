package reflected

var typeHandlers = make(map[string]Payload)

func RegisterHandler(t Payload) {
	typeHandlers[t.TypeName()] = t
}

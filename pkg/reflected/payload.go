package reflected

type Payload interface {
	TypeName() string
	Cast(interface{}) Payload
}

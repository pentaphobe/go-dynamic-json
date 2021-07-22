package reflected

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type TypedContainer struct {
	Type    string `export:"type"`
	Payload Payload
}

func (t *TypedContainer) MarshalJSON() ([]byte, error) {
	// TODO: (not really a TODO) error handling is removed because I can't seem to actually
	// cause an error

	// if t == nil {
	// 	return nil, NewReflectedError(t, "marshal envelope", fmt.Errorf("can't retrieve type from nil container"))
	// }
	// if t.Payload == nil {
	// 	return nil, NewReflectedError(t, "marshal envelope", fmt.Errorf("can't retrieve type from nil payload"))
	// }
	env := struct {
		Type string `json:"type"`
	}{
		Type: t.Payload.TypeName(),
	}
	envJson, _ := json.Marshal(env)

	payloadJson, _ := json.Marshal(t.Payload)
	// if err != nil {
	// 	return nil, NewReflectedError(t, "marshal payload", err)
	// }
	result := append([]byte{'{'}, envJson[1:len(envJson)-1]...)
	result = append(result, []byte{','}...)
	result = append(result, payloadJson[1:]...)
	return result, nil
}

func (t *TypedContainer) UnmarshalJSON(in []byte) error {
	env := struct {
		Type string `json:"type"`
	}{}
	err := json.Unmarshal(in, &env)
	if err != nil {
		return NewReflectedError(t, "unmarshal envelope", err)
	}
	t.Type = env.Type
	handler, ok := typeHandlers[env.Type]
	if !ok {
		return NewReflectedError(t, fmt.Sprintf("unmarshal unhandled type %s", env.Type), err)
	}
	structType := reflect.ValueOf(handler).Elem().Type()
	v := handler.Cast(reflect.New(structType).Interface())
	err = json.Unmarshal(in, v)
	if err != nil {
		return NewReflectedError(t, "unmarshal payload", err)
	}
	t.Payload = v

	return nil
}

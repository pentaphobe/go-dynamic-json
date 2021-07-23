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

// MarshalJSONContainer is the core functionality behind TypedContainer.MarshalJSON
// It performs some cheekiness in order to marshal type fields alongside dynamic ones
//
// Unless you're writing custom containers, you probably don't need this
func MarshalJSONContainer(p Payload) ([]byte, error) {
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
		Type: p.TypeName(),
	}
	envJson, _ := json.Marshal(env)

	payloadJson, _ := json.Marshal(p)
	// if err != nil {
	// 	return nil, NewReflectedError(t, "marshal payload", err)
	// }
	result := append([]byte{'{'}, envJson[1:len(envJson)-1]...)
	result = append(result, []byte{','}...)
	result = append(result, payloadJson[1:]...)
	return result, nil
}

// MarshalJSON implements the Marshaler interface from encoding/json
// It is just a wrapper for MarshalJSONContainer
func (t *TypedContainer) MarshalJSON() ([]byte, error) {
	return MarshalJSONContainer(t.Payload)
}

// MarshalJSONContainer is the core functionality behind TypedContainer.UnmarshalJSON
// It performs some cheekiness in order to marshal type fields alongside dynamic ones
//
// Unless you're writing custom containers, you probably don't need this
func UnmarshalJSONContainer(context ContainerContext, in []byte) (string, Payload, error) {
	env := struct {
		Type string `json:"type"`
	}{}
	err := json.Unmarshal(in, &env)
	if err != nil {
		return "", nil, NewReflectedError(context, "unmarshal envelope", err)
	}
	// t.Type = env.Type
	handler, ok := typeHandlers[env.Type]
	if !ok {
		return "", nil, NewReflectedError(context, fmt.Sprintf("unmarshal unhandled type %s", env.Type), err)
	}
	structType := reflect.ValueOf(handler).Elem().Type()
	v := handler.Cast(reflect.New(structType).Interface())
	err = json.Unmarshal(in, v)
	if err != nil {
		return "", nil, NewReflectedError(context, "unmarshal payload", err)
	}
	// t.Payload = v

	return env.Type, v, nil

}

func (t *TypedContainer) UnmarshalJSON(in []byte) error {
	ty, pl, err := UnmarshalJSONContainer(t, in)
	if err != nil {
		return err
	}
	t.Type = ty
	t.Payload = pl
	return nil
}

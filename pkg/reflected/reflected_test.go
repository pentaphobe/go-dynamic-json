package reflected

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

//////////////////////////////////////////////////////////////

type TypedBasic struct {
	Name string `json:"name"`
}

func (t *TypedBasic) TypeName() string {
	return "basic"
}

func (t *TypedBasic) Cast(i interface{}) Payload {
	return i.(*TypedBasic)
}

//////////////////////////////////////////////////////////////

type TypedAnother struct {
	Age int `json:"age"`
}

func (t *TypedAnother) TypeName() string {
	return "another"
}

func (t *TypedAnother) Cast(i interface{}) Payload {
	return i.(*TypedAnother)
}

//////////////////////////////////////////////////////////////

func (s *Suite) SetupTest() {
	RegisterHandler(&TypedBasic{})
	RegisterHandler(&TypedAnother{})
}

func (s *Suite) TestReflect() {
	jsonData := `
	{
		"foo": [
			{
				"type": "basic",
				"name": "foo"
			},
			{
				"type": "another",
				"age": 23
			}
		]
	}
`
	t := struct {
		Foo []TypedContainer `json:"foo"`
	}{
		Foo: make([]TypedContainer, 0),
	}
	err := json.Unmarshal([]byte(jsonData), &t)
	s.Nil(err)
	result, err := json.Marshal(t)
	s.Nil(err)
	s.JSONEq(jsonData, string(result))
}

func (s *Suite) Test_UnmarshalErrors() {
	inputs := map[string]string{
		"invalid envelope": `{ "type": [] }`,
		"missing handler":  `{ "type": "YOU CAN'T HANDLE THIS" }`,
		"invalid payload":  `{ "type": "another", "age": [] }`,
	}
	for name, input := range inputs {
		s.Run(name, func() {
			t := &TypedContainer{}
			err := json.Unmarshal([]byte(input), t)
			s.NotNil(err)
			s.IsType(ReflectedError{}, err)
			re := err.(ReflectedError)
			s.Equal(t, re.Container)
			s.Contains(re.Error(), re.Msg)
		})
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

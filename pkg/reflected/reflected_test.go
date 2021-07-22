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

func (t *TypedBasic) Cast(i interface{}) Typed {
	return i.(*TypedBasic)
}

//////////////////////////////////////////////////////////////

type TypedAnother struct {
	Age int `json:"age"`
}

func (t *TypedAnother) TypeName() string {
	return "another"
}

func (t *TypedAnother) Cast(i interface{}) Typed {
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

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

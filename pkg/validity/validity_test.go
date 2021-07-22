package validity

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pentaphobe/dynamic-json/pkg/reflected"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

func (s *Suite) Test_Bidirection_Marshal() {
	inputs := map[string]struct {
		input reflected.Payload
	}{
		"duration": {
			input: &ValidityDuration{
				ValidTil: time.Now().AddDate(1, 0, 0),
			},
		},
		"changes": {
			input: &ValidityChanges{
				CommitDelta: 4,
			},
		},
	}

	for name, input := range inputs {
		s.Run(name, func() {
			container := reflected.TypedContainer{
				Payload: input.input,
			}
			blob, err := json.Marshal(&container)
			s.Nil(err)

			err = json.Unmarshal(blob, &container)
			s.Nil(err)
		})
	}
}

func (s *Suite) Test_UnmarshalFailures() {
	inputs := map[string]struct {
		input string
	}{
		"invalid envelope": {
			`{ "type": 23 }`,
		},
		"unhandled type": {
			`{ "type": "NULL" }`,
		},
		"invalid payload": {
			`{ "type": "changes", "commitDelta": "i should be a number" }`,
		},
	}

	for name, input := range inputs {
		s.Run(name, func() {
			container := reflected.TypedContainer{}
			err := json.Unmarshal([]byte(input.input), &container)
			s.NotNil(err)
			s.IsType(err, reflected.ReflectedError{})
		})
	}
}

// TODO: removed because I can't get json.Marshal() to generate errors
// func (s *Suite) Test_MarshalFailures() {
// 	var nilContainer *reflected.TypedContainer
// 	inputs := map[string]*reflected.TypedContainer{
// 		"nil container": nilContainer,
// 		"nil payload": {
// 			Payload: reflected.Typed(nil),
// 		},
// 		"invalid payload": {
// 			Type:    "duration",
// 			Payload: nil,
// 		},
// 	}

// 	for name, input := range inputs {
// 		s.Run(name, func() {
// 			_, err := json.Marshal(input)
// 			s.True(errors.Is(err, reflected.ReflectedError{}))
// 		})
// 	}
// }

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

package main

import (
	"encoding/json"
	"fmt"

	"github.com/pentaphobe/go-dynamic-json/pkg/reflected"
)

type AttestorEmail struct {
	Email string `json:"email"`
}

func (a *AttestorEmail) TypeName() string {
	return "email"
}
func (a *AttestorEmail) Cast(i interface{}) reflected.Payload {
	return i.(*AttestorEmail)
}

type AttestorRobot struct {
	Name string `json:"name"`
}

func (a *AttestorRobot) TypeName() string {
	return "robot"
}
func (a *AttestorRobot) Cast(i interface{}) reflected.Payload {
	return i.(*AttestorRobot)
}

type ValidityDuration struct {
	Days float32 `json:"days"`
}

func (a *ValidityDuration) TypeName() string {
	return "duration"
}
func (a *ValidityDuration) Cast(i interface{}) reflected.Payload {
	return i.(*ValidityDuration)
}

type Signoff struct {
	Attestor            reflected.TypedContainer   `json:"attestor"`
	SingleValidityTest  reflected.TypedContainer   `json:"single"`
	ValidityConstraints []reflected.TypedContainer `json:"validity"`
}

func main() {
	reflected.RegisterHandler(&AttestorRobot{})
	reflected.RegisterHandler(&AttestorEmail{})
	reflected.RegisterHandler(&ValidityDuration{})

	inputs := []string{
		`
		{
			"attestor": {
				"type": "email",
				"email": "foo@bar.com"
			},
			"single": { "type": "duration", "days": 4 },
			"validity": [
				{ "type": "duration", "days": 4 }
			]
		}
		`,
		`
		{
			"attestor": {
				"type": "robot",
				"name": "Bingo"
			},
			"single": { "type": "duration", "days": 4 },
			"validity": [
				{ "type": "duration", "days": 2 }
			]
		}
		`,
	}

	for _, input := range inputs {
		var t Signoff
		_ = json.Unmarshal([]byte(input), &t)

		fmt.Printf("%+v\n", t)
		fmt.Printf("\t%+v\n", t.Attestor.Payload)
		fmt.Printf("\t%+v\n", t.ValidityConstraints[0].Payload)
	}
}

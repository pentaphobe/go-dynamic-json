package validity

import (
	"time"

	"github.com/pentaphobe/go-dynamic-json/pkg/reflected"
)

type ValidityDuration struct {
	ValidTil time.Time `json:"validTil"`
}

func (v *ValidityDuration) TypeName() string {
	return "duration"
}

func (v *ValidityDuration) Cast(i interface{}) reflected.Payload {
	return i.(*ValidityDuration)
}

type ValidityChanges struct {
	CommitDelta int `json:"commitDelta"`
}

func (v *ValidityChanges) TypeName() string {
	return "changes"
}

func (v *ValidityChanges) Cast(i interface{}) reflected.Payload {
	return i.(*ValidityChanges)
}

func init() {
	reflected.RegisterHandler(&ValidityDuration{})
	reflected.RegisterHandler(&ValidityChanges{})
}
